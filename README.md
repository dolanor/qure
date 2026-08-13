# The Qure

The Qure is a URL shortener and a QRCode generator written in Go.

## Endpoints

- `/admin/urls/`: manage the URLs
- `/<slug>`: short URL that redirects to the target URL
- `/qr/<slug>`: QRCode that links to the slug ([in UPPERCASE to reduce the density of the QR Code to make it more scannable by phones](https://shkspr.mobi/blog/2025/02/why-are-qr-codes-with-capital-letters-smaller-than-qr-codes-with-lower-case-letters/))

## Slug

The slug must be only alpha letters, numbers and `$ % * + - . :` to take advantage of the Uppercase QRCode optimization.
Also, reduce the slug length as much as possible to not jump into the greater QRCode density (maximum URL size &lt; 29 characters).

So, a short domain name + a short slug = low density QRCode, better scan.

## Admin panel

It is clearly not great, as it was a proof of concept for my https://rip.rest/ module, but it was useful enough and not too buggy enough that I didn't invest more time in it (there is PR in the work, but I've been way to busy to finish and merge it). And clearly, I'm no UI designer and lack the HTML/CSS skills to make it look good fast.

Anyway, the admin panel will look better once I improve RIP.

Also the pagination is not handled in the HTML admin panel, yet. So you need to add the `?page=x` yourself.

## Run

### With Go

```
go run github.com/dolanor/qure@latest
```

### Configuration

It uses environment variable to run

| Env var | format | role |
| --- | --- | --- |
| QURE_DOMAIN | link.example.com | The domain it's gonna run on. The shorter the better | 
| QURE_HOST | 0.0.0.0 | The IP address it's gonna run at. |
| QURE_PORT | 4444 | The port it's gonna run on. |
| QURE_DB_DIR | /var/lib/qure | The directory in which the database will be saved |

There is no TLS configuration because I tend to run my services behind a reverse proxy like Traefik or Caddy which auto update my certificates with Let's Encrypt.

### Docker

There is a ready to use `docker-compose.yaml` file.

- copy `.env.sample` and edit to your needs
- run `docker compose up`

You're done.

## Dependencies

- [RIP](https://rip.rest): a module that generates REST inspired endpoints based on pure Go types. It reduces the pain to create web handlers for mostly CRUD parts and handle content negotiation which, in that case gives us a free admin panel.
- [modernc.org/sqlite](https://modernc.org/sqlite): It works with sqlite with a pure Go implementation, which reduces performance compared to pure C wrappers with CGO, but allows for Go simple cross-compilation (just change GOOS and GOARCH to what you want/need).

## Design

Currently, the `URL ID` != `slug`. Maybe it's confusing and I should just make the slug the ID as it shouldn't be duplicable.
