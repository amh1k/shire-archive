# ShireArchive

ShireArchive is a secure, monolithic web application written in Go that allows users to create, share, and view text snippets. It implements a layered middleware chain, interface-driven database models, CSRF protection, and a comprehensive test suite against in-memory mocks.

## Features

- **User Authentication**: Registration, login, and logout backed by `bcrypt` password hashing (cost factor 12).
- **Account Management**: Authenticated users can view their profile and update their password in-app.
- **Session Management**: MySQL-backed sessions via `alexedwards/scs` with a 12-hour lifetime and secure cookies.
- **Snippet Management**: Create snippets with a title, content, and expiration of 1, 7, or 365 days.
- **Redirect-After-Login**: Unauthenticated users are redirected to their intended destination after successfully logging in.
- **Embedded UI Assets**: All HTML templates and static files are bundled into the binary using `go:embed`.
- **Debug Mode**: A `-debug` flag exposes server-side error details and stack traces in the browser.
- **Security-First Design**:
  - TLS-only server with restricted curve preferences (`X25519`, `CurveP256`).
  - CSRF protection on every state-mutating request via `justinas/nosurf`.
  - `Cache-Control: no-store` on all protected routes to prevent back-button leaks.
  - Strict security headers: `Content-Security-Policy`, `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`.
  - Context-based authentication propagation via a typed context key.

## Architecture

```
                    HTTPS Request
                         |
                         v
          +-----------------------------+
          |  Standard Middleware Chain  |
          |  recoverPanic               |
          |  logRequest                 |
          |  secureHeaders              |
          +-----------------------------+
                         |
                         v
          +-----------------------------+
          |       httprouter Mux        |
          +-----------------------------+
                         |
             +-----------+-----------+
             |                       |
             v                       v
    [Public routes]         [Dynamic middleware]
    /ping                   sessionManager.LoadAndSave
    /static/*               noSurf (CSRF)
                            authenticate
                                    |
                       +------------+------------+
                       |                         |
                       v                         v
             [Open dynamic]            [Protected] (+requireAuthentication)
             GET  /                    GET/POST /snippet/create
             GET  /snippet/view/:id    POST     /user/logout
             GET  /user/signup         GET      /account/view
             POST /user/signup         GET/POST /account/password/update
             GET  /user/login
             POST /user/login
             GET  /about
                       |
                       v
          +-----------------------------+
          |       application struct    |
          |  snippets SnippetModelIface |
          |  users    UserModelIface    |
          |  templateCache              |
          |  formDecoder                |
          |  sessionManager             |
          +-----------------------------+
                    /       \
                   v         v
          +-----------+  +-----------+
          | SnippetModel| | UserModel |
          | (MySQL)     | | (MySQL)   |
          +-----------+  +-----------+
```

### Request Lifecycle

Every request passes through the **standard chain** (`recoverPanic -> logRequest -> secureHeaders`) before reaching the router. Dynamic routes additionally pass through `sessionManager.LoadAndSave -> noSurf -> authenticate`. Routes requiring a logged-in user append `requireAuthentication`, which stores the attempted path in the session and redirects to `/user/login` if unauthenticated.

## Tech Stack

| Concern             | Library / Tool                          |
|---------------------|-----------------------------------------|
| Language            | Go 1.26                                 |
| Database            | MySQL                                   |
| Router              | `julienschmidt/httprouter`              |
| Middleware chaining | `justinas/alice`                        |
| Session management  | `alexedwards/scs/v2` + `mysqlstore`     |
| CSRF protection     | `justinas/nosurf`                       |
| Form decoding       | `go-playground/form/v4`                 |
| Password hashing    | `golang.org/x/crypto/bcrypt`            |
| TLS driver          | Standard library `crypto/tls`           |

## Project Structure

```
.
├── cmd/
│   └── web/
│       ├── main.go               # Dependency wiring, server bootstrap
│       ├── routes.go             # Route definitions and middleware chains
│       ├── handlers.go           # HTTP handler functions and form structs
│       ├── middleware.go         # secureHeaders, logRequest, recoverPanic,
│       │                         # requireAuthentication, noSurf, authenticate
│       ├── helpers.go            # render, newTemplateData, clientError, serverError
│       ├── templates.go          # Template cache construction, humanDate func
│       ├── context.go            # Typed context key for auth propagation
│       ├── handlers_test.go      # Integration tests: Ping, SnippetView,
│       │                         # UserSignup, SnippetCreate
│       ├── middleware_test.go    # Middleware-level tests
│       ├── templates_test.go     # Template cache tests
│       └── testutils_test.go     # newTestApplication, newTestServer,
│                                 # extractCSRFToken helpers
├── internal/
│   ├── models/
│   │   ├── snippets.go           # SnippetModel + SnippetModelInterface
│   │   ├── users.go              # UserModel + UserModelInterface
│   │   ├── errors.go             # Sentinel errors (ErrNoRecord, ErrDuplicateEmail, ...)
│   │   ├── testutils_test.go     # Model-layer test helpers
│   │   ├── users_test.go         # UserModel unit tests
│   │   ├── testdata/             # SQL seed scripts for model tests
│   │   └── mocks/                # In-memory mock implementations of both interfaces
│   ├── assert/
│   │   └── assert.go             # Equal, StringContains, NilError test helpers
│   └── validator/
│       └── validator.go          # Validator struct, CheckField, NotBlank, MinChars,
│                                 # MaxChars, PermittedValue, Matches, EmailRX
├── tls/
│   ├── cert.pem                  # TLS certificate
│   └── key.pem                   # TLS private key
├── ui/
│   ├── efs.go                    # go:embed declaration
│   ├── html/
│   │   ├── base.tmpl             # Base layout
│   │   ├── pages/                # home, view, create, signup, login,
│   │   │                         # about, account, password templates
│   │   └── partials/             # nav partial
│   └── static/                   # CSS, JS, images
├── go.mod
├── go.sum
└── main.go                       # Entry point (delegates to cmd/web)
```

## Test Suite

The project includes a layered test suite with no external test database required.

### Architecture

- **`internal/models/mocks`** provides in-memory implementations of `SnippetModelInterface` and `UserModelInterface`, seeded with fixture data (e.g. snippet ID 1, user `alice@example.com`).
- **`internal/assert`** provides generic `Equal[T]`, `StringContains`, and `NilError` helpers.
- **`cmd/web/testutils_test.go`** provides `newTestApplication` (wires mocks), `newTestServer` (wraps `httptest.NewTLSServer` with a persistent cookie jar and redirect control), and `extractCSRFToken` (regex-based CSRF scraper).

### Test Coverage

| Test                     | What it verifies                                                                 |
|--------------------------|----------------------------------------------------------------------------------|
| `TestPing`               | Health check endpoint returns `200 OK` with body `"OK"`                          |
| `TestSnippetView`        | Valid ID serves content; invalid/negative/decimal/string/empty IDs return 404    |
| `TestUserSignup`         | Valid submission redirects 303; invalid CSRF returns 400; field validation returns 422 and re-renders the form |
| `TestSnippetCreate`      | Unauthenticated request redirects to `/user/login`; authenticated session serves the create form |

The `postForm` helper sets a `Referer` header matching the test server origin to satisfy `nosurf`'s HTTPS referrer validation.

### Running Tests

```bash
# Run all tests with verbose output
go test -v ./...

# Run only handler tests
go test -v ./cmd/web/

# Run only model tests
go test -v ./internal/models/
```

## Getting Started

### Prerequisites

- Go 1.20 or later
- MySQL server (local or remote)

### Database Setup

Create the database and user, then apply the schema:

```sql
CREATE DATABASE shirearchive CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'web'@'localhost' IDENTIFIED BY 'abc';
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
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

### TLS Certificates

For local development, generate a self-signed certificate pair:

```bash
mkdir tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
mv cert.pem tls/cert.pem
mv key.pem tls/key.pem
```

### Running the Application

```bash
go run ./cmd/web
```

The server starts on `:4000` by default. Navigate to `https://localhost:4000`.

#### CLI Flags

| Flag    | Default                                    | Description                                    |
|---------|--------------------------------------------|------------------------------------------------|
| `-addr` | `:4000`                                    | Network address and port                       |
| `-dsn`  | `web:abc@/shirearchive?parseTime=true`     | MySQL data source name                         |
| `-debug`| `false`                                    | Show error details and stack traces in browser |

Example with custom flags:

```bash
go run ./cmd/web -addr=":8080" -dsn="user:pass@/dbname?parseTime=true" -debug
```

## Security Notes

- All session cookies are configured with `Secure`, `HttpOnly`, and `SameSite=Lax`.
- Session tokens are rotated on login and logout via `RenewToken` to prevent session fixation.
- Duplicate email registration is caught at the database constraint level and surfaced as a typed sentinel error (`ErrDuplicateEmail`).
- The `-debug` flag must never be enabled in production as it exposes internal error details to the client.
