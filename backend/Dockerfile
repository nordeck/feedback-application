FROM golang:1.19.3-alpine AS builder

RUN apk update && apk add --no-cache gcc && apk add --no-cache libc-dev && apk add --no-cache make
COPY . /build/

RUN make -C /build build

FROM alpine:3.16.2
COPY --from=builder /build/out/feedback-api /feedback-api
EXPOSE 8080
CMD ["/feedback-api"]
