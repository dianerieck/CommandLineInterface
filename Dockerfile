FROM golang:1.22-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /stress ./cmd

FROM scratch
COPY --from=build /stress /stress
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/stress"]