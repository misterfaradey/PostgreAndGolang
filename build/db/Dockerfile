FROM  postgres:12

COPY migrations/migrations.sql /docker-entrypoint-initdb.d/migrations.sql
#COPY ./migrate.sh migrate.sh

#CMD  ./migrate.sh