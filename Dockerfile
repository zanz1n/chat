FROM docker.io/library/golang:alpine AS builder

WORKDIR /go/src/izanr.com/chat

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN apk add --no-cache just git

COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    just BIN=/go/bin build

FROM gcr.io/distroless/static-debian13

COPY --from=builder /go/bin/chat /bin/chat

ENTRYPOINT [ "/bin/chat" ]
CMD [ "--migrate" ]
