FROM alpine:latest AS build

RUN apk add go

WORKDIR /app
COPY . .

RUN go get
RUN go build -ldflags="-w -s"

FROM alpine:latest AS run

WORKDIR /app

COPY --from=build /app/eloiserver ./eloiserver
COPY --from=build /app/config.json ./config.json

EXPOSE 8080

CMD ["./eloiserver"]
