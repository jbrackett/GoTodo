# GoTodo

A simple Todo application that allows users to create, read, update, and delete todos.

## Technologies
- Go
- Echo
- Gorm
- PostgreSQL
- Docker

## Architecture

The application consists of only one file which contains both http and repository functions. In a production application this would be split into multiple layers.

### Database
The application uses PostgreSQL for the production database. Gorm is used to manage the database schema.

## Running the application
```bash
go build cmd/app/main.go
./run.sh
```

