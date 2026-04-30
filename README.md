# Core Ticket

A Go-based ticketing system API service that manages support tickets with robust features for authentication, database operations, and external integrations.

## Features

- **Ticket Management**: Create, read, update, and manage support tickets
- **JWT Authentication**: Secure API endpoints with JWT token-based authentication
- **Database Support**: Multi-database support (Main DB and Ticketing DB)
- **Email Notifications**: Integrated SMTP support for email notifications
- **File Management**: S3 compatibility for file uploads and storage
- **Error Tracking**: Sentry integration for error monitoring
- **API Validation**: Request validation using Gin and validator/v10
- **Database Migrations**: Automated database schema management using golang-migrate

## Prerequisites

- **Go**: 1.22 or higher
- **MySQL**: 5.7 or higher (for both Main DB and Ticketing DB)
- **Git**: For version control
- **Environment Configuration**: `app.local.env` file for local development

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd core-ticket
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Set Up Environment Variables

Copy and configure the `app.local.env` file with your local development settings:

```bash
# Copy the template
cp app.local.env app.local.env
```

Edit `app.local.env` with your local database credentials:

```env
APP_ENV=development
APP_PORT=3005
APP_HOST=localhost

# Main Database Configuration
MAIN_DB_HOST=localhost
MAIN_DB_NAME=esb_main
MAIN_DB_USERNAME=root
MAIN_DB_PASSWORD=your_password

# Ticketing Database Configuration
TICKETING_DB_HOST=localhost
TICKETING_DB_NAME=esb_support
TICKETING_DB_USERNAME=root
TICKETING_DB_PASSWORD=your_password

# JWT Configuration
APP_JWT_SECRET=your_secret_key
APP_REFRESH_SECRET=your_refresh_secret
APP_JWT_TOKEN_LIFE_SPAN=480

# SMTP Configuration (for email notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_ENCRYPTION=tls
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Timezone
TZ=Asia/Jakarta
```

### 4. Set Up Databases

Ensure both `MAIN_DB` and `TICKETING_DB` databases exist in your MySQL server:

```sql
CREATE DATABASE IF NOT EXISTS esb_main;
CREATE DATABASE IF NOT EXISTS esb_support;
```

### 5. Run Database Migrations

The application automatically runs migrations on startup using the migration files in the `migrations/` directory:

```bash
# Migrations are applied automatically when the app starts
```

## Running the Application

### Development Mode

```bash
go run main.go
```

The server will start on **http://localhost:3005**

### Building for Production

```bash
go build -o core-ticket \
  -ldflags="-X main.Version=v1.0.0 -X main.CommitSHA=$(git rev-parse --short HEAD) -X main.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
  .
```

Run the compiled binary:

```bash
./core-ticket
```

## Project Structure

```
core-ticket/
├── main.go                 # Entry point
├── router.go              # API route definitions
├── validator.go           # Request validators
├── go.mod                 # Go module definition
├── app.local.env          # Local environment configuration
│
├── config/                # Configuration loading
├── constants/             # Application constants
├── database/              # Database setup and connections
├── migrations/            # Database migration files
│
├── models/                # Data models
│   ├── ticket/
│   ├── ticket_history/
│   └── ticket_product_type/
│
├── modules/               # Feature modules
│   ├── common/
│   ├── ip_whitelist/
│   └── return_ticket/
│
├── base/                  # Base utilities and helpers
│   ├── dto/               # Data Transfer Objects
│   ├── helpers/           # Utility helpers
│   ├── token/             # Token management
│   └── utility/
│
├── middlewares/           # HTTP middleware
├── helpers/               # Application-level helpers
└── assets/                # Static assets

```

## API Endpoints

The application provides RESTful API endpoints for ticket management. Key endpoints include:

- **Tickets**: CRUD operations for tickets
- **Ticket History**: Track ticket changes and history
- **Health Check**: `/healthcheck` endpoint for monitoring

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment | development |
| `APP_PORT` | Server port | 3005 |
| `APP_HOST` | Server host | localhost |
| `MAIN_DB_HOST` | Main database host | localhost |
| `MAIN_DB_NAME` | Main database name | esb_main |
| `MAIN_DB_USERNAME` | Main database user | root |
| `MAIN_DB_PASSWORD` | Main database password | - |
| `TICKETING_DB_HOST` | Ticketing database host | localhost |
| `TICKETING_DB_NAME` | Ticketing database name | esb_support |
| `TICKETING_DB_USERNAME` | Ticketing database user | root |
| `TICKETING_DB_PASSWORD` | Ticketing database password | - |
| `APP_JWT_SECRET` | JWT signing secret | - |
| `APP_JWT_TOKEN_LIFE_SPAN` | JWT token lifetime (minutes) | 480 |
| `SMTP_HOST` | SMTP server host | smtp.gmail.com |
| `SMTP_PORT` | SMTP server port | 587 |
| `TZ` | Timezone | Asia/Jakarta |

## Database Schema

Database migrations are managed by `golang-migrate`. Migration files are located in the `migrations/` directory:

- `230828001_init_database.up.sql` - Schema initialization
- `230828001_init_database.down.sql` - Schema rollback

## Development

### Wire (Dependency Injection)

This project uses Google Wire for dependency injection. To regenerate wire files after modifying providers:

```bash
cd modules/ip_whitelist && wire
cd ../return_ticket && wire
```

### Code Structure

- **Models**: Defined in `models/` directory with database schema
- **DTOs**: Data transfer objects in `base/dto/` for request/response validation
- **Repositories**: Data access layer for each module
- **Services**: Business logic layer for each module
- **Handlers**: HTTP request handlers (API endpoints)
- **Middlewares**: Authentication, logging, and request processing

## License

[Add your license information here]

## Support

For issues, questions, or contributions, please refer to the project repository.
