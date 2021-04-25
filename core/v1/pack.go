package v1

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/jiexun/minijkstool/jks"

	cli "gopkg.in/urfave/cli.v2"
)

var CertFilePackCommand = &cli.Command{
	Name:      "packfile",
	Usage:     "pack a certfile into a keystore jks file",
	ArgsUsage: "ca.crt out.jks",
	Action:    CertFilePack,
}

func CertFilePack(c *cli.Context) error {
	switch c.NArg() {
	case 0:
		cli.ShowSubcommandHelp(c)
		return errors.New("need input a cert file and output jks file name")

	case 2:
		// OK

	default:
		return errors.New("need input a cert file and output jks file name")
	}

	certFile := c.Args().Get(0)
	outFn := c.Args().Get(1)

	st, err := os.Stat(certFile)
	if err != nil {
		return err
	} else if st.IsDir() {
		return fmt.Errorf("%q must be a file", certFile)
	}

	f, err := os.OpenFile(outFn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}

	if err = pack(f, certFile); err != nil {
		_ = f.Close()
		_ = os.Remove(outFn)
		return err
	}

	if err = f.Close(); err != nil {
		_ = os.Remove(outFn)
		return err
	}

	return nil
}

func pack(out io.Writer, certFile string) error {

	var (
		err  error
		ks   jks.Keystore
		opts = jks.Options{
			KeyPasswords: make(map[string]string),
		}
	)
	// TODO  put a stable password
	opts.Password = "changeit"

	if _, err = os.Stat(certFile); err == nil {
		if err = packCerts(&opts, &ks, certFile); err != nil {
			return err
		}
	}


	raw, err := ks.Pack(&opts)
	if err != nil {
		return err
	}

	_, err = out.Write(raw)
	return err
}

func packPassword(dirname string) (string, error) {
	fn := filepath.Join(dirname, "password")
	p, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	// strip a possible trailing newline
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}

	// ensure it's valid UTF-8
	if !utf8.Valid(p) {
		return "", fmt.Errorf("%s: not valid UTF-8", fn)
	}
	return string(p), nil
}

func packCerts(opts *jks.Options, ks *jks.Keystore, certFile string) error {

	file, err := os.Stat(certFile)
	if err != nil {
		return err
	}

	if file.IsDir() || file.Name()[0] == '.' ||
		filepath.Ext(file.Name()) != ".pem" {
		fmt.Fprintf(os.Stderr, "ignoring %q (must be "+
			"non-dot-file ending .pem)\n", file.Name())

		return errors.New("the cert file is not avilabel")
	}

	cert, err := packLoadCert(certFile)
	if err != nil {
		return err
	}

	alias := filepath.Base(file.Name())
	alias = alias[:len(alias)-4] // strip ".pem"
	ks.Certs = append(ks.Certs, &jks.Cert{
		Alias:     alias,
		Timestamp: file.ModTime(),
		Cert:      cert,
	})

	return nil
}


func packLoadPem(fname string) (*pem.Block, error) {
	pemraw, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(pemraw)
	if block == nil {
		return nil, fmt.Errorf("%q: not a PEM file", fname)
	} else if len(rest) != 0 {
		return nil, fmt.Errorf("%q: has data beyond first PEM block",
			fname)
	}
	return block, nil
}

func packLoadCert(fname string) (*x509.Certificate, error) {
	block, err := packLoadPem(fname)
	if err != nil {
		return nil, err
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("%q: expected CERTIFICATE but found %q",
			fname, block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", fname, err)
	}
	return cert, nil
}
