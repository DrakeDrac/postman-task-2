# Campus Leave & Attendance Management System

## Tech Stack
- **Language:** Go 1.25 + Gin
- **Database:** PostgreSQL 
- **ORM:** GORM
- **Authentication:** JWT
- **Containerisation:** Docker & Docker Compose

---

## Setup & Usage

### 1 · Clone the repo
```bash
git clone https://github.com/DrakeDrac/postman-task-2
cd postman-task-2
```

### 2 · Start the stack with Docker Compose
Builds the Go API and starts both the API and Postgres database.
```bash
docker-compose up --build
```
* `db`  – Postgres on port **5432** (internal) → localhost:5432
* `app` – Go API on port **8080**            → http://localhost:8080
The API waits for the database health-check before launching.

### 3 · Verify the service
```bash
curl http://localhost:8080/health
# → {"message":"Server is running","status":"ok"}
```

### 4 · Run sample tests (optional)
A quick script that exercises typical flows.
```bash
chmod +x test_api.sh
./test_api.sh
```

### 5 · Stopping / cleaning up
```bash
# stop containers
docker-compose down

# stop + remove the database volume
docker-compose down -v
```