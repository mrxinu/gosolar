FROM golang:1.13-buster as builder

WORKDIR /go/solar
COPY . .

# Download our modules.
WORKDIR /go/solar/command
RUN go build -o /solarcmd .

#FROM alpine:latest
#COPY --from=builder /solarcmd /solarcmd
ENTRYPOINT ["/solarcmd"]
#RUN apt-get install -y bash
#ENTRYPOINT ["/bin/bash"]