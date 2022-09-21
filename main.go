package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var args struct {
	Debug      *bool  `arg:"-D,--debug"`
	Force      *bool  `arg:"-f,--force"`
	File       string `arg:"positional" default:"/etc/ssl/java/cacerts"`
	CertBundle string `arg:"-c,--certificate-bundle" default:"/etc/ssl/certs/ca-certificates.crt"`
	Password   string `arg:"-p,--password" default:"changeit"`
}

func main() {
	logger := logrus.New()
	arg.MustParse(&args)

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	_, err := os.Stat(args.File)
	if err == nil {
		if args.Force == nil || *args.Force == false {
			logger.Fatalf("%s already exists and --force was not provided", args.File)
		}
	}

	p := NewParser(args.CertBundle, logger)
	certs, err := p.getCertificates()
	if err != nil {
		logger.Fatalf("unable to get certificates from bundle: %v", err)
	}

	outputFile, err := os.OpenFile(args.File, os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()
	if err != nil {
		logger.Fatalf("unable to open %s: %v", args.File, err)
	}

	ks := keystore.New()
	ksCreation := time.Now()

	// For each certificate in the system trust store, add them to the keyStore
	for _, c := range certs {
		fingerprint := fmt.Sprintf("%02x", sha256.Sum256(c.Raw))
		certName := fmt.Sprintf("C=%s,O=%s,CN=%s,fingerprint=%s\n",
			c.Subject.Country,
			c.Subject.Organization,
			c.Subject.CommonName,
			fingerprint,
		)
		err = ks.SetTrustedCertificateEntry(fingerprint, keystore.TrustedCertificateEntry{
			CreationTime: ksCreation,
			Certificate: keystore.Certificate{
				Type:    "X509",
				Content: c.Raw,
			},
		})
		if err != nil {
			logger.Errorf("unable to add %s: %v", certName, err)
		}
	}

	err = ks.Store(outputFile, []byte(args.Password))
	if err != nil {
		logger.Fatalf("unable to store KeyStore: %v", err)
	}
}
