# [Zetahedron](https://l.thisworddoesnotexist.com/Uq7L)
Simple URL shortener build with [Go](https://go.dev/), [Echo](https://echo.labstack.com/) and [Redis](https://redis.io/)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
# Build and run
## Using Docker Compose
```
# Clone the repository
git clone https://github.com/weirdvic/zetahedron.git

# Run the app
docker compose up --build --detach
```
By default the app will be listening on http://localhost:1323
# API description
## /shorten
Try to shorten an URL.

Request params:

`url` — an URL to shorten

`url_slug` — optional short URL slug

`expiry` — optional short link expiration time in hours
Example request:
```
curl -s -X POST -H 'Content-Type: application/json' \
-d '{"url":"https://thebestmotherfucking.website/"}' \
http://localhost:1323/shorten
```
Example output:
```
{
  "url": "https://thebestmotherfucking.website/",
  "short_url": "http://localhost:1323/mpdN0iJlZ0h",
  "expiry": 24
}
```
## /:url
Try to resolve short URL slug to the original URL

Request params: none, short URL provided in request path

Example request:
```
curl -Ls -w %{url_effective} -o /dev/null http://localhost:1323/mpdN0iJlZ0h
```
Example response:
```
https://thebestmotherfucking.website/
```
