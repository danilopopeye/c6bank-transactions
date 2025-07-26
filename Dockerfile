FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build -trimpath -ldflags="-s -w" \
  -o /app/c6bank-transactions ./cmd/c6bank-transactions

FROM alpine:3
EXPOSE 4500
RUN apk add --no-cache tesseract-ocr tesseract-ocr-data-eng tesseract-ocr-data-por
USER nobody
COPY --from=builder /app/c6bank-transactions /c6bank-transactions
CMD ["/c6bank-transactions"]
