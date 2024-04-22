# This container is meant to be used as a test environment
# It's not intended for use by end-users, just development
FROM golang:latest

RUN apt update
RUN apt install sudo

COPY . /gimme

WORKDIR /gimme/gimme-core
RUN go install gimme-core.go # prep the binary why not, go nuts

RUN useradd -m -s /bin/bash gimmetest
RUN passwd -d gimmetest
RUN adduser gimmetest sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN echo "PATH="${PATH}"" >> /etc/environment
RUN chmod -R o+rwx /go
RUN chmod -R o+rwx /gimme

USER gimmetest
RUN echo ${PATH}
WORKDIR /gimme
RUN ./INSTALL.sh

WORKDIR /gimme/gimme-core
ENTRYPOINT ["go", "test", "-v", "./..."]