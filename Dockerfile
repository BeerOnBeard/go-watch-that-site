FROM golang:1.13-alpine

ENV GWTS_EMAIL_TO=""
ENV GWTS_EMAIL_FROM=""
ENV GWTS_EMAIL_PASSWORD=""

RUN apk update \
 && apk upgrade \
 && apk add --no-cache git

WORKDIR /go/src/github.com/beeronbeard/go-watch-that-site
COPY . .

RUN go get -v ./...
RUN go test -v -vet=off ./...
RUN go install -v ./...

CMD [ "go-watch-that-site" ]
