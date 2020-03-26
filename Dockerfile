FROM golang:1.13.5-alpine AS build
RUN apk add --no-cache tzdata
#ENV GO111MODULE=on
#WORKDIR /app
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#ADD . .
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service .

#FROM scratch as prod
#COPY --from=build /app /
#CMD ["./service"]
#EXPOSE 8080
