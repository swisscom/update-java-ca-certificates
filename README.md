# update-java-ca-certificates

This small utility takes care of creating a system-wide trust store
starting from your Linux CA trust store.

This command is supposed to be run after running 
[`update-ca-certificates (8)`](https://manpages.ubuntu.com/manpages/jammy/man8/update-ca-certificates.8.html),
so that the Java Keystore is in sync with the system trust store.  
  
The issue that this tool is trying to solve is already solved by Arch Linux's 
[update-ca-trust (8)](https://man.archlinux.org/man/update-ca-trust.8). Sadly not all the Linux distributions
have solved the issue (yet), thus this is a tool to help standardize the mess that's currently out there in terms
of path standardization and ca-certificates location.

## Usage

```
Usage: update-java-ca-certificates [--debug] [--force] [--certificate-bundle CERTIFICATE-BUNDLE] [--password PASSWORD] FILE

Positional arguments:
  FILE

Options:
  --debug, -D
  --force, -f
  --certificate-bundle CERTIFICATE-BUNDLE, -c CERTIFICATE-BUNDLE [default: /etc/ssl/certs/ca-certificates.crt]
  --password PASSWORD, -p PASSWORD [default: changeit]
  --help, -h             display this help and exit
```

### Example

```
update-java-ca-certificates -c /etc/ssl/certs/ca-certificates.crt /etc/ssl/java/cacerts
```

## Building

### Requirements

- Golang (1.17+)
- Make

### Steps

```
make
./bin/update-java-ca-certificates -h
```

## Paths

This tool assumes the directories are set up according to what
[update-ca-trust (8)](https://man.archlinux.org/man/update-ca-trust.8) uses.

### `/etc/ssl/certs`

This directory should contain individual CA certificates trusted for TLS authentication usage.
The format to be used is the `BEGIN CERTIFICATE` / `END CERTIFICATE` one.

If you are able to parse the certificate with:
```
openssl x509 -in /etc/ssl/certs/your-certificate.pem  -noout -text
```

then you're good.

### `/etc/ssl/ca-certificates.crt`

This file contains a bundle that is updated by `update-ca-trust` / `update-ca-certificates`.

### `/etc/ssl/java/cacerts`

This file contains the trust anchor for Java. Its format is the 
[Java Key Store (`JKS`)](https://docs.oracle.com/javase/7/docs/technotes/guides/security/crypto/CryptoSpec.html#KeyManagement).