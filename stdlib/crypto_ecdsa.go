
package stdlib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"encoding/asn1"

	"github.com/2dprototype/tender"
)

var cryptoECDSAModule = &tender.ImmutableMap{
	Value: map[string]tender.Object{
		"generate_key": &tender.UserFunction{Name: "generate_key", Value: ecdsaGenerateKey},
		"sign":         &tender.UserFunction{Name: "sign", Value: ecdsaSign},
		"verify":       &tender.UserFunction{Name: "verify", Value: ecdsaVerify},
		"export_key":   &tender.UserFunction{Name: "export_key", Value: ecdsaExportKey},
		"import_key":   &tender.UserFunction{Name: "import_key", Value: ecdsaImportKey},
	},
}

func ecdsaGenerateKey(args ...tender.Object) (ret tender.Object, err error) {
	curveName := "p256"
	if len(args) > 0 {
		curveNameArg, ok := tender.ToString(args[0])
		if !ok {
			return nil, tender.ErrInvalidArgumentType{
				Name:     "curve",
				Expected: "string",
				Found:    args[0].TypeName(),
			}
		}
		curveName = curveNameArg
	}

	var curve elliptic.Curve
	switch curveName {
	case "p224":
		curve = elliptic.P224()
	case "p256":
		curve = elliptic.P256()
	case "p384":
		curve = elliptic.P384()
	case "p521":
		curve = elliptic.P521()
	default:
		// return nil, &tender.Error{Value: &tender.String{Value: "unsupported curve: " + curveName}}
		return nil, nil
	}

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return wrapError(err), nil
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: privateKeyBytes}, nil
}

func ecdsaSign(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	data, _ := tender.ToByteSlice(args[0])
	privateKeyBytes, _ := tender.ToByteSlice(args[1])

	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return wrapError(err), nil
	}

	hashed := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])
	if err != nil {
		return wrapError(err), nil
	}

	// Encode signature as ASN.1
	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: signature}, nil
}

func ecdsaVerify(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 3 {
		return nil, tender.ErrWrongNumArguments
	}

	data, _ := tender.ToByteSlice(args[0])
	signature, _ := tender.ToByteSlice(args[1])
	publicKeyBytes, _ := tender.ToByteSlice(args[2])

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return wrapError(err), nil
	}

	ecdsaPublicKey := publicKey.(*ecdsa.PublicKey)

	// Decode ASN.1 signature
	var sig struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(signature, &sig)
	if err != nil {
		return wrapError(err), nil
	}

	hashed := sha256.Sum256(data)
	valid := ecdsa.Verify(ecdsaPublicKey, hashed[:], sig.R, sig.S)

	if valid {
		return tender.TrueValue,  nil
	} else {
		return tender.FalseValue,  nil
	}
	// return &tender.Bool{Value: valid}, nil
}

func ecdsaExportKey(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	keyBytes, _ := tender.ToByteSlice(args[0])
	keyType, ok := tender.ToString(args[1])
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "keyType",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	var pemBlock *pem.Block
	switch keyType {
	case "private":
		privateKey, err := x509.ParseECPrivateKey(keyBytes)
		if err != nil {
			return wrapError(err), nil
		}
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return wrapError(err), nil
		}
		pemBlock = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
	case "public":
		publicKey, err := x509.ParsePKIXPublicKey(keyBytes)
		if err != nil {
			return wrapError(err), nil
		}
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return wrapError(err), nil
		}
		pemBlock = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
	default:
		// return nil, &tender.Error{Value: &tender.String{Value: "keyType must be 'private' or 'public'"}}
		return nil, nil
	}

	pemData := pem.EncodeToMemory(pemBlock)
	return &tender.Bytes{Value: pemData}, nil
}

func ecdsaImportKey(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	pemData, _ := tender.ToByteSlice(args[0])
	keyType, ok := tender.ToString(args[1])
	if !ok {
		return nil, tender.ErrInvalidArgumentType{
			Name:     "keyType",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil {
		// return nil, &tender.Error{Value: &tender.String{Value: "invalid PEM data"}}
		return nil, nil
	}

	switch keyType {
	case "private":
		privateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
		if err != nil {
			return wrapError(err), nil
		}
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return wrapError(err), nil
		}
		return &tender.Bytes{Value: privateKeyBytes}, nil
	case "public":
		publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
		if err != nil {
			return wrapError(err), nil
		}
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return wrapError(err), nil
		}
		return &tender.Bytes{Value: publicKeyBytes}, nil
	default:
		// return nil, &tender.Error{Value: &tender.String{Value: "keyType must be 'private' or 'public'"}}
		return nil, nil
	}
}

// Helper function for ASN.1 marshaling
func asn1Marshal(value interface{}) ([]byte, error) {
	return asn1.Marshal(value)
}

// Helper function for ASN.1 unmarshaling  
func asn1Unmarshal(b []byte, value interface{}) ([]byte, error) {
	return asn1.Unmarshal(b, value)
}