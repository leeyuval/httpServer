# GitHub Repositories Collector - Program Design and Considerations

The GitHub Repositories Collector is a Go application designed to fetch and display GitHub repository information based on organization and query phrase. This document outlines the design of the program, highlighting key components and considerations.

## Program Design Overview

The program is designed using the following main components:

1. **Repositories Collector Interface (`ReposCollector`):** This interface defines the contract for any repository information collector. It specifies methods to configure the collector and handle HTTP requests to retrieve repositories.

2. **GitHub Repositories Collector (`GitHubReposCollector`):** This implementation of the `ReposCollector` interface is responsible for fetching and managing GitHub repository information. It uses the Gorilla Mux router to handle HTTP requests and Redis for caching fetched data.

3. **Utilities (`utils`):** This package contains utility functions used by the GitHubReposCollector, including functions for formatting creation time and rendering HTML templates.

4. **Main Program (`main.go`):** This is the entry point of the application. It configures the GitHubReposCollector, sets up the HTTP routes using Gorilla Mux, and starts the HTTP server.

5. **Docker Configuration (`Dockerfile`, `docker-compose.yml`):** These files define the Docker configuration for building and running the application along with its dependencies (Redis).

## Considerations and Design Decisions

### 1. Redis Caching:

The program utilizes Redis for caching fetched GitHub repository data. Caching helps reduce the load on the GitHub API and improves response times for frequently requested data. Cached data is stored with a predefined cache key based on the organization and query phrase, and it expires after a certain period (12 hours in this case).

### 2. Pagination:

GitHub API responses can be paginated, so the program includes pagination support. The `paginate` function splits the fetched data into pages of a predefined size (`PerPage` constant) to allow users to navigate through the results more easily.

### 3. Error Handling:

The program incorporates robust error handling. It handles errors related to fetching data from the GitHub API, decoding JSON responses, and interacting with Redis. Proper HTTP status codes and error messages are returned to the client in case of errors.

### 4. Dockerization:

The program includes Docker files (`Dockerfile` and `docker-compose.yml`) to simplify setup and deployment. This approach ensures consistent environments across different systems and avoids potential compatibility issues.

### 5. Separation of Concerns:

The program adheres to the principle of separation of concerns. The `GitHubReposCollector` focuses on fetching and managing GitHub repository data, while the `utils` package contains utility functions. This separation enhances code organization and maintainability.

### 6. Configuration and Extensibility:

The use of interfaces and a well-defined configuration approach (`ConfigureCollector` method) makes the program extensible. It allows for the addition of new repository collectors or data sources by implementing the `ReposCollector` interface and configuring them similarly.

### 7. Testing:

The program design encourages testability. Unit tests can be written for individual functions, and integration tests can be conducted to ensure that different components work together as expected.

### 8. Route Definitions:

The program defines two main routes for fetching repositories: one with only the organization name and another with both the organization name and a query phrase. This design provides flexibility for users to search repositories based on different criteria.

## Conclusion

The GitHub Repositories Collector is designed with modularity, caching, error handling, and extensibility in mind. It leverages the power of Gorilla Mux, Redis caching, and Docker to provide a scalable and efficient solution for fetching and displaying GitHub repository information. The considerations outlined above ensure a robust and maintainable application.