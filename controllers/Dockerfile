FROM golang:1.10-alpine AS build-stage
ARG example
WORKDIR /go/src/github.com/siddhardhan/k8s
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/pod-annotator --ldflags "-w -extldflags '-static'"  github.com/siddhardhan/k8s/controllers

# Final image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-stage /bin/pod-annotator /usr/local/bin/pod-annotator
ENTRYPOINT ["/usr/local/bin/pod-annotator"]