# Golang Pet Care Service

## Project Overview

This project is a comprehensive pet care service backend built with Go, featuring a robust server-side architecture that supports a React Native mobile application. The system provides a complete solution for pet clinics and pet owners, offering features like appointment management, medical records tracking, pet scheduling, and e-commerce capabilities.

## Documentation

   https://deepwiki.com/quanganh247-qa/golang-pet-care

### Component Diagram

![Component Diagram](https://github.com/user-attachments/assets/1471ba84-6f95-4c29-b92c-b16bdeba51dc)

## Technology Stack

### Backend
- **Go (Golang)**: Main programming language
- **Gin Web Framework**: HTTP web framework
- **PostgreSQL**: Primary database
- **Redis**: Caching and queue management
- **Docker**: Containerization and deployment
- **Swagger**: API documentation

### Third-Party Services
- **JWT**: Secure token-based authentication
- **Goong Maps**: Location services (alternative to Google Maps for Vietnam)
- **Gmail**: Email notifications
- **VietQR**: Payment integration
- **Google Gemini**: AI-powered chatbot functionality
- **OpenFDA**: Drug and medical information

## Features

### User Management
- User registration and authentication
- Role-based access control
- Profile management
- Email verification

### Pet Management
- Pet profile creation and management
- Medical history tracking
- Vaccination records
- Pet scheduling

### Appointment System
- Appointment booking and management
- Doctor availability scheduling
- Queue management
- Appointment history

### Medical Services
- Medical records management
- Treatment plans
- Prescription management
- Disease tracking and management

### E-commerce
- Product catalog
- Shopping cart functionality
- Order management
- Payment processing

### Location Services
- Clinic location finder
- Directions and navigation
- Distance calculation

### AI-Powered Features
- Chatbot for pet health inquiries
- Drug information lookup
- Side effect analysis
- Treatment recommendations

## Setup Instructions

### Prerequisites

- Docker and Docker Compose
- Go 1.23 or higher
- PostgreSQL
- Redis
- Elasticsearch
- MinIO

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd golang-pet-care
   ```

2. Configure environment variables:
   - Create an `app.env` file with necessary configuration
   - Set API keys for Gmail, VietQR, Goong Maps, Google Gemini, and OpenFDA

3. Install dependencies:
   ```
   go mod tidy
   go mod verify
   ```

4. Initialize the database:
   ```
   make mup
   ```

5. Start the services using Docker Compose:
   ```
   make docker-up
   ```
   
   Or run individual services:
   ```
   make postgres
   make redis
   ```

6. Run the backend server:
   ```
   go run main.go
   ```

### Development Tools

- Hot reload with gin:
  ```
  make server
  ```

- Database migrations:
  ```
  make migrateup    # Apply migrations
  make migratedown  # Revert migrations
  make new_migration name=migration_name  # Create new migration
  ```

- Generate Swagger documentation:
  ```
  make swagger
  ```

## API Documentation

### Swagger Documentation
The API is documented using Swagger. Once the server is running, you can access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

To regenerate the Swagger documentation after making API changes:
```
make swagger
```

### Authentication Endpoints
- `POST /api/v1/user/create` - Register a new user
- `POST /api/v1/user/login` - User login
- `POST /api/v1/user/verify_email` - Verify email
- `GET /api/v1/user/refresh_token` - Refresh access token
- `POST /api/v1/user/logout` - Logout user

### Pet Management
- `POST /api/v1/pet/create` - Create a new pet
- `GET /api/v1/pet/:pet_id` - Get pet details
- `GET /api/v1/pet/list` - List all pets
- `GET /api/v1/pet/` - List pets by username
- `PUT /api/v1/pet/:pet_id` - Update pet information
- `DELETE /api/v1/pet/delete/:pet_id` - Delete a pet

### Appointment Management
- `POST /api/v1/appointment/` - Create an appointment
- `GET /api/v1/appointment/user` - Get user's appointments
- `GET /api/v1/appointment/:id` - Get appointment details
- `PUT /api/v1/appointment/:id` - Update appointment
- `GET /api/v1/doctor/:doctor_id/time-slot` - Get available time slots

### Medical Records
- `GET /api/v1/pets/:pet_id/medical-records` - Get pet's medical records
- `POST /api/v1/pets/:pet_id/medical-records` - Create medical record
- `GET /api/v1/pets/:pet_id/medical-histories` - List medical history

### E-commerce
- `GET /api/v1/products/` - Get all products
- `POST /api/v1/cart/cart` - Add to cart
- `GET /api/v1/cart/cart` - Get cart items
- `POST /api/v1/cart/order` - Create order
- `GET /api/v1/cart/order/:order_id` - Get order details

### Payment
- `GET /api/v1/payment/banks` - Get list of banks
- `POST /api/v1/payment/generate-qr` - Generate payment QR code

### Location Services
- `GET /api/v1/location/places/autocomplete` - Place autocomplete
- `GET /api/v1/location/places/detail` - Get place details
- `GET /api/v1/location/directions` - Get directions

### Chatbot
- `POST /api/v1/chatbot/chat` - Chat with AI assistant
- `GET /api/v1/chatbot/drug-info` - Get drug information

## Docker Deployment

The application can be deployed using Docker Compose:

```
make docker-build  # Build Docker images
make docker-up     # Start all services
make docker-down   # Stop all services
make docker-logs   # View logs
```

