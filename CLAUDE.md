# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **security-hardened web forum** written in Go, featuring HTTPS/TLS, rate limiting, session management, and bcrypt password hashing. The application is designed to run on both PaaS platforms (Heroku, Render, Railway, Fly.io, etc.) and VPS/bare metal servers with automatic environment detection.

## Build and Run

```bash
# Build the application
go build -o forum ./cmd/web

# Run locally (HTTP mode)
PORT=8080 ./forum

# Run with autocert (requires VPS with ports 80/443 and domain)
sudo USE_AUTOCERT=true DOMAIN=yourdomain.com ./forum

# Build and run with Docker
docker build -t forum .
docker run -p 8080:8080 -e PORT=8080 forum
```

## Database

The application uses **SQLite** with the database file at `./forum.db`. The schema is auto-created on startup via `forum.InitDB()` in `forum/setupDB.go`.

### Schema:
- **users**: username (PK), image, email, password (bcrypt hash), access, loggedIn, liked/disliked posts/comments
- **sessions**: sessionID (UUID, PK), username (FK to users)
- **posts**: postID (auto-increment PK), author (FK), title, content, category, postTime, likes, dislikes, ips (view tracking)
- **comments**: commentID (auto-increment PK), author (FK), postID (FK), content, commentTime, likes, dislikes

## Architecture

### Deployment Mode Detection (cmd/web/main.go)

The application automatically detects the environment:

**PaaS Mode (Heroku, Render, etc.):**
- Triggered when `PORT` env var exists and `USE_AUTOCERT != "true"`
- Listens on HTTP only (platform handles TLS termination at edge/load balancer)
- Uses `runOnPaaS()` function
- Works with any PaaS that sets `PORT` and provides TLS termination

**Autocert Mode (VPS/Bare Metal):**
- Triggered when no `PORT` or `USE_AUTOCERT=true`
- Binds to ports 80 (ACME challenges) and 443 (HTTPS)
- Uses Let's Encrypt via `golang.org/x/crypto/acme/autocert`
- Uses `runWithAutocert()` function

### Request Flow

1. **Rate Limiting**: All routes wrapped with `forum.RateLimiter()` middleware (forum/rateLimiter.go)
   - Global token bucket: 3 initial tokens, refills 1 every 200ms
   - **Note**: Current implementation is global for all users (not per-IP)

2. **Session Management** (forum/session.go):
   - UUIDv4 session tokens stored in `sessions` table
   - `loggedIn(r)` checks if request has valid session cookie
   - `obtainCurUserFormCookie(r)` retrieves current user from session
   - Sessions expire after 15 minutes (login) or 30 minutes (register)

3. **Handlers** (forum/handlers.go):
   - All handlers check `r.Method` and return 400 for invalid methods
   - `loggedIn(r)` redirects to home if already authenticated (login/register)
   - Template execution with `.gohtml` files in `templates/`

### Like/Dislike System (forum/LikesAndDislikes.go, forum/processPostsAndComments.go)

Uses a **string-based toggle mechanism**:
- User likes/dislikes stored as delimiter-separated strings: `"-postID1-postID2-postID3"`
- Clicking like/dislike appends to string
- `CountLikesByUser()` counts odd occurrences (toggle logic: odd = active, even = cancelled)
- `SumOfAllLikes()` aggregates across all users
- `DistLikesToPosts()` distributes counts to post objects for display

### Post Filtering (forum/filter.go)

Three filter types:
- **Category**: Uses SQL `LIKE '%($category)%'` on category column
- **Author**: Direct username match
- **Liked**: Filters posts by current user's liked posts

Categories are stored as `(Category1)(Category2)` format.

## Security Considerations

### Current Implementation:
- ✅ Bcrypt password hashing (cost 10)
- ✅ Parameterized SQL queries (SQL injection protection)
- ✅ UUID session tokens
- ✅ TLS cipher suite configuration (when using autocert)
- ✅ Go's `html/template` auto-escapes template variables

### Known Issues to Address:
- ❌ **No CSRF protection** - all state-changing forms vulnerable
- ❌ **Session cookies lack security flags** - need `HttpOnly`, `Secure`, `SameSite`
- ❌ **XSS via unquoted HTML attributes** - template attributes should be quoted
- ❌ **Plaintext password/session logging** - remove debug logs in production
- ❌ **Global rate limiter** - should be per-IP or per-user
- ❌ **Weak email validation** - only checks for `@` and `.`

### Cookie Configuration (forum/session.go, forum/register.go):

Current cookies need security flags added:
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    sid.String(),
    MaxAge:   900,
    HttpOnly: true,                    // ADD: Prevent JavaScript access
    Secure:   true,                    // ADD: HTTPS only
    SameSite: http.SameSiteStrictMode, // ADD: CSRF protection
})
```

## Template Structure

Templates use Go's `html/template` with composition:
- `header.gohtml` / `header2.gohtml`: Navigation (logged out / logged in)
- `footer.gohtml`: Footer
- `index.gohtml` / `index2.gohtml`: Home page (logged out / logged in)
- `post.gohtml` / `post2.gohtml`: Post detail page (logged out / logged in)

Template data passed via structs in `forum/structs.go`.

## Dependencies

Key packages (see `go.mod`):
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/satori/go.uuid` - UUID generation for sessions
- `golang.org/x/crypto/bcrypt` - Password hashing
- `golang.org/x/crypto/acme/autocert` - Let's Encrypt integration

## Deployment

### Heroku
```bash
git push heroku master
# Automatically detects PORT env var → PaaS mode
# Enable free SSL: heroku certs:auto:enable
```

### Render
```bash
# Connect your Git repository in Render dashboard
# Build Command: go build -o app ./cmd/web
# Start Command: ./app
# Automatically detects PORT env var → PaaS mode
# SSL enabled by default (free)
```

### VPS/Bare Metal
```bash
# Set domain and enable autocert
export USE_AUTOCERT=true
export DOMAIN=yourdomain.com
sudo ./forum  # Requires sudo for ports 80/443
```

## Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `PORT` | Server port (set by PaaS platforms) | Triggers autocert mode if unset |
| `USE_AUTOCERT` | Force autocert mode | `false` |
| `DOMAIN` | Domain for Let's Encrypt | `www.elephorum.com` |

## File Organization

```
cmd/web/
  main.go     - Entry point with environment detection
  mainS.go    - Legacy HTTPS config (can be deleted)
forum/
  setupDB.go              - Database initialization
  structs.go              - Core data structures
  session.go              - Authentication & session management
  register.go             - User registration
  handlers.go             - HTTP route handlers
  rateLimiter.go          - Rate limiting middleware
  processPostsAndComments.go - POST request processing
  LikesAndDislikes.go     - Like/dislike logic
  displayPostsAndComments.go - Data retrieval
  filter.go               - Post filtering
  server.go               - Server configuration (timeouts)
templates/                - HTML templates (.gohtml)
assets/                   - Static files (CSS, images)
```
