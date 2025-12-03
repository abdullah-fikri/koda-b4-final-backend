# README – Backend (Golang + Gin + PostgreSQL + Redis)

Backend ini adalah layanan REST API untuk aplikasi KODA URL Shortener yang menyediakan fitur:

- Pembuatan short link

- Redirect berdasarkan slug

- Statistik klik (database + Redis)

- Autentikasi user (JWT)

- CRUD user & management link

- Redis caching & flushing mechanism

## 🚀 How to Run Backend

Masuk ke directory backend:

- cd backend
- go mod tidy
- go run main.go

| Method | Endpoint                    | Description       |
| ------ | --------------------------- | ----------------- |
| POST   | `/api/v1/links`             | Create short link |
| GET    | `/:slug`                    | Redirect          |
| POST   | `/api/v1/auth/login`        | Login             |
| POST   | `/api/v1/auth/register`     | Register          |
| PUT   | `/api/v1/links/:slug`     | Edit          |
| DELETE   | `/api/v1/links/:slug`     | Delete          |


## 🔁 Redis Flushing Mechanism

Counter klik real-time

Cache short link (slug → URL)



### ERD 
```mermaid
erDiagram

    users {
        INT id PK
        VARCHAR email
        TEXT password
        TEXT role
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    profile {
        INT id PK
        INT users_id FK
        VARCHAR username
        VARCHAR phone
        VARCHAR address
        TEXT profile_picture
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    short_links {
        INT id PK
        INT user_id FK
        VARCHAR slug
        TEXT url
        INT clicks
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    daily_analytics {
        INT id PK
        INT user_id FK
        DATE date
        INT total_links
        INT total_visits
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    sessions {
        INT id PK
        INT user_id FK
        TEXT refresh_token_hash
        TIMESTAMP revoked_at
        TIMESTAMP expires_at
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }


    users ||--|{ profile : "has one"
    users ||--|{ short_links : "creates"
    users ||--|{ daily_analytics : "analytics"
    users ||--|{ sessions : "login sessions"
```