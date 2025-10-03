FROM golang:1.24-alpine
RUN apk add --no-cache git make build-base
WORKDIR /go/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 go build -a -o app ./cmd/telebackup/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /
COPY --from=0 /go/app /

ENTRYPOINT ["./app"]