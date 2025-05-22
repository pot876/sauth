
```sh
HTTP_LISTEN_ADDR=:8085 \
CP_PG_ENABLED=true \
CP_PG_CONNECTION_STRING="postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable" \
    go run ./cmd/auth_server

```
