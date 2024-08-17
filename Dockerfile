# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

ENV GOMODCACHE=/root/.cache/gocache

WORKDIR /app

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=bind,target=. \
    go install .

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/qure /bin/

ENTRYPOINT [ "/bin/qure" ]
