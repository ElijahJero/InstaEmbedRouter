## Configurable InstaEmbedRouter fork

This fork keeps the upstream app structure intact while adding runtime configuration so you can run it on your own domain instead of `zzinstagram.com`.

## Features supported

| Host | Post Description | Video handling | Image Index | Notes |
| --- | :-: | :-: | :-: | --- |
| `<your-domain>` |   | ✅ | ✅ | Base proxy host |
| `<gallery-subdomain>.<your-domain>` |   | ✅ | ✅ | Embed the post without a description |
| `<direct-subdomain>.<your-domain>` |   | ✅ |   | Embed the post without the integration frame |
| `<normal-subdomain>.<your-domain>` | ✅ | ✅ | ✅ | The normal way to embed posts (description + username) |

Default values remain compatible with upstream behavior:

- `PROXY_BASE_DOMAIN=zzinstagram.com`
- `PROXY_GALLERY_SUBDOMAIN=g`
- `PROXY_DIRECT_SUBDOMAIN=d`
- `PROXY_NORMAL_SUBDOMAIN=n`

## Build

Make sure you have [Go](https://go.dev/) installed, then run:

```bash
go build .
```

## Run locally

```bash
./InstagramEmbedResolver \
  -p 8080 \
  -proxy-domain embeds.example.com \
  -gallery-subdomain media \
  -direct-subdomain raw \
  -normal-subdomain post
```

You can also use environment variables:

```bash
PROXY_PORT=8080
PROXY_BASE_DOMAIN=embeds.example.com
PROXY_GALLERY_SUBDOMAIN=media
PROXY_DIRECT_SUBDOMAIN=raw
PROXY_NORMAL_SUBDOMAIN=post
RESOLVERS_FILE=resolvers.json
```

## Docker

Build and run with Docker:

```bash
docker build -t insta-embed-router .
docker run --rm -p 8080:8080 \
  -e PROXY_BASE_DOMAIN=embeds.example.com \
  -e PROXY_GALLERY_SUBDOMAIN=media \
  -e PROXY_DIRECT_SUBDOMAIN=raw \
  -e PROXY_NORMAL_SUBDOMAIN=post \
  insta-embed-router
```

Or use Compose:

```bash
docker compose up --build
```

## Soft fork notes

This fork only adds:

- runtime domain/subdomain configuration
- Docker launch files
- fork-specific documentation

That keeps the code changes small and makes future upstream syncs easier.

## Credits

As this app is only acting as a proxy, it relies on other Instagram embedding softwares such as [Instafix](https://github.com/Wikidepia/InstaFix/), [vxinstagram](https://github.com/Lainmode/InstagramEmbed-vxinstagram), [OGInstagram](https://github.com/seirenkr/OGInstagram) and [InstafixRevived](https://github.com/Bl0ck154/InstaFix-Revived).
