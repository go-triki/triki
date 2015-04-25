#!/bin/sh

echo Console for MongoDB.
echo MongoDB should be running.
echo ---------------------------------------------------------------------------

source ./mongo_opts.txt

exec mongo --ssl --sslPEMKeyFile="./trikipedia.pem" --sslPEMKeyPassword="pass" --sslCAFile="./trikipedia.pem" \
	-u "${TRIKI_USR}" -p "${TRIKI_PASS}" "localhost:27017/${TRIKI_DB}"
