/*
 * Nuts node
 * Copyright (C) 2021 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package crypto

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nuts-foundation/nuts-node/audit"
	"github.com/nuts-foundation/nuts-node/crypto/log"
	"github.com/nuts-foundation/nuts-node/crypto/storage/spi"
	"github.com/shengdoushi/base58"
)

// ErrUnsupportedSigningKey is returned when an unsupported private key is used to sign. Currently only ecdsa and rsa keys are supported
var ErrUnsupportedSigningKey = errors.New("signing key algorithm not supported")

var supportedAlgorithms = []jwa.SignatureAlgorithm{jwa.PS256, jwa.PS384, jwa.PS512, jwa.ES256, jwa.ES384, jwa.ES512}

func isAlgorithmSupported(alg jwa.SignatureAlgorithm) bool {
	for _, curr := range supportedAlgorithms {
		if curr == alg {
			return true
		}
	}
	return false
}

// SignJWT creates a JWT from the given claims and signs it with the given key.
func (client *Crypto) SignJWT(ctx context.Context, claims map[string]interface{}, key interface{}) (string, error) {
	privateKey, kid, err := client.getPrivateKey(key)
	if err != nil {
		return "", err
	}

	audit.Log(ctx, log.Logger(), audit.CryptoSignJWTEvent).
		Infof("Signing a JWT with key: %s (issuer: %s, subject: %s)", kid, claims["iss"], claims["sub"])

	keyAsJWK, err := jwkKey(privateKey)
	if err != nil {
		return "", err
	}

	if err = keyAsJWK.Set(jwk.KeyIDKey, kid); err != nil {
		return "", err
	}

	return signJWT(keyAsJWK, claims, nil)
}

// SignJWS creates a signed JWS using the indicated key and map of headers and payload as bytes.
func (client *Crypto) SignJWS(ctx context.Context, payload []byte, headers map[string]interface{}, key interface{}, detached bool) (string, error) {
	privateKey, kid, err := client.getPrivateKey(key)
	if err != nil {
		return "", err
	}

	audit.Log(ctx, log.Logger(), audit.CryptoSignJWSEvent).Infof("Signing a JWS with key: %s", kid)

	return signJWS(payload, headers, privateKey, detached)
}

func jwkKey(signer crypto.Signer) (key jwk.Key, err error) {
	key, err = jwk.New(signer)
	if err != nil {
		return nil, err
	}

	switch k := signer.(type) {
	case *rsa.PrivateKey:
		key.Set(jwk.AlgorithmKey, jwa.PS256)
	case *ecdsa.PrivateKey:
		var alg jwa.SignatureAlgorithm
		alg, err = ecAlg(k)
		key.Set(jwk.AlgorithmKey, alg)
	default:
		err = errors.New("unsupported signing private key")
	}
	return
}

// signJWT signs claims with the signer and returns the compacted token. The headers param can be used to add additional headers
func signJWT(key jwk.Key, claims map[string]interface{}, headers map[string]interface{}) (token string, err error) {
	var sig []byte
	t := jwt.New()

	for k, v := range claims {
		if err := t.Set(k, v); err != nil {
			return "", err
		}
	}
	hdr := convertHeaders(headers)

	sig, err = jwt.Sign(t, jwa.SignatureAlgorithm(key.Algorithm()), key, jwt.WithHeaders(hdr))
	token = string(sig)

	return
}

// JWTKidAlg parses a JWT, does not validate it and returns the 'kid' and 'alg' headers
func JWTKidAlg(tokenString string) (string, jwa.SignatureAlgorithm, error) {
	j, err := jws.ParseString(tokenString)
	if err != nil {
		return "", "", err
	}

	if len(j.Signatures()) != 1 {
		return "", "", errors.New("incorrect number of signatures in JWT")
	}

	sig := j.Signatures()[0]
	hdrs := sig.ProtectedHeaders()
	return hdrs.KeyID(), hdrs.Algorithm(), nil
}

// PublicKeyFunc defines a function that resolves a public key based on a kid
type PublicKeyFunc func(kid string) (crypto.PublicKey, error)

// ParseJWT parses a token, validates and verifies it.
func ParseJWT(tokenString string, f PublicKeyFunc, options ...jwt.ParseOption) (jwt.Token, error) {
	kid, alg, err := JWTKidAlg(tokenString)
	if err != nil {
		return nil, err
	}

	key, err := f(kid)
	if err != nil {
		return nil, err
	}

	if !isAlgorithmSupported(alg) {
		return nil, fmt.Errorf("token signing algorithm is not supported: %s", alg)
	}

	options = append(options, jwt.WithVerify(alg, key))
	options = append(options, jwt.WithValidate(true))

	return jwt.ParseString(tokenString, options...)
}

// ParseJWS parses a JWS byte array object, validates and verifies it.
// This method returns the value of the payload as byte array, or an error if
// the parsing fails at any level.
func ParseJWS(token []byte, f PublicKeyFunc) (payload []byte, err error) {
	message, err := jws.Parse(token)
	if err != nil {
		return nil, err
	}
	headers, body, _, err := jws.SplitCompact(token)
	if err != nil {
		return nil, err
	}
	signatures := message.Signatures()
	for i := range signatures {
		signature := signatures[i]
		// Get and check the algorithm
		alg := signature.ProtectedHeaders().Algorithm()
		if !isAlgorithmSupported(alg) {
			return nil, fmt.Errorf("token signing algorithm is not supported: %s", alg)
		}
		// Get the verifier for the algorithm
		verifier, err := jws.NewVerifier(alg)
		if err != nil {
			return nil, err
		}
		// Get the key id, and get the associated key
		kid := signature.ProtectedHeaders().KeyID()
		key, err := f(kid)
		if err != nil {
			return nil, err
		}
		// This seems an awkward way of appending 3 arrays.
		var payload []byte
		parts := [][]byte{headers, []byte("."), body}
		for _, part := range parts {
			payload = append(payload, part...)
		}
		err = verifier.Verify(payload, signature.Signature(), key)
		if err != nil {
			return nil, err
		}
	}

	body = message.Payload()
	return body, nil
}

func signJWS(payload []byte, protectedHeaders map[string]interface{}, privateKey crypto.Signer, detachedPayload bool) (string, error) {
	headers := jws.NewHeaders()
	for key, value := range protectedHeaders {
		if err := headers.Set(key, value); err != nil {
			return "", fmt.Errorf("unable to set header %s: %w", key, err)
		}
	}
	privateKeyAsJWK, err := jwkKey(privateKey)
	if err != nil {
		return "", err
	}
	// The JWX library is fine with creating a JWK for a private key (including the private exponents), so
	// we want to make sure the `jwk` header (if present) does not (accidentally) contain a private key.
	// That would lead to the node leaking its private key material in the resulting JWS which would be very, very bad.
	if headers.JWK() != nil {
		var jwkAsPrivateKey crypto.Signer
		if err := headers.JWK().Raw(&jwkAsPrivateKey); err == nil {
			// `err != nil` is good in this case, because that means the key is not assignable to crypto.Signer,
			// which is the interface implemented by all private key types.
			return "", errors.New("refusing to sign JWS with private key in JWK header")
		}
	}
	algo := jwa.SignatureAlgorithm(privateKeyAsJWK.Algorithm())

	var (
		data []byte
	)
	if detachedPayload {
		// Sign JWS with detached payload
		data, err = jws.Sign(nil, algo, privateKey, jws.WithHeaders(headers), jws.WithDetachedPayload(payload))
	} else {
		// Sign normal JWS
		data, err = jws.Sign(payload, algo, privateKey, jws.WithHeaders(headers))
	}
	if err != nil {
		return "", fmt.Errorf("unable to sign JWS %w", err)
	}
	return string(data), nil
}

func (client *Crypto) getPrivateKey(key interface{}) (crypto.Signer, string, error) {
	var kid string
	switch k := key.(type) {
	case exportableKey:
		return k.Signer(), k.KID(), nil
	case Key:
		kid = k.KID()
	case string:
		kid = k
	default:
		return nil, "", errors.New("provided key must be either string or Key")
	}

	privateKey, err := client.storage.GetPrivateKey(kid)
	if err != nil {
		if errors.Is(err, spi.ErrNotFound) {
			return nil, "", ErrPrivateKeyNotFound
		}
		return nil, "", err
	}
	return privateKey, kid, nil
}

func convertHeaders(headers map[string]interface{}) (hdr jws.Headers) {
	hdr = jws.NewHeaders()

	for k, v := range headers {
		hdr.Set(k, v)
	}
	return
}

func ecAlg(key *ecdsa.PrivateKey) (alg jwa.SignatureAlgorithm, err error) {
	alg, err = ecAlgUsingPublicKey(key.PublicKey)
	return
}

func ecAlgUsingPublicKey(key ecdsa.PublicKey) (alg jwa.SignatureAlgorithm, err error) {
	switch key.Params().BitSize {
	case 256:
		alg = jwa.ES256
	case 384:
		alg = jwa.ES384
	case 521:
		alg = jwa.ES512
	default:
		err = ErrUnsupportedSigningKey
	}
	return
}

// SignatureAlgorithm determines the jwa.SigningAlgorithm for ec/rsa/ed25519 keys.
func SignatureAlgorithm(key crypto.PublicKey) (jwa.SignatureAlgorithm, error) {
	if key == nil {
		return "", errors.New("no key provided")
	}

	var ptr interface{}
	switch v := key.(type) {
	case rsa.PrivateKey:
		ptr = &v
	case rsa.PublicKey:
		ptr = &v
	case ecdsa.PrivateKey:
		ptr = &v
	case ecdsa.PublicKey:
		ptr = &v
	default:
		ptr = v
	}

	switch k := ptr.(type) {
	case *rsa.PrivateKey:
		return jwa.PS256, nil
	case *rsa.PublicKey:
		return jwa.PS256, nil
	case *ecdsa.PrivateKey:
		return ecAlgUsingPublicKey(k.PublicKey)
	case *ecdsa.PublicKey:
		return ecAlgUsingPublicKey(*k)
	case ed25519.PrivateKey:
		return jwa.EdDSA, nil
	case ed25519.PublicKey:
		return jwa.EdDSA, nil
	default:
		return "", fmt.Errorf(`could not determine signature algorithm for key type '%T'`, key)
	}
}

// Thumbprint generates a Nuts compatible thumbprint: Base58(SHA256(rfc7638-json))
func Thumbprint(key jwk.Key) (string, error) {
	pkHash, err := key.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", err
	}
	return base58.Encode(pkHash[:], base58.BitcoinAlphabet), nil
}
