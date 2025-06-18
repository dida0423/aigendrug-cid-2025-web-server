# AIGENDRUG Platform Server

## Introduction

**Dedicated Backend Server** for**AIGENDRUG Tool Selection and Interaction Platform**. This platform provides an integrated web interface for interacting with AI tools developed for drug discovery and experimental optimization. It allows users to chat, receive tool recommendations, input data, and view tool results — all from a single, streamlined dashboard. You can find the client [here](https://github.com/khinwaiyan/aigendrug-cid-2025-web-client/blob/main/README.md).

## Getting Started

### Prerequisites

To use this application, make sure you have Docker and Docker Compose installed.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/dida0423/aigendrug-cid-2025-web-server.git
   cd aigendrug-cid-2025-web-server
   ```

2. Set up environmental variables in a new `.env` file.

```.env
PORT=8080
RUN_MODE=debug

DB_CONNECTION_STRING="postgres://postgres:postgres@db:5432/aigendrug?sslmode=disable&options=-c%20search_path=ks_admin"
MAIN_DB_HOST=db
MAIN_DB_PORT=5432
MAIN_DB_USER=postgres
MAIN_DB_PASSWORD=postgres
MAIN_DB_NAME=aigendrug
MAIN_DB_SCHEMA=ks_admin

TOOL_ROUTER_HOST=https://router-aigendrug-cid-2025.luidium.com

OPENAI_API_KEY=
```

3. Build and run with Docker Compose

   ```bash
   docker-compose up --build [-d]
   ```

---
