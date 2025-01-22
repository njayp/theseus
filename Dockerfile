FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o output/theseus ./cmd/main.go

FROM alpine
COPY --from=build /app/output/theseus /bin/theseus
ENTRYPOINT ["theseus"]
