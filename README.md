# About

Mb is as a micro blog server that aims to be small, fast and self contained. It was heavly inspired by [this project](https://github.com/l1mey112/me.l-m.dev) that was writen in V. I rewrote it in Go in a much cleaner fashion and added features. The main similarities that are left is the philosophy and the visuals (which are improved for mobile).

## Features

- Easy deploy with Docker or single binary
- Docker image of just 8 MB!
- Simple authentication
- Full markdown support
- Fully self contained with no external dependencies
- 100% Javascript free
- RSS feed
- Custom Youtube and Spotify embeds
- built in media manager for images, video and audio
- built in backups
- Mobile friendly
- Light weight and performant
- Tag and text search
- Private/draft posts
- ...

# Dev Setup

## Dev dependencies

Latest version of Go and the following codegen tools.
Use the latest version or the version currently used in the repo.

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

This is an example to run mb with Docker. You can you can use it as a base for your own config.

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
MB_HOST=https://mb.pine32.be
```

### Compose file

```yml
# docker-compose.yml
services:
  mb:
    image: pineapple217/mb:latest
    container_name: mb
    restart: unless-stopped
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
| `MB_HOST`          | The full host url where the blog is hosted. For ex. `https://mb.pine32.be`        | `http://localhost:3000`                                                 |

# Performance

TODO

profiling util

https://github.com/sevennt/echo-pprof

http://127.0.0.1:3000/debug/pprof/

go tool pprof -http=:8080 .\profile
