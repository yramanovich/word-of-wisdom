FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o server -trimpath cmd/server/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /app/server .
ENTRYPOINT ["/app/server"]