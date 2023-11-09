FROM golang:latest

WORKDIR /app
EXPOSE 1323

ENV PATH="${PATH}:/go/bin"

RUN set -x && go install github.com/cosmtrek/air@latest
RUN set -x && go install github.com/kyoh86/richgo@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 1323

CMD ["air"]