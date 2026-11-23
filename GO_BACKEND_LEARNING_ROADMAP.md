# Go Backend Services: Learning-by-Building Roadmap

You have finished the core Go language and syntax. This roadmap is for learning backend development by building progressively larger services. The goal is understanding how backend systems work, not building production-ready products.

## How to use this roadmap

- Build every project yourself before moving to the next stage.
- Start with the Go standard library whenever it is capable of solving the problem.
- Add one major library or concept at a time.
- Keep each project small enough to rewrite after learning something new.
- Write tests for the important behavior, but do not spend time on deployment, scaling, or advanced security yet.
- Keep a short `README.md` for each project containing how to run it and example requests.

## Recommended learning order

| Stage | Focus | Project |
|---|---|---|
| 0 | Backend foundations | HTTP and JSON service |
| 1 | Routing and CRUD | In-memory task API |
| 2 | Persistence | SQLite task or notes API |
| 3 | API structure | Layered book/library API |
| 4 | Authentication | User and session service |
| 5 | Background work | Job queue and notification service |
| 6 | External services | URL shortener with an external API |
| 7 | Real-time communication | Chat or live notifications |
| 8 | Service architecture | Small e-commerce or expense system |
| 9 | Capstone | Multi-feature backend service |

---

## Stage 0: Build an HTTP and JSON service

### Concepts

- HTTP methods, status codes, headers, and request bodies
- URL paths and query parameters
- JSON encoding and decoding
- Request validation and error responses
- Starting and stopping an HTTP server
- Packages, modules, and environment configuration

### Use these first

- `net/http`
- `encoding/json`
- `fmt`
- `log/slog`
- `os`
- `strconv`
- `time`

### Project: Quote API

Build an API that stores a small list of quotes in memory.

Suggested endpoints:

```text
GET    /health
GET    /quotes
GET    /quotes/{id}
POST   /quotes
DELETE /quotes/{id}
```

Practice:

- Returning JSON consistently
- Returning `400`, `404`, and `500` responses
- Parsing path values and query parameters
- Separating HTTP handling from quote logic

Do not use a web framework yet. Learn what `net/http` provides first.

---

## Stage 1: Build a CRUD service with routing

### Concepts

- REST-style resource design
- Routers and route parameters
- In-memory repositories
- Service and handler responsibilities
- Table-driven unit tests
- Middleware

### Standard library

- `net/http`
- `testing`
- `httptest`
- `context`

### Optional routing library

- `github.com/go-chi/chi/v5`

Use `chi` only after implementing a few routes with `net/http`. Compare the two approaches.

### Project: Task Manager API

Features:

- Create, list, update, and delete tasks
- Filter tasks by completion status
- Add due dates and tags
- Validate missing or invalid fields
- Add request logging and a request ID middleware

Testing goals:

- Handler tests using `httptest`
- Service tests using a fake repository
- Tests for invalid JSON and missing resources

---

## Stage 2: Add database persistence

### Concepts

- Relational data and tables
- Primary keys and foreign keys
- SQL queries and parameters
- Database migrations
- Transactions
- Repository interfaces
- Pagination and sorting

### Start with SQLite

- `database/sql`
- `modernc.org/sqlite` (pure Go SQLite driver)
- `github.com/golang-migrate/migrate/v4` (migrations)

Alternative driver to learn later:

- `github.com/mattn/go-sqlite3`

### Project: Notes API

Features:

- Users can create and edit notes
- Notes have titles, content, and timestamps
- Search notes by title
- Paginate the list endpoint
- Store data in SQLite instead of memory

Practice:

- Writing SQL manually
- Using `QueryRowContext`, `QueryContext`, and `ExecContext`
- Scanning rows into structs
- Closing rows correctly
- Using a transaction for a multi-step operation

Do not hide SQL behind an ORM at the beginning. SQL knowledge is more valuable here.

---

## Stage 3: Build a maintainable API

### Concepts

- Package boundaries
- Dependency injection
- Interfaces at the consumer side
- DTOs versus database models
- Centralized error handling
- Configuration loading
- Graceful server shutdown

### Useful libraries

- `github.com/go-chi/chi/v5` for routing
- `github.com/joho/godotenv` for local `.env` files
- `github.com/google/uuid` for IDs
- `github.com/stretchr/testify` for readable test assertions

### Project: Library Management API

Resources:

- Books
- Authors
- Members
- Loans

Learn to:

- Model relationships
- Prevent a book from being loaned twice
- Use transactions for checkout and return
- Add service-level business rules
- Organize code into `handler`, `service`, and `repository` packages

Suggested structure:

```text
cmd/server/main.go
internal/config/
internal/handler/
internal/service/
internal/repository/
internal/model/
internal/middleware/
migrations/
```

This is a learning structure, not a rule that every Go project must follow.

---

## Stage 4: Add authentication and authorization

### Concepts

- Password hashing
- Registration and login
- Sessions versus tokens
- Authentication middleware
- Roles and permissions
- Protecting resources owned by a user

### Useful libraries

- `golang.org/x/crypto/bcrypt` for password hashing
- `github.com/golang-jwt/jwt/v5` for learning JWTs
- `github.com/google/uuid` for user and session IDs

### Project: Personal Expense API

Features:

- User registration and login
- Create, update, and delete personal expenses
- Categorize expenses
- Monthly summary endpoint
- Admin-only endpoint for listing users

Learning rule:

- First implement cookie-based sessions with a database.
- Then implement JWT authentication as a comparison exercise.
- Do not treat JWTs as automatically better; understand the trade-offs.

For learning, keep secrets in environment variables and use deliberately simple authorization rules.

---

## Stage 5: Learn background work

### Concepts

- Goroutines and channels in a backend context
- Worker pools
- Job states and retries
- Scheduled tasks
- Context cancellation
- Idempotency at a basic level

### Standard library

- `context`
- `sync`
- `time`
- `net/http`

### Optional libraries

- `github.com/robfig/cron/v3` for scheduled jobs
- `github.com/hibiken/asynq` for a Redis-backed task queue
- `github.com/redis/go-redis/v9` for Redis access

### Project: Report Generation Service

Features:

- User submits a report request
- API returns a job ID
- Worker processes the report in the background
- User checks job status
- Store completed reports as local files

Start with an in-process channel worker. Add Redis and Asynq only after understanding the basic implementation.

---

## Stage 6: Integrate external services

### Concepts

- Calling another HTTP API
- Timeouts and cancellation
- Decoding external responses
- Handling upstream errors
- Caching a response
- Faking external services in tests

### Useful libraries

- `net/http` with a custom `http.Client`
- `github.com/stretchr/testify`
- `github.com/patrickmn/go-cache` for a simple in-memory cache

### Project: URL Shortener with Link Preview

Features:

- Create short URLs
- Redirect to the original URL
- Fetch and cache a page title
- Count visits
- Add an endpoint for link statistics

Practice configuring an HTTP client with a timeout and testing it with `httptest.Server`.

---

## Stage 7: Build a real-time service

### Concepts

- WebSockets
- Connection management
- Broadcasting messages
- Concurrent access to shared state
- Heartbeats and disconnects

### Useful library

- `github.com/coder/websocket`

### Project: Room Chat Service

Features:

- Join a named room
- Send and receive messages
- List active rooms
- Store recent messages in SQLite
- Notify clients when someone joins or leaves

Keep the first version single-process and in memory. Distributed WebSocket infrastructure is outside this learning roadmap.

---

## Stage 8: Combine the concepts

### Project: Small e-commerce backend

Build a learning version with:

- Users and authentication
- Products and categories
- Shopping cart
- Order creation
- Inventory updates
- Background email-like notifications
- Admin endpoints

Important exercises:

- Use a transaction when creating an order and reducing inventory.
- Make the order state explicit: `pending`, `paid`, `cancelled`, `completed`.
- Add tests for insufficient inventory and duplicate requests.
- Use a fake notification service rather than integrating a real email provider.

### Alternative project: Habit tracker

If e-commerce feels too large, build a habit tracker with users, recurring habits, daily check-ins, reminders, statistics, and a background worker.

---

## Stage 9: Capstone project

Choose one service and rebuild it cleanly from the earlier lessons:

- Learning platform API
- Issue tracker API
- Personal finance API
- Event booking API
- Blog platform API

The capstone should include:

- REST API with documented endpoints
- SQLite or PostgreSQL persistence
- Authentication and authorization
- Validation and consistent errors
- Unit tests and handler tests
- One background job
- One external HTTP integration, mocked in tests
- Structured logging
- A simple CLI or seed script

Use PostgreSQL only after you are comfortable with SQLite and SQL.

Useful PostgreSQL tools:

- `github.com/jackc/pgx/v5`
- `github.com/golang-migrate/migrate/v4`
- Docker Compose for running PostgreSQL locally

---

## Library and module checklist

Learn these in roughly this order:

### Go standard library

- `net/http`
- `encoding/json`
- `database/sql`
- `context`
- `testing`
- `net/http/httptest`
- `log/slog`
- `sync`
- `time`
- `os`
- `flag`

### Common backend libraries

- `github.com/go-chi/chi/v5` - routing
- `github.com/stretchr/testify` - test assertions and mocks
- `github.com/google/uuid` - UUIDs
- `golang.org/x/crypto/bcrypt` - password hashing
- `github.com/golang-jwt/jwt/v5` - JWT learning exercise
- `github.com/golang-migrate/migrate/v4` - migrations
- `modernc.org/sqlite` - SQLite driver
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/robfig/cron/v3` - scheduled jobs
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/coder/websocket` - WebSockets

### Tools

- `go mod`
- `go fmt`
- `go vet`
- `go test`
- `go test -race`
- `go doc`
- `curl` or an API client such as Postman/Insomnia
- Docker and Docker Compose for local databases

---

## Topics to postpone

Do not make these the focus while learning the fundamentals:

- Kubernetes
- Microservices and service meshes
- Complex cloud deployment
- Distributed tracing platforms
- High-availability database design
- Premature performance optimization
- Advanced event-driven architecture
- Building a custom ORM or framework
- Production-grade secrets and compliance systems

Return to these after you can comfortably build, test, debug, and explain a complete small service.

## Definition of progress

You are ready for the next stage when you can:

1. Explain the request flow from HTTP handler to database and back.
2. Add a feature without putting all logic in the handler.
3. Write a test that fails before the feature and passes afterward.
4. Handle invalid input and expected failures explicitly.
5. Run and debug the service without copying an unexplained solution.
6. Describe at least one limitation of your current design.
