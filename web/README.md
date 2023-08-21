# c-share/web <img src="/web/public/assets/footer.png" align="right" width="150" />

![last-commit](https://img.shields.io/github/last-commit/cyan903/c-share)

The frontend for `c-share`. Make sure `config.yaml` matches up with `.env`. This frontend is self-contained and can be moved if necessary. This frontend is not required and the API can be accessed programmatically if needed.

### Install

```sh
$ npm i
$ cp .env.example .env
$ nano .env # or any editor...
```

### Production/Development

```sh
$ npm run build # Production
$ npm run dev # Development
```

### Lint/Format

```sh
$ npm run lint
$ npm run format
$ npm run type-check
```

## License

[MIT](LICENSE)
