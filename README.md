# Job Apply Automator

> A production-ready career automation platform that helps college students and working professionals find verified job opportunities, understand companies, match jobs with the right resume, apply efficiently, reach out for referrals, automate follow-ups, and track the complete job-search journey.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Product Flow](#product-flow)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Folder Structure](#folder-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Docker Setup](#docker-setup)
- [Environment Variables](#environment-variables)
- [API Reference](#api-reference)
- [Database Migrations](#database-migrations)
- [Development Workflow](#development-workflow)
- [Roadmap](#roadmap)

---

## Project Overview

The job search process is broken for students and early professionals. It means:

- Manually checking LinkedIn, Wellfound, Naukri, company career pages every day
- Maintaining multiple tailored resumes with no central place to manage them
- Writing cold emails and follow-ups from scratch every time
- Having no visibility into whether a job listing is real, expired, or a scam

**Job Apply Automator** solves all of this in one platform:

```
FIND → VERIFY → SELECT RESUME → MATCH → APPLY → OUTREACH → FOLLOW UP → TRACK
```

---

## Product Flow

```
User
  │
  ▼
Career Profile + Saved Resumes
  │
  ▼
Job Discovery  ──────────────────────────────────────────────────────────────────┐
(Public job fetching APIs · Company Career Pages · ATS feeds)              │
  │                                                                              │
  ▼                                                                         Feedback Loop
Normalize + Deduplicate                                                          │
  │                                                                              │
  ▼                                                                              │
Job Evidence / Source Verification                                               │
(Source URL · Application destination · Last checked · Salary evidence)         │
  │                                                                              │
  ▼                                                                              │
Company Passport                                                                 │
(Identity · Size · Workplace signals · Review evidence · Source links)          │
  │                                                                              │
  ▼                                                                              │
Select or Upload Resume for THIS Job                                             │
  │                                                                              │
  ▼                                                                              │
Resume-to-Job Matching                                                           │
(Skills · Experience · Role · Location · Salary · Education · Gaps)             │
  │                                                                              │
  ▼                                                                              │
Save / Apply / Ask for Referral                                                  │
  │                                                                              │
  ▼                                                                              │
Application + Outreach Automation                                                │
(Email templates · AI drafts · Scheduled send · Gmail/Microsoft OAuth)          │
  │                                                                              │
  ▼                                                                              │
Follow-ups + Notifications                                                       │
  │                                                                              │
  ▼                                                                              │
Application Outcome  ────────────────────────────────────────────────────────────┘
```

---

## Features

| Feature | Description |
|---|---|
| **Smart Job Discovery** | Aggregate jobs from company career pages, ATS feeds, and permitted job platforms |
| **Job Evidence** | Show original source URL, last verified time, application destination, and salary evidence |
| **Company Passport** | Company identity, size, workplace signals, review evidence, and source links per job |
| **Resume-first Matching** | User selects or uploads the resume for each specific job; system scores the match |
| **Application Automation** | User-controlled automation rules — prepare and submit applications when criteria are met |
| **Outreach Automation** | AI-assisted referral emails, recruiter outreach, scheduled sending, automatic follow-ups |
| **Application Tracker** | Track every application across states: Applied → OA → Interview → Offer → Rejection |
| **Feedback Loop** | Application outcomes feed back into future matching to improve recommendations |

---

## Architecture

The first version is a **modular monolith with separate background workers**. This is easier to build and test than starting with microservices. High-load modules (scraper, matching, outreach) can be extracted later when needed.

```
                          ┌─────────────────────────────────┐
                          │           WEB BROWSER            │
                          └────────────────┬────────────────┘
                                           │ HTTPS
                          ┌────────────────▼────────────────┐
                          │     Next.js + Tailwind CSS       │
                          │        (Frontend — SSR)          │
                          └────────────────┬────────────────┘
                                           │ REST API
                          ┌────────────────▼────────────────┐
                          │         Load Balancer            │
                          └────────────────┬────────────────┘
                                           │
                          ┌────────────────▼────────────────┐
                          │     Go API Server (Gin)          │
                          │   cmd/api/main.go   │
                          └──────┬──────────────┬───────────┘
                                 │              │
               ┌─────────────────▼──┐    ┌──────▼──────────────┐
               │   Business Logic   │    │   Kafka Event Bus    │
               │  (service layer)   │    │  (async events)      │
               └─────────┬──────────┘    └──────┬───────────────┘
                         │                      │
          ┌──────────────▼───────┐   ┌──────────▼──────────────┐
          │   PostgreSQL (GORM)  │   │   Background Workers     │
          │  Primary source of   │   │  cmd/worker/main.go      │
          │       truth          │   │  · Job ingestion         │
          └──────────────────────┘   │  · Resume matching       │
                                     │  · Email scheduler       │
          ┌───────────────────────┐  │  · Notifications         │
          │   Redis               │  └──────────────────────────┘
          │  · Session cache      │
          │  · Rate limiting      │  ┌──────────────────────────┐
          │  · Distributed locks  │  │   OpenSearch              │
          └───────────────────────┘  │  · Job/company search    │
                                     │  · Filters + facets      │
          ┌───────────────────────┐  └──────────────────────────┘
          │   Vector DB           │
          │  · Semantic resume    │  ┌──────────────────────────┐
          │    ↔ JD matching      │  │   AWS S3                 │
          └───────────────────────┘  │  · Resume file storage   │
                                     └──────────────────────────┘
```

### Application State Machine

Every job application moves through these states:

```
DISCOVERED → MATCHED → QUALIFIED → RESUME_SELECTED → APPLICATION_READY → AWAITING_APPROVAL → SUBMITTED
                                                                                                    │
                                                                                          Recruiter reaches out
```

Automation rules are user-controlled. A fully automatic action only fires when the user has explicitly enabled it and the source permits it.

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | Next.js 14 (App Router) · TypeScript · Tailwind CSS · Zustand |
| **Backend API** | Go 1.23 · Gin · GORM · JWT |
| **Background Workers** | Go · Kafka consumers |
| **Primary Database** | PostgreSQL 16 |
| **Cache / Locks** | Redis |
| **Event Bus** | Apache Kafka |
| **Search** | OpenSearch |
| **Semantic Matching** | Embeddings + Vector DB (pgvector / Qdrant) |
| **File Storage** | AWS S3 |
| **Email** | Gmail OAuth · Microsoft Graph (no password stored) |
| **Containerization** | Docker · Docker Compose |
| **Observability** | OpenTelemetry · Prometheus · Grafana · Zap (structured logs) |

---

## Folder Structure

```
job-apply-automator/
│
├── frontend/                              # Next.js 14 + TypeScript + Tailwind
│   └── src/
│       ├── app/                           # App Router pages
│       │   ├── (auth)/                    # Login, Signup (no dashboard layout)
│       │   └── (dashboard)/               # Protected pages with sidebar
│       │       ├── discover/              # Job discovery feed
│       │       ├── jobs/[id]/             # Job detail + evidence + match
│       │       ├── companies/[id]/        # Company Passport
│       │       ├── resumes/               # Resume management
│       │       ├── applications/          # Application tracker
│       │       ├── outreach/              # Email campaigns
│       │       └── profile/              # Career profile + preferences
│       ├── components/
│       │   ├── ui/                        # Base design system (Button, Input, Modal…)
│       │   ├── layout/                    # Sidebar, Navbar, PageWrapper
│       │   └── features/                  # Domain components (JobCard, MatchBreakdown…)
│       ├── hooks/                         # Custom React hooks
│       ├── services/                      # HTTP calls to Go backend
│       ├── store/                         # Zustand global state
│       ├── types/                         # TypeScript interfaces
│       └── utils/                         # cn(), formatDate(), constants
│
├── backend/                               # Go — Modular Monolith
│   ├── cmd/
│   │   ├── api/main.go                    # API server entry point
│   │   └── worker/main.go                 # Background worker entry point
│   │
│   ├── config/
│   │   └── config.go                      # Typed env config (loaded once at startup)
│   │
│   ├── internal/                          # Private app code (not importable externally)
│   │   ├── domain/                        # Core entities — pure structs, no DB/HTTP logic
│   │   │   ├── user.go
│   │   │   ├── job.go
│   │   │   ├── company.go
│   │   │   ├── resume.go
│   │   │   ├── match.go
│   │   │   ├── application.go
│   │   │   └── outreach.go
│   │   │
│   │   ├── repository/                    # Data access layer
│   │   │   ├── interface.go               # Repository interfaces (contracts)
│   │   │   ├── postgres/                  # PostgreSQL implementations
│   │   │   └── mock/                      # Mock implementations for testing
│   │   │
│   │   ├── service/                       # Business logic layer
│   │   │   └── auth_service.go
│   │   │
│   │   ├── api/                           # HTTP layer
│   │   │   ├── router.go                  # Route registration
│   │   │   ├── handlers/                  # Request → Service → Response
│   │   │   └── middleware/                # JWT auth, CORS, rate limiting, logging
│   │   │
│   │   ├── db/
│   │   │   └── postgres.go                # DB connection + pool config
│   │   │
│   │   ├── worker/                        # Background worker processors
│   │   ├── scraper/                       # Job source adapters + normalizer
│   │   │   └── sources/                   # linkedin.go, wellfound.go, naukri.go…
│   │   ├── evidence/                      # Source verification logic
│   │   ├── matching/                      # Resume ↔ JD scoring engine
│   │   ├── outreach/                      # Email composer + follow-up scheduler
│   │   │   └── provider/                  # gmail.go, microsoft.go
│   │   └── workflow/                      # Application state machine + automation rules
│   │
│   ├── pkg/                               # Reusable shared packages
│   │   ├── logger/                        # Zap structured logger wrapper
│   │   ├── response/                      # Standard API response helpers
│   │   ├── validator/                     # Input validation + error formatting
│   │   ├── pagination/                    # Page/limit/offset parsing
│   │   └── events/                        # Kafka producer/consumer wrapper
│   │
│   └── migrations/                        # Raw SQL migration files (run in order)
│       ├── 001_create_users.sql
│       ├── 002_create_jobs.sql
│       └── ...
│
├── docker/
│   ├── frontend.Dockerfile
│   ├── backend.Dockerfile
│   └── worker.Dockerfile
│
├── docker-compose.yml                     # Full local dev stack
├── .gitignore
└── README.md
```

---

## Getting Started

### Prerequisites

Make sure you have these installed:

| Tool | Version | Check |
|---|---|---|
| Go | 1.23+ | `go version` |
| Node.js | 18+ | `node --version` |
| PostgreSQL | 15+ | `psql --version` |
| Redis | 7+ | `redis-cli --version` |
| Docker | 24+ | `docker --version` |

---

### Backend Setup

**1. Clone the repo**

```bash
git clone https://github.com/Snehaaa-Kri/job-apply-automator.git
cd job-apply-automator
```

**2. Install Go dependencies**

```bash
cd backend
go mod tidy
```

**3. Set up environment variables**

```bash
cp .env.example .env.local
# Open .env.local and fill in your values (DB password, JWT secret, etc.)
```

**4. Set up PostgreSQL**

Option A — Local install (Linux):
```bash
sudo apt update && sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql

# Create database and set password
sudo -i -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'yourpassword';"
sudo -i -u postgres psql -c "CREATE DATABASE job_automator;"
```

Option B — Docker (quicker):
```bash
docker run --name job-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=job_automator \
  -p 5432:5432 -d postgres:16
```

**5. Run database migrations**

```bash
# Always run in order
psql -U postgres -d job_automator -f migrations/001_create_users.sql
```

**6. Start the API server**

```bash
go run cmd/api/main.go
```

Server starts at: `http://localhost:8080`

Verify it's running:
```bash
curl http://localhost:8080/health
# → {"status":"ok"}
```

---

### Frontend Setup

> Frontend is built with Next.js and will be set up separately after the backend is stable.

```bash
cd frontend
npm install
cp .env.example .env.local
# Set NEXT_PUBLIC_API_URL=http://localhost:8080
npm run dev
```

Frontend runs at: `http://localhost:3000`

---

### Docker Setup

To run the full stack (API + Workers + Postgres + Redis + Kafka) in one command:

```bash
# From the repo root
docker compose up --build
```

Services started:
| Service | Port |
|---|---|
| API Server | `8080` |
| Frontend | `3000` |
| PostgreSQL | `5432` |
| Redis | `6379` |
| Kafka | `9092` |

To stop:
```bash
docker compose down
```

To stop and wipe all data (careful):
```bash
docker compose down -v
```

---

## Environment Variables

All environment variables live in `backend/.env.local` (never committed to git).

Use `backend/.env.example` as the template.

| Variable | Description | Example |
|---|---|---|
| `APP_ENV` | Environment | `development` |
| `APP_PORT` | API server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | DB username | `postgres` |
| `DB_PASSWORD` | DB password | `yourpassword` |
| `DB_NAME` | Database name | `job_automator` |
| `DB_SSL_MODE` | SSL mode | `disable` (dev) |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `JWT_SECRET` | JWT signing secret | `min 32 chars, random` |
| `JWT_EXPIRY_HOURS` | Token expiry | `24` |
| `KAFKA_BROKERS` | Kafka broker addresses | `localhost:9092` |
| `AWS_REGION` | S3 region | `ap-south-1` |
| `AWS_BUCKET_NAME` | Resume storage bucket | `job-automator-resumes` |
| `GMAIL_CLIENT_ID` | Gmail OAuth client ID | from Google Cloud Console |
| `GMAIL_CLIENT_SECRET` | Gmail OAuth secret | from Google Cloud Console |

---

## API Reference

Base URL: `http://localhost:8080/api/v1`


**Rules:**
- Never edit a migration file once it has been run in any environment
- To fix a mistake, create a new migration: `002_fix_users_add_phone.sql`
- Each file is idempotent (`CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`)

---

### Build & Check

```bash
# Compile everything — if no output, no errors
go build ./...

# Run the API server
go run cmd/api/main.go

# Run the background worker
go run cmd/worker/main.go

# Run tests
go test ./...
```
---

## License

MIT — see [LICENSE](LICENSE) for details.
