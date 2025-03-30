# Remote Dictionary

A lightweight, high-performance, in-memory key-value store service with a clean HTTP API interface.

## Installation

### Prerequisites
- Go 1.23 or higher

### From Source
```bash
# Clone the repository
git clone https://github.com/2k4sm/remoteDictionary.git
cd remoteDictionary

# Build the binary
go build -o remoteDictionary

# Run the server
./remoteDictionary
```

### Using Go Install
```bash
go install github.com/2k4sm/remoteDictionary@latest

# Run the server
./remoteDictionary
```

### Using Docker
```bash
# Build using
docker build -t remotedictionary .
```

```bash
# Run using
docker run -p 7171:7171 remotedictionary
```

### You can also pull the image from dockerhub and use it.
```bash
# Pull the image using
docker pull sm2k4/remotedictionary:main

# Run it using
docker run -p 7171:7171  sm2k4/remotedictionary:main
```

## Configuration
Given below are the default values but you can change it by creating a `.env` file in the project root:
```
PORT=7171
MAX_KEY_SIZE=256
MAX_VALUE_SIZE=256
```

## API Usage

### Store a Value
```
POST /put
```

Request body:
```json
{
  "key": "user:1234",
  "value": "John Doe"
}
```

Response:
```json
{
  "status": "OK",
  "message": "Key inserted/updated successfully."
}
```

### Retrieve a Value
```
GET /get?key=user:1234
```

Success Response:
```json
{
  "status": "OK",
  "key": "user:1234",
  "value": "John Doe",
  "message": "Key retrieved successfully."
}
```

Key Not Found Response:
```json
{
  "status": "ERROR",
  "message": "Key not found."
}
```

## Design Choices and Optimizations

### Memory Efficiency

- **LRU Cleanup**: We remove the oldest, least-used items first when memory gets tight. This keeps only the stuff you're actually using.

- **Usage Monitor**: The system keeps an eye on memory usage and automatically cleans up when it reaches 70% full, so it never crashes from running out of memory.

- **Gradual Cleanup**: When clearing memory, we start by removing just a few items at a time, then gradually remove more if needed. This keeps the system responsive while freeing up space.

### Speed Enhancements

- **Buffer Reuse**: We reuse memory buffers instead of creating new ones each time, which makes everything run faster during heavy usage.

- **Concurrent Access**: The system allows multiple users to read data at the same time, while making sure updates happen safely without conflicts.

- **Request Limits**: We put caps on request sizes to prevent anyone from crashing the system by sending extremely large requests.

### Growth & Expansion

- **Horizontal Scaling**: The core service can be easily replicated across multiple servers with a load balancer since it doesn't rely on shared state (except for the cache itself).

- **Workload Adaptation**: The system constantly monitors resource usage and adjusts to handle changing workloads.

Built by [2k4sm](https://github.com/2k4sm)
