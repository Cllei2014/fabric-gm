package gm

import (
	"encoding/pem"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pkg/errors"
	"github.com/tw-bc-group/aliyun-kms/sm2"
	"github.com/tw-bc-group/fabric-gm/bccsp"
	"os"
)

type kmsSm2PrivateKey struct {
	adapter *sm2.KeyAdapter
}

func (sm2 *kmsSm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.Errorf("Unsupported")
}

func (sm2 *kmsSm2PrivateKey) SKI() []byte {
	return []byte(sm2.adapter.KeyID())
}

func (sm2 *kmsSm2PrivateKey) Symmetric() bool {
	return false
}

func (sm2 *kmsSm2PrivateKey) Private() bool {
	return true
}

func (sm2 *kmsSm2PrivateKey) PublicKey() (bccsp.Key, error) {
	pemPubKey, err := sm2.adapter.GetPublicKey()
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(pemPubKey))
	pubKey, err := x509.ParseSm2PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &gmsm2PublicKey{pubKey: pubKey}, nil
}

func createKmsSm2PrivateKey() (*kmsSm2PrivateKey, error) {
	client, err := kms.NewClientWithAccessKey(os.Getenv("ALIBABA_CLOUD_REGION"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY"), os.Getenv("ALIBABA_CLOUD_ACCESS_SECRET"))
	if err != nil {
		return nil, err
	}

	adapter, err := sm2.CreateSm2KeyAdapter(client, sm2.SignAndVerify, "")
	if err != nil {
		return nil, err
	}

	return &kmsSm2PrivateKey{
		adapter: adapter,
	}, nil
}
