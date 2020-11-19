package gm

import (
	"github.com/pkg/errors"
	"github.com/Hyperledger-TWGC/tjfoc-gm/sm2"
	"github.com/tw-bc-group/fabric-gm/bccsp"
)

type zhSm2PrivateKey struct {
	sm2PrivateKey *sm2.PrivateKey
	sm2PublicKey *sm2.PublicKey
}

func (sm2 *zhSm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.Errorf("Unsupported")
}

func (sm2 *zhSm2PrivateKey) SKI() []byte {
	return []byte("")
}

func (sm2 *zhSm2PrivateKey) Private() bool {
	return true
}

func (sm2 *zhSm2PrivateKey) Symmetric() bool {
	return false
}

func (sm2 *zhSm2PrivateKey) PublicKey() (bccsp.Key, error) {
	return &gmsm2PublicKey{pubKey: sm2.sm2PublicKey}, nil
}

func createZhSm2PrivateKey() (*zhSm2PrivateKey, error) {
	return &zhSm2PrivateKey{}, nil
}
