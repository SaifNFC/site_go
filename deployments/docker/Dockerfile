# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

# ---- Final stage ----
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/api .
COPY --from=builder /app/web/static ./web/static

USER nonroot:nonroot

EXPOSE 8081

ENTRYPOINT ["/app/api"]
