# Remote Dictionary

A lightweight, high-performance, in-memory key-value store service with a clean HTTP API interface.

## How to Run Locally

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
docker pull sm2k4/remotedictionary:latest

# Run it using
docker run -p 7171:7171  sm2k4/remotedictionary:latest
```

## Configuration

Given below are the default values but you can change it by creating a `.env` file in the project root:

```
PORT=7171
MAX_KEY_SIZE=256
MAX_VALUE_SIZE=256
```

### Secrets Configuration (CI/CD)

To enable the GitHub Actions pipeline, configure the following **Repository Secrets**:

| Secret Name                      | Description                                               |
| :------------------------------- | :-------------------------------------------------------- |
| `GCP_PROJECT_ID`                 | Your Google Cloud Project ID.                             |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | The full identifier of the Workload Identity Provider.    |
| `GCP_SERVICE_ACCOUNT`            | The email address of the Google Service Account.          |
| `GKE_CLUSTER_NAME`               | The name of your GKE cluster.                             |
| `GKE_ZONE`                       | The compute zone of your cluster (e.g., `us-central1-a`). |
| `DOCKER_USERNAME`                | Your Docker Hub username.                                 |
| `DOCKER_PASSWORD`                | Your Docker Hub access token or password.                 |

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

## CI/CD Pipeline

The project utilizes a robust CI/CD pipeline powered by GitHub Actions.

### Continuous Integration (CI)

Run on every push to `main` and Pull Requests:

1.  **Linting**: Uses `golangci-lint` to ensure code quality.
2.  **Security Scanning**:
    - **SAST**: `gosec` scans for security vulnerabilities in code.
    - **SCA**: `trivy` scans dependencies for known vulnerabilities.
3.  **Unit Tests**: Runs `go test -race` to verification logic and concurrency.
4.  **Build**: Verifies the application compiles successfully.

### Continuous Delivery (CD)

Triggered on successful runs on the `main` branch:

1.  **Container Build**: Builds the Docker image.
2.  **Container Scan**: Scans the image with `trivy` for OS-level vulnerabilities.
3.  **Deployment**: Pushes the image to Docker Hub and deploys to the configured GKE cluster.
4.  **Dynamic Scan**: Runs OWASP ZAP (DAST) against the live deployment to check for runtime security issues.

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
