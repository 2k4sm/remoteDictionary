# Project Proposal: CI/CD Pipeline for remoteDictionary Application

**Student Name:** Shrinibas Mahanta
**Scaler Student ID:** 23bcs10067
**Project Title:** CI/CD Pipeline for remoteDictionary Application

## 1. Application Description

- The **remoteDictionary** is a lightweight, high-performance distributed key-value store written in **Go (Golang)**.

- It supports RESTful API endpoints for setting (`POST /put`) and retrieving (`GET /get`) key-value pairs.

- The application features an in-memory LRU cache with configurable eviction policies based on memory usage, designed to handle high-throughput read/write operations efficiently.

**GitHub Repository:** `https://github.com/2k4sm/remoteDictionary`

## 2. Problem Statement

Manual deployments and lack of automated quality gates lead to:

- **Security Vulnerabilities**: Introduction of insecure dependencies or coding flaws.

- **Regressions**: Breaking changes deployed to production due to insufficient testing.

- **Inconsistent Environments**: "It works on my machine" issues.

- **Slow Feedback Loops**: Developers waiting too long to know if their code works.

The goal is to eliminate these risks by implementing a robust, automated pipeline that enforces quality, security, and reliability standards before any code reaches production.

## 3. Chosen CI/CD Stages & Justification

The pipeline is designed with a "Shift-Left" security mindset, ensuring issues are caught early.

| Stage                      | Tool            | Justification                                                                                                   |
| :------------------------- | :-------------- | :-------------------------------------------------------------------------------------------------------------- |
| **Linting**                | `golangci-lint` | Enforces Go coding standards and idiomatic usage, preventing technical debt accumulation.                       |
| **SAST (Static Analysis)** | `gosec`         | Scans source code for security flaws (e.g., hardcoded credentials, SQL injection risks) without executing code. |
| **SCA (Dependency Scan)**  | `trivy`         | Detects known vulnerabilities (CVEs) in third-party libraries (`go.mod`) to prevent supply-chain attacks.       |
| **Unit Testing**           | `go test`       | Validates business logic and ensures new changes do not break existing functionality (regression testing).      |
| **Build**                  | `go build`      | Ensures the application compiles correctly on a clean CI environment.                                           |
| **Containerization**       | `Docker`        | Packages the application and its runtime into an immutable artifact for consistent deployment.                  |
| **Image Scanning**         | `trivy`         | Scans the built Docker image for OS-level vulnerabilities (e.g., in Alpine Linux base) before registry push.    |
| **Runtime Smoke Test**     | `Docker + curl` | Simply validates that the container starts and listens on the expected port.                                    |
| **Registry Push**          | `Docker Hub`    | Stores the verified, secure artifact for consumption by the deployment pipeline.                                |

## 4. Expected Outcomes

- Fully automated build and test process on every commit.
- Zero critical vulnerabilities in deployed container images.
- Standardized code quality report.
- Published Docker image ready for Kubernetes deployment.
