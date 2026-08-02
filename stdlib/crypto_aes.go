package stdlib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/2dprototype/tender"
)

var cryptoAESModule = &tender.ImmutableMap{
	Value: map[string]tender.Object{
		"encrypt":    &tender.NativeFunction{Name: "encrypt", Value: aesEncrypt},
		"decrypt":    &tender.NativeFunction{Name: "decrypt", Value: aesDecrypt},
		"block_size": &tender.Int{Value: aes.BlockSize},
		// New addition exposing full AES potential
		"new_cipher": &tender.NativeFunction{Name: "new_cipher", Value: aesNewCipher},
		// EVP and PKCS7 utility functions
		"evp_bytes_to_key": &tender.NativeFunction{Name: "evp_bytes_to_key", Value: aesEVPBytesToKey},
		"pkcs7_pad":        &tender.NativeFunction{Name: "pkcs7_pad", Value: aesPKCS7Pad},
		"pkcs7_unpad":      &tender.NativeFunction{Name: "pkcs7_unpad", Value: aesPKCS7Unpad},
	},
}

// -----------------------------------------------------------------------------
// Default Encrypt/Decrypt (Maintained for Compatibility)
// -----------------------------------------------------------------------------

func aesEncrypt(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}

	plaintext, _ := tender.ToByteSlice(args[0])
	key, _ := tender.ToByteSlice(args[1])

	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.String{Value: ciphertext}, nil
}

func aesDecrypt(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}

	ciphertext, _ := tender.ToString(args[0])
	key, _ := tender.ToByteSlice(args[1])

	plaintext, err := DecryptAES(ciphertext, key)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.String{Value: plaintext}, nil
}


// aesEVPBytesToKey wraps EVPBytesToKey for Tender
func aesEVPBytesToKey(args ...tender.Object) (tender.Object, error) {
	if len(args) != 4 {
		return nil, tender.ErrWrongNumArguments
	}
	password, _ := tender.ToByteSlice(args[0])
	salt, _ := tender.ToByteSlice(args[1])
	keyLen, _ := tender.ToInt(args[2])
	ivLen, _ := tender.ToInt(args[3])
	key, iv := EVPBytesToKey(password, salt, keyLen, ivLen)
	// Return as a map with "key" and "iv" fields
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"key": &tender.Bytes{Value: key},
			"iv":  &tender.Bytes{Value: iv},
		},
	}, nil
}

// aesPKCS7Pad wraps pkcs7Pad for Tender
func aesPKCS7Pad(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	data, _ := tender.ToByteSlice(args[0])
	blockSize, _ := tender.ToInt(args[1])
	padded := pkcs7Pad(data, blockSize)
	return &tender.Bytes{Value: padded}, nil
}

// aesPKCS7Unpad wraps pkcs7Unpad for Tender
func aesPKCS7Unpad(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	data, _ := tender.ToByteSlice(args[0])
	unpadded, err := pkcs7Unpad(data)
	if err != nil {
		return wrapError(err), nil
	}
	return &tender.Bytes{Value: unpadded}, nil
}


func EVPBytesToKey(password, salt []byte, keyLen, ivLen int) (key, iv []byte) {
	derived := make([]byte, 0, keyLen+ivLen)
	hash := make([]byte, 0)
	for len(derived) < keyLen+ivLen {
		// MD5(hash + password + salt)
		h := md5.New()
		h.Write(hash)
		h.Write(password)
		h.Write(salt)
		hash = h.Sum(nil)
		derived = append(derived, hash...)
	}
	return derived[:keyLen], derived[keyLen : keyLen+ivLen]
}

// pkcs7Pad pads data to a multiple of blockSize (16 for AES)
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// pkcs7Unpad removes PKCS#7 padding
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}

// EncryptAES encrypts a plaintext string with a passphrase.
// Returns a Base64 string compatible with CryptoJS.AES.encrypt().
func EncryptAES(plaintext, passphrase []byte) (string, error) {
	// Generate random 8-byte salt
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// Derive 32-byte key and 16-byte IV using EVP_BytesToKey
	key, iv := EVPBytesToKey(passphrase, salt, 32, 16)

	// Pad plaintext (PKCS#7)
	paddedPlain := pkcs7Pad(plaintext, aes.BlockSize)

	// AES CBC encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(paddedPlain))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlain)

	// Build final data: "Salted__" + salt + ciphertext
	final := make([]byte, 0, 8+8+len(ciphertext))
	final = append(final, []byte("Salted__")...)
	final = append(final, salt...)
	final = append(final, ciphertext...)

	// Base64 encode
	return base64.StdEncoding.EncodeToString(final), nil
}

// DecryptAES decrypts a CryptoJS-compatible Base64 ciphertext.
// Returns the original plaintext string.
func DecryptAES(ciphertextBase64 string, passphrase []byte) (string, error) {
	// Decode Base64
	final, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// Validate and extract salt (must start with "Salted__")
	if len(final) < 16 || string(final[:8]) != "Salted__" {
		return "", errors.New("invalid ciphertext: missing Salted__ header")
	}
	salt := final[8:16]
	ciphertext := final[16:]

	// Derive key and IV
	key, iv := EVPBytesToKey(passphrase, salt, 32, 16)

	// AES CBC decryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of block size")
	}
	plainPadded := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainPadded, ciphertext)

	// Remove PKCS#7 padding
	plaintext, err := pkcs7Unpad(plainPadded)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// -----------------------------------------------------------------------------
// Advanced AES API 
// -----------------------------------------------------------------------------

func aesNewCipher(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}

	key, _ := tender.ToByteSlice(args[0])
	block, err := aes.NewCipher(key)
	if err != nil {
		return wrapError(err), nil
	}

	return makeBlockWrapper(block), nil
}

func makeBlockWrapper(block cipher.Block) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"block_size": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					return &tender.Int{Value: int64(block.BlockSize())}, nil
				},
			},
			"new_gcm": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					gcm, err := cipher.NewGCM(block)
					if err != nil {
						return wrapError(err), nil
					}
					return makeAEADWrapper(gcm), nil
				},
			},
			"new_cbc_encrypter": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					iv, _ := tender.ToByteSlice(args[0])
					if len(iv) != block.BlockSize() {
						return wrapError(errors.New("IV length must equal block size")), nil
					}
					mode := cipher.NewCBCEncrypter(block, iv)
					return makeBlockModeWrapper(mode), nil
				},
			},
			"new_cbc_decrypter": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					iv, _ := tender.ToByteSlice(args[0])
					if len(iv) != block.BlockSize() {
						return wrapError(errors.New("IV length must equal block size")), nil
					}
					mode := cipher.NewCBCDecrypter(block, iv)
					return makeBlockModeWrapper(mode), nil
				},
			},
			"new_ctr": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					iv, _ := tender.ToByteSlice(args[0])
					if len(iv) != block.BlockSize() {
						return wrapError(errors.New("IV length must equal block size")), nil
					}
					stream := cipher.NewCTR(block, iv)
					return makeStreamWrapper(stream), nil
				},
			},
		},
	}
}

// Wrapper for Authenticated Encryption with Associated Data (GCM)
func makeAEADWrapper(aead cipher.AEAD) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"nonce_size": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					return &tender.Int{Value: int64(aead.NonceSize())}, nil
				},
			},
			"seal": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) < 2 || len(args) > 3 {
						return nil, tender.ErrWrongNumArguments
					}
					nonce, _ := tender.ToByteSlice(args[0])
					plaintext, _ := tender.ToByteSlice(args[1])
					var additionalData []byte
					if len(args) == 3 {
						additionalData, _ = tender.ToByteSlice(args[2])
					}
					
					ciphertext := aead.Seal(nil, nonce, plaintext, additionalData)
					return &tender.Bytes{Value: ciphertext}, nil
				},
			},
			"open": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) < 2 || len(args) > 3 {
						return nil, tender.ErrWrongNumArguments
					}
					nonce, _ := tender.ToByteSlice(args[0])
					ciphertext, _ := tender.ToByteSlice(args[1])
					var additionalData []byte
					if len(args) == 3 {
						additionalData, _ = tender.ToByteSlice(args[2])
					}
					
					plaintext, err := aead.Open(nil, nonce, ciphertext, additionalData)
					if err != nil {
						return wrapError(err), nil
					}
					return &tender.Bytes{Value: plaintext}, nil
				},
			},
		},
	}
}

// Wrapper for Block Modes (CBC)
func makeBlockModeWrapper(mode cipher.BlockMode) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"crypt_blocks": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					src, _ := tender.ToByteSlice(args[0])
					if len(src)%mode.BlockSize() != 0 {
						return wrapError(errors.New("input not full blocks")), nil
					}
					dst := make([]byte, len(src))
					mode.CryptBlocks(dst, src)
					return &tender.Bytes{Value: dst}, nil
				},
			},
		},
	}
}

// Wrapper for Stream Modes (CTR, CFB, OFB)
func makeStreamWrapper(stream cipher.Stream) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"xor_key_stream": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					src, _ := tender.ToByteSlice(args[0])
					dst := make([]byte, len(src))
					stream.XORKeyStream(dst, src)
					return &tender.Bytes{Value: dst}, nil
				},
			},
		},
	}
}