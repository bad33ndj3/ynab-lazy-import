FROM golang:alpine3.12

RUN apk add --no-cache git

WORKDIR /src/ynab-lazy-import

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/ynab-lazy-import .

ENTRYPOINT ["./out/ynab-lazy-import"]