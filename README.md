# =============================================================================
# Starting the server
# =============================================================================
```
go run cmd/server/main.go
```
# =============================================================================
# Migrations
# =============================================================================

## To install golang migrate
```
go get -u -d github.com/golang-migrate/migrate/cmd/migrate
```

## To create a migration
```
migrate create -ext sql -dir migrations -seq init_table
```

## To run a migration
```
migrate -source migrations -database migrate -path migrations -database postgres://postgres:root@localhost:5432/microauth?sslmode=disable up
```

# =============================================================================
