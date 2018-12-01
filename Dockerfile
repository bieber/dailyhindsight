FROM golang:1.10-alpine AS builder
WORKDIR /go/src/github.com/bieber/dailyhindsight
RUN apk add git
COPY . /go/src/github.com/bieber/dailyhindsight
RUN go get -d github.com/bieber/dailyhindsight
RUN mkdir /app
RUN go build -o /app/server github.com/bieber/dailyhindsight

FROM alpine:latest
LABEL maintainer="docker@biebersprojects.com"
EXPOSE 80

WORKDIR /app
RUN apk add ca-certificates
COPY --from=builder /app/server /app/server
ENTRYPOINT ["/app/server"]
