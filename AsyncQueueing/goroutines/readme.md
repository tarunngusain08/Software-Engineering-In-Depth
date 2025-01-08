# Producer-Consumer Example in Go

This is a simple Go implementation of the producer-consumer pattern using goroutines and channels. It demonstrates how to use synchronization mechanisms like `sync.WaitGroup` and idiomatic channel handling with `range` for graceful termination.

## How It Works

1. **Producer:**
   - Generates integers from 0 to 9 and sends them to a shared channel (`dataChannel`).
   - Closes the channel after sending all values to signal that no more data will be sent.

2. **Consumer:**
   - Reads integers from the shared channel.
   - Prints each integer to the console.
   - Exits gracefully when the channel is closed.

3. **Synchronization:**
   - A `sync.WaitGroup` is used to ensure both the producer and consumer complete before the `main` function exits.

## Code Overview

### `main` Function
- Initializes the channel and `WaitGroup`.
- Starts the producer and consumer as separate goroutines.
- Waits for both goroutines to finish using the `WaitGroup`.

### `producer` Function
- Sends integers to the channel.
- Closes the channel after all integers are sent.

### `consumer` Function
- Uses a `for range` loop to read from the channel.
- Terminates gracefully when the channel is closed.

## Example Output
```
0 1 2 3 4 5 6 7 8 9
```

## Key Features
- Uses Go's idiomatic `range` to consume channel data.
- Clean synchronization with `sync.WaitGroup`.
- No need for sentinel values for termination.

## How to Run
1. Save the code in a file, e.g., `main.go`.
2. Run the program using:
   ```bash
   go run main.go
   ```