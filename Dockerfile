#build
FROM golang AS build

WORKDIR /go/src/github.com/zs5460/ping

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build .

CMD ["./ping"]

#production
FROM scratch AS prod

COPY --from=build /go/src/github.com/zs5460/ping/default.config.json ./config.json
COPY --from=build /go/src/github.com/zs5460/ping/ping .

CMD ["./ping"]