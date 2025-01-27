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
RUN mkdir -p \
    src/localization src/jwtKeys \
    src/application/communication/emailService/templates/acceptInvitation \
    src/application/communication/emailService/templates/activateAccount \
    src/application/communication/emailService/templates/forgotPassword \
    src/application/communication/emailService/templates/remindToActivateAccount

COPY --from=build /go/bin/app /go/bin/app
COPY --from=build /go/src/app/src/localization/*.json src/localization
COPY --from=build /go/src/app/jwt-decoder.sh .
COPY --from=build /go/src/app/src/application/communication/emailService/templates/acceptInvitation/*.html src/application/communication/emailService/templates/acceptInvitation
COPY --from=build /go/src/app/src/application/communication/emailService/templates/activateAccount/*.html src/application/communication/emailService/templates/activateAccount
COPY --from=build /go/src/app/src/application/communication/emailService/templates/forgotPassword/*.html src/application/communication/emailService/templates/forgotPassword
COPY --from=build /go/src/app/src/application/communication/emailService/templates/remindToActivateAccount/*.html src/application/communication/emailService/templates/remindToActivateAccount

CMD ["/bin/sh", "-c", "./jwt-decoder.sh && /go/bin/app"]
