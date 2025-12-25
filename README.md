# Literary Opinions Graph

A web application that visualizes relationships between writers and literary works based on documented opinions.

## Concept

Writers often expressed views about each other's works — admiration, criticism, or outright disdain. This project captures those connections as a graph where:

- **Nodes** are writers and their works
- **Edges** are documented opinions (positive or negative)
- **Each edge includes proof**: a quote, source, and context

## Tech Stack

- **Frontend**: Next.js (React + TypeScript)
- **Backend**: Go
- **Database**: PostgreSQL

## Data Model

Three entities: `Writer`, `Work`, and `Opinion`. Writers create works; writers express opinions about other writers' works. Each opinion is backed by a verifiable source.

## Quick Start with Docker

### Prerequisites

- Docker and Docker Compose installed

### Setup

1. Copy environment file:
```bash
cp .env.example .env
```

2. Start all services:
```bash
docker compose up -d
```

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api/v1
   - PostgreSQL: localhost:5432

### Stop Services

```bash
docker compose down
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### Rebuild Services

```bash
docker compose up -d --build
```

### Development Notes

- The frontend connects to the backend using the service name `backend` within Docker network
- For local development outside Docker, set `NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1`
- Database data persists in the `postgres_data` volume