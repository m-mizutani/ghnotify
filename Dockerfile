FROM golang:1.20 AS build-go
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -o ghnotify .

FROM gcr.io/distroless/base
COPY --from=build-go /src/ghnotify /ghnotify
WORKDIR /
ENTRYPOINT ["/ghnotify"]
