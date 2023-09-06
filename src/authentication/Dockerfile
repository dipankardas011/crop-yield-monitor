FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 go build -v -o demo cmd/addersvr/*

FROM scratch

WORKDIR /app

COPY --from=builder /app/demo /app/demo
COPY --from=builder /app/swaggerui /app/swaggerui
COPY --from=builder /app/gen /app/gen

ENTRYPOINT [ "/app/demo" ]

EXPOSE 8080
