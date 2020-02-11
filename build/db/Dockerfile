FROM  postgres:12

COPY migrations/migrations.sql migrations.sql
COPY ./migrate.sh migrate.sh

CMD  ./migrate.sh