package gm

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/tw-bc-group/fabric-gm/bccsp"
	"github.com/tw-bc-group/zhonghuan-ce/sm2"
)

type zhSm2PrivateKey struct {
	adapter *sm2.KeyAdapter
}

func (sm2 *zhSm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.Errorf("Unsupported")
}

func (sm2 *zhSm2PrivateKey) SKI() []byte {
	publicKey := sm2.adapter.PublicKey()
	raw := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (sm2 *zhSm2PrivateKey) Private() bool {
	return true
}

func (sm2 *zhSm2PrivateKey) Symmetric() bool {
	return false
}

func (sm2 *zhSm2PrivateKey) PublicKey() (bccsp.Key, error) {
	return &gmsm2PublicKey{pubKey: sm2.adapter.PublicKey()}, nil
}

func createZhSm2PrivateKey() (*zhSm2PrivateKey, error) {
	adapter, err := sm2.CreateSm2KeyAdapter("")
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

type zhsm2Decryptor struct{}

func (s *zhsm2Decryptor) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {
	return k.(*zhSm2PrivateKey).adapter.Decrypt(rand.Reader, ciphertext, opts)
}
