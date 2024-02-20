FROM golang:1.20-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main cmd/main.go

FROM scratch
COPY ./migrations /migrations
COPY --from=builder /bin/main /bin/main
ENTRYPOINT ["/bin/main"]