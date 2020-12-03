package gm

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"github.com/Hyperledger-TWGC/tjfoc-gm/sm2"

	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	x509GM "github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"github.com/pkg/errors"
	zhsm2 "github.com/tw-bc-group/zhonghuan-ce/sm2"
	"strings"
)

func LoadZHX509KeyPair(certPEMBlock, keyBytes []byte) (gmtls.Certificate, error) {
	fail := func(err error) (gmtls.Certificate, error) { return gmtls.Certificate{}, err }

	var cert gmtls.Certificate
	var skippedBlockTypes []string
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
		}
	}

	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			return fail(errors.New("tls: failed to find any PEM data in certificate input"))
		}
		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			return fail(errors.New("tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
		}
		return fail(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
	}

	priAdapter, err := zhsm2.CreateSm2KeyAdapter(strings.TrimSpace(string(keyBytes)))
	if err != nil {
		return fail(fmt.Errorf("tls: failed to load zh sm2 key adapter, Got err: %s", err))
	}

	x509Cert, err := x509GM.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fail(err)
	}

	cert.PrivateKey = priAdapter

	//check pub and private match
	switch pub := x509Cert.PublicKey.(type) {
	case *rsa.PublicKey:
		priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return fail(errors.New("tls: private key type does not match public key type"))
		}
		if pub.N.Cmp(priv.N) != 0 {
			return fail(errors.New("tls: private key does not match public key"))
		}
	case *ecdsa.PublicKey:
		pub, _ = x509Cert.PublicKey.(*ecdsa.PublicKey)
		switch pub.Curve {
		case sm2.P256Sm2():
			zhPubKey := priAdapter.PublicKey()
			if pub.X.Cmp(zhPubKey.X) != 0 || pub.Y.Cmp(zhPubKey.Y) != 0 {
				return fail(errors.New("tls: zh sm2 private key does not match public key"))
			}
		default:
			return fail(errors.New("tls: private key type does not match public key type"))
		}
	default:
		return fail(errors.New("tls: unknown public key algorithm"))
	}

	return cert, nil
}
