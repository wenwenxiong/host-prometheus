#FROM alpine:latest
FROM scratch

#RUN apk add --no-cache ca-certificates 

COPY  ./host-prometheus /usr/local/bin/
COPY  ./conf.yaml /etc/host-prometheus/conf.yaml

EXPOSE 9111
ENTRYPOINT ["/usr/local/bin/host-prometheus"]
