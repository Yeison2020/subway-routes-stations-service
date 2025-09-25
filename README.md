# Subway Routing Service

A **RESTful API** for computing subway routes and stations using the MBTA API.  
This service models the subway network as a **graph**, supports caching, and provides **human-readable route directions**.

---

## Table of Contents

- [Features](#features)  
- [Architecture](#architecture)  
- [API Endpoints](#api-endpoints)  
- [Installation](#installation)  
- [Environment Variables](#environment-variables)  
- [Usage](#usage)  
- [Example Response](#example-response)  
- [Docker](#docker)  
- [Testing](#testing)  
- [How the Graph and BFS Works](#how-the-graph-and-bfs-works)  

---

## Features

- Fetch all subway routes and stations from the MBTA API.  
- Cache station data to reduce API calls and improve performance.  
- Compute valid subway routes between two stations using **BFS graph traversal**.  
- Generate human-readable directions with transfers.  
- Handles invalid stations and disconnected routes gracefully.  
- Containerized with Docker for easy deployment.  

---

## Architecture

```
+---------------------+
|    MBTA API Client  |
|---------------------|
| Fetch routes & stops|
| Caching (in-memory) |
+----------+----------+
           |
           v
+---------------------+
|      Graph Module   |
|---------------------|
| Nodes = stations    |
| Edges = connections |
| BFS traversal       |
+----------+----------+
           |
           v
+-------------------------------------------+
|               API Handlers                |
+-------------------+-----------------------+--------------------------------------+
| Method / Endpoint  | Handler Function      | Description                          |
+-------------------+-----------------------+--------------------------------------+
| GET /api/v1/health | HealthHandler        | Simple health check (returns 200 OK) |
| GET /api/v1/subways| GetSubwaysHandler    | List all subway lines and stations   |
| GET /api/v1/routes | GetRouteHandler      | Compute a route between two stations |
+-------------------+-----------------------+--------------------------------------+
```

---

## API Endpoints

## GET api/v1/heatlhz

- Returns ok numbers of routes and status of API
- **Example Request:**  
```
GET /api/v1/heatlhz
```
- **Response:**
```json
{
"mbta": "ok, 8 routes",
"message": "API is running",
"status": "healthy"
}
```
---
## GET api/v1/subways

- Returns all subway routes and their stations.  
- **Example Request:**  

```
GET /api/v1/subways
```

- **Response:**

```json
[
  {
    "id": "Red",
    "name": "Red Line",
    "stations": ["Alewife", "Davis", "Porter", "..."]
  }
]
```

---

## GET /api/v1/routes?start=<station>&end=<station>

- Computes a valid subway route between two stations, including transfers.  
- **Query Parameters:**  
  - `start` – starting station name  
  - `end` – destination station name  

- **Example Request:**

 - Accepts stations Ids
```
GET /api/v1/routes?start=place-bmmnl&end=place-coecl
```

- **Response:**

```json
{   "description": [
        "Path 1",
        "Path 2",
        "Path 3"
    ],
"lines": [
        "Blue Line",
        "Orange Line",
        "Red Line",
    ],
    "stations": [
        "Beachmont",
        "Suffolk Downs",
        "Orient Heights",
    ]}
```

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/subway-routes-stations-service.git
cd subway-routes-stations-service
```

2. Install dependencies (Go modules):

```bash
go mod tidy
```

3. Set environment variables (see below).  

4. Run the server:

```bash
go run cmd/server/main.go
```

The server runs on `http://localhost:8080` by default.  

---

## Environment Variables

| Variable       | Description                     |
|----------------|---------------------------------|
| `MBTA_API_KEY` | Your MBTA API key               |

---

## Usage

- Health endpoint `GET api/v1/healthz`
- List subway lines: `GET api/v1/subway`  
- Compute a route: `GET /route?start=<station>&end=<station>`  

-- 

## Swagger 

Navigate to ```http://localhost:8080/swagger/index.html```

---

## Docker

Build and run the app using Docker:

```bash
docker build -t subway-routes-stations-service .
docker run -p 8080:8080 --env MBTA_API_KEY=your_api_key subway-routes-stations-service
```

---

## Testing

- Unit tests for graph traversal, caching, and API handlers.  
- Run tests with:

```bash
go test ./internal/tests
```
---

## Concerns to consider if deploying to production

1. **Cache Mechanisms:**

- Issue: Expired keys are never removed; they just stay in the map until overwritten.
  - Improvement: Add a cleanup mechanism (background job) to cleanup expired items, or consider a library like Ristretto

- Issue: RWMutex locks the entire map, which can become a bottleneck under high traffic.
  - Improvement: Use sharded locks (lock per key)

- Issue: TTL is the same for all keys.
  - Improvement: Support per-key TTLs by attaching an expiration time to each entry. Default TTL can remain as a fallback.

- Issue: No monitoring or visibility into cache performance.
  - Improvement: Expose custom metrics (e.g. memory usage) to track cache effectiveness.

2. **Factors outside of the code:**
- Rate Limiting & Throttling: Protect the service from overload and abusive traffic.
- Load Balancing: Distribute requests across multiple server instances to prevent bottlenecks.



## How the Graph and BFS Works

1. **Graph Modeling:**  
   - Each station = node  
   - Each connection = edge with line name  
   - Bidirectional edges for both directions  

2. **BFS Traversal:**  
   - Explore stations level by level  
   - Track Path (stations) and Lines (lines taken) for each potential route.
   - Continue exploring even after finding one path to collect multiple possible routes (limit can be applied to avoid infinite results).

3. **Handling Transfers and Errors:**  
   - Detect line changes to generate transfer points for each route.
   - If a station does not exist or no path exists → return 404.
   - Return all valid routes in the response, each with its stations, lines, and human-readable description.

4. **Route Description:**
   - For each route, generate a description:
     - Start station → line taken → transfer points → end station.
   - Each route has its own stations, lines, and description.
---

**Author:** Yeison Casado

