# GitHub Repositories Collector

The GitHub Repositories Collector is a Go application that fetches and displays GitHub repository information based on organization and query phrase. 
It uses the Gorilla Mux router for handling HTTP requests and Redis for caching the fetched data.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Caching](#caching)

## Features

- Fetches GitHub repositories based on organization and query phrase.
- Displays repository information including name, owner, URL, creation time, and stars.
- Provides pagination for browsing repositories.
- Caches fetched data using Redis for improved performance.

## Project Structure

The project is organized as follows:

- `cacheDBs`: Contains code related to caching data using Redis.
- `httpServer_tests`: Test files for the HTTP server.
- `repositoriesCollectors`: Implements the repositories information collector.
- `templates`: HTML templates for rendering repository information.
- `utils`: Utility functions and helpers.
- `docker-compose.yml`: Docker Compose configuration for setting up Redis and the application.
- `Dockerfile`: Docker configuration for building and running the application.
- `go.mod` and `go.sum`: Go module files.

## Getting Started

### Prerequisites

Before running the application, ensure you have the following dependencies:

- [Docker](https://www.docker.com/get-started)
- [Go] (https://go.dev/dl/)

### Running the Application with Docker (Recommended)

1. Clone this repository:

   ```sh
   git clone https://github.com/leeyuval/repositoriesCollector.git
   ```

2. Navigate to the project directory:

   ```sh
   cd repositoriesCollector
   ```

3. Start the application using Docker Compose:

   ```sh
   docker-compose up -d
   ```

   This will start a Redis instance and the Go application.

4. Access the application in your browser:

   Open a web browser and go to [http://localhost:8080](http://localhost:8080) to access the application.
   
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

The application provides a web interface to search and display GitHub repositories. 
You can enter an organization name and an optional search phrase to filter repositories. 
The repositories are displayed with their details, and pagination is available at the bottom of the page.

## API Endpoints

The application provides the following API endpoints:

- `GET /repositories/org/{org}`: Fetches repositories for the specified organization.
- `GET /repositories/org/{org}/q/{q}`: Fetches repositories for the specified organization with the given search phrase.

## Caching

The application uses Redis as a caching mechanism to store fetched data. 
Cached data is stored for 12 hours (configurable in `redisDB.go`). 
This improves performance by reducing the number of requests to the GitHub API.
