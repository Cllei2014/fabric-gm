package gm

import (
	"github.com/pkg/errors"
	"github.com/tw-bc-group/aliyun-kms/sm2"
	"github.com/tw-bc-group/fabric-gm/bccsp"
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
	pubKey, err := sm2.adapter.GetPublicKey()
	if err != nil {
		return nil, err
	}
	return &gmsm2PublicKey{pubKey: pubKey}, nil
}

func createKmsSm2PrivateKey() (*kmsSm2PrivateKey, error) {
	adapter, err := sm2.CreateSm2KeyAdapter("", sm2.SignAndVerify)
	if err != nil {
		return nil, err
	}

	return &kmsSm2PrivateKey{
		adapter: adapter,
	}, nil
}
