ARG HTTP_PROXY=
ARG HTTPS_PROXY=

#build stage
FROM harbor.cicd-p-az1.l12m.nl/docker-hub-proxy/library/golang:alpine AS builder

ARG HTTP_PROXY
ARG HTTPS_PROXY

#COPY ./certs/* /usr/local/share/ca-certificates/
#RUN cat /usr/local/share/ca-certificates/logius* >> /etc/ssl/certs/ca-certificates.crt

RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./main.go

#final stage
FROM harbor.cicd-p-az1.l12m.nl/docker-hub-proxy/library/alpine:latest

ARG HTTP_PROXY
ARG HTTPS_PROXY

ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV HTTP_PROXY=${HTTP_PROXY}

#COPY --from=builder /etc/ssl/certs /etc/ssl/certs
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/app"]
CMD ["serve", "--addr", "0.0.0.0:80"]
LABEL Name=iac-assets Version=0.0.1
EXPOSE 80
