# Agnos Test

Hospital middleware backend for staff authentication and hospital-scoped patient search.<br>
The service is built with Go, Gin, GORM, PostgreSQL, Docker Compose, and Nginx.

## Features

- Create hospital staff accounts with bcrypt-hashed passwords.
- Authenticate staff members with hospital-scoped credentials.
- Issue signed JWT access tokens containing staff and hospital IDs.
- Restrict patient searches to the authenticated staff member's hospital.
- Search local patient records using optional identity and contact fields.
- Query Hospital A when a matching local patient is not found.

## Architecture

```text
Client
  |
  | HTTP :8080
  v
Nginx
  |
  v
Go/Gin application :8080
  |                 |
  |                 +--> Hospital A API
  |
  +--> PostgreSQL :5432
```

The application runs database migrations during startup. The database contains `hospitals`, `staffs`, and `patients` tables. Staff and patients are linked to a hospital through `hospital_id`; patient queries always include the authenticated hospital ID.

On startup, the application idempotently seeds one record for each table:

| Table | Seed data |
|---|---|
| `hospitals` | `Hospital A` |
| `staffs` | Username `demo`, password `demo123`, assigned to `Hospital A` |
| `patients` | Somchai Jaidee, national ID `1234567890123`, patient HN `HN001`, assigned to `Hospital A` |

The seed uses lookup keys and does not create duplicate records on restart. Change or remove the demo credentials before using the service outside local development.

## Project structure

```text
Agnos-test/
├── cmd/server/main.go          # Application entrypoint and route registration
├── internal
    ├── auth/                   # JWT claims, token handling, and middleware
    ├── config/                 # Environment-based configuration
    ├── db/                     # PostgreSQL connection and migrations
    ├── hospital/               # Hospital A HTTP client and response mapping
    ├── models/                 # Hospital, Staff, and Patient models
    ├── patient/                # Patient search service and handler
    └── staff/                  # Staff creation, login, and handlers
├── nginx/nginx.conf            # Reverse proxy configuration
├── docker-compose.yml          # PostgreSQL, app, and Nginx services
├── Dockerfile                  # Multi-stage Go container build
└── README.md
```

## Configuration

| Variable | Description | Example |
|---|---|---|
| `APP_NAME` | Application name used in startup logs | `agnos-backend` |
| `APP_PORT` | Port used by the Go application inside the Docker network | `8080` |
| `DB_HOST` | PostgreSQL hostname; Compose overrides this to `postgres` | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL username | `agnos` |
| `DB_PASSWORD` | PostgreSQL password | Set a local development value |
| `DB_NAME` | PostgreSQL database name | `agnos` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` for local Compose |
| `JWT_SECRET` | Secret used to sign and validate JWTs | Long random secret |
| `HOSPITAL_A_URL` | Hospital A API base URL | `https://hospital-a.api.co.th` |

## Start the application

Build and start all services:

```powershell
docker compose up -d --build
```

Check service status:

```powershell
docker compose ps
```

View application logs:

```powershell
docker compose logs -f app
```

The public API is available through Nginx at `http://localhost:8080`. PostgreSQL is available on `localhost:5432` for local administration.

To remove the PostgreSQL data volume and reset local data:

```powershell
docker compose down -v
```

## API

All examples below use the Nginx URL.

### `POST /staff/create`

Creates a staff account. If the named hospital does not exist, it is created.

Request:

```http
POST /staff/create
Content-Type: application/json
```

```json
{
  "username": "alice",
  "password": "secret123",
  "hospital": "Hospital A"
}
```

Successful response: `201 Created`

```json
{
  "id": 1,
  "username": "alice",
  "hospital_id": 1
}
```

Possible errors:

- `400 Bad Request`: malformed JSON or missing required fields.
- `409 Conflict`: username already exists in the specified hospital.
- `500 Internal Server Error`: database or server failure.


### `POST /staff/login`

Authenticates staff credentials and returns an eight-hour JWT.

Request:

```http
POST /staff/login
Content-Type: application/json
```

```json
{
  "username": "alice",
  "password": "secret123",
  "hospital": "Hospital A"
}
```

Successful response: `200 OK`

```json
{
  "token": "<jwt>"
}
```

Possible errors:

- `400 Bad Request`: malformed JSON or missing required fields.
- `401 Unauthorized`: username, password, or hospital is incorrect.
- `500 Internal Server Error`: server failure.


### `GET /patient/search`

Searches patients belonging to the hospital encoded in the JWT. At least one query parameter is recommended.

Supported optional query parameters:

| Parameter | Description |
|---|---|
| `national_id` | National identification number |
| `passport_id` | Passport identification number |
| `first_name` | Thai or English first name, partial match |
| `middle_name` | Thai or English middle name, partial match |
| `last_name` | Thai or English last name, partial match |
| `date_of_birth` | Date in `YYYY-MM-DD` format |
| `phone_number` | Phone number |
| `email` | Email address |

Request:

```http
GET /patient/search?national_id=123
Authorization: Bearer <jwt>
```

Successful response: `200 OK`

```json
{
  "patients": [
    {
      "id": 1,
      "hospital_id": 1,
      "first_name_en": "Alice",
      "last_name_en": "Example",
      "national_id": "123",
      "gender": "F"
    }
  ]
}
```

Possible errors:

- `401 Unauthorized`: missing, malformed, or invalid JWT.
- `400 Bad Request`: invalid search parameters or invalid date format.
- `500 Internal Server Error`: database or external integration failure.

If no local patient matches and the hospital is `Hospital A`, the service uses `national_id` or `passport_id` to call:

```text
GET {HOSPITAL_A_URL}/patient/search/{id}
```

The returned patient is mapped into the local schema and stored under the authenticated hospital.

## Testing

Run all Go tests:

```powershell
go test ./...
```

The test suite covers JWT handling, staff creation, password hashing, Hospital A response parsing, patient search hospital scoping, and validation errors.

Validate the Compose file without starting containers:

```powershell
docker compose config
```
