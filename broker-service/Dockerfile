FROM golang:1.24-alpine as builder 

# installing packages 
RUN apk update && apk add --no-cache git gcc musl-dev

# build go module
ARG APP_SRC=/home/app/broker-service
WORKDIR ${APP_SRC}

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /brokerApp ./cmd/api

# run stage
FROM alpine:3.21 as runner 
COPY --from=builder /brokerApp .
COPY --from=builder ${APP_SRC} ${APP_SRC}
CMD ["/brokerApp"]
