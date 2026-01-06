# Integration Test Suite Implementation Plan

## Overview
Build a comprehensive HTTP-based integration test suite for the Pulsar backend, covering all 54 API endpoints with full HTTP stack testing (routing, middleware, JSON serialization).

## Requirements
- **Coverage**: Full API coverage - all 54 endpoints
- **Database**: Docker Compose with separate test PostgreSQL database
- **Approach**: HTTP client tests (full stack testing)

---

## Directory Structure

```
backend/
├── tests/
│   └── integration/
│       ├── testutils/
│       │   ├── config.go       # Test configuration
│       │   ├── database.go     # DB connection, migrations, cleanup
│       │   ├── server.go       # Test server setup
│       │   ├── client.go       # HTTP client with auth helpers
│       │   └── fixtures.go     # Test data factories
│       ├── main_test.go        # Suite setup/teardown
│       ├── auth_test.go        # 5 endpoints
│       ├── alert_test.go       # 9 endpoints
│       ├── team_test.go        # 9 endpoints
│       ├── user_test.go        # 1 endpoint
│       ├── schedule_test.go    # 20 endpoints
│       ├── escalation_test.go  # 13 endpoints
│       ├── notification_test.go # 15 endpoints
│       ├── incident_test.go    # 14 endpoints
│       └── webhook_test.go     # 8 endpoints
docker-compose.test.yml         # Test database (port 5434, tmpfs)
```

---

## Implementation Phases

### Phase 1: Infrastructure Setup
1. Create `docker-compose.test.yml` with PostgreSQL on port 5434 (tmpfs for speed)
2. Create `tests/integration/testutils/config.go` - test environment config
3. Create `tests/integration/testutils/database.go` - migrations, truncate helpers
4. Create `tests/integration/testutils/server.go` - httptest server with full Gin router
5. Update `Makefile` with test targets

### Phase 2: Test Utilities
1. Create `tests/integration/testutils/client.go` - HTTP client with auth token support
2. Create `tests/integration/testutils/fixtures.go` - factories for users, teams, alerts, etc.
3. Create `tests/integration/main_test.go` - TestMain setup

### Phase 3: Auth Tests
1. Implement `auth_test.go` - Register, Login, RefreshToken, GetMe, Logout
2. Test cases: success, validation errors, duplicate email, invalid credentials, unauthorized

### Phase 4: Core Entity Tests
1. Implement `alert_test.go` - CRUD, List, Acknowledge, Close, Snooze, Assign
2. Implement `team_test.go` - CRUD, member management
3. Implement `user_test.go` - ListOrganizationUsers

### Phase 5: Complex Entity Tests
1. Implement `schedule_test.go` - schedules, rotations, participants, overrides, on-call
2. Implement `escalation_test.go` - policies, rules, targets
3. Implement `notification_test.go` - channels, preferences, send, logs

### Phase 6: Advanced Tests
1. Implement `incident_test.go` - incidents, responders, timeline, alert linking
2. Implement `webhook_test.go` - endpoints, deliveries, incoming tokens

---

## Key Files to Create/Modify

### New Files
| File | Purpose |
|------|---------|
| `docker-compose.test.yml` | Test PostgreSQL (port 5434, tmpfs) |
| `tests/integration/testutils/config.go` | Test config loader |
| `tests/integration/testutils/database.go` | DB helpers, migrations, truncate |
| `tests/integration/testutils/server.go` | Test server with full router |
| `tests/integration/testutils/client.go` | HTTP client with auth |
| `tests/integration/testutils/fixtures.go` | Test data factories |
| `tests/integration/main_test.go` | TestMain setup |
| `tests/integration/*_test.go` | 9 test files for handlers |

### Modified Files
| File | Changes |
|------|---------|
| `Makefile` | Add test-db-up, test-db-down, test-integration targets |
| `go.mod` | Add golang-migrate dependency |

---

## Test Patterns

### Test Structure
```go
func TestEndpoint_Scenario(t *testing.T) {
    cleanDatabase(t)                           // Reset DB
    ctx := context.Background()
    client := newTestClient(t)

    user, _ := testFixtures.CreateUniqueUser(ctx)  // Create test user
    client.SetAuthToken(user.AccessToken)          // Authenticate

    resp := client.Post("/api/v1/endpoint", body)  // Make request
    client.ExpectStatus(resp, http.StatusCreated)  // Assert status

    var result map[string]interface{}
    client.ParseJSON(resp, &result)                // Parse response
    // Assert response body...
}
```

### Database Isolation
- Each test calls `cleanDatabase(t)` to truncate all tables
- Tables truncated in reverse dependency order (FK-safe)
- Unique user/org created per test via `CreateUniqueUser()`

---

## Makefile Targets

```makefile
test-db-up:          # Start test PostgreSQL
test-db-down:        # Stop and remove test database
test-integration:    # Run all integration tests
test-coverage:       # Run with coverage report
```

---

## Dependencies to Add

```
github.com/golang-migrate/migrate/v4
```

---

## Reference Files

- `cmd/api/main.go` - Route setup to replicate in test server
- `internal/middleware/auth.go` - JWT token structure
- `internal/repository/interfaces.go` - Repository interfaces for fixtures
- `internal/service/auth_service.go` - Registration/login patterns
- `docker-compose.yml` - Pattern for test compose file
