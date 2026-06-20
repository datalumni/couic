FROM golang:bookworm AS builder

ARG TARGET=linux

RUN apt-get update -qq && apt-get install -y -qq \
    libgl1-mesa-dev \
    xorg-dev \
    libglfw3-dev \
    && rm -rf /var/lib/apt/lists/*

RUN if [ "$TARGET" = "windows" ]; then \
        apt-get update -qq && apt-get install -y -qq gcc-mingw-w64-x86-64 && rm -rf /var/lib/apt/lists/*; \
    fi

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN if [ "$TARGET" = "windows" ]; then \
        go install github.com/tc-hib/go-winres@latest && \
        go-winres simply --icon ui/couic-icon.png --out scripts/rsrc; \
    fi

RUN go test ./core/ ./export/ -v && \
    case "$TARGET" in \
      windows) \
        cp scripts/rsrc_windows_amd64.syso . && \
        CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
        go build -ldflags="-s -w" -o /build/couic.exe . ;; \
      *) \
        go build -ldflags="-s -w" -o /build/couic . ;; \
    esac

FROM scratch

COPY --from=builder /build/couic* /

ENTRYPOINT ["/couic"]
