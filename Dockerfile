FROM ubuntu
COPY dist/lin app
COPY db db
COPY config.json config.json
EXPOSE 80
EXPOSE 8080
EXPOSE 443
ENTRYPOINT ["./app"]