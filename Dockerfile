# Stage 1: build the binary
FROM docker.io/balenalib/raspberrypi4-64-golang:1.20.1 as builder
WORKDIR /app
COPY . .
RUN go mod download
# Use all available cores for the build process
RUN GOMAXPROCS=$(nproc) go build -o bot -p $(nproc)

# Stage 2: inimal runtime image
FROM alpine:latest
RUN addgroup -S botgroup && adduser -S -G botgroup botuser
WORKDIR /app
COPY --from=builder /app/bot .
USER botuser

# Run the bot
CMD ["./bot"]