#!/bin/sh

set -e

service=$SERVICE_NAME
host=$SERVICE_HOST
port=$SERVICE_PORT
url_path=$URL_PATH

#cmd=$@

>&2 echo "Checking service $service for UP on: $host:$port"

until curl http://$host:$port/$url_path; do
  >&2 echo $host:$port" is unavailable, sleeping"
  sleep 1
done

>&2 echo $host:$port" is up"

#exec $cmd