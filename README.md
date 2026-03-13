# 📦 Microservice Ecommerce DDD & Hexagonal
## 💻 Teckstack
  - golang
  - redis
  - keycloak
## 🚀 Services
  - API Gateway Service
  - Auth Service
  - Product Service
## 🚀 Quick Start
 - docker compose start
  ```bash
     docker compose up -f docker-compose.{stage}.yaml up -d
  ```
 - go install package
  ```bash
    go mod dowload
  ```
 - start with dotenvx
  ```bash
    dotenvx run -f ./envs/{stage}/.env -- go run ./cmd/{service}/main.go
  ```
