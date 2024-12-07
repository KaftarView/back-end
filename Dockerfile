FROM golang:1.23.4-alpine3.20 AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download \
    && CGO_ENABLED=0 go build -o /go/bin/app

# FROM gcr.io/distroless/static-debian12:nonroot AS deploy
FROM gcr.io/distroless/static-debian12:debug AS deploy

COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/src/localization /src/localization
CMD ["/app"]
