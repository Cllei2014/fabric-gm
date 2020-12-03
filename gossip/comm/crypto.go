/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"context"
	"crypto/rand"
	"encoding/pem"
	"github.com/tw-bc-group/fabric-gm/bccsp/gm"
	zhsm2 "github.com/tw-bc-group/zhonghuan-ce/sm2"
	"math/big"

	tls "github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	credentials "github.com/Hyperledger-TWGC/tjfoc-gm/gmtls/gmcredentials"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"github.com/tw-bc-group/fabric-gm/common/util"
	"google.golang.org/grpc/peer"
)

// GenerateCertificatesOrPanic generates a a random pair of public and private keys
// and return TLS certificate
func GenerateCertificatesOrPanic() tls.Certificate {
	privateKey, err := zhsm2.CreateSm2KeyAdapter("")
	if err != nil {
		panic(err)
	}

	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		SerialNumber: sn,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	rawBytes, err := x509.CreateCertificate(&template, &template, privateKey.PublicKey(), privateKey)
	if err != nil {
		panic(err)
	}

	encodedCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawBytes})

	cert, err := gm.LoadZHX509KeyPair(encodedCert, []byte(privateKey.KeyID()))
	if err != nil {
		panic(err)
	}
	if len(cert.Certificate) == 0 {
		panic("Certificate chain is empty")
	}

	//FIXME: matrix

	////privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//privateKey, err := sm2.GenerateKey(nil)
	//if err != nil {
	//	panic(err)
	//}
	//sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	//if err != nil {
	//	panic(err)
	//}
	//template := x509.Certificate{
	//	KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	//	SerialNumber: sn,
	//	ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	//}
	//rawBytes, err := x509.CreateCertificate(&template, &template, &privateKey.PublicKey, privateKey)
	//if err != nil {
	//	panic(err)
	//}
	////privBytes, err := x509.MarshalECPrivateKey(privateKey)
	//privBytes, err := x509.MarshalSm2UnecryptedPrivateKey(privateKey)
	//if err != nil {
	//	panic(err)
	//}
	//encodedCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawBytes})
	//encodedKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	//cert, err := tls.X509KeyPair(encodedCert, encodedKey)
	//if err != nil {
	//	panic(err)
	//}
	//if len(cert.Certificate) == 0 {
	//	panic("Certificate chain is empty")
	//}
	return cert
}

func certHashFromRawCert(rawCert []byte) []byte {
	if len(rawCert) == 0 {
		return nil
	}
	return util.ComputeSHA256(rawCert)
}

// ExtractCertificateHash extracts the hash of the certificate from the stream
func extractCertificateHashFromContext(ctx context.Context) []byte {
	pr, extracted := peer.FromContext(ctx)
	if !extracted {
		return nil
	}

	authInfo := pr.AuthInfo
	if authInfo == nil {
		return nil
	}

	tlsInfo, isTLSConn := authInfo.(credentials.TLSInfo)
	if !isTLSConn {
		return nil
	}
	certs := tlsInfo.State.PeerCertificates
	if len(certs) == 0 {
		return nil
	}
	raw := certs[0].Raw
	return certHashFromRawCert(raw)
}
