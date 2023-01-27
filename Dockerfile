FROM alpine:latest as build

RUN apk add go

WORKDIR /app
COPY . .

RUN go get
RUN go build -ldflags="-w -s"

FROM alpine:latest as run

RUN apk add cairo

WORKDIR /app

COPY --from=build /app/eloiserver ./run
COPY --from=build /app/config.json ./config.json

EXPOSE 8080

CMD ["./run"]
