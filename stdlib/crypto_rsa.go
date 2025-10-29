
package stdlib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/2dprototype/tender"
)

var cryptoRSAModule = &tender.ImmutableMap{
	Value: map[string]tender.Object{
		"generate_key": &tender.UserFunction{Name: "generate_key", Value: rsaGenerateKey},
		"encrypt":      &tender.UserFunction{Name: "encrypt", Value: rsaEncrypt},
		"decrypt":      &tender.UserFunction{Name: "decrypt", Value: rsaDecrypt},
		"sign":         &tender.UserFunction{Name: "sign", Value: rsaSign},
		"verify":       &tender.UserFunction{Name: "verify", Value: rsaVerify},
		"export_key":   &tender.UserFunction{Name: "export_key", Value: rsaExportKey},
		"import_key":   &tender.UserFunction{Name: "import_key", Value: rsaImportKey},
	},
}

func rsaGenerateKey(args ...tender.Object) (ret tender.Object, err error) {
	bits := 2048
	if len(args) > 0 {
		bitsArg, ok := tender.ToInt(args[0])
		if !ok {
			return nil, tender.ErrInvalidArgumentType{
				Name:     "bits",
				Expected: "int",
				Found:    args[0].TypeName(),
			}
		}
		bits = bitsArg
	}

	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: x509.MarshalPKCS1PrivateKey(key)}, nil
}

func rsaEncrypt(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	data, _ := tender.ToByteSlice(args[0])
	publicKeyBytes, _ := tender.ToByteSlice(args[1])

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		// Try parsing as PKIX format
		pubInterface, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return wrapError(err), nil
		}
		publicKey = pubInterface.(*rsa.PublicKey)
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: ciphertext}, nil
}

func rsaDecrypt(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	ciphertext, _ := tender.ToByteSlice(args[0])
	privateKeyBytes, _ := tender.ToByteSlice(args[1])

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return wrapError(err), nil
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: plaintext}, nil
}

func rsaSign(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 2 {
		return nil, tender.ErrWrongNumArguments
	}

	data, _ := tender.ToByteSlice(args[0])
	privateKeyBytes, _ := tender.ToByteSlice(args[1])

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return wrapError(err), nil
	}

	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hashed[:])
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: signature}, nil
}

func rsaVerify(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 3 {
		return nil, tender.ErrWrongNumArguments
	}

	data, _ := tender.ToByteSlice(args[0])
	signature, _ := tender.ToByteSlice(args[1])
	publicKeyBytes, _ := tender.ToByteSlice(args[2])

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		pubInterface, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return wrapError(err), nil
		}
		publicKey = pubInterface.(*rsa.PublicKey)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(publicKey, 0, hashed[:], signature)
	if err != nil {
		return tender.FalseValue, nil
	}

	return tender.TrueValue, nil
}

func rsaExportKey(args ...tender.Object) (ret tender.Object, err error) {
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
		privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			return wrapError(err), nil
		}
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
	case "public":
		publicKey, err := x509.ParsePKCS1PublicKey(keyBytes)
		if err != nil {
			pubInterface, err := x509.ParsePKIXPublicKey(keyBytes)
			if err != nil {
				return wrapError(err), nil
			}
			publicKey = pubInterface.(*rsa.PublicKey)
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
		// return &tender.Error{Value: &tender.String{Value: "keyType must be 'private' or 'public'"}}, nil
		return nil, nil
	}

	pemData := pem.EncodeToMemory(pemBlock)
	return &tender.Bytes{Value: pemData}, nil
}

func rsaImportKey(args ...tender.Object) (ret tender.Object, err error) {
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
		privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return wrapError(err), nil
		}
		return &tender.Bytes{Value: x509.MarshalPKCS1PrivateKey(privateKey)}, nil
	case "public":
		publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
		if err != nil {
			publicKey, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
			if err != nil {
				return wrapError(err), nil
			}
			publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
			if err != nil {
				return wrapError(err), nil
			}
			return &tender.Bytes{Value: publicKeyBytes}, nil
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