FROM golang:1.9.1-alpine

LABEL maintainer="Anthony Josiah Dodd <Dodd.AnthonyJosiah@gmail.com>"

EXPOSE 4004

# Fetch needed package deps.
WORKDIR /go/src/gitlab.com/project-leaf/mq-service-go
RUN apk add --no-cache git gcc musl-dev libsasl cyrus-sasl-dev && \
    go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/skelterjohn/rerun

# Copy over API source & needed files.
COPY ./Gopkg.lock Gopkg.lock
COPY ./Gopkg.toml Gopkg.toml
COPY ./src src
COPY ./main.go main.go

# Build the API.
RUN dep ensure && go install

# Use a CMD here (instead of ENTRYPOINT) for easy overwrite in docker ecosystem.
CMD ["mq-service-go"]
