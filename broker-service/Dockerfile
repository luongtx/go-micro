FROM golang:1.24-bullseye as debugger 

# installing packages 
RUN apt-get update && apt-get install -y \
    gcc \
    git \
    musl-dev

# build go module
WORKDIR /home/app

COPY . .
# RUN go mod download
# RUN go build -o /brokerApp ./cmd/api
# EXPOSE 8000
# CMD ["/brokerApp"]

# run debugging with dlv
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN chmod +x dlv.sh
EXPOSE 2345
CMD [ "./dlv.sh" ]

FROM alpine:3.21 as runner 
COPY --from=debugger /home/app/brokerApp .
CMD ["./brokerApp"]
