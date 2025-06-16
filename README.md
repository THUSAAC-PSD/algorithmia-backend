# Algorithmia Backend

Algorithmia Backend is the server-side application for the Algorithmia platform, designed to manage CP contests, problems, users, and related functionalities. It provides a RESTful API and WebSocket support for real-time communication.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Clone the Repository](#clone-the-repository)
  - [Install Tools](#install-tools)
  - [Install Dependencies](#install-dependencies)
  - [Configuration](#configuration)
    - [JSON Configuration Files](#json-configuration-files)
    - [The .env file](#the-env-file)
  - [Running the Application](#running-the-application)
- [Makefile Targets](#makefile-targets)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
- [Code Quality](#code-quality)
- [License](#license)

## Features

*   **User Management:**
    *   User registration with email verification.
    *   User login and session management.
    *   Get current user profile.
    *   Role-based access control (RBAC) with permissions.
*   **Contest Management:**
    *   Create, list, and delete contests.
    *   Define problem count limits for contests.
    *   Set contest deadlines.
    *   Assign/unassign problems to/from contests.
*   **Problem Management:**
    *   **Problem Drafts:** Create, update, list, and delete problem drafts with multi-language support for details (title, background, statement, etc.) and examples.
    *   **Problem Submission:** Submit drafts for review.
    *   **Problem Lifecycle:**
        *   Review (approve, reject, needs revision).
        *   Assign testers.
        *   Testing (passed, failed).
        *   Mark as complete.
    *   **Problem Details:** View problem versions, details, examples, reviews, and test results.
    *   **Problem Chat:** Real-time WebSocket-based chat for discussing problems, including notifications for submissions, reviews, tests, and completions.
*   **Problem Difficulty:** Manage and list problem difficulties with multi-language display names.
*   **Media Management:** Support for uploading media related to problem drafts and chat messages.

## Tech Stack

*   **Language:** Go
*   **Web Framework:** [Echo](https://echo.labstack.com/)
*   **ORM:** [GORM](https://gorm.io/)
*   **Database:** PostgreSQL
*   **Configuration:** [Viper](https://github.com/spf13/viper) (JSON files & environment variables)
*   **Logging:** [Zap](https://github.com/uber-go/zap)
*   **Command-line Interface:** [Cobra](https://github.com/spf13/cobra)
*   **Dependency Injection:** [Dig](https://github.com/uber-go/dig)
*   **Real-time Communication:** WebSockets
*   **Containerization:** Docker, Docker Compose
*   **Development Tooling:**
    *   Node.js (for Husky and Commitlint)
    *   Make / [Task](https://taskfile.dev/)
    *   Linters: `golangci-lint`, `revive`, `staticcheck`
    *   Formatters: `gofmt`, `goimports`

## Prerequisites

*   Go (version specified in `go.mod`, typically latest stable)
*   Docker and Docker Compose
*   Node.js and npm (for commit hooks and linting setup)
*   (Optional) Task (`taskfile.dev`) if you prefer it over Make.

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/THUSAAC-PSD/algorithmia-backend.git
cd algorithmia-backend
```

### Install Tools

This project uses various Go tools for development, linting, and formatting. Install them by running:

```bash
make install-tools
# or
task install-tools
```

This will execute `./scripts/install-tools.sh`.

If you plan to commit code, ensure Node.js and npm are installed, then run:
```bash
npm install
```
This sets up Husky and Commitlint for commit message linting.

### Install Dependencies

Install Go module dependencies:

```bash
make install-dependencies
# or
task install-dependencies
```
This executes `./scripts/install-dependencies.sh`, which runs `go mod tidy`.

### Configuration

The application uses a combination of JSON configuration files and environment variables. Viper is used to manage configuration, with environment variables taking precedence.

#### JSON Configuration Files

Configuration files are located in the `config/` directory.

1.  **Copy the example configuration:**
    ```bash
    cp config/config.development.json.example config/config.development.json
    ```

2.  **Edit `config/config.development.json`:**
    *   Update `gormOptions` if your PostgreSQL setup differs from the default (localhost:5432, user: postgres, pass: postgres, db: algorithmia).
    *   Set a `sessionSecret` under `echoHttpOptions`.
        This is crucial for session security. To generate a strong secret, use a cryptographically secure random string generator (e.g., `openssl rand -base64 32` or `head -c 32 /dev/urandom | base64` on Linux/macOS). Aim for at least 32 bytes of randomness (which becomes a longer Base64 string).
    *   Configure `gomailOptions` if you need email sending functionality (e.g., for email verification).

#### The .env file

1.  **Create a `.env` file** in the project root directory (e.g., `algorithmia-backend/.env`).
2.  **Add your environment variables.** For example:

    ```env
    # .env
    APP_ENV=development # Application Environment (development, test, production). This determines which config file to use (config.*.json)
    PROJECT_NAME=algorithmia-backend # Ensure this matches the root folder name
    ```

### Running the Application

The recommended way to run the application is with Docker Compose, which manages both the Go application container and the PostgreSQL database container.

1.  **Build and Run with Docker Compose:**
    Ensure Docker is running, then execute the following command from the project root:
    ```bash
    docker-compose up --build -d
    ```
    * `--build`: This flag tells Docker Compose to build the application image from the `Dockerfile` before starting the services.
    * `-d`: This runs the containers in detached mode (in the background).

    The application will now be running. The Go app's port (e.g., 9090) will be mapped to the same port on your host machine, ready to receive requests from a reverse proxy like Nginx.

2.  **Stopping the Application:**
    To stop both the application and the database containers:
    ```bash
    docker-compose down
    ```

## Makefile Targets

The `Makefile` provides several useful targets for development:

*   `install-tools`: Installs necessary Go tools.
*   `run-app`: Runs the application locally.
*   `build`: Builds the application binary.
*   `install-dependencies`: Installs Go module dependencies.
*   `format`: Formats the Go codebase.
*   `lint`: Runs linters (`golangci-lint`, `revive`, `staticcheck`).
*   `update-dependencies`: Updates Go dependencies.

A `taskfile.yml` is also provided with similar commands that can be run using `task <task-name>`.

## Project Structure

The project follows a modular structure, primarily within the `internal/` directory. This structure is designed to separate concerns and promote maintainability.

```
algorithmia-backend/
├── Makefile                  # Main build and task runner
├── Dockerfile                # Instructions to build the application container
├── docker-compose.yml        # Defines and runs the multi-container setup
├── README.md                 # This file
├── cmd/
│   └── app/
│       └── main.go           # Application entry point
├── config/
│   └── config.development.json.example # Example configuration
├── deployments/
│   └── docker-compose/
│       └── docker-compose.infrastructure.yaml # Docker Compose for PostgreSQL
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── golangci.yml              # GolangCI-Lint configuration
├── revive-config.toml        # Revive linter configuration
├── staticcheck.conf          # Staticcheck linter configuration
├── taskfile.yml              # Alternative task runner (Task)
├── package.json              # Node.js dependencies (for dev tools)
├── package-lock.json         # Node.js lock file
├── scripts/                  # Utility shell scripts for Makefile/Task
│   ├── build.sh
│   ├── format.sh
│   ├── install-dependencies.sh
│   ├── install-tools.sh
│   ├── lint.sh
│   ├── run.sh
│   └── update-dependencies.sh
└── internal/                 # Core application logic (not for external import)
    ├── contest/              # Contest module
    │   ├── feature/          # Feature-sliced (CQRS-like) sub-packages
    │   │   ├── assignproblem/
    │   │   ├── createcontest/
    │   │   ├── ...           # (e.g., command.go, handler.go, endpoint.go, repository.go)
    │   └── endpoint_params.go # Shared Echo group parameters for contest endpoints
    ├── problem/              # Problem module (similar structure to contest)
    ├── problemdifficulty/    # Problem Difficulty module
    ├── problemdraft/         # Problem Draft module
    ├── user/                 # User module
    │   ├── feature/
    │   │   ├── login/
    │   │   ├── register/
    │   │   └── ...
    │   ├── constant/         # User-specific constants
    │   ├── infrastructure/   # User-specific infrastructure (e.g., password hasher)
    │   └── endpoint_params.go
    └── pkg/                  # Shared internal packages (utilities, core infrastructure)
        ├── app/              # Application bootstrapping and core
        │   ├── application/      # Application struct, lifecycle, DI resolution
        │   └── applicationbuilder/ # Builder for application setup
        ├── config/           # Configuration loading helpers
        ├── constant/         # Global constants (permissions, statuses, etc.)
        ├── contract/         # Interfaces defining contracts between components
        ├── customerror/      # Custom error types
        ├── database/         # GORM models, DB connection, Unit of Work
        ├── environment/      # Environment handling
        ├── http/             # HTTP related utilities
        │   ├── echoweb/      # Echo framework setup, middleware, helpers
        │   └── httperror/    # Custom HTTP error handling
        ├── logger/           # Logging setup and interface (Zap)
        ├── mailing/          # Email sending setup (Gomail)
        ├── reflection/       # Reflection utilities
        ├── tzinit/           # Timezone initialization (sets to UTC)
        └── websocket/        # WebSocket hub, client, protocol, broadcaster
```

**Key Principles:**

*   **Modularity:** Features are organized into modules (e.g., `user`, `contest`, `problem`).
*   **Feature Slicing (CQRS-like):** Within each module, features (e.g., `registeruser`, `createcontest`) are often organized into their own sub-packages. These typically contain:
    *   `command.go` / `query.go`: Defines the input data structure for the operation.
    *   `command_handler.go` / `query_handler.go`: Contains the business logic for the operation.
    *   `endpoint.go`: Maps the HTTP/WebSocket route and handles request/response binding.
    *   `gorm_repository.go`: Implements the data access logic for the specific feature using GORM.
    *   `response.go`: Defines the output data structure.
*   **Dependency Injection:** The application uses `go.uber.org/dig` for dependency injection, managed primarily in `internal/pkg/app/applicationbuilder/` and `internal/pkg/app/application/`.
*   **Shared Packages (`internal/pkg/`):** Common utilities, infrastructure setup (database, HTTP, logging, etc.), and core contracts are placed here.
*   **Clear Entry Point:** `cmd/app/main.go` is the single entry point that initializes and runs the application.
*   **Configuration Separation:** Configuration (`config/`) and deployment scripts (`deployments/`) are kept separate from application code.
*   **Scripts:** Repetitive tasks are automated via `Makefile` and shell scripts in `scripts/`.

**Contributing New Features:**

1.  **Identify the Module:** Determine which existing module your feature belongs to (e.g., `user`, `problem`). If it's a new domain, create a new directory under `internal/`.
2.  **Create a Feature Package:** Inside the module, create a new directory for your feature (e.g., `internal/user/feature/updateprofile/`).
3.  **Define Command/Query & Response:** Create `command.go` (or `query.go`) and `response.go` to define the data structures.
4.  **Implement Handler:** Create `command_handler.go` (or `query_handler.go`) with the business logic. Define a repository interface here for data access.
5.  **Implement Repository:** Create `gorm_repository.go` (or similar) to implement the repository interface using GORM.
6.  **Create Endpoint:** Create `endpoint.go` to handle HTTP/WebSocket requests, bind data, call the handler, and format the response.
7.  **Register Dependencies:**
    *   Add your handler, repository implementation, and endpoint to the dependency injection container in `internal/pkg/app/applicationbuilder/application_builder_features.go` and `internal/pkg/app/application/application_handlers.go`.
    *   If you introduce new infrastructure components (e.g., a new type of mailer), register them in `internal/pkg/app/applicationbuilder/application_builder_infrastructure.go`.
8.  **Add Database Migrations (if needed):** Update `internal/pkg/app/application/application_infrastructure.go` in the `migrateDatabase` function if you add or change GORM models.
9.  **Add Tests:** (None yet, but crucial) Write unit and/or integration tests for your new feature.
10. **Update Documentation:** If your feature adds new API endpoints or significantly changes behavior, update the APIDog documentation and/or this README file.

## API Overview

The backend exposes a RESTful API, primarily under the `/api/v1/` prefix. Key resource groups include:

*   `/api/v1/auth`: Authentication-related endpoints (register, login, logout, email verification).
*   `/api/v1/users`: User-related endpoints (e.g., get current user).
*   `/api/v1/contests`: Contest management.
*   `/api/v1/problems`: Problem management.
*   `/api/v1/problem-drafts`: Problem draft management.
*   `/api/v1/problem-difficulties`: Problem difficulty listing.
*   `/api/v1/testers`: List users who can be testers.
*   `/api/v1/ws/chat`: WebSocket endpoint for real-time problem chat.

Refer to the `internal/**/endpoint.go` files for specific route definitions and handlers.

## Code Quality

*   **Formatting:** Code is formatted using `gofmt` and `goimports`. Run `make format`.
*   **Linting:** The project uses a combination of linters:
    *   `golangci-lint` (configured in `golangci.yml`)
    *   `revive` (configured in `revive-config.toml`)
    *   `staticcheck` (configured in `staticcheck.conf`)
    Run `make lint` to check the codebase.
*   **Commit Messages:** Commit messages are linted using `commitlint` with the `config-conventional` standard, enforced by Husky git hooks.

## License

Except as otherwise noted in individual files, this project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.
