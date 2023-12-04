FROM golang:1.19.2-alpine as build

WORKDIR /go/src
COPY . .

RUN apk add --update --no-cache build-base
RUN make build

FROM scratch as run

WORKDIR /go
COPY --from=build /go/src/config.json ./
COPY --from=build /go/src/be ./

ENTRYPOINT [ "./be" ]