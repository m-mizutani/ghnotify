FROM golang:1.18 AS build-go
COPY . /src
WORKDIR /src
RUN go build -o ghnotify .

FROM gcr.io/distroless/base
COPY --from=build-go /src/ghnotify /ghnotify
WORKDIR /
ENTRYPOINT ["/ghnotify"]
