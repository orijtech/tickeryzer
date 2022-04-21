FROM golang:latest

# Add Code and install
ADD . /go/src/github.com/orijtech/tickeryzer/
WORKDIR /go/src/github.com/orijtech/tickeryzer/cmd/tickeryzer
RUN go install -v

# Setup work path
WORKDIR /data
ENTRYPOINT ["/go/bin/tickeryzer"]
