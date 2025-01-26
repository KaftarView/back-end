FROM golang:1.23.4-alpine3.20 AS build

WORKDIR /go/src/app

COPY . .

RUN apk add --no-cache sed \
    && sed -i 's/\r$//' jwt-decoder.sh \
    && chmod +x jwt-decoder.sh \
    && go mod download \
    && CGO_ENABLED=0 go build -o /go/bin/app

# FROM gcr.io/distroless/static-debian12:nonroot AS deploy
FROM golang:1.23.4-alpine3.20 AS deploy

WORKDIR /go/src/app
RUN mkdir -p src/localization src/jwtKeys
COPY --from=build /go/bin/app /go/bin/app
COPY --from=build /go/src/app/src/localization/*.json src/localization
COPY --from=build /go/src/app/jwt-decoder.sh .

CMD ["/bin/sh", "-c", "./jwt-decoder.sh && /go/bin/app"]
