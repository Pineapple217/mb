# mb

micro-blog

insiperd by https://github.com/l1mey112/me.l-m.dev

TODO: better description

# Features

TODO

# Dev Setup

TODO

## Dev dependencies

Latest version of Go and the following codegen tools.

```sh
go install github.com/a-h/templ/cmd/templ@latest
```

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Air is optional but strongly recommended.

```sh
go install github.com/cosmtrek/air@latest
```

# Deploying

## Docker

This is an expample to run mb with Docker. You can you can use it as a base for your own config.

### Env

See [Configuration](#configuration) for more info.

```env
# .env
MB_AUTH_PASSWORD=test123
MB_TIMEZONE=Europe/Brussels
MB_LOGO=" ▄▄▄·▪   ▐ ▄ ▄▄▄ .\n▐█ ▄███ •█▌▐█▀▄.▀·\n ██▀·▐█·▐█▐▐▌▐▀▀▪▄\n▐█▪·•▐█▌██▐█▌▐█▄▄▌\n.▀   ▀▀▀▀▀ █▪ ▀▀▀ "
MB_RIGHTS=pine32.be
MB_LINK=https://pine32.be
MB_MESSAGE="A funny little cycle."
```

### Compose file

```yml
# docker-compose.yml
services:
  mb:
    image: pineapple217/mb:latest
    container_name: mb
    volumes:
      - ./data:/app/data
      # this makes sure the datetimes are right on the container
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - .env
    ports:
      - 3000:3000
```

```console
docker compose up -d
```

# Configuration

## Environment variables

| variable           | info                                                                              | default                                                                 |
| ------------------ | --------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `MB_AUTH_PASSWORD` | The password used when login in `/auth`                                           | random generated password that will be printed on startup               |
| `MB_TIMEZONE`      | Timezone used to when showing dates of posts                                      | Timezone of machine the webserver is running on                         |
| `MB_LOGO`          | Logo displayed on the top of the main page. Make sure to use `\n` for line breaks | MB logo, see: https://patorjk.com/software/taag/#p=display&f=Elite&t=MB |
| `MB_RIGHTS`        | Name of the right holder thingy                                                   | `mb.dev` just a placeholder                                             |
| `MB_LINK`          | Link displayed at the top fo the page                                             | `https://mb.dev` link does not exist                                    |
| `MB_MESSAGE`       | Message displayed at the top of the page                                          | `Created without any JS.`                                               |

# Performance

profiling util

https://github.com/sevennt/echo-pprof

http://127.0.0.1:3000/debug/pprof/

go tool pprof -http=:8080 .\profile
