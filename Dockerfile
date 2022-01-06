FROM golang:1.17-alpine3.15 AS builder
RUN apk add --no-cache make
COPY / /app
WORKDIR /app
RUN make

FROM scratch
COPY --from=builder /app/bin/update-java-ca-certificates /update-java-ca-certificates
ENTRYPOINT ["/update-java-ca-certificates"]