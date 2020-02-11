#!/bin/sh
#set -e

dbExist=$(psql -U postgres -lqt | cut -d \| -f 1 | grep -qw test)

if [ -z "$dbExist" ]; then
  psql -U postgres -a -f /migrations.sql
  echo "migrate complete"
fi