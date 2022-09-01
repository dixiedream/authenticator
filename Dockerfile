FROM golang:alpine AS dev
WORKDIR /app
COPY . .
RUN go mod download && \
    go install github.com/cespare/reflex@latest
EXPOSE 3000
CMD reflex -g '**/*.go' go run main.go --start-service

FROM dev AS test
RUN apk update && \
    apk add --no-cache libc-dev gcc 
RUN go test ./...

FROM test AS audit 
COPY --from=aquasec/trivy:latest /usr/local/bin/trivy /usr/local/bin/trivy
RUN trivy filesystem --no-progress /

FROM golang:alpine AS build
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
RUN apk update && \
    apk add --no-cache upx
COPY . .
RUN go mod download && \
    go build -ldflags="-w -s" -o app && \
    upx app

FROM alpine:latest AS prod
EXPOSE 3000
RUN apk add --no-cache ca-certificates bash tini
COPY --from=build app .
COPY ./docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT [ "docker-entrypoint.sh" ]
CMD ./app
