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

## Configuration
Given below are the default values but you can change it by creating a `.env` file in the project root:
```
PORT=7171
MAX_KEY_SIZE=256
MAX_VALUE_SIZE=256
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

## You can also pull the image from dockerhub and use it.
```bash
# Pull the image using
docker pull sm2k4/remotedictionary:main

# Run it using
docker run -p 7171:7171  sm2k4/remotedictionary:main
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

Built by [2k4sm](https://github.com/2k4sm)
