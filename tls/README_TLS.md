# TLS Certs

This directory should hold the TLS key pair and nothing else. Generate them using something like this for local development.

```sh
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

You'll end up with a private key (`key.pem`) and public key (`cert.pem`) that won't go into VCS.
