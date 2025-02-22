FROM golang:1.24-alpine as builder 

# build go module
ARG APP_SRC=/home/app/authentication-service
WORKDIR ${APP_SRC}

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /authApp ./cmd/api

# run stage
FROM alpine:3.21 as runner 
COPY --from=builder /authApp .
COPY --from=builder ${APP_SRC} ${APP_SRC}
CMD ["/authApp"]
