# ###
# Build tailwind final css
FROM node:18 AS tailwind

WORKDIR /app

COPY package.json .
COPY tailwind.config.js .
RUN npm i

COPY assets/ ./assets
RUN ls -alh

RUN npx @tailwindcss/cli -i ./assets/tailwind/app.css -o ./assets/css/app.css


# ###
# Build server
FROM golang:1.24 AS go

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=tailwind /app/assets ./assets/

RUN CGO_ENABLED=0 go build \
    -ldflags "-X 'github.com/indeedhat/barista/internal/version.Version=${VERSION}' \
    -X 'github.com/indeedhat/barista/internal/version.BuildTime=$(date +%s)'" \
    -o barista ./cmd/barista/main.go


# ###
# Create final minimal container to expose
FROM alpine:latest

WORKDIR /app

COPY --from=go /app/barista .
RUN mkdir -p ./data/uploads

EXPOSE 8087

CMD ["./barista"]
