# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make build

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11:nonroot AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/rinha /rinha

ENV TZ="America/Sao_Paulo"

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/rinha"]