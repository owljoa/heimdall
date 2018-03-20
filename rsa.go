package heimdall

import (
	"crypto/rsa"
	"errors"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
	"crypto/sha256"
)

type RsaSigner struct{}

func (s *RsaSigner) Sign(key Key, digest []byte, opts SignerOpts) ([]byte, error) {

	if opts == nil {
		return nil, errors.New("invalid options")
	}

	return key.(*RsaPrivateKey).priv.Sign(rand.Reader, digest, opts)

}

type RsaVerifier struct{}

func (v *RsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if opts == nil {
		return false, errors.New("invalid options")
	}

	switch opts.(type) {
	case *rsa.PSSOptions:
		err := rsa.VerifyPSS(key.(*rsaPublicKey).pub,
			(opts.(*rsa.PSSOptions)).Hash,
				digest, signature, opts.(*rsa.PSSOptions))

		if err != nil {
			return false, errors.New("verifying error occurred")
		}

		return true, nil
	default:
		return false, errors.New("invalid options")
	}

}

type rsaKeyMarshalOpt struct {
	N *big.Int
	E int
}

type RsaPrivateKey struct {
	priv *rsa.PrivateKey
}

func (key *RsaPrivateKey) SKI() ([]byte) {

	if key.priv == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		key.priv.N, key.priv.E,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *RsaPrivateKey) Algorithm() string {
	return RSA
}

func (key *RsaPrivateKey) PublicKey() (pub Key, err error) {
	return &rsaPublicKey{&key.priv.PublicKey}, nil
}

type rsaPublicKey struct {
	pub *rsa.PublicKey
}

func (key *rsaPublicKey) SKI() ([]byte) {

	if key.pub == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		big.NewInt(123), 57,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (key *rsaPublicKey) Algorithm() string {
	return RSA
}