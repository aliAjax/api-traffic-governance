FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/control-api ./cmd/control-api
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/control-api /control-api
EXPOSE 8085
ENTRYPOINT ["/control-api"]
