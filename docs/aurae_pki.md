# Creating an Aurae PKI

## Usage

**Create a Root CA**

Get a new CA certificate and key pair. This is the root of trust for all other certificates. The output is a json string.

```bash
$ ./bin/ae pki create unsafe.aurae.io
{
    "cert": "<certificate>",
    "key": "<key>"
}
```

With the flag `-d` corresponding files can be created:

```bash
$ ./bin/ae pki create unsafe.aurae.io -d ./pki/
```

Which is the equivalent of 

```bash
$ openssl req \
    -new \
    -x509 \
    -nodes \
    -days    9999 \
    -addext  "subjectAltName = DNS:unsafe.aurae.io" \
    -subj    "/C=IS/ST=aurae/L=aurae/O=Aurae/OU=Runtime/CN=unsafe.aurae.io" \
    -keyout  "./pki/ca.key" \
    -out     "./pki/ca.crt" 2>/dev/null
```

**Create CSR**

Get a new certificate signing request and key pair. The output is a json string.

```
$ ./bin/ae pki create unsafe.aurae.io --user christoph
{
    "csr": "<certificate request>",
    "key": "<key>",
    "user": "christoph"
}
```

With the flag `-d` corresponding files can be created:

```bash
$ ./bin/ae pki create unsafe.aurae.io --user christoph -d ./pki/
```

Which is the equivalent of

```
$ openssl genrsa -out "./pki/client.${NAME}.key" 4096 2>/dev/null
$ openssl req \
    -new \
    -addext  "subjectAltName = DNS:${NAME}.unsafe.aurae.io" \
    -subj    "/C=IS/ST=aurae/L=aurae/O=Aurae/OU=Runtime/CN=${NAME}.unsafe.aurae.io" \
    -key     "./pki/client.${NAME}.key" \
    -out     "./pki/client.${NAME}.csr" 2>/dev/null
```

**Create client certificate**

*to be implemented*

<!-- 
TODO:
- can we have less arguements?
    - can we derive CA PK from ca.crt?
    - can we derive user from user.csr?

```
$ ./bin/ae pki create unsafe.aurae.io \
    --user christoph \
    --csr ./_tmp/pki/client.christoph.csr
    --ca ./_tmp/pki/ca.crt \
    --caKey ./_tmp/pki/ca.key
```

-->

<!-- **Sign with existing CSR**

```
ae sign <domain> --user=nova --root=ca.crt
$ ./bin/ae sign \
    --user christoph \
    --csr ./_tmp/pki/client.christoph.csr \
    --ca ./_tmp/pki/ca.crt \
    --caKey ./_tmp/pki/ca.key
```
-->

<!--
## Check certificate contents

openssl x509 -noout -ext keyUsage < test.crt
X509v3 Key Usage: critical
    Digital Signature, Key Encipherment
-->
