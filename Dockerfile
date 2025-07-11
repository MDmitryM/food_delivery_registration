FROM golang:1.23 AS build

ENV GOPATH=/
WORKDIR /src/
COPY . .

RUN go mod download; CGO_ENABLED=0 go build -o /food_delivery_registtration ./cmd/main.go

FROM alpine:3.17
COPY --from=build /food_delivery_registtration /food_delivery_registtration

EXPOSE 8000
CMD [ "./food_delivery_registtration" ]