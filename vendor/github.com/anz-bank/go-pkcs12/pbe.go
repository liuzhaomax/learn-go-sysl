package pkcs12

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
)

const DefaultPBKDFIterations = 1000

type Cipher interface {
	Encrypt(decrypted []byte) []byte
	Decrypt(encrypted []byte) ([]byte, error)
}

func sha256Sum(in []byte) []byte {
	sum := sha256.Sum256(in)
	return sum[:]
}

type sha256With128AESCBC struct {
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
	blockSize int
}

func create(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

func deriveKey(salt, password []byte, iterations int) []byte {
	return pbkdf(sha256Sum, 32, 64, salt, password, iterations, 1, 16)
}

func deriveIV(salt, password []byte, iterations int) []byte {
	return pbkdf(sha256Sum, 32, 64, salt, password, iterations, 2, 16)
}

func newSHA256With128AESCBCCipher(password, salt []byte, iter int) (*sha256With128AESCBC, error) {
	encodedPass, err := bmpString(string(password))
	if err != nil {
		return nil, err
	}

	key := deriveKey(salt, encodedPass, iter)
	iv := deriveIV(salt, encodedPass, iter)

	block, err := create(key)
	if err != nil {
		return nil, err
	}

	pbe := &sha256With128AESCBC{
		encrypter: cipher.NewCBCEncrypter(block, iv),
		decrypter: cipher.NewCBCDecrypter(block, iv),
		blockSize: block.BlockSize(),
	}

	return pbe, nil
}

func PBEDecrypt(encrypted, password, salt []byte, iter int) ([]byte, error) {
	c, err := newSHA256With128AESCBCCipher(password, salt, iter)
	if err != nil {
		return nil, err
	}
	if len(encrypted)%c.blockSize != 0 {
		return nil, errors.New("pkcs12: input is not a multiple of the block size")
	}
	decrypted := make([]byte, len(encrypted))
	c.decrypter.CryptBlocks(decrypted, encrypted)

	psLen := int(decrypted[len(decrypted)-1])
	if psLen == 0 || psLen > c.blockSize {
		return nil, ErrDecryption
	}

	if len(decrypted) < psLen {
		return nil, ErrDecryption
	}
	ps := decrypted[len(decrypted)-psLen:]
	decrypted = decrypted[:len(decrypted)-psLen]
	if !bytes.Equal(ps, bytes.Repeat([]byte{byte(psLen)}, psLen)) {
		return nil, ErrDecryption
	}

	return decrypted, nil
}

func PBEEncrypt(decrypted, password, salt []byte, iter int) ([]byte, error) {
	c, err := newSHA256With128AESCBCCipher(password, salt, iter)
	if err != nil {
		return nil, err
	}

	psLen := c.blockSize - len(decrypted)%c.blockSize
	encrypted := make([]byte, len(decrypted)+psLen)
	copy(encrypted[:len(decrypted)], decrypted)
	copy(encrypted[len(decrypted):], bytes.Repeat([]byte{byte(psLen)}, psLen))
	c.encrypter.CryptBlocks(encrypted, encrypted)

	return encrypted, nil
}
