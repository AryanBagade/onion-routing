package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

type OnionLayer struct {
	NodeID    string
	PublicKey *rsa.PublicKey
	AESKey    []byte
	GCM       cipher.AEAD
}

type OnionPacket struct {
	Layers [][]byte
	Data   []byte
}

func CreateOnionLayers(nodeKeys []*rsa.PublicKey, nodeIDs []string) ([]OnionLayer, error) {
	if len(nodeKeys) != len(nodeIDs) {
		return nil, errors.New("mismatched node keys and IDs")
	}

	layers := make([]OnionLayer, len(nodeKeys))
	
	for i, pubKey := range nodeKeys {
		aesKey := make([]byte, 32)
		if _, err := rand.Read(aesKey); err != nil {
			return nil, err
		}

		block, err := aes.NewCipher(aesKey)
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		layers[i] = OnionLayer{
			NodeID:    nodeIDs[i],
			PublicKey: pubKey,
			AESKey:    aesKey,
			GCM:       gcm,
		}
	}

	return layers, nil
}

func EncryptOnion(data []byte, layers []OnionLayer) (*OnionPacket, error) {
	payload := data
	
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		
		nonce := make([]byte, layer.GCM.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}

		encrypted := layer.GCM.Seal(nonce, nonce, payload, nil)
		
		encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, layer.PublicKey, layer.AESKey, nil)
		if err != nil {
			return nil, err
		}

		payload = append(encryptedKey, encrypted...)
	}

	return &OnionPacket{Data: payload}, nil
}

func DecryptOnionLayer(packet []byte, privateKey *rsa.PrivateKey) ([]byte, []byte, error) {
	keySize := privateKey.Size()
	if len(packet) < keySize {
		return nil, nil, errors.New("packet too small")
	}

	encryptedKey := packet[:keySize]
	encryptedData := packet[keySize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKey, nil)
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, nil, errors.New("encrypted data too small")
	}

	nonce := encryptedData[:nonceSize]
	ciphertext := encryptedData[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, nil, err
	}

	return plaintext, aesKey, nil
}