FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS builder
LABEL authors="gucooing"

RUN apk add --no-cache bash
WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download && go mod verify
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags "rel" netgo -o /app/BaPs ./cmd/BaPs/BaPs.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /usr/ba
COPY --from=builder /app/BaPs .
COPY --from=builder /app/data/ ./data/
RUN chmod +x BaPs
EXPOSE 5000/tcp
ENTRYPOINT ["./BaPs"]