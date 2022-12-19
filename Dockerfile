FROM golang:1.19-bullseye AS build
WORKDIR /build
COPY . .
RUN go test -v ./internal/helpers && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags="-w -s" \
    -o /build/api-server \
    ./cmd/api-server

FROM scratch AS api-server
ARG API_PORT=1323
COPY --from=build /build/api-server /zetahedron
COPY .env .env
EXPOSE ${API_PORT}
ENTRYPOINT [ "/zetahedron" ]
