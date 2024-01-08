FROM mongo:6

RUN openssl rand -base64 756 > /mongodb.key
RUN chmod 400 /mongodb.key
RUN chown 999:999 /mongodb.key

CMD ["exec", "docker-entrypoint.sh", "$$@"]
