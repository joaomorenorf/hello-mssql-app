# Simple hello world app to test mssql connection

Required environment variables:
```shell
export DB_USER="SA"
export DB_PASSWORD="password"
export DB_SERVER="URLorIP"
export DB_NAME="master"
```

optional:
```shell
export PORT="8000"
```

Run using docker:
```shell
docker run -e DB_USER -e DB_PASSWORD -e DB_SERVER -e DB_NAME -e PORT joaomorenorf/hello-mssql-app:1.0.1
```

Run directly:
```shell
go run main.go
```