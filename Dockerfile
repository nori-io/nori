FROM golang:1.10.1 AS builder
RUN go version

RUN go get -v github.com/secure2work/nori
WORKDIR /go/src/github.com/secure2work/nori/

#RUN set -x && \
#    go get github.com/golang/dep/cmd/dep && \
#    dep ensure -v

RUN go build -o nori .

# Stage 2

FROM ubuntu:18.04
# RUN apk --no-cache add ca-certificates

RUN mkdir nori
RUN mkdir /nori/config
RUN mkdir /nori/plugins
RUN mkdir /nori/certs
RUN chmod 644 /nori

WORKDIR /nori/
COPY --from=builder /go/src/github.com/secure2work/nori .
RUN touch /nori/config/nori.json
RUN echo '{\n\
    "nori": {\n\
        "grpc": {\n\
          "enable": true,\n\
          "tls": {\n\
            "ca": "/nori/certs/ca.pem",\n\
            "private": "/nori/certs/private.pem"\n\
          }\n\
        },\n\
      "storage": {\n\
        "type": "none"\n\
      }\n\
    },\n\
    "plugins": {\n\
      "dir": "/nori/plugins"\n\
    }\n\
}' >> /nori/config/nori.json

EXPOSE 8080
EXPOSE 29876
ENTRYPOINT ["/nori/nori", "--config=/nori/config/nori.json", "server"]