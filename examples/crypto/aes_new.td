// tender@v0.0.8
import "crypto"

// Legacy/Default AES Encryption
aes_key := "0123456789abcdef0123456789abcdef" // 32 bytes
encrypted := crypto.aes.encrypt("secret message", aes_key)
decrypted := crypto.aes.decrypt(encrypted, aes_key)
println("Default Encrypted:", encrypted)
println("Default Decrypted:", decrypted)

// Create block cipher instance for advanced modes
block := crypto.aes.new_cipher(aes_key)
println("AES Block Size:", block.block_size())

// Advanced AES-GCM (Authenticated Encryption with Associated Data)
gcm := block.new_gcm()
nonce := crypto.random.bytes(gcm.nonce_size())
plaintext := "authenticated secret"
additional_data := "metadata-context"

gcm_ciphertext := gcm.seal(nonce, plaintext, additional_data)
gcm_decrypted := gcm.open(nonce, gcm_ciphertext, additional_data)
println("GCM Decrypted:", string(gcm_decrypted))

// Advanced AES-CBC (Cipher Block Chaining)
iv_cbc := crypto.random.bytes(block.block_size())
cbc_encrypter := block.new_cbc_encrypter(iv_cbc)

// Note: CBC requires manual padding to match block multiples; raw block example:
raw_block_data := "1234567890abcdef" // Must be exact multiple of block size (16 bytes)
cbc_encrypted := cbc_encrypter.crypt_blocks(raw_block_data)

cbc_decrypter := block.new_cbc_decrypter(iv_cbc)
cbc_decrypted := cbc_decrypter.crypt_blocks(cbc_encrypted)
println("CBC Decrypted:", string(cbc_decrypted))

// Advanced AES-CTR (Counter Stream Mode)
iv_ctr := crypto.random.bytes(block.block_size())
ctr_stream := block.new_ctr(iv_ctr)

stream_ciphertext := ctr_stream.xor_key_stream("stream me up!")
// CTR mode is symmetric: running XOR again decrypts the stream
ctr_stream_decrypt := block.new_ctr(iv_ctr) // Re-initialize stream with same IV/nonce state
stream_decrypted := ctr_stream_decrypt.xor_key_stream(stream_ciphertext)
println("CTR Decrypted:", string(stream_decrypted))