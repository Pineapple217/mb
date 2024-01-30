# mb

micro-blog

insiperd by https://github.com/l1mey112/me.l-m.dev

TODO: better description

# Features

TODO

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

# Dev Setup

TODO

# Deploying

TODO

# Performance

profiling util

https://github.com/sevennt/echo-pprof

http://127.0.0.1:3000/debug/pprof/

go tool pprof -http=:8080 .\profile
