package stdlib

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512" 
    "crypto/sha1"
    "golang.org/x/crypto/sha3"
	"encoding/hex"
	"github.com/2dprototype/tender"
)

var cryptoModule = map[string]tender.Object{
    "md5":          &tender.UserFunction{Name: "md5", Value: cryptoMd5Hash},
    "sha1":         &tender.UserFunction{Name: "sha1", Value: cryptoSha1Hash},
    "sha256":       &tender.UserFunction{Name: "sha256", Value: cryptoSha256Hash},
    "sha512":       &tender.UserFunction{Name: "sha512", Value: cryptoSha512Hash},
    "sha3_256":     &tender.UserFunction{Name: "sha3_256", Value: cryptoSha3_256Hash},
	"aes":          cryptoAESModule,
}

func cryptoMd5Hash(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}

	input, _ := tender.ToByteSlice(args[0])
	data := input
	hash := md5.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	return &tender.String{Value: hashString}, nil
}

func cryptoSha256Hash(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}

	input, _ := tender.ToByteSlice(args[0])
	data := input
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	return &tender.String{Value: hashString}, nil
}


func cryptoSha1Hash(args ...tender.Object) (ret tender.Object, err error) {
    if len(args) != 1 {
        return nil, tender.ErrWrongNumArguments
    }

    input, _ := tender.ToByteSlice(args[0])
    data := input
    hash := sha1.Sum(data)
    hashString := hex.EncodeToString(hash[:])
    return &tender.String{Value: hashString}, nil
}

func cryptoSha512Hash(args ...tender.Object) (ret tender.Object, err error) {
    if len(args) != 1 {
        return nil, tender.ErrWrongNumArguments
    }

    input, _ := tender.ToByteSlice(args[0])
    data := input
    hash := sha512.Sum512(data)
    hashString := hex.EncodeToString(hash[:])
    return &tender.String{Value: hashString}, nil
}

func cryptoSha3_256Hash(args ...tender.Object) (ret tender.Object, err error) {
    if len(args) != 1 {
        return nil, tender.ErrWrongNumArguments
    }

    input, _ := tender.ToByteSlice(args[0])
    data := input
    hash := sha3.Sum256(data)
    hashString := hex.EncodeToString(hash[:])
    return &tender.String{Value: hashString}, nil
}
