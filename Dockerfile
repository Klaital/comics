FROM alpine

RUN apk update && apk add ca-certificates && update-ca-certificates && apk --no-cache add tzdata
COPY ./comics /comics

ENTRYPOINT ["/comics"]
