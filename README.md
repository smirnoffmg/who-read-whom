# Literary Opinions Graph

A web application that visualizes relationships between writers and literary works based on documented opinions.

## Concept

Writers often expressed views about each other's works — admiration, criticism, or outright disdain. This project captures those connections as a graph where:

- **Nodes** are writers and their works
- **Edges** are documented opinions (positive or negative)
- **Each edge includes proof**: a quote, source, and context

## Tech Stack

- **Frontend**: React
- **Backend**: Go
- **Database**: PostgreSQL

## Data Model

Three entities: `Writer`, `Work`, and `Opinion`. Writers create works; writers express opinions about other writers' works. Each opinion is backed by a verifiable source.