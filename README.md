# cPanel Matrix webhook handler

Simple Matrix webhook handler for cPanel notifications

## Compiling

`go build bartoostveen.nl/cpanel-matrix/cmd/cpanel-matrix`

## Setup

Pass the path to the config file (if any), using the `-c` flag, defaults to `config.yaml`.

## Configuration overview

```yaml
port: 8080
matrix:
  homeserver_url: "https://matrix.example.com/"
  access_token: verysecrettoken
  mx_id: ""@bot:example.com"
  rooms:
    - id: "!roomidhere"
      token: "supersecrettoken"
```

### Environment variables

```shell
CPANEL_PORT=8080
CPANEL_HOMESERVER_URL=https://matrix.example.com/
CPANEL_ACCESS_TOKEN=verysecrettoken
CPANEL_MX_ID=@bot:example.com
```

## License

This software comes as-is without any warranty, and is licensed by the GPLv3 license. Read more [here](./LICENSE).
