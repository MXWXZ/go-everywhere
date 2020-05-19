FROM golang:latest AS build
WORKDIR /build
COPY . .

ENV GOFLAGS -mod=vendor
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -o go-everywhere main.go

FROM scratch AS prod
WORKDIR /app
COPY --from=build /build/go-everywhere .
ENV GIN_MODE release

expose 8080
CMD ["./go-everywhere"]
