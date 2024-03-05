FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src

RUN go build ./cmd/unollm

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/unollm /app/unollm

WORKDIR /app

EXPOSE 11451
EXPOSE 19198

CMD ["./unollm"]
