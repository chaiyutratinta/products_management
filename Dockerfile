FROM golang:1.13.1 AS builder

ARG app_name=products_management

WORKDIR /go/src/${app_name}
COPY . /go/src/${app_name}
RUN go get . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  

ARG app_name=products_management

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/${app_name}/app .

CMD ["./app"]  
