FROM golang:1.21 as builder

WORKDIR /app

LABEL MAINTAINER="dipankar das"

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -v -o image .

FROM scratch

LABEL MAINTAINER="dipankar das"

WORKDIR /app

COPY --from=builder /app/image /app/image

ENTRYPOINT [ "/app/image" ]

EXPOSE 8090
