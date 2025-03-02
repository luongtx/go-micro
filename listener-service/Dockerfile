ARG APP_SRC=/home/app/listener-service

FROM golang:1.24-alpine as builder 

# build go module
ARG APP_SRC
WORKDIR ${APP_SRC}

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o ./listenerApp .
# run stage
FROM alpine:3.21 as runner 
ARG APP_SRC
WORKDIR ${APP_SRC}
COPY --from=builder ${APP_SRC} . 
CMD ["./listenerApp"]
