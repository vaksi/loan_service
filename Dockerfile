##
# Multi‑stage build for the loan service. The builder stage compiles
# the Go binary with minimal dependencies, and the final stage
# produces a small container containing only the binary. This design
# improves security by eliminating unused tooling and reduces image
# size.
FROM golang:1.20-alpine AS builder

# install git (required to fetch dependencies) and build tools
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
RUN go mod download
COPY . .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loan_service ./cmd

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /app/loan_service /loan_service

EXPOSE 8080
ENTRYPOINT ["/loan_service"]