package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

/*
 * VerifySignature: Verify Signature
 */
func VerifySignature(data, signature, publicKey string, algorithm string, hash string) (bool, error) {
	if algorithm != "RSA" && hash != "SHA256" {
		return false, errors.New("invalid algorithm or hash")
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, errors.New("invalid Base64 public key")
	}
	pubKey, err := parsePublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.New("invalid Base64 signature")
	}

	hashBytes := sha256.Sum256([]byte(data))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashBytes[:], signatureBytes)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func parsePublicKeyFromBytes(keyBytes []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func SignMessage(message string, privateKey string, algorithm string, hash string) (string, error) {
	if algorithm != "RSA" && hash != "SHA256" {
		return "", errors.New("invalid algorithm or hash")
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", errors.New("failed to decode private key from Base64")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func GenerateRSAKeys(bits int) (string, string, error) {

	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKeyBytes)

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKeyBytes)

	return privKeyBase64, pubKeyBase64, nil
}
