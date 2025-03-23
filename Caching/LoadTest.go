package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests int
	Successful    int
	Failed        int
	TotalLatency  time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	mu            sync.Mutex
}

func (m *Metrics) RecordSuccess(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
	m.Successful++
	m.TotalLatency += latency
	if m.MinLatency == 0 || latency < m.MinLatency {
		m.MinLatency = latency
	}
	if latency > m.MaxLatency {
		m.MaxLatency = latency
	}
}

func (m *Metrics) RecordFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
	m.Failed++
}

func (m *Metrics) AverageLatency() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Successful == 0 {
		return 0
	}
	return m.TotalLatency / time.Duration(m.Successful)
}

func generateKeyValuePairs(count int) map[string]string {
	keyValues := make(map[string]string, count)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		keyValues[key] = value
	}
	return keyValues
}

func performSetOperations(keyValues map[string]string, rps int, duration time.Duration, metrics *Metrics) {
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()

	endTime := time.Now().Add(duration)
	var wg sync.WaitGroup

	for key, value := range keyValues {
		if time.Now().After(endTime) {
			break
		}
		<-ticker.C
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			start := time.Now()
			data := map[string]string{k: v}
			body, _ := json.Marshal(data)
			resp, err := http.Post("http://localhost:8080/set", "application/json", bytes.NewBuffer(body))
			latency := time.Since(start)
			if err != nil {
				fmt.Printf("Failed to set key %s: %v\n", k, err)
				metrics.RecordFailure()
				return
			}
			resp.Body.Close()
			metrics.RecordSuccess(latency)
		}(key, value)
	}
	wg.Wait()
}

func performGetOperations(keys []string, rps int, duration time.Duration, metrics *Metrics) {
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()

	endTime := time.Now().Add(duration)
	var wg sync.WaitGroup

	for _, key := range keys {
		if time.Now().After(endTime) {
			break
		}
		<-ticker.C
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			start := time.Now()
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/get?key=%s", k))
			latency := time.Since(start)
			if err != nil {
				fmt.Printf("Failed to get key %s: %v\n", k, err)
				metrics.RecordFailure()
				return
			}
			resp.Body.Close()
			metrics.RecordSuccess(latency)
		}(key)
	}
	wg.Wait()
}

func main() {
	keyValues := generateKeyValuePairs(100000)
	keys := make([]string, 0, len(keyValues))
	for key := range keyValues {
		keys = append(keys, key)
	}

	setMetrics := &Metrics{}
	getMetrics := &Metrics{}

	fmt.Println("Starting SET operations...")
	performSetOperations(keyValues, 1200, 10*time.Second, setMetrics)
	fmt.Println("Completed SET operations.")
	fmt.Printf("SET Metrics: Total Requests: %d, Successful: %d, Failed: %d, Average Latency: %v, Min Latency: %v, Max Latency: %v\n",
		setMetrics.TotalRequests, setMetrics.Successful, setMetrics.Failed, setMetrics.AverageLatency(), setMetrics.MinLatency, setMetrics.MaxLatency)

	fmt.Println("Starting GET operations...")
	performGetOperations(keys, 1000, 10*time.Second, getMetrics)
	fmt.Println("Completed GET operations.")
	fmt.Printf("GET Metrics: Total Requests: %d, Successful: %d, Failed: %d, Average Latency: %v, Min Latency: %v, Max Latency: %v\n",
		getMetrics.TotalRequests, getMetrics.Successful, getMetrics.Failed, getMetrics.AverageLatency(), getMetrics.MinLatency, getMetrics.MaxLatency)
}
