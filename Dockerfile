FROM --platform=$BUILDPLATFORM golang:1.24.5-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o megafon-bot ./cmd/bot

FROM gcr.io/distroless/static-debian11

COPY --from=build /src/megafon-bot /bin/megafon-bot

COPY --from=build /src/config       /config
COPY --from=build /src/migrations   /migrations

ENV CONFIG_PATH=/config/config.dev.yaml

ENTRYPOINT ["/bin/megafon-bot"]