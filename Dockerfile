FROM alpine
COPY dist/lin app
COPY docs_old docs
COPY default-config.json config.json
EXPOSE 80
EXPOSE 8080
EXPOSE 443
ENTRYPOINT ["./app"]
