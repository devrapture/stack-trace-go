# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25.5

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS development
ARG AIR_VERSION=v1.67.3
RUN go install github.com/air-verse/air@${AIR_VERSION}

COPY . .

EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM base AS build
ARG TARGETOS
ARG TARGETARCH

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /bin/server ./cmd/api

FROM alpine:3.22 AS production

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app -H -s /sbin/nologin app

COPY --from=build --chown=app:app /bin/server /bin/server

USER app
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
