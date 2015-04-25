#!/bin/sh

host="trikipedia"
fqdn="localhost"

rm -i -R RootCA SigningCA1 SigningCA2
rm -i *.pem *.crt *.key *.csr

echo For easy testing and compatability with other scripts set all passwords to \"pass\".
echo ---------------------------------------------------------------------------
echo Create the self signed root certificate
echo
openssl genrsa -aes256 -out root-ca.key 2048
openssl req -new -x509 -days 3650 -key root-ca.key -out root-ca.crt \
	-subj "/C=UK/ST=England/L=Cambridge/O=Trikipedia/CN=${fqdn}"
mkdir RootCA
mkdir RootCA/ca.db.certs
echo "01" >> RootCA/ca.db.serial
touch RootCA/ca.db.index
#echo $RANDOM >> RootCA/ca.db.rand
openssl rand -out RootCA/ca.db.rand 1024
mv root-ca* RootCA/
cp root-CA* RootCA/

echo ---------------------------------------------------------------------------
echo Generate the 2 signing certificates.
echo
for index in `seq 1 1`;
do
	openssl genrsa -aes256 -out signing-ca-${index}.key 2048
	openssl req -new -days 1460 -key signing-ca-${index}.key \
	            -out signing-ca-${index}.csr \
				-subj "/C=UK/ST=England/L=Cambridge/O=Trikipedia/CN=${fqdn}"
	openssl ca -name RootCA -config root-CA.cfg -extensions v3_ca \
	           -out signing-ca-${index}.crt \
	           -infiles signing-ca-${index}.csr
	
	mkdir SigningCA${index}
	mkdir SigningCA${index}/ca.db.certs
	echo "01" >> SigningCA${index}/ca.db.serial
	touch SigningCA${index}/ca.db.index
	openssl rand -out SigningCA${index}/ca.db.rand 1024
	mv signing-ca-${index}* SigningCA${index}/
done

echo ---------------------------------------------------------------------------
echo Generate certificates for clients
echo
openssl genrsa -aes256 -out ${host}.key 2048
openssl req -new -days 365 -key ${host}.key -out ${host}.csr \
            -subj "/C=UK/ST=England/L=Cambridge/O=Trikipedia/CN=${fqdn}"
openssl ca -name SigningCA1 -config root-CA.cfg -out ${host}.crt \
           -infiles ${host}.csr
# Create the .pem file with the certificate and private key
cat ${host}.crt ${host}.key >> ${host}.pem

########

# generate certificates
#openssl req -newkey rsa:2048 -new -x509 -days 365 -nodes -out mongodb-cert.crt -keyout mongodb-cert.key
#cat mongodb-cert.key mongodb-cert.crt > mongodb.pem
