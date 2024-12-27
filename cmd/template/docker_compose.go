package template

var DockerComposeTemplate = `
version: "3.8"

services:
  app:
    container_name: api
    build:
      context: .
      dockerfile: Dockerfile
    tty: true
    ports:
      - 8080127.0.0.1:8080
    depends_on:
      - db
    volumes:
      - .:/app
`
