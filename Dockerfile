FROM golang:1.21 as dependencies
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
ENV TZ=Europe/Berlin
WORKDIR /app/
VOLUME ["/app"]
RUN apk add --no-cache tzdata
COPY --from=builder /go/src/app/awb-kh-api .
CMD ["./awb-kh-api"]
