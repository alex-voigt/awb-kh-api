FROM golang:1.17.7 as dependencies
WORKDIR /go/src/app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM dependencies as builder
COPY . .
RUN go test ./... -timeout 30s -cover
RUN CGO_ENABLED=0 go build -o awb-kh-api

FROM alpine:latest
LABEL maintainer="Alex Voigt <mail@alexander-voigt.info>"
WORKDIR /app/
VOLUME ["/app"]
COPY --from=builder /go/src/app/awb-kh-api .
CMD ["./awb-kh-api"]
