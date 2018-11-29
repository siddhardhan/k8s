
# Final image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY /bin/pod-annotator /usr/local/bin/pod-annotator
ENTRYPOINT ["/usr/local/bin/pod-annotator"]