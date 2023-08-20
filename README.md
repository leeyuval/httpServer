# GitHub Repositories Collector

The GitHub Repositories Collector is a Go application that fetches and displays GitHub repository information based on organization and query phrase. It uses the Gorilla Mux router for handling HTTP requests and Redis for caching the fetched data.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Configuration](#configuration)

## Prerequisites

Before running this application, you need to have the following software installed:

- Docker (optional, if you're using Docker - recommended)
- Go
- Redis

## Getting Started

Follow these steps to get the application up and running:

### Running with Docker

1. Make sure you have Docker installed.
2. Clone this repository to your local machine.
3. Open a terminal and navigate to the root directory of the cloned repository.
4. Run the following command to build and start the application along with Redis:

```bash
docker-compose up
```

### Running Locally

1. Clone this repository to your local machine.
2. Install Go and Redis if you haven't already.
3. Open a terminal and navigate to the root directory of the cloned repository.
4. Run the following command to build the application:

```bash
go build -o main .
```

5. Start the Redis server.
6. Run the application:

```bash
./main
```

The application will be available at `http://localhost:8080`.

## Usage

The application allows you to fetch and display GitHub repository information based on organization and query phrase.

You can access the endpoints using the following URLs:

- `http://localhost:8080/repositories/org/{org}`: Fetches and displays repositories for a specific organization.
- `http://localhost:8080/repositories/org/{org}/q/{q}`: Fetches and displays repositories for a specific organization with a query phrase.

Replace `{org}` with the organization name and `{q}` with the query phrase.

## Endpoints

- `GET /repositories/org/{org}`: Fetches and displays repositories for a specific organization.
- `GET /repositories/org/{org}/q/{q}`: Fetches and displays repositories for a specific organization with a query phrase.

## Configuration

The application's configuration can be found in the `main.go` file. You can configure the Redis address and other settings there.

---