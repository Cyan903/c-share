# c-share <img src="/web/public/assets/landing.png" align="right" width="150" />

![go-mod](https://img.shields.io/github/go-mod/go-version/cyan903/c-share?style=flat) ![go-report](https://goreportcard.com/badge/github.com/cyan903/c-share?style=flat) ![node-version](https://img.shields.io/node/v-lts/vite?style=flat)

A secure password protected file hosting service with built in support for applications such as [ShareX](https://getsharex.com/). Includes developer API documentation and file management.

```sh
$ make build # build for production
$ make dev # run for development
$ make format # format & lint code
$ make update # update and validate dependencies
```

## Features

- Self contained API
- Hashed file passwords
- File sharing support
- Development API support
- Frontend web interface
- Email verification
- Password reset

## Install

Clone the source code, import the DB and install dependencies.

```sh
$ git clone https://github.com/Cyan903/c-share.git
$ cd c-share
$ cp config.example.yaml config.yaml
$ mysql -u MY_DB_USER -p MY_DB_NAME < ext/db.sql
$ make update
```

Run or build the app.

```sh
$ make build
$ make dev
```

## License

[MIT](LICENSE)
