FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather

FROM scratch
WORKDIR /app
COPY --from=build /app/weather .
ENTRYPOINT ["./weather"]
