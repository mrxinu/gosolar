FROM golang:1.13-buster as builder

WORKDIR /go/solar
COPY . .

# Download our modules.
WORKDIR /go/solar/command
RUN go build -o /solarcmd .

FROM debian:buster
COPY --from=builder /solarcmd /solarcmd
ENTRYPOINT ["/solarcmd"]