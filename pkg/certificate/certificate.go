package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GenerateCA(caPath string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Orbital OSS",
			Organization: []string{"Orbital OSS"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour), // Valid 10 years
		KeyUsage:              x509.KeyUsageCRLSign | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, privateKey.Public(), privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	certFile := filepath.Join(caPath, "ca.crt")
	keyFile := filepath.Join(caPath, "ca.key")

	if err = savePEMFile(certFile, "CERTIFICATE", certBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to save CA certificate: %w", err)
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err = savePEMFile(keyFile, "EC PRIVATE KEY", privateKeyBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to save CA private key: %w", err)
	}

	return caCert, privateKey, nil
}

func LoadCA(caPath string) (*x509.Certificate, *ecdsa.PrivateKey, error) {

	certFile := filepath.Join(caPath, "ca.crt")
	keyFile := filepath.Join(caPath, "ca.key")

	certBytes, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA certificate file: %w", err)
	}

	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("failed to decode CA certificate")
	}

	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA private key file: %w", err)
	}

	// Decode the CA private key
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil || keyBlock.Type != "EC PRIVATE KEY" {
		return nil, nil, fmt.Errorf("failed to decode CA private key")
	}

	caKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA private key: %w", err)
	}

	return caCert, caKey, nil
}

func GenerateServerCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, serverCertPath, ip string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate server private key: %w", err)
	}

	serverCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: ip,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(1 * 365 * 24 * time.Hour), // Available for 1 year
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP(ip)},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, serverCert, caCert, privateKey.Public(), caKey)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %w", err)
	}

	certFile := fmt.Sprintf("%s/server.crt", serverCertPath)
	keyFile := fmt.Sprintf("%s/server.key", serverCertPath)

	if err = savePEMFile(certFile, "CERTIFICATE", certBytes); err != nil {
		return fmt.Errorf("failed to save server certificate: %w", err)
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err = savePEMFile(keyFile, "EC PRIVATE KEY", privateKeyBytes); err != nil {
		return fmt.Errorf("failed to save server private key: %w", err)
	}

	return nil
}

func savePEMFile(filePath, blockType string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer file.Close()

	block := &pem.Block{
		Type:  blockType,
		Bytes: data,
	}

	return pem.Encode(file, block)
}
