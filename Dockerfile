FROM alpine

RUN apk update && apk add ca-certificates && update-ca-certificates && apk --no-cache add tzdata
RUN mkdir -p /web
COPY ./web/* /web
COPY ./comics /comics

ENTRYPOINT ["/comics"]
