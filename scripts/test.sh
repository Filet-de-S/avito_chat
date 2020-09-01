#!/bin/sh

set -e

printf "Waiting for API up:\n"

printf "\n Running entrypoint.sh \n"
sh $SCRIPTS_FOLDER/entrypoint.sh

printf "\n Running psql \n"
PGPASSFILE=$PGPASSFILE psql -a -h $PG_SNAME -d msg_service -U service_client \
                            -f $SCRIPTS_FOLDER/001_relations.clear.sql

printf "\n Running Newman Tests \n"
newman run $POSTMAN_FOLDER/i_test.postman_collection.json --folder scenario \
        -r html \
        --reporter-html-export $POSTMAN_FOLDER/report.html \
        --reporter-html-template $POSTMAN_FOLDER/template.hbs \
        --bail \
        --global-var baseUrl=$SERVICE_HOST:$SERVICE_PORT
printf "\n Passed Newman Tests \n"
