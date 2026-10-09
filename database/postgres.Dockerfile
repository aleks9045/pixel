FROM postgres:18.6-alpine

COPY ./postgresql.conf /etc/postgresql.conf

EXPOSE 5432

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["postgres", "-c", "config_file=/etc/postgresql.conf"]