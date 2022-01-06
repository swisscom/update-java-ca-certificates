package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"os"
)

var args struct {
	Debug      *bool  `arg:"-D,--debug"`
	Force      *bool  `arg:"-f,--force"`
	File       string `arg:"positional" default:"/etc/ssl/java/cacerts"`
	CertBundle string `arg:"-c,--certificate-bundle" default:"/etc/ssl/certs/ca-certificates.crt"`
}

func main() {
	logger := logrus.New()
	arg.MustParse(&args)

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	outputFileStat, err := os.Stat(args.File)
	if os.IsNotExist(err) {

	}

	p := NewParser(args.CertBundle, logger)
	certs, err := p.getCertificates()
	if err != nil {
		logger.Fatalf("unable to get certificates from bundle: %v" ,err)
	}

	// For each certificate in the system trust store, add them to the keyStore
	for _, c := range certs {
		logger.Debugf("C=%s,O=%s,CN=%s\n", c.Subject.Country, c.Subject.Organization, c.Subject.CommonName)
	}
}