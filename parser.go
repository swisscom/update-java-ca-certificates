package main

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Parser struct {
	bundlePath string
	logger     *logrus.Logger
}

func NewParser(certBundle string, logger *logrus.Logger) *Parser {
	p := Parser{
		bundlePath: certBundle,
		logger:     logger,
	}
	return &p
}

func (p *Parser) getCertificates() ([]x509.Certificate, error) {
	f, err := os.Open(p.bundlePath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to open certificate bundle file: %v", err)
	}

	s := bufio.NewScanner(f)
	s.Split(splitByBeginCertificate)

	var certs []x509.Certificate

	for s.Scan() {
		// Read PEM block
		b, _ := pem.Decode(s.Bytes())
		if b.Type != "CERTIFICATE" {
			p.logger.Errorf("invalid type found: %s", b.Type)
			continue
		}
		c, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			p.logger.Errorf("unable to parse certificate: %v", err)
			continue
		}

		certs = append(certs, *c)
	}
	return certs, nil
}

func splitByBeginCertificate(data []byte, eof bool) (advance int, token []byte, err error) {
	if eof && len(data) == 0 {
		return 0, nil, nil
	}
	endCertificate := "-----END CERTIFICATE-----"
	idx := bytes.Index(data, []byte(endCertificate))
	if idx >= 0 {
		trimmed := strings.TrimSpace(string(data[0 : idx+len(endCertificate)]))
		return idx + len(endCertificate), []byte(trimmed), nil
	}

	// EOF but non terminated string (well, we can't terminate with a BEGIN CERTIFICATE)
	if eof {
		return len(data), nil, nil
	}

	// Get more data
	return 0, nil, nil
}
