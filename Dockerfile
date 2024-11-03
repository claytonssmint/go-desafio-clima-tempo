FROM golang:1.23 as build
WORKDIR /app
COPY . .
COPY  .env .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-desafio-clima-tempo ./cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/go-desafio-clima-tempo .
COPY --from=build /app/.env .env
ENTRYPOINT ["/app/go-desafio-clima-tempo"]
