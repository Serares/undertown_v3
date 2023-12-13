FROM public.ecr.aws/docker/library/golang:1.21-alpine3.18 AS builder
RUN mkdir /addproperty
WORKDIR /addproperty
# TODO this looks like a bad idea
RUN mkdir -p ./internal/database
RUN mkdir -p ./repositories
COPY go.mod go.sum ./
COPY services/api/addProperty/ ./
COPY internal/database/ ./internal/database
COPY repositories/ ./repositories
RUN ls -lah
RUN tree ./
RUN CGO_ENABLED=0 GOOS=linux go build -o addproperty -ldflags="-s -w" .
FROM public.ecr.aws/docker/library/alpine:latest
RUN mkdir /app && adduser -h /app -D addproperty
WORKDIR /app
COPY --chown=addproperty --from=builder /addproperty/addproperty .
CMD ["/app/addproperty"]