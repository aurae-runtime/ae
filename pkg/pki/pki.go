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
	"time"
)

type Certificate struct {
	Certificate string `json:"cert" yaml:"cert"`
	PrivateKey  string `json:"key" yaml:"key"`
}

func (c *Certificate) GetCertificate() (*x509.Certificate, error) {
	crtPem, _ := pem.Decode([]byte(c.Certificate))
	if crtPem == nil || crtPem.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	crt, err := x509.ParseCertificate(crtPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return crt, nil
}

func (c *Certificate) GetCertAsString() string {
	return c.Certificate
}

func (c *Certificate) GetPrivateKey() (*rsa.PrivateKey, error) {
	pkPem, _ := pem.Decode([]byte(c.PrivateKey))
	if pkPem == nil || pkPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode private key")
	}

	pk, err := x509.ParsePKCS1PrivateKey(pkPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return pk, nil
}

func (c *Certificate) GetPrivateKeyAsString() string {
	return c.PrivateKey
}

func (c *Certificate) WriteCertificateToFile(path, filename string) error {
	err := createFile(path, filename, c.GetCertAsString())
	if err != nil {
		return err
	}
	return nil
}

func (c *Certificate) WritePrivateKeyToFile(path, filename string) error {
	err := createFile(path, filename, c.GetPrivateKeyAsString())
	if err != nil {
		return err
	}
	return nil
}

type CertificateRequest struct {
	CSR        string `json:"csr" yaml:"csr"`
	PrivateKey string `json:"key" yaml:"key"`
	User       string `json:"user" yaml:"user"`
}

func (c *CertificateRequest) GetCsr() (*x509.CertificateRequest, error) {
	csrPem, _ := pem.Decode([]byte(c.CSR))
	if csrPem == nil || csrPem.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("failed to decode certificate request")
	}

	csr, err := x509.ParseCertificateRequest(csrPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate request: %w", err)
	}

	return csr, nil
}

func (c *CertificateRequest) GetCsrAsString() string {
	return c.CSR
}

func (c *CertificateRequest) GetPrivateKey() (*rsa.PrivateKey, error) {
	pkPem, _ := pem.Decode([]byte(c.PrivateKey))
	if pkPem == nil || pkPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode private key")
	}

	pk, err := x509.ParsePKCS1PrivateKey(pkPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return pk, nil
}

func (c *CertificateRequest) GetPrivateKeyAsString() string {
	return c.PrivateKey
}

func (c *CertificateRequest) WriteCsrToFile(path, filename string) error {
	err := createFile(path, filename, c.GetCsrAsString())
	if err != nil {
		return err
	}
	return nil
}

func (c *CertificateRequest) WritePrivateKeyToFile(path, filename string) error {
	err := createFile(path, filename, c.GetPrivateKeyAsString())
	if err != nil {
		return err
	}
	return nil
}

func HandleCreateAuraeRootCA(path string, domainName string) (*Certificate, error) {
	crtPem, keyPem, err := createCA(domainName)
	if err != nil {
		return nil, err
	}

	ca := &Certificate{
		Certificate: string(crtPem),
		PrivateKey:  string(keyPem),
	}

	if path != "" {
		err = ca.WriteCertificateToFile(path, "ca.crt")
		if err != nil {
			return ca, err
		}
		err = ca.WritePrivateKeyToFile(path, "ca.key")
		if err != nil {
			return ca, err
		}
	}

	return ca, nil
}

func createCA(domainName string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
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
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	crtPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crtBytes,
	})

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: crtBytes,
	})

	return crtPem, keyPem, nil
}

func HandleCreateClientCSR(path, domain, user string) (*CertificateRequest, error) {
	csrPem, keyPem, err := createClientCSR(domain, user)
	if err != nil {
		return &CertificateRequest{}, err
	}

	csr := &CertificateRequest{
		CSR:        string(csrPem),
		PrivateKey: string(keyPem),
		User:       user,
	}

	if path != "" {
		err = csr.WriteCsrToFile(path, fmt.Sprintf("client.%s.csr", csr.User))
		if err != nil {
			return csr, err
		}
		err = csr.WritePrivateKeyToFile(path, fmt.Sprintf("client.%s.key", csr.User))
		if err != nil {
			return csr, err
		}
	}

	return csr, nil
}

func createClientCSR(domain, user string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("failed to generate private key: %w", err)
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
		return []byte{}, []byte{}, fmt.Errorf("could not create certificate request: %w", err)
	}

	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	return csrPem, keyPem, nil
}
