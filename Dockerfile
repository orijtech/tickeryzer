FROM golang as builder

# Add Code and install
ADD . /go/src/github.com/orijtech/tickeryzer/
WORKDIR /go/src/github.com/orijtech/tickeryzer/cmd/tickeryzer
RUN go install -v

FROM alpine:latest
COPY --from=builder /go/bin/tickeryzer /go/bin/tickeryzer
# Setup work path
WORKDIR /data
ENTRYPOINT ["/go/bin/tickeryzer"]
