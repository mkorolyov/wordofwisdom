#Build stage

FROM golang:1.19.4-alpine3.16 as BuildStage

ARG path

WORKDIR /usr/src/app

COPY . .

RUN go build -v -o /app cmd/${path}/main.go

# Deploy stage

FROM alpine:3.16

COPY --from=BuildStage /app /app

ARG addr
ENV addr_env=$addr

ENTRYPOINT "/app" "-addr" "${addr_env}"