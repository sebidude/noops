FROM scratch

COPY noops /noops

ENTRYPOINT ["/noops"]
CMD ["-c","/config.yaml"]


