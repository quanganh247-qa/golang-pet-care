# Golang Pet Care Service

## Project Overview

This project is a comprehensive pet care service backend built with Go, featuring a robust server-side architecture that supports a React Native mobile application. The system provides a complete solution for pet clinics and pet owners, offering features like appointment management, medical records tracking, pet scheduling, and e-commerce capabilities.

## Architecture

![Architecture Diagram](https://github.com/user-attachments/assets/ec053782-9055-42d0-9a1d-59407e2f5ff3)

### Component Diagram

![Component Diagram](https://github.com/user-attachments/assets/1471ba84-6f95-4c29-b92c-b16bdeba51dc)

## Technology Stack

### Backend
- **Go (Golang)**: Main programming language
- **Gin Web Framework**: HTTP web framework
- **PostgreSQL**: Primary database
- **Redis**: Caching and queue management
- **Elasticsearch**: Search functionality
- **MinIO**: Object storage for files and images
- **Docker**: Containerization and deployment

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

## API Documentation

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

## Elasticsearch Commands

```
curl -X GET "localhost:9200/_cat/indices?v"  # List all indices
curl -X DELETE "localhost:9200/petclinic_medicines"  # Delete medicines index
curl -X DELETE "localhost:9200/petclinic_diseases"  # Delete diseases index
```
