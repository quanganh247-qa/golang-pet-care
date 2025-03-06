<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
## Project Documentation

### Overview

This project is a React Native mobile application backed by a robust server-side architecture built with Go, PostgreSQL, and Redis, deployed using Docker. The architecture ensures efficient user management, secure authentication, and high-performance data processing.

The system integrates third-party services like JWT for token management, Gmail for email notifications, and a custom REST API for seamless client-server communication.

### Component

![image](https://github.com/user-attachments/assets/1471ba84-6f95-4c29-b92c-b16bdeba51dc)

<<<<<<< HEAD
<<<<<<< HEAD
### Architecture Diagram

![image](https://github.com/user-attachments/assets/ec053782-9055-42d0-9a1d-59407e2f5ff3)

### Features

- Frontend (React Native)
- Backend (Go , Gin web framework)
- Database (PostgreSQL)
- Caching Layer, manage queue (Redis)
- Third-Party Services
  - JWT: Secure token-based authentication for users.
  - Goong Maps: The perfect alternative to Google Maps API in Vietnam
  - Gmail: Automates email notifications for critical user actions (e.g., email confirmation).
  - Notifee: Hanlde notifications function
  - VietQR : Handle payment intergration

### Setup Instructions

#### Prerequisites

- Docker: Ensure Docker is installed on your machine.
- Node.js: Required for React Native development.
- Go: Backend server development.
- PostgreSQL: Database setup.
- Redis : Caching setup

### Installation

1. Clone the repository:
2. Start Docker containers:
   ` cd go-lang-petcare`
   `make postgres`
   `make redis`
3. Configure environment variables:
   - Backend: .env file for the Go server.
   - Frontend: .env file for React Native app.
   - Add API keys for Gmail, VietQr, and Goong maps...
4. Install dependencies by running:
   `go mod tidy`
   `go mod verify`
5. Initialize the database:
   `make mup`
6. Run the backend server:
   `go run main.go`

go install github.com/codegangsta/gin@latest
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.zshrc

curl -X GET "localhost:9200/\_cat/indices?v"
curl -X DELETE "localhost:9200/petclinic_medicines"
curl -X DELETE "localhost:9200/petclinic_diseases"
<<<<<<< HEAD
<<<<<<< HEAD
=======
keploy record -c "gin -p 8089 -i run main.go"

keploy test -c "gin -p 8089 -i run main.go" --delay 20
>>>>>>> 3980627 (generated test cases with keploy)
=======
## Project Documentation

### Overview

=======
## Project Documentation

### Overview
>>>>>>> 2968ee5 (Update readme.md)
=======
## Project Documentation

### Overview
>>>>>>> 2968ee5 (Update readme.md)
This project is a React Native mobile application backed by a robust server-side architecture built with Go, PostgreSQL, and Redis, deployed using Docker. The architecture ensures efficient user management, secure authentication, and high-performance data processing.

The system integrates third-party services like JWT for token management, Gmail for email notifications, and a custom REST API for seamless client-server communication.

### Component

![image](https://github.com/user-attachments/assets/1471ba84-6f95-4c29-b92c-b16bdeba51dc)

<<<<<<< HEAD
<<<<<<< HEAD
### Architecture Diagram

![image](https://github.com/user-attachments/assets/ec053782-9055-42d0-9a1d-59407e2f5ff3)
<<<<<<< HEAD
>>>>>>> 2968ee5 (Update readme.md)
=======

### Features

- Frontend (React Native)
- Backend (Go , Gin web framework)
- Database (PostgreSQL)
- Caching Layer, manage queue (Redis)
- Third-Party Services
  - JWT: Secure token-based authentication for users.
  - Goong Maps: The perfect alternative to Google Maps API in Vietnam
  - Gmail: Automates email notifications for critical user actions (e.g., email confirmation).
  - Notifee: Hanlde notifications function
  - VietQR : Handle payment intergration

### Setup Instructions

#### Prerequisites

- Docker: Ensure Docker is installed on your machine.
- Node.js: Required for React Native development.
- Go: Backend server development.
- PostgreSQL: Database setup.
- Redis : Caching setup

### Installation
<<<<<<< HEAD
  1. Clone the repository:
  2. Start Docker containers:
     ``` cd go-lang-petcare```
     ``` make postgres ```
     ``` make redis ```
  4. Configure environment variables:
     - Backend: .env file for the Go server.
     - Frontend: .env file for React Native app.
     - Add API keys for Gmail, VietQr, and Goong maps...
  5. Install dependencies by running:
      ``` go mod tidy ```
      ``` go mod verify ```
  7. Initialize the database:
       ``` make mup ```
  8. Run the backend server:
       ``` go run main.go ```




<<<<<<< HEAD
>>>>>>> 2c765c9 (Update readme.md)
=======
=======
>>>>>>> e859654 (Elastic search)

1. Clone the repository:
2. Start Docker containers:
   ` cd go-lang-petcare`
   `make postgres`
   `make redis`
3. Configure environment variables:
   - Backend: .env file for the Go server.
   - Frontend: .env file for React Native app.
   - Add API keys for Gmail, VietQr, and Goong maps...
4. Install dependencies by running:
   `go mod tidy`
   `go mod verify`
5. Initialize the database:
   `make mup`
6. Run the backend server:
   `go run main.go`

go install github.com/codegangsta/gin@latest
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.zshrc
<<<<<<< HEAD
>>>>>>> 3bf345d (happy new year)
=======

curl -X GET "localhost:9200/\_cat/indices?v"
curl -X DELETE "localhost:9200/petclinic_medicines"
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> ada3717 (Docker file)
=======
keploy record -c "gin -p 8089 -i run main.go"

keploy test -c "gin -p 8089 -i run main.go" --delay 20
>>>>>>> 3980627 (generated test cases with keploy)
=======

=======
>>>>>>> e859654 (Elastic search)
### Architecture Diagram

![image](https://github.com/user-attachments/assets/ec053782-9055-42d0-9a1d-59407e2f5ff3)
<<<<<<< HEAD
>>>>>>> 2968ee5 (Update readme.md)
=======

### Features

- Frontend (React Native)
- Backend (Go , Gin web framework)
- Database (PostgreSQL)
- Caching Layer, manage queue (Redis)
- Third-Party Services
  - JWT: Secure token-based authentication for users.
  - Goong Maps: The perfect alternative to Google Maps API in Vietnam
  - Gmail: Automates email notifications for critical user actions (e.g., email confirmation).
  - Notifee: Hanlde notifications function
  - VietQR : Handle payment intergration

### Setup Instructions

#### Prerequisites

- Docker: Ensure Docker is installed on your machine.
- Node.js: Required for React Native development.
- Go: Backend server development.
- PostgreSQL: Database setup.
- Redis : Caching setup

### Installation
<<<<<<< HEAD
  1. Clone the repository:
  2. Start Docker containers:
     ``` cd go-lang-petcare```
     ``` make postgres ```
     ``` make redis ```
  4. Configure environment variables:
     - Backend: .env file for the Go server.
     - Frontend: .env file for React Native app.
     - Add API keys for Gmail, VietQr, and Goong maps...
  5. Install dependencies by running:
      ``` go mod tidy ```
      ``` go mod verify ```
  7. Initialize the database:
       ``` make mup ```
  8. Run the backend server:
       ``` go run main.go ```




<<<<<<< HEAD
>>>>>>> 2c765c9 (Update readme.md)
=======
=======
>>>>>>> e859654 (Elastic search)

1. Clone the repository:
2. Start Docker containers:
   ` cd go-lang-petcare`
   `make postgres`
   `make redis`
3. Configure environment variables:
   - Backend: .env file for the Go server.
   - Frontend: .env file for React Native app.
   - Add API keys for Gmail, VietQr, and Goong maps...
4. Install dependencies by running:
   `go mod tidy`
   `go mod verify`
5. Initialize the database:
   `make mup`
6. Run the backend server:
   `go run main.go`

go install github.com/codegangsta/gin@latest
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.zshrc
<<<<<<< HEAD
>>>>>>> 3bf345d (happy new year)
=======

curl -X GET "localhost:9200/\_cat/indices?v"
curl -X DELETE "localhost:9200/petclinic_medicines"
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> ada3717 (Docker file)
=======
keploy record -c "gin -p 8089 -i run main.go"

keploy test -c "gin -p 8089 -i run main.go" --delay 20
>>>>>>> 3980627 (generated test cases with keploy)
=======

=======
>>>>>>> e859654 (Elastic search)
### Architecture Diagram

![image](https://github.com/user-attachments/assets/ec053782-9055-42d0-9a1d-59407e2f5ff3)
<<<<<<< HEAD
>>>>>>> 2968ee5 (Update readme.md)
=======

### Features

- Frontend (React Native)
- Backend (Go , Gin web framework)
- Database (PostgreSQL)
- Caching Layer, manage queue (Redis)
- Third-Party Services
  - JWT: Secure token-based authentication for users.
  - Goong Maps: The perfect alternative to Google Maps API in Vietnam
  - Gmail: Automates email notifications for critical user actions (e.g., email confirmation).
  - Notifee: Hanlde notifications function
  - VietQR : Handle payment intergration

### Setup Instructions

#### Prerequisites

- Docker: Ensure Docker is installed on your machine.
- Node.js: Required for React Native development.
- Go: Backend server development.
- PostgreSQL: Database setup.
- Redis : Caching setup

### Installation
<<<<<<< HEAD
  1. Clone the repository:
  2. Start Docker containers:
     ``` cd go-lang-petcare```
     ``` make postgres ```
     ``` make redis ```
  4. Configure environment variables:
     - Backend: .env file for the Go server.
     - Frontend: .env file for React Native app.
     - Add API keys for Gmail, VietQr, and Goong maps...
  5. Install dependencies by running:
      ``` go mod tidy ```
      ``` go mod verify ```
  7. Initialize the database:
       ``` make mup ```
  8. Run the backend server:
       ``` go run main.go ```




<<<<<<< HEAD
>>>>>>> 2c765c9 (Update readme.md)
=======
=======
>>>>>>> e859654 (Elastic search)

1. Clone the repository:
2. Start Docker containers:
   ` cd go-lang-petcare`
   `make postgres`
   `make redis`
3. Configure environment variables:
   - Backend: .env file for the Go server.
   - Frontend: .env file for React Native app.
   - Add API keys for Gmail, VietQr, and Goong maps...
4. Install dependencies by running:
   `go mod tidy`
   `go mod verify`
5. Initialize the database:
   `make mup`
6. Run the backend server:
   `go run main.go`

go install github.com/codegangsta/gin@latest
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.zshrc
<<<<<<< HEAD
>>>>>>> 3bf345d (happy new year)
=======

curl -X GET "localhost:9200/\_cat/indices?v"
curl -X DELETE "localhost:9200/petclinic_medicines"
>>>>>>> e859654 (Elastic search)
