FROM golang:1.15-alpine as build
WORKDIR /src/printMeAt
COPY . .
RUN go mod vendor
RUN go build -mod=vendor

FROM alpine:latest
COPY --from=build /src/printMeAt/printMeAt .
CMD ["./printMeAt"]
EXPOSE 8080