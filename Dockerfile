FROM alpine
RUN  apk add --no-cache ca-certificates tzdata
COPY templates /templates
COPY main /main
CMD ["/main"]
