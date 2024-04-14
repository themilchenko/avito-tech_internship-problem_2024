# syntax=docker/dockerfile:1

FROM golang:1.22.2-alpine3.18 AS BUILDER

RUN go version

COPY . /github.com/themilchenko/avito_internship-problem_2024
WORKDIR /github.com/themilchenko/avito_internship-problem_2024

RUN go mod download
RUN GOOS=linux go build -o ./bin/server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=BUILDER /github.com/themilchenko/avito_internship-problem_2024/bin/server .
COPY --from=BUILDER /github.com/themilchenko/avito_internship-problem_2024/configs/ configs/

EXPOSE 8080

CMD ["./server", "-ConfigPath", "configs/app/deploy.yml"]
