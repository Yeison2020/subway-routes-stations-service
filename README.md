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

### GET /subway

- Returns all subway routes and their stations.  
- **Example Request:**  

```
GET /subway
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

### GET /route?start=<station>&end=<station>

- Computes a valid subway route between two stations, including transfers.  
- **Query Parameters:**  
  - `start` – starting station name  
  - `end` – destination station name  

- **Example Request:**

```
GET /route?start=Alewife&end=Government+Center
```

- **Response:**

```json
{
  "stations": ["Alewife","Davis","Porter","Harvard","Central","Kendall/MIT","Charles/MGH","Park Street","Government Center"],
  "lines": ["Red Line","Red Line","Red Line","Red Line","Red Line","Red Line","Red Line","Green Line B"],
  "description": "Start at Alewife, transfer at Park Street to Green Line B, take Green Line B to Government Center."
}
```

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/subway-routing-service.git
cd subway-routing-service
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

- List subway lines: `GET /subway`  
- Compute a route: `GET /route?start=<station>&end=<station>`  

---

## Docker

Build and run the app using Docker:

```bash
docker build -t subway-routing-service .
docker run -p 8080:8080 --env MBTA_API_KEY=your_api_key subway-routing-service
```

---

## Testing

- Unit tests for graph traversal, caching, and API handlers.  
- Run tests with:

```bash
go test ./internal/tests
```
---

## How the Graph and BFS Works

1. **Graph Modeling:**  
   - Each station = node  
   - Each connection = edge with line name  
   - Bidirectional edges for both directions  

2. **BFS Traversal:**  
   - Explore stations level by level  
   - Track `Path` (stations) and `Lines` (lines taken)  
   - Stops when destination is reached  

3. **Handling Transfers and Errors:**  
   - Detect line changes to generate transfers  
   - If station not found or no path exists → return 404  

---

**Author:** Yeison Casado

