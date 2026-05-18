# MAXAT Social Network Microservices 

This is a backend system for a social network built with Go, structured as individual microservices that talk to each other using gRPC. It also has a frontend built with React / Next.js.

## Project Structure

* **`api-gateway`** – The single entry point for all frontend requests. It translates incoming HTTP REST requests into internal **gRPC** calls to the backend services.
* **`auth-service`** – Handles user authentication, registration, issuing ***JWT tokens***, and checking credentials.
* **`user-service`** – Manages user profiles, lookups, and account details. Stores data in ***MongoDB***.
* **`content-service`** – Manages social feed elements like posts, comments, and likes. Uses ***Redis*** for caching frequently accessed feed data. 
* **`chat-service`** – Handles real-time messaging, group chat configurations, and chat histories.
* **`notification-service`** – Listens for asynchronous events (like someone liking a post or registering) using NATS to send out automated emails or alerts.
* **`frontend`** – The user interface built with ***React/Next.js*** to interact with the backend services via the API Gateway.
* **`prometheus` & `grafana`** – Configured to gather ***metrics and monitor*** system performance across containers.

**Protocol Buffers:** The shared gRPC definitions (`.proto` files and compiled `.pb.go` stubs) are stored in a separate dedicated repository named `maxat-protobuf`. 

---

## Technical Stack
* **Language:** Go (Golang)
* **Frontend:** React / Next.js
* **Communication:** gRPC (Internal service communication) & HTTP REST (API Gateway to Frontend)
* **Databases:** MongoDB (with replication for main storage) & Redis (for caching feed items)
* **Event Broker:** NATS (for async notifications)
* **Monitoring:** Prometheus & Grafana

---

## Project Setup

### 1. Requirements
* Docker 
* Go 1.23+
* Node.js

### 2. Environment Setup
Every service expects its own configuration variables (`api-gateway/.env`, `user-service/.env`, `content-service/.env`, etc.). 

An example for service `.env` file:
```env
PORT=8080
MONGO_URI=mongodb://mongo-replica:27017/dbname
REDIS_ADDR=redis:6379
NATS_URL=nats://nats:4222
JWT_SECRET=secret_key
```
### 3. Running the Application
From root folder, run:  
``` 
docker-compose up -d --build
``` 
To run individual service locally (example for content-service):
```
cd content-service
go run cmd/content-service/main.go
``` 
## Database Migrations
Project uses JSON configurations inside migrations/ folder setups instead of regular SQL files.  
It is used in:  
- **user-service/migrations/**  
  Contains schema definitions such as `000001_create_users.up.json` and rollback steps like `000001_create_users.down.json`. These json files apply strict administrative Mongo collection criteria ("$jsonSchema" rule structures matching properties like username, email, etc.).  
- **content-service/migrations/**   
  Handles the setup and indexes configuration for posts, comments, and likes collections (`000001_create_content_collections.up.json` and `000001_create_content_collections.down.json`).

It executes in:  
-  `internal/repository/db/migration.go` inside both services.

## Transactions
To enforce strict database atomicity, the services use explicit MongoDB Client Sessions to process relational mutations safely.  
It is used in:  
- **user-service/internal/repository/dao/user-dao.go**  
  Utilizes sessions to handle atomic changes to a user profile cleanly.  
- **content-service/internal/repository/dao/** (post-dao.go, comment-dao.go, like-dao.go)   
  Manages structural adjustments safely. For example, creating a like entry while simultaneously updating a counter state or associated cache entry.

## NATS event streaming

- **auth-service/internal/event/**  
  Publisher implementation pushes event payload to the "user.registered" subject channel as soon as a registration successfully completes.  
- **content-service/internal/events/**   
  `post-publisher.go` and `like-queue.go` contain wrappers wrapping a `*nats.Conn` pointer object instance. When actions complete, they serialize details into structured JSON byte streams and hit `nc.Publish(...)` targeting subjects such as `"post.created"` or `"like.created"`.
- **notification-service/internal/subscriber/**   
  Acts as a downstream event listener. It monitors the unified NATS broker stream (subscriber.go), parsing message streams asynchronously to perform notification pipelines (like sending automated emails via its mailer package) entirely detached from the primary API request loop.

## Clean Archiecture. Microservice's workflow: 
   ``` 
   cmd (entry point) ➔ internal/app (wiring dependencies) ➔ internal/transport (handlers/gRPC mappers) ➔ internal/usecase (business logic) ➔ internal/repository (database access objects)
   ```
