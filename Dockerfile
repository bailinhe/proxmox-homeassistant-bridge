FROM golang:1.21 as build
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN mkdir -p /build/bin
COPY . /build
WORKDIR /build
RUN go build -o bin/main main.go

FROM alpine:3.18
COPY --from=build build/bin/main /usr/local/bin/proxmox-ha-bridge
ENTRYPOINT ["proxmox-ha-bridge"]
