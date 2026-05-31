# Mirage

A lightweight HTTP server that captures and logs incoming probes to PostgreSQL, with a real-time admin dashboard for reviewing attack data.

Built in Go using only the standard library and minimal dependencies. No frameworks.

**Live demo:** [mirage-oz0l.onrender.com](https://mirage-oz0l.onrender.com)  
**Admin panel:** [mirage-oz0l.onrender.com/admin](https://mirage-oz0l.onrender.com/admin)

---

## What it does

Mirage runs as an HTTP server that responds convincingly to every incoming request while silently logging the source IP, HTTP method, path, and timestamp to a PostgreSQL database. A protected admin panel provides a real-time view of all captured traffic with filtering, pagination, and severity classification.

Paths like `/etc/passwd`, `/.env`, and `/wp-config` are flagged as dangerous. Paths like `/admin`, `/login`, and `/wp-admin` are flagged as suspicious. Everything else is logged as normal traffic.

---

## Features

- Captures IP address, HTTP method, path, and timestamp for every request
- Responds with a convincing HTML page to avoid detection
- Admin dashboard with live metrics — total probes, last 24 hours, unique IPs, top path
- Filter by IP, HTTP method, or path
- Pagination — 25 rows per page
- Severity colour-coding — red for dangerous paths, amber for suspicious, grey for normal
- Basic Auth protection on the admin panel
- Separate honeypot and admin servers to prevent cross-contamination
- Environment variable configuration — no secrets in source code
- Cloud PostgreSQL via Neon — data persists across deployments

---

## Tech stack

- **Language:** Go 1.22
- **Database:** PostgreSQL (local) / Neon (production)
- **Driver:** github.com/lib/pq
- **Config:** github.com/joho/godotenv
- **Deployment:** Render
- **Templates:** Go `html/template` — server-side rendered
- **Auth:** HTTP Basic Auth middleware — no external packages

---

## Project structure

```
mirage/
├── main.go          — entry point, database connection, server startup
├── server.go        — HTTP handlers, middleware, route registration
├── storage.go       — Store interface, PostgreSQL implementation
├── config.go        — environment variable loading
├── templates/
│   └── admin.html   — admin dashboard template
└── static/
    └── favicon.ico
```

---

## What I learned building this

This project was built across 10 deliberate versions, each targeting one new concept:

| Version | Concept |
|---|---|
| v1 | `database/sql`, connection strings, `db.Ping`, `db.Exec`, `db.QueryRow` |
| v2 | Struct-based DB models, `rows.Scan`, thinking in types |
| v3 | Go interfaces — definition, implicit satisfaction, dependency injection |
| v4 | `net/http` — live request capture, `r.RemoteAddr`, method and path extraction |
| v5 | `encoding/json`, `json.NewEncoder`, serving data over HTTP |
| v6 | `html/template`, server-side rendering, separating logic from presentation |
| v7 | URL query params, dynamic SQL `WHERE` clauses, pagination |
| v8 | HTTP middleware pattern — wrapping handlers, `http.Handler`, `http.HandlerFunc` |
| v9 | Environment variables, 12-factor app config, `.env` vs production injection |
| v10 | Cloud PostgreSQL, Render deployment, production config |

---

## Running locally

**Prerequisites:** Go 1.22+, PostgreSQL

**1. Clone the repo**

```bash
git clone https://github.com/t-Shack/mirage.git
cd mirage
```

**2. Create the database**

```sql
CREATE DATABASE mirage;

\c mirage

CREATE TABLE IF NOT EXISTS requests (
    id        SERIAL PRIMARY KEY,
    ip        TEXT NOT NULL,
    method    TEXT NOT NULL,
    path      TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);
```

**3. Create a `.env` file**

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mirage
SERVER_PORT=8080
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password
```

**4. Run**

```bash
go run .
```

Visit `http://localhost:8080` to trigger the honeypot.  
Visit `http://localhost:8080/admin` for the dashboard.

---

## Deployment

Mirage is deployed on [Render](https://render.com) with a [Neon](https://neon.tech) PostgreSQL database.

Required environment variables on the host:

| Variable | Description |
|---|---|
| `DATABASE_URL` | Full PostgreSQL connection string |
| `ADMIN_USERNAME` | Admin panel username |
| `ADMIN_PASSWORD` | Admin panel password |

The `render.yaml` in the repo root handles build and start configuration automatically.

---

## Other projects

- **Sniffer** — phishing URL detection tool. Live at [sniffer-3g7n.onrender.com](https://sniffer-3g7n.onrender.com)
- **Port Scanner** — concurrent TCP port scanner with goroutines, channels, and file I/O. [GitHub](https://github.com/t-Shack/port-scanner)
