package proof

import (
	"encoding/json"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-node/crypto"
	"github.com/nuts-foundation/nuts-node/vcr/signature"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLDProof_Verify(t *testing.T) {

}

func TestLDProofVerifier_Verify(t *testing.T) {
	t.Run("ok - JSONWebSignature2020 test vector", func(t *testing.T) {
		vc_0 := `{
			"@context": [
				 "https://www.w3.org/2018/credentials/v1",
				 "https://www.w3.org/2018/credentials/examples/v1",
				 "https://w3c-ccg.github.io/lds-jws2020/contexts/lds-jws2020-v1.json"
			],
			"id": "http://example.gov/credentials/3732",
			"type": ["VerifiableCredential", "UniversityDegreeCredential"],
			"issuer": { "id": "https://example.com/issuer/123" },
			"issuanceDate": "2020-03-10T04:24:12.164Z",
			"credentialSubject": {
				 "id": "did:example:456",
				 "degree": {
					  "type": "BachelorDegree",
					  "name": "Bachelor of Science and Arts"
				 }
			},
			"proof": {
				 "type": "JSONWebSignature2020",
				 "created": "2019-12-11T03:50:55Z",
				 "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..MJ5GwWRMsadCyLNXU_flgJtsS32584MydBxBuygps_cM0sbU3abTEOMyUvmLNcKOwOBE1MfDoB1_YY425W3sAg",
				 "proofPurpose": "assertionMethod",
				 "verificationMethod": "https://example.com/issuer/123#ovsDKYBjFemIy8DVhc-w2LSi8CvXMw2AYDzHj04yxkc"
			}
		 }`

		rawVerificationMethod := `{
			"id": "did:key:abc#ovsDKYBjFemIy8DVhc-w2LSi8CvXMw2AYDzHj04yxkc",
			"type": "JsonWebKey2020",
			"controller": "did:key:z6Mkf5rGMoatrSj1f4CyvuHBeXJELe9RPdzo2PKGNCKVtZxP",
			"publicKeyJwk": {
				 "kty": "OKP",
				 "crv": "Ed25519",
				 "x": "CV-aGlld3nVdgnhoZK0D36Wk-9aIMlZjZOK2XhPMnkQ"
			}
		}`
		signedDocument := SignedDocument{}
		if !assert.NoError(t, json.Unmarshal([]byte(vc_0), &signedDocument)) {
			return
		}

		verificationMethod := did.VerificationMethod{}
		if !assert.NoError(t, json.Unmarshal([]byte(rawVerificationMethod), &verificationMethod)) {
			return
		}
		pk, err := verificationMethod.PublicKey()
		if !assert.NoError(t, err) {
			return
		}

		ldProof, err := NewLDProofFromDocumentProof(signedDocument.FirstProof())
		assert.NoError(t, err)
		err = ldProof.Verify(signedDocument.DocumentWithoutProof(), signature.JSONWebSignature2020{}, pk)
		assert.NoError(t, err, "expected no error when verifying the JSONWebSignature2020 test vector")
	})
}

func TestLDProof_Sign(t *testing.T) {
	t.Run("sign a document", func(t *testing.T) {
		now := time.Now()
		expires := now.Add(20 * time.Hour)
		challenge := "stand on 1 leg for 2 hours"
		domain := "chateau Torquilstone"

		pOptions := ProofOptions{
			Created:        now,
			Domain:         &domain,
			Challenge:      &challenge,
			ExpirationDate: &expires,
			ProofPurpose:   "assertion",
		}

		ldProof := NewLDProof(pOptions)

		document := map[string]interface{}{
			"@context": []interface{}{
				map[string]interface{}{"title": "https://schema.org#title"},
			},
			"title": "Hello world!",
		}

		kid := "did:nuts:123#abc"
		testKey := crypto.NewTestKey(kid)

		result, err := ldProof.Sign(document, signature.JSONWebSignature2020{}, testKey)
		if !assert.NoError(t, err) || !assert.NotNil(t, result) {
			return
		}
		signedDocument := result.(SignedDocument)
		docProof := signedDocument.FirstProof()
		assert.Equal(t, challenge, docProof["challenge"])
		assert.Equal(t, domain, docProof["domain"])
		t.Logf("%+v", signedDocument)
	})
}
