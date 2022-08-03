FROM golang:1.16-alpine as builder
ARG version=beta11
RUN apk add make
RUN mkdir /go/src/app
ADD . /go/src/app
RUN cd /go/src/app \
    && make Version=${version} linux_build


FROM alpine:latest
RUN apk add --no-cache ca-certificates  bash
RUN mkdir -p /app/logs
RUN mkdir -p /app/profile
COPY --from=builder /go/src/app/dist/peep  /app/

COPY entrypoint.sh /bin
RUN chmod +x /bin/entrypoint.sh
RUN chmod +x /app/peep
WORKDIR /app
EXPOSE 8000
CMD ["entrypoint.sh","-"]