# ShireArchive

ShireArchive is a secure, monolithic web application written in Go that allows users to create, share, and view text snippets. It features user authentication, active session management, template rendering with an embedded filesystem, and CSRF protection.

## Features

- **User Authentication**: Secure user registration and login using `bcrypt` password hashing.
- **Session Management**: Persistent sessions backed by MySQL using `scs/v2` with a 12-hour expiration time.
- **Snippet Management**: Create snippets with titles, content, and expiration dates (1, 7, or 365 days).
- **Embedded UI Assets**: All HTML templates and static assets (CSS, JS, images) are bundled securely into the Go binary using `go:embed`.
- **Security-First Approach**:
  - Secure TLS server configuration with restricted curve preferences.
  - Built-in Cross-Site Request Forgery (CSRF) protection using `nosurf`.
  - Comprehensive form validation for user input.
  - Context-based authentication middleware to protect sensitive routes.

## Disclaimer on Testing

> [!WARNING]
> Please note that **testing has not been covered yet** in the scope of this project. Comprehensive unit and integration test coverage will be iteratively added in future updates.

## Tech Stack

- **Language**: Go 1.20+
- **Database**: MySQL
- **Router**: `julienschmidt/httprouter`
- **Middleware Chaining**: `justinas/alice`
- **Session Management**: `alexedwards/scs/v2`
- **Form Decoder**: `go-playground/form/v4`

## Project Structure

```text
├── cmd
│   └── web/            # Application configuration, handlers, routing, and middleware
├── internal
│   ├── models/         # Database models for Users and Snippets
│   └── validator/      # Custom form validation logic
├── tls/                # TLS certificates for secure HTTPS connections
├── ui/                 # Embedded frontend assets
│   ├── html/           # Go HTML templates (pages, partials, base)
│   ├── static/         # Public CSS, JS, and image files
│   └── efs.go          # Embedded File System configurations
├── go.mod / go.sum     # Go module definitions
└── main.go             # Application entry pont
```

## Getting Started

### Prerequisites

- Go installed on your machine
- MySQL database running locally or remotely

### Configuration

1. Clone the repository.
2. Initialize your MySQL database with the `users` and `snippets` tables.
3. Configure your database connection utilizing the `-dsn` flag or hardcode your DSN into the main application.
4. Ensure your TLS certificates are located under `./tls/cert.pem` and `./tls/key.pem`.

### Running the App

Run the following command at the root of the project:

```bash
go run ./cmd/web
```

By default, the server will start on port `4000`. Open your browser and navigate to:
`https://localhost:4000`

## Contributing

Contributions are welcome. Please ensure that all new features adhere to the project's embedded structure and strict validation patterns. As mentioned, testing frameworks have not been integrated yet, so thorough manual testing is expected before submitting a pull request.
