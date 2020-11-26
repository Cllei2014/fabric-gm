package gm

import (
	"github.com/pkg/errors"
	"github.com/tw-bc-group/fabric-gm/bccsp"
	"github.com/tw-bc-group/mock-collaborative-encryption-lib/sm2"
)

type zhSm2PrivateKey struct {
	adapter *sm2.KeyAdapter
}

func (sm2 *zhSm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.Errorf("Unsupported")
}

func (sm2 *zhSm2PrivateKey) SKI() []byte {
	return []byte(sm2.adapter.KeyID())
}

func (sm2 *zhSm2PrivateKey) Private() bool {
	return true
}

func (sm2 *zhSm2PrivateKey) Symmetric() bool {
	return false
}

func (sm2 *zhSm2PrivateKey) PublicKey() (bccsp.Key, error) {
	pubKey, err := sm2.adapter.GetPublicKey()
	if err != nil {
		return nil, err
	}
	return &gmsm2PublicKey{pubKey: pubKey}, nil
}

func createZhSm2PrivateKey() (*zhSm2PrivateKey, error) {
	adapter, err := sm2.CreateSm2KeyAdapter(sm2.SignAndVerify, "")
	if err != nil {
		return nil, err
	}

	return &zhSm2PrivateKey{
		adapter: adapter,
	}, nil
}

type zhSm2PrivateKeySigner struct{}

func (s *zhSm2PrivateKeySigner) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	return k.(*zhSm2PrivateKey).adapter.AsymmetricSign(digest)
}
