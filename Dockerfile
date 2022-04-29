FROM scratch
COPY --from=alpine:latest /tmp /tmp
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dist/lin app
COPY default-config.json config.json
EXPOSE 80
EXPOSE 8080
EXPOSE 443
ENTRYPOINT ["./app"]
