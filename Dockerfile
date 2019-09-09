FROM golang:1.13 as builder

WORKDIR /usr/src
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o /go/bin/dafang-exporter .

FROM gcr.io/distroless/base
COPY --from=builder /go/bin/dafang-exporter /
CMD ["/dafang-exporter"]