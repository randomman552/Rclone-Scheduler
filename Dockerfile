FROM golang:1.21.5-alpine3.19 AS build
WORKDIR /build
COPY . .
RUN go build -o entrypoint .

FROM alpine:3.19 AS final

RUN apk add curl

COPY ./templates /templates
COPY --from=build /build/entrypoint /

HEALTHCHECK CMD curl --fail http://localhost:3000/health || exit 1
ENTRYPOINT [ "/entrypoint" ]