ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} AS build
WORKDIR /src

ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -ldflags='-s -w -extldflags "-static"' -o /bin/server -tags sqlite_math_functions ./cmd/server
    # static linking is necessary because of CGO dependency
    # -s -w removes debug info for smaller bin

FROM alpine:latest AS final

ARG GIT_COMMIT=unspecified
LABEL org.opencontainers.image.version=$GIT_COMMIT
LABEL org.opencontainers.image.source=https://github.com/Pineapple217/mb

# Removed user because of file permisions
# ARG UID=10001
# RUN adduser \
#     --disabled-password \
#     --gecos "" \
#     --home "/nonexistent" \
#     --shell "/sbin/nologin" \
#     --no-create-home \
#     --uid "${UID}" \
#     appuser
# USER appuser

WORKDIR /app
COPY --from=build /bin/server /app/server

EXPOSE 3000

CMD [ "/app/server", "-listen", "0.0.0.0" ]
