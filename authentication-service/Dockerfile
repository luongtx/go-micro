FROM golang:1.24-alpine as builder 

# installing packages 
RUN apk update && apk add --no-cache git gcc musl-dev

# build go module
ARG APP_SRC=/home/app/authentication-service
WORKDIR ${APP_SRC}

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /authApp ./cmd/api
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN chmod +x dlv.sh

FROM alpine:3.21 as runner 
COPY --from=builder /authApp .
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder ${APP_SRC} ${APP_SRC}
CMD ["/authApp"]
