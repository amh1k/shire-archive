# ShireArchive

ShireArchive is a server-rendered Go web application for creating and sharing expiring text snippets. It includes user registration, authentication, account management, MySQL-backed sessions, CSRF protection, and an embedded frontend.

> [!IMPORTANT]
> This repository is a learning project with production-oriented security practices. Review the [production checklist](#production-checklist) before exposing it to the public internet.

## Features

- Create snippets that expire after 1, 7, or 365 days
- Register, log in, log out, and change account passwords
- Return users to their intended protected page after login
- Persist application data and sessions in MySQL
- Protect state-changing requests with CSRF tokens
- Hash passwords with bcrypt
- Serve exclusively over HTTPS
- Embed templates and static assets in the application binary
- Apply security headers and disable caching on authenticated pages
- Recover from panics and log requests and server errors

## Technology

| Area | Choice |
| --- | --- |
| Language | Go 1.26.2 |
| Database | MySQL |
| Router | `julienschmidt/httprouter` |
| Middleware | `justinas/alice` |
| Sessions | `alexedwards/scs/v2` with `mysqlstore` |
| CSRF protection | `justinas/nosurf` |
| Forms | `go-playground/form/v4` |
| Password hashing | `golang.org/x/crypto/bcrypt` |

## Requirements

- Go 1.26.2 or a compatible newer version
- MySQL 8.x (or a compatible MySQL server)
- OpenSSL, or another way to generate a local TLS certificate

## Quick start

### 1. Clone the repository

```bash
git clone git@github.com:amh1k/shire-archive.git
cd shire-archive
go mod download
```

### 2. Create the database

Open a MySQL session as an administrative user and run:

```sql
CREATE DATABASE shirearchive
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE USER 'web'@'localhost' IDENTIFIED BY 'replace-with-a-strong-password';
GRANT SELECT, INSERT, UPDATE, DELETE ON shirearchive.* TO 'web'@'localhost';

USE shirearchive;

CREATE TABLE snippets (
    id      INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title   VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
    id              INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created         DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);

CREATE TABLE sessions (
    token  CHAR(43) PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);
```

There is not yet a migration command in this repository. Schema changes must currently be applied manually.

### 3. Generate a development certificate

The server expects `tls/cert.pem` and `tls/key.pem`, relative to the directory from which it is started. The `tls/` directory is intentionally ignored by Git.

```bash
mkdir -p tls
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout tls/key.pem \
  -out tls/cert.pem \
  -days 365 \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
```

Your browser will warn about this self-signed certificate. That is expected for local development.

### 4. Start the server

Pass the database credentials at runtime instead of relying on the example default in the source code:

```bash
go run ./cmd/web \
  -addr=":4000" \
  -dsn="web:replace-with-a-strong-password@/shirearchive?parseTime=true"
```

Visit <https://localhost:4000>. A health check is available at <https://localhost:4000/ping>.

## Configuration

Configuration is currently provided through command-line flags. The application does not read `.env` or `.config` files.

| Flag | Default | Purpose |
| --- | --- | --- |
| `-addr` | `:4000` | HTTPS listen address |
| `-dsn` | `web:abc@/shirearchive?parseTime=true` | MySQL data source name |
| `-debug` | `false` | Display internal errors and stack traces in HTTP responses |

Run `go run ./cmd/web -help` to see the available flags.

> [!WARNING]
> Never enable `-debug` in production. Avoid putting real credentials in source control, screenshots, shell history, issue reports, or documentation. For a production deployment, inject the DSN through a secret manager or protected runtime configuration.

## Testing

Run the fast test suite, which skips database integration tests:

```bash
go test -short ./...
```

Run all tests with:

```bash
go test ./...
```

The `internal/models` integration test expects a local MySQL database with the following fixed test credentials:

```sql
CREATE DATABASE test_snippetbox
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE USER 'test_web'@'localhost' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON test_snippetbox.* TO 'test_web'@'localhost';
```

The test setup and teardown scripts create and remove their own tables. Use a dedicated test database only—never point the tests at development or production data.

Useful development checks:

```bash
go test -race -short ./...
go vet ./...
gofmt -w cmd internal ui
```

## Application structure

```text
.
├── cmd/web/                 HTTP server, routes, middleware, handlers, and tests
├── internal/assert/         Test assertion helpers
├── internal/models/         MySQL models, interfaces, mocks, and model tests
├── internal/validator/      Reusable form validation
├── ui/html/                 Page layouts, pages, and partial templates
├── ui/static/               CSS, JavaScript, and images
├── ui/efs.go                Embedded UI filesystem
├── go.mod                   Module and dependency declarations
└── main.go                  Earlier standalone example; not the ShireArchive server
```

The supported application entry point is `./cmd/web`. The root `main.go` is a separate, earlier example and does not start the full application.

## Request flow

```text
HTTPS request
  -> panic recovery
  -> request logging
  -> security headers
  -> router
      -> static files and /ping
      -> session loading + CSRF protection + authentication
          -> public dynamic routes
          -> authentication requirement -> protected routes
  -> handler
  -> MySQL model
```

Protected routes include snippet creation, logout, account viewing, and password changes. Unauthenticated requests are redirected to the login page, and the original path is stored so the user can return after authenticating.

## Security model

The application currently includes:

- HTTPS with explicit curve preferences
- `Secure`, `HttpOnly`, and `SameSite=Lax` CSRF cookies
- MySQL-backed server-side sessions with a 12-hour lifetime
- Session token rotation during login and logout
- Bcrypt password hashing
- CSRF protection on dynamic routes
- Content Security Policy, clickjacking protection, MIME sniffing protection, and referrer policy headers
- `Cache-Control: no-store` on authenticated routes
- Database uniqueness enforcement for email addresses

Security reports should not include credentials, session tokens, private keys, or personal data. If this repository is used by others, add a private vulnerability-reporting method in `SECURITY.md`.

## Production checklist

Before deploying publicly:

- Replace the example database credentials and remove the insecure default DSN from application code.
- Load secrets from protected runtime configuration or a secret manager.
- Use a certificate issued by a trusted certificate authority, or terminate TLS at a properly configured reverse proxy.
- Run the service as an unprivileged user and restrict access to the TLS private key.
- Add versioned database migrations and a documented backup/restore process.
- Add graceful shutdown and deployment health/readiness checks.
- Send structured logs to a persistent logging system and configure monitoring and alerting.
- Pin a supported Go toolchain and keep dependencies updated.
- Review the Content Security Policy and all cookie settings for the deployment domain.
- Put rate limiting and request-size limits at the application or reverse-proxy layer.
- Run `go test -race -short ./...`, `go vet ./...`, and a vulnerability scanner in CI.

## Troubleshooting

### `dial unix /var/run/mysqld/mysqld.sock: connect: no such file or directory`

MySQL is not running, is using a different socket, or the DSN needs an explicit TCP address. For example:

```bash
-dsn="web:password@tcp(127.0.0.1:3306)/shirearchive?parseTime=true"
```

### `access denied for user`

Confirm that the username, password, host, and grants match the DSN. In MySQL, `'web'@'localhost'` and `'web'@'127.0.0.1'` can be treated as different accounts.

### `open ./tls/cert.pem: no such file or directory`

Start the application from the repository root and generate the certificate files described above.

### Browser certificate warning

Self-signed certificates are not trusted automatically. Use them only for local development; use a trusted certificate in production.

## Contributing

1. Create a focused branch.
2. Make the change and add or update tests.
3. Run formatting, tests, and static checks.
4. Open a pull request explaining the motivation, behavior change, and verification performed.

Please keep commits focused and never commit database credentials, private keys, generated certificates, or local environment files.

## License

No license file is currently included. Unless a license is added, standard copyright restrictions apply and reuse is not automatically granted.
