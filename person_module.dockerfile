FROM golang:1.22 as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o person_module


#BUILD a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/person_module /app

CMD ["/app/person_module"]
