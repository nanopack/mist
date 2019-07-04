FROM golang:1.12

# Prepare external Go packages:
# we're going to add go.mod and go.sum into /tmp and download all of our go dependencies.
ENV GOPATH=/go
ENV GO111MODULE=on
ADD go.mod go.sum /tmp/go/
WORKDIR /tmp/go
RUN go mod download
# reset WORKDIR and remove tmp files
WORKDIR /
RUN rm -rf /tmp/go

ADD . /go/src/github.com/nanopack/mist
WORKDIR /go/src/github.com/nanopack/mist
RUN CGO_ENABLED=1 go build -o /mist main.go

# We don't want to run as root in production. Create a "mist" user.
RUN addgroup mist && useradd -g mist mist

# Tell docker that all future commands should run as the appuser user
USER mist

ENTRYPOINT [ "/mist" ]