FROM golang:1.10
WORKDIR /go/src/github.com/eduboard/backend
ADD . /go/src/github.com/eduboard/backend
RUN CGO_ENABLED=0 go build ./cmd/server/main.go

FROM alpine:latest
COPY --from=0 /go/src/github.com/eduboard/backend/main /

EXPOSE 8080
CMD [ "./main" ]
