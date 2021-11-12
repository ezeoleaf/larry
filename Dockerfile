FROM golang:alpine AS compiler

# Build directory.
WORKDIR /go/src/larry

# Add go modules and env files to the WORKDIR and install dependencies.
ADD go.mod go.sum ./

# Add code to the WORKDIR and trigger the build process which will assess code quality
# and check if unit tests are passing. Golang binary will be found under /bin/goapp
ADD . ./
RUN go build -o /bin/larry -ldflags="-w -s" cmd/larry/main.go

# Create final application image.
FROM alpine:3.12

RUN apk update && apk upgrade && apk add tzdata

COPY --from=compiler /bin/larry /larry

ENTRYPOINT ["/larry"]
