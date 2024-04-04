FROM golang:1.22.1 as builder
WORKDIR /app
COPY main.go ./
COPY go.mod ./
COPY go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /hello-mssql-app

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /hello-mssql-app /hello-mssql-app

ENV PORT 8000
ARG DB_USER
ARG DB_PASSWORD
ARG DB_SERVER
ARG DB_NAME

USER nonroot:nonroot
EXPOSE 8000
CMD ["/hello-mssql-app"]
