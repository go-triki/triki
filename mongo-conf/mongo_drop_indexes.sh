#!/bin/sh

echo "Drop (delete) all indexes in the triki database."
echo MongoDB should be running.
echo ---------------------------------------------------------------------------

source ./mongo_opts.txt

mongo --ssl --sslPEMKeyFile="./trikipedia.pem" --sslPEMKeyPassword="pass" --sslCAFile="./trikipedia.pem" localhost:27017 <<EOF
use ${TRIKI_DB}
db.auth("${TRIKI_USR}", "${TRIKI_PASS}")

db.getCollectionNames().forEach(function(collName) { 
	db.runCommand({dropIndexes: collName, index: "*"});
});

exit
EOF
