FROM golang:1.24-bullseye as debugger 

# installing packages 
RUN apt-get update && apt-get install -y \
    gcc \
    git \
    musl-dev

# build go module
WORKDIR /home/app/broker-service

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /brokerApp ./cmd/api

# run debugging with dlv
RUN go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 80
CMD ["dlv", "exec", "/brokerApp", "--headless", "--log", "--listen=:80", "--api-version=2"]


FROM alpine:3.21 as runner 
COPY --from=debugger /brokerApp .
CMD ["./brokerApp"]
