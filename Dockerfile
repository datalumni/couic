FROM golang:bookworm AS builder

RUN apt-get update -qq && apt-get install -y -qq \
    libgl1-mesa-dev \
    xorg-dev \
    libglfw3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./core/ ./export/ -v && \
    go build -ldflags="-s -w" -o /build/couic .

FROM scratch

COPY --from=builder /build/couic /couic

ENTRYPOINT ["/couic"]
