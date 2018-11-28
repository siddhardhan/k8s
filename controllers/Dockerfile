FROM golang:1.10-alpine AS build-stage
WORKDIR /go/src/github.com/siddhardhan/k8s
COPY . .
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /bin/pod-annotator .

# Final image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-stage /bin/pod-annotator /usr/local/bin/pod-annotator
ENTRYPOINT ["/usr/local/bin/pod-annotator"]