# Envoy Gateway Setup on Minikube

This guide covers setting up Envoy Gateway on Minikube, configuring Load Balancing, and implementing Rate Limiting using Redis.

## Prerequisites
- **Minikube** installed (`minikube start` to launch a cluster)
- **Kubectl** installed
- **Helm** installed

---

## 1️⃣ Install Envoy Gateway
We use Helm to install Envoy Gateway in the `envoy-gateway-system` namespace.
```sh
helm install eg oci://docker.io/envoyproxy/gateway-helm --version v0.0.0-latest -n envoy-gateway-system --create-namespace
kubectl wait --timeout=5m -n envoy-gateway-system deployment/envoy-gateway --for=condition=Available
```

Verify the installation:
```sh
kubectl get pods -n envoy-gateway-system
kubectl get svc -n envoy-gateway-system
```

---

## 2️⃣ Apply Quickstart Configuration
```sh
kubectl apply -f https://github.com/envoyproxy/gateway/releases/download/latest/quickstart.yaml -n default
```
Check resources:
```sh
kubectl get gatewayclass,gateway,httproute -A
```

Expose Envoy Gateway:
```sh
kubectl port-forward -n envoy-gateway-system svc/envoy-default-eg-XXXXX 8080:80
```
Test:
```sh
curl -H "Host: www.example.com" http://127.0.0.1:8080/
```

---

## 3️⃣ Enable Load Balancing
Scale backend service to multiple replicas:
```sh
kubectl scale deployment backend --replicas=3
kubectl get pods -n default -l app=backend
kubectl get endpoints -n default backend
```
Test load balancing:
```sh
for i in {1..10}; do curl -H "Host: www.example.com" http://127.0.0.1:80/; echo ""; done
```

Modify `envoy-config.yaml` to customize Load Balancing:
```yaml
clusters:
  - name: backend_service
    load_balancing_policy: ROUND_ROBIN
```
Apply and restart:
```sh
kubectl apply -f envoy-config.yaml
kubectl rollout restart deployment/envoy-gateway -n envoy-gateway-system
```

---

## 4️⃣ Implement Rate Limiting
### **Step 1: Deploy Redis**
```sh
kubectl apply -f https://raw.githubusercontent.com/helm/charts/main/stable/redis/templates/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/helm/charts/main/stable/redis/templates/service.yaml
```

### **Step 2: Deploy Rate Limit Service**
Create `ratelimit.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ratelimit
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ratelimit
  template:
    metadata:
      labels:
        app: ratelimit
    spec:
      containers:
      - name: ratelimit
        image: envoyproxy/ratelimit:latest
        ports:
        - containerPort: 8081
        env:
        - name: REDIS_SOCKET_TYPE
          value: "tcp"
        - name: REDIS_URL
          value: "redis:6379"
```
Apply:
```sh
kubectl apply -f ratelimit.yaml
```

### **Step 3: Define Rate Limit Rules**
Create `ratelimit-config.yaml`:
```yaml
domain: "default"
descriptors:
  - key: "generic_key"
    value: "limited"
    rate_limit:
      unit: second
      requests_per_unit: 5
```
Apply:
```sh
kubectl apply -f ratelimit-config.yaml
```

### **Step 4: Modify Envoy Configuration**
Update `envoy-config.yaml`:
```yaml
http_filters:
- name: envoy.filters.http.ratelimit
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit
    domain: "default"
    failure_mode_deny: true
    rate_limit_service:
      grpc_service:
        envoy_grpc:
          cluster_name: ratelimit_service
```
Apply and restart:
```sh
kubectl apply -f envoy-config.yaml
kubectl rollout restart deployment/envoy-gateway -n envoy-gateway-system
```

### **Step 5: Test Rate Limiting**
```sh
for i in {1..10}; do curl -H "Host: www.example.com" http://127.0.0.1:80/; echo ""; done
```
Expected output: 
✔️ First 5 requests pass (200 OK) 
✔️ Next requests are blocked (429 Too Many Requests)

---

## 🎯 Next Steps
- Implement **JWT-based rate limiting**
- Add **mTLS security**
- Integrate **Grafana/Prometheus for monitoring**
