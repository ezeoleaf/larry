FROM alpine:3.12 AS base
WORKDIR /go/src/larry
RUN apk update && apk upgrade && apk add tzdata

FROM golang:alpine AS dev
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg mkdir -p go/bin && cd /go/bin && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	go install golang.org/x/tools/gopls@latest &&\
    apk add --no-cache bash git make

FROM dev AS compiler

# Add go modules and env files to the WORKDIR and install dependencies.
ADD go.mod go.sum ./

# Add code to the WORKDIR and trigger the build process which will assess code quality
# and check if unit tests are passing. Golang binary will be found under /bin/goapp
ADD . ./
RUN go build -o /bin/larry -ldflags="-w -s" cmd/larry/main.go

# Create final application image.
FROM base AS final

COPY --from=compiler /bin/larry /larry

ENTRYPOINT ["/larry"]
