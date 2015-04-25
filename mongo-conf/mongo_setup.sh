#!/bin/bash

source ./mongo_opts.txt

mkdir db

echo
echo Start MongoDB and use localhost exception to create admin user account.
echo Log in as admin and create triki database and user.
echo ---------------------------------------------------------------------------

sh ./mongo_start.sh &
PID=$!

sleep 5

mongo --ssl --sslPEMKeyFile="./trikipedia.pem" --sslPEMKeyPassword="pass" --sslCAFile="./trikipedia.pem" localhost:27017/admin <<EOF
use admin
db.createUser(
  {
    user: "${ADMIN_USR}",
    pwd: "${ADMIN_PASS}",
    roles: [ { role: "root", db: "admin" } ]
  }
)

db.auth("${ADMIN_USR}", "${ADMIN_PASS}")
use ${TRIKI_DB}
db.createUser(
  {
    user: "${TRIKI_USR}",
    pwd: "${TRIKI_PASS}",
    roles: [ { role: "readWrite", db: "${TRIKI_DB}" } ]
  }
)

exit
EOF

echo
echo Shutdown MongoDB.
echo ---------------------------------------------------------------------------

kill -s INT ${PID}

wait
echo DONE
