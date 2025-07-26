ARG GO_VERSION="1.24"

FROM golang:${GO_VERSION} AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build ./cmd/sink


FROM alpine

COPY --from=builder /src/sink /

ENTRYPOINT [ "/sink" ]