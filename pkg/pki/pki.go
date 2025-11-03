/* -------------------------------------------------------------------------- *\
 *             Apache 2.0 License Copyright © 2022 The Aurae Authors          *
 *                                                                            *
 *                +--------------------------------------------+              *
 *                |   █████╗ ██╗   ██╗██████╗  █████╗ ███████╗ |              *
 *                |  ██╔══██╗██║   ██║██╔══██╗██╔══██╗██╔════╝ |              *
 *                |  ███████║██║   ██║██████╔╝███████║█████╗   |              *
 *                |  ██╔══██║██║   ██║██╔══██╗██╔══██║██╔══╝   |              *
 *                |  ██║  ██║╚██████╔╝██║  ██║██║  ██║███████╗ |              *
 *                |  ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝ |              *
 *                +--------------------------------------------+              *
 *                                                                            *
 *                         Distributed Systems Runtime                        *
 *                                                                            *
 * -------------------------------------------------------------------------- *
 *                                                                            *
 *   Licensed under the Apache License, Version 2.0 (the "License");          *
 *   you may not use this file except in compliance with the License.         *
 *   You may obtain a copy of the License at                                  *
 *                                                                            *
 *       http://www.apache.org/licenses/LICENSE-2.0                           *
 *                                                                            *
 *   Unless required by applicable law or agreed to in writing, software      *
 *   distributed under the License is distributed on an "AS IS" BASIS,        *
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 *   See the License for the specific language governing permissions and      *
 *   limitations under the License.                                           *
 *                                                                            *
\* -------------------------------------------------------------------------- */

package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

type Certificate struct {
	Certificate string `json:"cert" yaml:"cert"`
	PrivateKey  string `json:"key" yaml:"key"`
}

type CertificateRequest struct {
	CSR        string `json:"csr" yaml:"csr"`
	PrivateKey string `json:"key" yaml:"key"`
	User       string `json:"user" yaml:"user"`
}

func CreateAuraeRootCA(path string, domainName string) (*Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return &Certificate{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	subj := pkix.Name{
		Organization:       []string{"Aurae"},
		OrganizationalUnit: []string{"Runtime"},
		Province:           []string{"aurae"},
		Locality:           []string{"aurae"},
		Country:            []string{"IS"},
		CommonName:         domainName,
	}

	now := time.Now()

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return &Certificate{}, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		Subject:               subj,
		NotBefore:             now,
		NotAfter:              now.Add(24 * time.Hour * 9999),
		IsCA:                  true,
		BasicConstraintsValid: true,
		DNSNames:              []string{domainName},
		SerialNumber:          serialNumber,
	}

	// To get an AuthorityKeyId which is equal to SubjectKeyId, we are manually
	// setting it according to the rules of x509.go.
	//
	// From X509 docs:
	// > The AuthorityKeyId will be taken from the SubjectKeyId of parent, if any,
	// > unless the resulting certificate is self-signed. Otherwise the value from
	// > template will be used.
	//
	// > If SubjectKeyId from template is empty and the template is a CA, SubjectKeyId
	// > will be generated from the hash of the public key.
	//
	// We need a hash of the publickey, so hopefully this link is right
	// https://stackoverflow.com/questions/52502511/how-to-generate-bytes-array-from-publickey#comment92269419_52502639
	pubHash := sha1.Sum(priv.N.Bytes())
	template.SubjectKeyId = pubHash[:]
	template.AuthorityKeyId = template.SubjectKeyId

	crtBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return &Certificate{}, fmt.Errorf("failed to create certificate: %w", err)
	}

	crtPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crtBytes,
	})

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	ca := &Certificate{
		Certificate: string(crtPem),
		PrivateKey:  string(keyPem),
	}

	if path != "" {
		err = createCAFiles(path, ca)
		if err != nil {
			return ca, err
		}
	}

	return ca, nil
}

func CreateClientCSR(path, domain, user string) (*CertificateRequest, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return &CertificateRequest{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	subj := pkix.Name{
		Organization:       []string{"Aurae"},
		OrganizationalUnit: []string{"Runtime"},
		Province:           []string{"aurae"},
		Locality:           []string{"aurae"},
		Country:            []string{"IS"},
		CommonName:         fmt.Sprintf("%s.%s", user, domain),
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           []string{fmt.Sprintf("%s.%s", user, domain)},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return &CertificateRequest{}, fmt.Errorf("could not create certificate request: %w", err)
	}

	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	csr := &CertificateRequest{
		CSR:        string(csrPem),
		PrivateKey: string(keyPem),
		User:       user,
	}

	if path != "" {
		err = createCsrFiles(path, csr)
		if err != nil {
			return csr, err
		}
	}

	return csr, nil
}

func createCAFiles(path string, ca *Certificate) error {
	path = filepath.Clean(path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	crtPath := filepath.Join(path, "ca.crt")
	keyPath := filepath.Join(path, "ca.key")

	err = writeStringToFile(crtPath, ca.Certificate)
	if err != nil {
		return err
	}

	err = writeStringToFile(keyPath, ca.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func createCsrFiles(path string, ca *CertificateRequest) error {
	path = filepath.Clean(path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	csrPath := filepath.Join(path, fmt.Sprintf("client.%s.csr", ca.User))
	keyPath := filepath.Join(path, fmt.Sprintf("client.%s.key", ca.User))

	err = writeStringToFile(csrPath, ca.CSR)
	if err != nil {
		return err
	}

	err = writeStringToFile(keyPath, ca.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func writeStringToFile(p string, s string) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", p, err)
	}
	defer f.Close()

	_, err = f.WriteString(s)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", p, err)
	}

	return nil
}
