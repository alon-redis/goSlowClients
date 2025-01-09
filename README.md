# differentBufferSize

## Overview
`differentBufferSize` is a Go script that interacts with a Redis server using both connection pooling and TCP sockets. The script supports both secured TLS and non-TLS connections and has two primary stages:

1. **Population Stage:**
   - Uses the Redis Go client to populate the Redis database with keys of varying sizes.
   - The key size starts at an initial size and increases incrementally based on the specified delta.

2. **Fetch Stage:**
   - Opens multiple parallel TCP socket connections to fetch the data.
   - Simulates slow responses by introducing a sleep time during data fetching.

## Usage
### Command-line Parameters
```bash
./differentBufferSize <redis_host> <redis_port> <num_connections> <initial_key_size_MB> <delta_MB> <sleep_time_seconds> <noflush> <use_tls>
```

### Parameters:
- **redis_host:** Redis server hostname.
- **redis_port:** Redis server port.
- **num_connections:** Number of parallel connections to use.
- **initial_key_size_MB:** Initial key size in megabytes.
- **delta_MB:** Incremental increase in key size per connection in megabytes.
- **sleep_time_seconds:** Time to sleep between sending commands during the fetch stage.
- **noflush:** Prevents flushing the Redis database before starting if set to `true`.
- **use_tls:** Enables secured TLS connections to the Redis server if set to `true`.

### Example Command:
```bash
./differentBufferSize 127.0.0.1 6379 10 1 1 2 false true
```

## Installation
1. **Install Go**:
   - Follow the instructions [here](https://go.dev/doc/install) or use:
     ```bash
     sudo apt update
     sudo apt install golang-go -y
     ```
2. **Install Redis**:
   - Install Redis server:
     ```bash
     sudo apt install redis-server -y
     ```
   - Start Redis:
     ```bash
     sudo systemctl start redis
     ```

3. **Set Up the Project**:
   - Create the project directory:
     ```bash
     mkdir differentBufferSize
     cd differentBufferSize
     ```
   - Initialize the Go module:
     ```bash
     go mod init differentBufferSize
     ```
   - Add the Redis Go client dependency:
     ```bash
     go get github.com/go-redis/redis/v8
     ```

## Build and Run the Script
1. **Build the binary**:
   ```bash
   go build -o differentBufferSize
   ```
2. **Run the script**:
   ```bash
   ./differentBufferSize <parameters>
   ```

## Additional Tips
- Ensure the Redis server is running before executing the script.
- Adjust the parameters as needed for your test scenarios.
- Use `true` for the `use_tls` parameter if connecting to a Redis server with TLS enabled.
- Use `false` for the `noflush` parameter if you want to flush the Redis database before the script starts.

