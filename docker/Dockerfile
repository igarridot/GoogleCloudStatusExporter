ARG GOLANG_VERSION=1.25
FROM golang:$GOLANG_VERSION AS builder
ARG PROJECT_NAME
ARG GOARCH=amd64
ENV PROJECT_NAME=${PROJECT_NAME}
ENV GOARCH=${GOARCH}

RUN \
    apt update && \
    apt install -y make git && \
    rm -rf /var/lib/apt/lists/*

#ENV GO111MODULE="on"
WORKDIR /go/src
COPY ./src ./
RUN \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${GOARCH} \
    go build -a -installsuffix cgo -o bin/${PROJECT_NAME}

FROM scratch
ARG PROJECT_NAME
ENV PROJECT_NAME=${PROJECT_NAME}
COPY --from=builder /go/src/bin/${PROJECT_NAME} /docker-entrypoint
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "/docker-entrypoint" ]
EXPOSE 9118
