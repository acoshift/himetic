FROM acoshift/go-scratch

ADD himetic /
COPY view /view
COPY assets /assets
COPY static.yaml /static.yaml
EXPOSE 8080

ENTRYPOINT ["/himetic"]
