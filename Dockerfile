FROM golang:1.16-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go mod download
WORKDIR /app/cmd/

RUN go build -o /client


FROM alpine:3.14
LABEL maintainer=binozoworks
LABEL org.opencontainers.image.source = "https://github.com/InstantMinecraft/Client"
LABEL org.opencontainers.image.description="A standalone container running a minecraft server which is controllable through an http endpoint"
# Default to UTF-8 file.encoding
ENV LANG C.UTF-8

# Update system
RUN apk update \
  && apk --update-cache upgrade \
  && apk add ca-certificates \
  && update-ca-certificates

# Install JRE
RUN apk add --no-cache \
  --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing \
  --repository http://dl-cdn.alpinelinux.org/alpine/edge/main \
    openjdk18-jre-headless curl

# Create Server directory
RUN mkdir server
WORKDIR server

# Install Client Software
COPY --from=builder /client /client

# Install Minecraft Server
ADD https://piston-data.mojang.com/v1/objects/f69c284232d7c7580bd89a5a4931c3581eae1378/server.jar server.jar

# Accept eula
RUN echo "eula=true" > eula.txt

ENTRYPOINT ["/client"]
