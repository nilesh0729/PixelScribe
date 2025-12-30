# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go)
```bash
# Run backend server
go run ./cmd/api

# Run all tests with coverage
make Test
# or: go test -v -cover ./...

# Generate SQLC code (after modifying SQL queries)
make Sqlc
# or: sqlc generate

# Database migrations
make MigrateUp    # Apply migrations
make MigrateDown  # Rollback migrations

# Database setup (for local dev without Docker)
make Container    # Start PostgreSQL container
make CreateDB     # Create database
make DropDB       # Drop database
```

### Frontend (React + Vite)
```bash
cd web

# Start dev server
npm run dev

# Build for production
npm run build

# Lint TypeScript/React code
npm run lint

# Preview production build
npm run preview
```

### Docker
```bash
# Start entire stack (frontend + backend + postgres)
docker compose up --build

# Frontend: http://localhost:5173
# Backend: http://localhost:8080
```

## Architecture

### Backend Architecture (Go + Gin + SQLC)

**Entry Point**: `cmd/api/main.go` loads config, establishes DB connection, creates server, and starts HTTP listener.

**Core Pattern**: Repository pattern with SQLC-generated type-safe queries
- `db/query/*.sql` → Raw SQL queries (CREATE, SELECT, UPDATE, DELETE)
- `sqlc generate` → Generates type-safe Go code in `internal/db/sqlc/`
- `internal/db/sqlc/querier.go` → Interface with all query methods
- `internal/db/sqlc/transaction.go` → Store interface with transactional operations

**Store Pattern**: The `Store` interface (`internal/db/sqlc/transaction.go`) extends `Querier` with transaction methods:
- `SubmitAttemptTx`: Creates attempt + updates/creates performance summary (atomic)
- `CreateUserTx`: Creates user + default settings (atomic)
- `DeleteDictationTx`: Cascading delete for dictation + attempts + performance summaries (atomic)

All business logic transactions follow this pattern: define result struct → implement `execTx` function → handle rollback on error.

**API Layer** (`internal/api/`):
- `server.go`: Route definitions, middleware setup (CORS, JWT auth), server initialization
- Handler files: `user.go`, `dictation.go`, `attempt.go`, `performance.go`, `settings.go`, `tts.go`
- Each handler extracts request params → calls store methods → returns JSON response
- `middleware.go`: JWT authentication middleware (`authMiddleware`) validates tokens and extracts user context

**Authentication Flow**:
- JWT tokens managed by `internal/token/` (interface + implementation pattern)
- `token.Maker` interface with `jwt_maker.go` implementation
- Protected routes use `authMiddleware` which validates token and sets user context
- Token payload contains: `id`, `username`, `issued_at`, `expired_at`

**Configuration**: `internal/util/config.go` uses Viper to load `.env` file with DB credentials, server address, JWT key, OpenAI API key, and token duration.

### Frontend Architecture (React + TypeScript + Vite)

**Structure**:
- `web/src/pages/`: Page components (Dashboard, Dictation, History, etc.)
- `web/src/components/`: Reusable UI components
- `web/src/services/`: API client modules (axios-based)
  - `attempt.ts`, `dictation.ts`, `performance.ts`, `tts.ts`
- `web/src/context/`: React context for auth state management
- `web/src/hooks/`: Custom React hooks
- `web/src/layouts/`: Layout wrapper components
- `web/src/types/`: TypeScript type definitions

**API Communication**: Services in `web/src/services/` handle all backend communication. Each service exports functions that make axios requests to the Go backend at `http://localhost:8080`.

**Styling**: Tailwind CSS v4 with utility-first approach. Custom styles in `index.css` and component-specific CSS.

### Database Schema

**Core Tables**:
- `users`: User accounts (username, hashed_password, full_name, email, created_at)
- `dictations`: Text/audio dictations (title, content, audio_url, language, user_id, type)
- `attempts`: User typing attempts (user_id, dictation_id, input_text, accuracy, wpm, time_spent, missed_words, incorrect_words, extra_words)
- `performance_summary`: Aggregated stats per user+dictation (total_attempts, best_accuracy, average_accuracy, average_time, last_attempt_at)
- `settings`: User preferences (default_voice, default_speed, highlight colors)

**Key Relationships**:
- Users → Dictations (one-to-many)
- Dictations → Attempts (one-to-many)
- Users + Dictations → Performance Summary (composite unique key)
- Users → Settings (one-to-one)

### Testing Strategy

**Backend Testing**:
- Unit tests for all handlers in `internal/api/*_test.go`
- Mock database using `go.uber.org/mock` (generated mocks in `internal/db/mock/`)
- Test structure: Setup mock store → Create test server → Make HTTP request → Assert response
- Database tests in `internal/db/sqlc/*_test.go` use real postgres connection
- All tests run with `make Test` or `go test -v -cover ./...`

**Test Helper**: `internal/api/auth_test_helper.go` provides utilities for creating test JWT tokens and authenticated requests.

## Adding New Database Operations

1. Write SQL query in `db/query/<table>.sql`
2. Run `make Sqlc` to generate Go code
3. Use generated methods from `internal/db/sqlc/querier.go`
4. For multi-step operations requiring atomicity, add transaction method to `Store` interface in `internal/db/sqlc/transaction.go`
5. Implement transaction using `execTx` helper
6. Add tests in `internal/db/sqlc/*_test.go`

## Adding New API Endpoints

1. Create handler function in appropriate file in `internal/api/` (e.g., `user.go`, `dictation.go`)
2. Define request/response structs
3. Register route in `internal/api/server.go` (use `authRoutes` for protected endpoints)
4. Add tests in corresponding `*_test.go` file with mock store
5. Update frontend service in `web/src/services/` to call new endpoint

## Environment Configuration

Backend requires `.env` file (see `app.env.example`) with:
- `DB_SOURCE`: PostgreSQL connection string
- `SERVER_ADDRESS`: Backend listen address (default: `0.0.0.0:8080`)
- `TOKEN_SYMMETRIC_KEY`: 32-character key for JWT signing
- `ACCESS_TOKEN_DURATION`: Token expiry (e.g., `15m`)
- `OPENAI_API_KEY`: For TTS generation via `/tts/generate` endpoint
