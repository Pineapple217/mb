ARG GO_VERSION=1.23
FROM golang:${GO_VERSION} AS build
WORKDIR /src

ENV CGO_ENABLED=0

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -ldflags='-s -w -extldflags "-static"' -trimpath -o /bin/server ./cmd/server
    # static linking is necessary because of CGO dependency
    # -s -w removes debug info for smaller bin

FROM scratch AS final

ARG GIT_COMMIT=unspecified
LABEL org.opencontainers.image.version=$GIT_COMMIT
LABEL org.opencontainers.image.source=https://github.com/Pineapple217/mb

COPY --from=alpine:3.19 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/server /server

EXPOSE 3000

CMD [ "/server", "-listen", "0.0.0.0" ]
