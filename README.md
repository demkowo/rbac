# RBAC Service

An RBAC (Role-Based Access Control) service that manages routes, roles, and their associations, ensuring that specific roles have access to specific routes.

## Features
- **Role Management**: Create, update, delete, and fetch roles.
- **Route Management**: Create, update, delete, and fetch routes, including batch activation of all discovered routes.
- **RBAC Associations**: Link roles to routes, enabling fine-grained access control.
- **Database Integration**: Uses PostgreSQL for persistent data storage.
- **APIs**: RESTful endpoints for roles, routes, and their RBAC bindings.

## Usage Example

### Add a Route

```bash
curl -X POST -H "Content-Type: application/json" -d '{
    "method": "GET",
    "path": "/api/v1/users",
    "active": true
}' http://localhost:5000/api/v1/routes
```

### Add a Role

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "admin"
}' http://localhost:5000/api/v1/roles
```

### Bind Role to Route

```bash
curl -X POST -H "Content-Type: application/json" -d '[
  {
    "route_id": "<ROUTE_UUID>",
    "role_id": "<ROLE_UUID>"
  }
]' http://localhost:5000/api/v1/rbac
```

### Retrieve All Routes By Role

```bash
curl -X GET http://localhost:5000/api/v1/routes/role/<ROLE_UUID>
```

## How It Works

### Initialization
- A PostgreSQL database connection is established using the environment - ariable `DB_CLIENT`.
- Necessary tables (routes, roles, rbac) are created if they don’t already - xist.

### Routes Registration
- The system automatically registers all available Gin routes, marks existing routes as inactive, and then updates each route’s status to active in the DB.
    
### RBAC Logic
- Roles and routes are stored in the database.
- The `rbac` table holds references linking roles to routes.
- When an RBAC record is added, the service checks that both the role and - the route exist before inserting a record into `rbac`.

## Workflow Overview

### Startup
- Application starts and connects to the DB.
- Creates tables if needed.
- Initializes repositories and services.
- Registers all endpoints with Gin.
- Discovers the endpoints and updates the DB routes table (active routes).

### API Operations
- Clients consume RESTful endpoints for roles, routes, and RBAC bindings.
- Queries, updates, and deletions are performed using structured - repositories.

### Data Consistency
- Foreign-key constraints maintain referential integrity across the rbac, - routes, and roles tables.

## Potential Improvements

- **Granular Logging:** Enhance logging for better troubleshooting (e.g., structured - logs, correlation IDs).
- **Caching:** Employ caching on frequently accessed data (e.g., route-to-role - mappings) to reduce database load.
- **Automated Tests:** Add more test coverage (unit/integration) with mocking/- stubbing for repositories.
- **Multi-DB Support:** Abstract database layer to allow switching to other - databases (MySQL, SQLite, etc.) if needed.