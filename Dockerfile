# syntax=docker/dockerfile:1.4

FROM golang:1.22-alpine AS build-dev
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache upx || \
RUN go run github.com/steebchen/prisma-client-go prefetch
COPY ./ ./
RUN go run github.com/steebchen/prisma-client-go generate
RUN CGO_ENABLED=0 go install -buildvcs=false -trimpath -ldflags '-w -s -extldflags "-static"'
#RUN CGO_ENABLED=0 GOOS=linux go install -buildvcs=false -trimpath -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

RUN [ -e /usr/bin/upx ] && upx /go/bin/go-prisma-example || echo
FROM scratch
COPY --link --from=build-dev /go/bin/go-prisma-example /go/bin/go-prisma-example
COPY --from=build-dev /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/go/bin/go-prisma-example"]
