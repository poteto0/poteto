package template

var DockerTemplate = `
FROM golang:1.23

RUN mkdir app
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download && go mod verify

CMD ["go", "run", "/app/main.go"]
`
