## Project Overview

* **Go API (beanbag-backend):** A backend API written in Go.
* **db:** A PostgreSQL Database.

## Prerequisites

Before you begin, ensure you have the following installed:

* **Docker:** Docker Desktop or the Docker Engine.
* **Docker Compose:** Included with Docker Desktop.
* **Git:** For version control.
* **Make:** To run the provided Makefile commands.
* **Go:** For working on the Go API.

## Local Development Setup

### Starting the Development Environment

1. Start all services using Docker Compose:

    ```bash
    make compose-up-build
    ```

    This command does the following:

    * Builds the Docker images for each service.
    * Starts all services as defined in the `docker-compose-dev.yml` file.

    The following services will be available:

      * **beanbag-backend:** Go API - accessible at `http://localhost:8080`.
      * **db:** PostgreSQL database.

### Stopping the Development Environment

1. Stop all services:

    ```bash
    make compose-down
    ```

## Running Tests

1. Run unit tests for the `beanbag-backend` service:

    ```bash
    make run-tests
    ```

    This command executes tests within Docker containers, using the test definitions from `docker-compose-test.yml`.

    Equivalent commands for debugging:
    ```bash
    docker compose -f docker-compose-dev.yml -f docker-compose-test.yml run --build beanbag-backend
    ```

## Debugging

1. Start the debugging environment:

    ```bash
    make compose-up-debug-build
    ```

    This command starts the services with debuggers enabled, as configured in `docker-compose-debug.yml`.

    Equivalent command for debugging:
    ```bash
    docker compose -f docker-compose-dev.yml -f docker-compose-debug.yml up --build
    ```

    * **Go API Debugging:**
        * The Go API starts with the `dlv` (Delve) debugger, listening on port `4000`.
        * Attach a Go debugger (e.g., in VS Code) to `localhost:4000`.

## Git Workflow

The project uses a standard Git workflow:

1. **Clone the Repository:**

2. **Create Branches:** Create a new branch for your changes:

    ```bash
    git checkout -b feature/my-new-feature
    ```

3. **Make Changes:** Develop and test your code.

4. **Stage Changes:** Add modified files:

    ```bash
    git add .
    ```

5. **Commit Changes:** Commit with a descriptive message:

    ```bash
    git commit -m "Implement my new feature"
    ```

6. **Push Changes:** Push your branch to the remote:

    ```bash
    git push -u origin feature/my-new-feature
    ```

7. **Create a Pull Request (PR):** Open a PR from your branch to `main` on GitHub.

8. **Code Review:** Address any feedback from code reviews.

9. **Merge:** After the code is reviewed, and the CI tests pass, merge your PR into main.
