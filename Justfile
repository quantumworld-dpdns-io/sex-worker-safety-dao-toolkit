# Sex Worker Safety DAO Toolkit

# Rust
build-gateway:
    cd services/gateway && cargo build

build-attestation:
    cd services/attestation && cargo build

test-gateway:
    cd services/gateway && cargo test

test-attestation:
    cd services/attestation && cargo test

# Go
build-checkin:
    cd services/checkin-dao && go build -o bin/ ./...

test-checkin:
    cd services/checkin-dao && go test ./...

# Docker
docker-up:
    docker compose up

docker-build:
    docker compose build

docker-down:
    docker compose down

docker-clean:
    docker compose down -v

# All
test-all: test-gateway test-attestation test-checkin
