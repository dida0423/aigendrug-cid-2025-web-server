FROM postgres:16

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=aigendrug

COPY main-db-init.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]