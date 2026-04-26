FROM golang:1.24.0

WORKDIR /app
COPY . .

RUN go build -o /app/bin/c-share /app/cmd/c-share/main.go
ENTRYPOINT [ "/app/bin/c-share" ]
