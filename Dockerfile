FROM golang:1.16-alpine3.13 AS builder

RUN apk add --update --no-cache gcc make musl-dev
RUN mkdir -p /app

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

ARG version=dev
ENV VERSION=$version

COPY . /src

RUN make -B VERSION=$VERSION LDFLAGS='-s -w -extldflags "-static"' BIN=/app/biscuit-api

FROM alpine:3.13

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

# Add new user 'biscuit'
RUN adduser -D biscuit
USER biscuit

COPY --from=builder /app /home/biscuit/app

WORKDIR /home/biscuit/app

# Since running as a non-root user, port bindings < 1024 is not possible
# 8080 for HTTP; 8443 for HTTPS;
EXPOSE 8080
#EXPOSE 8443

CMD ["./biscuit-api"]
