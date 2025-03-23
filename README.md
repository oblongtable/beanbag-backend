## Project Overview

* **Go API (beanbag-backend):** A backend API written in Go.
* **db:** A PostgreSQL Database.

## Prerequisites

Before you begin, ensure you have the following installed:

* **Docker:** Docker Desktop or the Docker Engine.
* **Docker Compose:** Included with Docker Desktop.
* **Git:** For version control.
* **Make:** To run the provided Makefile commands.
* **Go:** For doing `go get` and stopping dev env from complaining.
* **Goose:** For database migrations. Install with `go install github.com/pressly/goose/v3/cmd/goose@latest`

Extra recommended nice-to-haves:

* **VSCode**
* **WSL** w/ Ubuntu if developing under Windows
* **The following VSCode Extensions** 
  * [Remote - SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh) (browse/edit your WSL files from host machine)
  * [DB Code](https://marketplace.visualstudio.com/items?itemName=dbcode.dbcode) - Connect to and browse DB from VSCode
  * [Postman](https://marketplace.visualstudio.com/items?itemName=Postman.postman-for-vscode) - Make REST API calls from inside VSCode
  * [Go](https://marketplace.visualstudio.com/items?itemName=golang.go) - Rich Golang support for VSCode (autocomplete, etc.)
  * [Makefile tools](https://marketplace.visualstudio.com/items?itemName=ms-vscode.makefile-tools) - Makefile support for VSCode (autocomplete, etc.)
  * [Gemini Code Assist](https://marketplace.visualstudio.com/items?itemName=google.geminicodeassist) - Google+StackOverflow, but better

## Local Development Setup

### Starting the Development Environment

1. Start all services using Docker Compose:

    NOTE: Before doing the below, you need a file called `app.env` in the root directory of this project that has the format:

    ```
    POSTGRES_HOST=db
    POSTGRES_USER=<username>
    POSTGRES_PASSWORD=<password>
    POSTGRES_DB=beanbag-backend-db
    POSTGRES_PORT=5432

    PORT=8080
    CLIENT_ORIGIN=http://localhost:3000
    ```

    Then you can do:

    ```bash
    make compose-up-build
    ```

    This command does the following:

    * Builds the Docker iamge for the go backend and the postgres db.
    * Starts all services as defined in the `docker-compose-dev.yml` file.
    * Runs the database migrations using Goose.

    The following services will be available:

      * **beanbag-backend:** Go API - accessible at `http://localhost:8080`.
      * **db:** PostgreSQL database

### Stopping the Development Environment

1. Stop all services:

    ```bash
    make compose-down
    ```

    [Optional]: Wipe database dir
    ```bash
    make compose-down-wipe
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

## Adding a New Model (Table) or Changing the Database Schema

This project uses Goose for database migrations. Here's how to add a new model or change the database schema:

1.  **Create a New Migration File:**
    *   Use the `goose create` command to create a new migration file in the `migrations` directory:

        ```bash
        goose -dir migrations create add_my_new_table sql
        ```

        *   Replace `add_my_new_table` with a descriptive name for your migration (e.g., `add_users_table`, `add_email_to_users`, etc.).
        *   The `sql` argument specifies that you want to create a SQL migration file.

2.  **Define Your Schema Changes:**
    *   Open the newly created migration file (e.g., `migrations/YYYYMMDDHHMMSS_add_my_new_table.sql`).
    *   Define your schema changes in the `-- +goose Up` section.
    *   Define the reverse changes (to undo the migration) in the `-- +goose Down` section.

    ```sql
    -- +goose Up
    -- +goose StatementBegin
    CREATE TABLE IF NOT EXISTS my_new_table (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    -- +goose StatementEnd

    -- +goose Down
    -- +goose StatementBegin
    DROP TABLE IF EXISTS my_new_table;
    -- +goose StatementEnd
    ```

3.  **Run the Migrations:**
    *   The migrations will automatically run when you start the development environment with `make compose-up-build`.
    *   The `make` commands automatically run `sqlc generate`, which will update the `db/` directory with new queries based on any files in the `queries` folder. SQLC will also use the combined `-- +goose Up` versions of all the migrations to determine the database schema.

4.  **Create New Queries (Optional):**
    *   If you've added a new table or need to interact with the database in a new way, you'll likely need to create new SQL queries.
    *   Create new `.sql` files in the appropriate subdirectory under the `queries/` directory (e.g., `queries/users/`, `queries/quizzes/`, etc.).
    *   Define your queries using the SQLC `-- name:` directive.

    ```sql
    -- queries/users/get_user_by_id.sql
    -- name: GetUserByID :one
    SELECT * FROM users WHERE user_id = $1;
    ```

## Docs

When the docker containers are running the docs are available in interactive format at http://localhost:8080/swagger/index.html