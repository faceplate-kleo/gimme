# This container is meant to be used as a test environment
# It's not intended for use by end-users, just development
FROM golang:latest

RUN apt update

COPY . /gimme

WORKDIR /gimme/gimme-core
RUN go install gimme-core.go # prep the binary why not, go nuts

ENTRYPOINT ["go", "test", "-v", "./..."]