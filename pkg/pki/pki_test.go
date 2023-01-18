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
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAuraeRootCA(t *testing.T) {
	t.Run("createAuraeCA", func(t *testing.T) {
		domainName := "unsafe.aurae.io"

		auraeCa, err := CreateAuraeRootCA("", "unsafe.aurae.io")
		if err != nil {
			t.Errorf("could not create auraeCA")
		}

		cert, _ := pem.Decode([]byte(auraeCa.Certificate))
		crt, err := x509.ParseCertificate(cert.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate")
		}
		if crt.Subject.CommonName != domainName {
			t.Errorf("certificate does not contain common name")
		}
	})

	t.Run("createAuraeCA with local files", func(t *testing.T) {
		path := "_tmp/pki"
		domainName := "unsafe.aurae.io"

		_, err := CreateAuraeRootCA(path, domainName)
		if err != nil {
			t.Errorf("could not create auraeCA")
		}

		auraeCaFile, err := os.ReadFile(filepath.Join(path, "ca.crt"))
		if err != nil {
			t.Errorf("could not read ca file")
		}

		// load ca.cert
		cert, _ := pem.Decode(auraeCaFile)

		crt, err := x509.ParseCertificate(cert.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate")
		}
		if crt.Subject.CommonName != domainName {
			t.Errorf("certificate does not contain common name")
		}

		// cleanup files
		err = os.Remove(filepath.Join(path, "ca.crt"))
		if err != nil {
			t.Errorf("could not delete %s", filepath.Join(path, "ca.crt"))
		}
		err = os.Remove(filepath.Join(path, "ca.key"))
		if err != nil {
			t.Errorf("could not delete %s", filepath.Join(path, "ca.key"))
		}
	})
}

func TestCreateCSR(t *testing.T) {
	t.Run("createCSR", func(t *testing.T) {
		clientCsr, err := CreateClientCSR("", "unsafe.aurae.io", "christoph")
		if err != nil {
			t.Errorf("could not create csr")
		}

		csrBytes := []byte(clientCsr.CSR)
		csrPem, _ := pem.Decode(csrBytes)
		_, err = x509.ParseCertificateRequest(csrPem.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate request")
		}

		keyBytes := []byte(clientCsr.PrivateKey)
		keyPem, _ := pem.Decode(keyBytes)
		_, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes)
		if err != nil {
			t.Errorf("could not parse private key")
		}
	})

	t.Run("createCSR with local files", func(t *testing.T) {
		path := "_tmp/pki"
		clientCsr, err := CreateClientCSR(path, "unsafe.aurae.io", "christoph")
		if err != nil {
			t.Errorf("could not create csr")
		}

		// read and load csr file
		csrFilePath := filepath.Join(path, fmt.Sprintf("client.%s.csr", clientCsr.User))
		csrBytes, err := os.ReadFile(csrFilePath)
		if err != nil {
			t.Errorf("could not read ca file")
		}

		csrPem, _ := pem.Decode(csrBytes)
		_, err = x509.ParseCertificateRequest(csrPem.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate request")
		}

		// read and load key file
		keyFilePath := filepath.Join(path, fmt.Sprintf("client.%s.key", clientCsr.User))
		keyBytes, err := os.ReadFile(keyFilePath)
		if err != nil {
			t.Errorf("could not read key file")
		}

		keyPem, _ := pem.Decode(keyBytes)
		_, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes)
		if err != nil {
			t.Errorf("could not parse private key")
		}

		// cleanup files
		err = os.Remove(csrFilePath)
		if err != nil {
			t.Errorf("could not delete %s", csrFilePath)
		}
		err = os.Remove(keyFilePath)
		if err != nil {
			t.Errorf("could not delete %s", keyFilePath)
		}
	})

	t.Run("test createCSR against reference", func(t *testing.T) {
		// Reference material generated with (NAME = "christoph"):
		//    openssl genrsa -out "./pki/client.${NAME}.key" 4096 2>/dev/null
		//    openssl req \
		//    -new \
		//    -addext  "subjectAltName = DNS:${NAME}.unsafe.aurae.io" \
		//    -subj    "/C=IS/ST=aurae/L=aurae/O=Aurae/OU=Runtime/CN=${NAME}.unsafe.aurae.io" \
		//    -key     "./pki/client.${NAME}.key" \
		//    -out     "./pki/client.${NAME}.csr" 2>/dev/null
		// ...and base64 encoded
		referenceCSR64 := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJRTd6Q0NBdGNDQVFBd2N6RUxNQWtHQTFVRUJoTUNTVk14RGpBTUJnTlZCQWdNQldGMWNtRmxNUTR3REFZRApWUVFIREFWaGRYSmhaVEVPTUF3R0ExVUVDZ3dGUVhWeVlXVXhFREFPQmdOVkJBc01CMUoxYm5ScGJXVXhJakFnCkJnTlZCQU1NR1dOb2NtbHpkRzl3YUM1MWJuTmhabVV1WVhWeVlXVXVhVzh3Z2dJaU1BMEdDU3FHU0liM0RRRUIKQVFVQUE0SUNEd0F3Z2dJS0FvSUNBUUMyY2hEc1A2VnZPd1ZybVlGbE1DVjc5Uk95Rk40Qkdxem01NFM5SHVhVAppK09HY3RRUHNzL3loWVhmaWNPN3JWaldEMW02ODgvRXFGeFlXZ2xKbHJRT240NXE0dWxxVitPK0JvZGhqTWlzClNxbU1zamEraGM5cXE0bnIyRDF3Vno1L1RjQnBlZkd5OTVoL05xSElUTStRU0xRbFZ2djA4SWs0QXRybmVmekQKbCtnWGdkSVFVQkpXTTJ5U3VXZ1ViSjlSNURrUm9vRnZOVzhjaU1zQlhOQlJrQTBrZjdhTEQ3NGxCanFYYXNwcwo5bWlxSUhiSUR6TjZLZGZIbDdNOU9UQnFxMURQZHpwMy9qeVNoR21YVFYvai80WVBMUnVtVFRHQmlnK1RLTzFpCm53VnlRUDdoYjVIb0dla2RGK1ZyZ3Z2THcvN0MxUDBhMUYrMXp6YVJITUJwbWQzUVFxZEtwRzlZeW1WelhHd2wKUmJhTXVDWG1Ea2hTY3RYaDJYU0tCM0ZSQVhlc2x5L25Mc0NjOWIzYXZpTThrVUFsTG9rdUlCYzlDanBiZTBJRQpDTlovZklZcnRoZ2dIWVZOaWJpbm1NVTE3b2J5RzRuTW1RWTh1cngwYi95VFdrdk0rbjBLMzljd3FxRGNoNnZLCm5OZU90WVhiY3dDM2djdExpYkZ3Vm1Pai83MTRtSHU0UnVBWWVQcUxtUmlCNVkvcnNQbWR0Q2p0SEJxcE40NzkKdUV3SXZaVUxNVlJZcWFMZ0dFejdDckNCUFoyU09IYzB5akxxb3lLS0tCVGU5OGJVeFdUay9YQUhIYVlibjBIdwoyS1kxWmxnZ1haME1Rd0w1VVcwYkc4VkdlQk9hc2VoQ3M1dzBSckpqcU9QQ3gwQmlUcTVNeXAzV0dtSFFydXdjCmlRSURBUUFCb0Rjd05RWUpLb1pJaHZjTkFRa09NU2d3SmpBa0JnTlZIUkVFSFRBYmdobGphSEpwYzNSdmNHZ3UKZFc1ellXWmxMbUYxY21GbExtbHZNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUNBUUJpSzZSbUdRbzd2NXozMXo3UAp5Z3BIS2V0Z2U4U3FvcFU2Yk5TWUdETTVWcmRBYkQ2YUdBOUxQREs3bngwaEx4VVhka3NrcHY0SkpORVdIZHhVCkN2WUFOY1o1RGVNaklHaWlTZi9RVzkvVm4rWWlLaDV0Nk5tdHplZ2xXYi9DN2lzdmwwdytyR3VyQVhBU0VZaGwKUXp0dUNvR3hWVlo2aUo3OGU3cGhVTGFRNjlySDZhd3FIVVV4MXkxKzBPU3JEVG15WkI3Z20yRjVqVlRYZDZyRwp5MjNwajZUd3B3UHc2b1pEditoSjAxVWY2VXhSZ09QR042aS93dWhjQ2E2RFBjbGh4WFZmYjRzR2RwaWZwMWl6Ck9WMFNTZlpRNUdRanNXeWNzQTFOcGRDeTlkNjVXLzRTRG9ZK0JPVTBGZVpoS0Y2ODQrVEhwR3FaWXNaQ1MySGsKOUhaWCtRYnA0WVhKYTQxSjJRb29OdjZ4VHY4VG5OVW5RdGpJaHQ1aEl1eFBTMnZTTWYxUnlpT1oxOWh4NWZoRQo0SmtXdk1nekNNZHh4RVQ1NUFCaWt2SVo0bUNFYm4xT3JiWlNvR2hWZkFjUG5hOHUvOEFvQkZqL1Y5dnF3TzdOCjYwWmFYK1Y4WkRlY3ZIMkUyRHVJdTQyWXBqZVNUMjF0L1FEdE5McHVBdTc1OUYwZjRUbUNJbm1wZG9HcDNFcUwKL0JMNjRNRDByZWRTaDhnUFR2SDUyODEyc2I5dmJxUHBFSjBWOG90amxQMlhKdVJDNFJPOVFoOEo0TWl3TnNUOQp5NFFpSWxaOTl3d0lpL1YvaEUrVzVWNWdybC9wdTk1LzJmTDUrL2RvdC83c0NMb0pnd0UwS2svNm91NUtjZFFICitmeS8rVyswK2JleEdVYTdLWVZ4YkNCMzhRPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUgUkVRVUVTVC0tLS0tCg=="
		referenceKey64 := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS0FJQkFBS0NBZ0VBdG5JUTdEK2xienNGYTVtQlpUQWxlL1VUc2hUZUFScXM1dWVFdlI3bWs0dmpobkxVCkQ3TFA4b1dGMzRuRHU2MVkxZzladXZQUHhLaGNXRm9KU1phMERwK09hdUxwYWxmanZnYUhZWXpJckVxcGpMSTIKdm9YUGFxdUo2OWc5Y0ZjK2YwM0FhWG54c3ZlWWZ6YWh5RXpQa0VpMEpWYjc5UENKT0FMYTUzbjh3NWZvRjRIUwpFRkFTVmpOc2tybG9GR3lmVWVRNUVhS0JielZ2SElqTEFWelFVWkFOSkgrMml3KytKUVk2bDJyS2JQWm9xaUIyCnlBOHplaW5YeDVlelBUa3dhcXRRejNjNmQvNDhrb1JwbDAxZjQvK0dEeTBicGsweGdZb1BreWp0WXA4RmNrRCsKNFcrUjZCbnBIUmZsYTRMN3k4UCt3dFQ5R3RSZnRjODJrUnpBYVpuZDBFS25TcVJ2V01wbGMxeHNKVVcyakxnbAo1ZzVJVW5MVjRkbDBpZ2R4VVFGM3JKY3Y1eTdBblBXOTJyNGpQSkZBSlM2SkxpQVhQUW82VzN0Q0JBaldmM3lHCks3WVlJQjJGVFltNHA1akZOZTZHOGh1SnpKa0dQTHE4ZEcvOGsxcEx6UHA5Q3QvWE1LcWczSWVyeXB6WGpyV0YKMjNNQXQ0SExTNG14Y0Zaam8vKzllSmg3dUViZ0dIajZpNWtZZ2VXUDY3RDVuYlFvN1J3YXFUZU8vYmhNQ0wyVgpDekZVV0ttaTRCaE0rd3F3Z1QyZGtqaDNOTW95NnFNaWlpZ1UzdmZHMU1WazVQMXdCeDJtRzU5QjhOaW1OV1pZCklGMmRERU1DK1ZGdEd4dkZSbmdUbXJIb1FyT2NORWF5WTZqandzZEFZazZ1VE1xZDFocGgwSzdzSElrQ0F3RUEKQVFLQ0FnRUFnSDBtMCtzakRJb0prRFRrdnluQVRHTldRcVdWa0R1RUozNUhxcFYzbDlQK0lqTCtqQ3ZIYmFxQgpsT1BHR0lmRnQ4UEowdk5na01SdGZMKzBLTUpjL3F0Nk5tYW1Nb0hCWDVQamhsMEsrdVArTXB0VUdLdk9YdlorClJMM2V6eDV5WWwrVXNmUHl0N0xPRUZHZWNKMC8xUUtPOUhrbEt1UzREdDFiNDRleTd1RXQwRmhhWTZpd3NVcTQKSFVFOFBwNGRPaVE3Mk9LVXU0aHJQekpMbmlNS2gxYW5HdHhpNTk3bmI5WEtMOWRDeHFobkgrR0xKZXdtdWROOApKeEg4WnBLL09YQjdraEVLK1hUd25kTnBOZWlGTHVKSFBLcnMvUnNDVVpPMDBsUVJrdElobU15VGRKc0pxK2VMCm1EUzdHeE45VjQwcC8zYlc1aTFKVnBhZmZHVStVREFBNlRCcTQvNHNNV1ZFL2FRUWdFRW1jMWQzK3d6MWkxQSsKRWF2dGZUU0VlaEIxemlXKzg0WFluVm96M3JQL0RRVE9qNXRIZzdqMzZYMDlQZ2J1LzFUUU9URmtCWDJhbElQQwpab295RkcvVXUrNkN0ZWtjdTVRZStodTYxbEt2MzE3SXlrWFVXMGV0MmEwQXFYZ0twMlFhNExQNDYwdUxRdW0vCjEzZGIvdnVibFNpMDBLQU1lVVRUaVM1N1Z6UU1nUVFUbGVCekZDbGxvZnFwTjNCeTJ6SVJIeHVMVmZPeVlsU0oKTVF3M1B2bTloTkduUkVDSHhIQ2MrT1BBa25WK1VCRTBVLzdTeWF0Rm9rdnFHNitQbGFUbHd6OE1uYXNBTm5HTAoyMW9HSWN5U0xISnpWek1ST3JiZ0Rua2QzREI0QWE4RXFwamN1NGtQdFFSbFlxb1djaTBDZ2dFQkFPTitlRjdDCmcvbnpKSVIzRWtvZGNYc2FsU2xjTjlKcVF6NExYbXdQUm9uYlhoeUlwOEdBR2xYSDJ5YUZ4Z0l3aENDYWtBQmgKT0grVGRGWlVteGFkRTdQRmxiU0lZL3hoK3JlVHNaL2lSSitDcXh4VkpqTHYwcXNtbGN4aWV0NTJSUklOVXJLegp1YVNCZDJwdTV1TStQWEtZSjM4VmdFdVdIVjlVRm9KTG00YVlKdnFQVU0zT1VicTdpTVVYODBGeklKRm93QUhJCnU0dHNvZ3JGdGxLdU5tREZnVEZWY0xNbXJrRG5CUzA2YjFROU9ROVVkWUF6OXBFSHRPcnFuT2ZvTGptQmRpZmgKQTVqQzNnL1NYdW5XcnV3SW13K1ZkTXdGSnVsaDZld01vd0VkRXdXcjUwK0lrWlFuSW84cUs5c0wxZFZJNEJHagp2WEVqY3ZDdUVvYjVZOE1DZ2dFQkFNMU9pVWtuK3lXWVNlV0NHRVhsR0RyYklNN0lLSERiQzg2d1B0bWVWNXVZCkhxeGh5akdhYll4WWZUSEI2RDJjVlUxYlNPSVZ0TEpaVVRhSlRPK0dEMCtKNU4vcStQbWloWldQcXJKNVg1WFMKSWRHQ1BFaEEyZERMc0ZoYlpQMmVyVXBqd1d3YmJBVjBoWU41QUoxdTNuMVQzdlVkTGNsT0xkeGRURkZSd0Q5VApyZG1xL1VMOStTVmprSFpqVEdpQ1UyWHVqVXdvYnNRNlgvUk0yeG52T3JjMDk0R2hRNVdMUmpZbThid0NMclA2CnVES1JLWUFHdkhSakJDYmRGRG5NSGQrMWMrVi95Nk96TGtVdlM3LzZZbkRlenB4Rk1WSGdJRGpvaytBZjFpRnAKc1FXNlEwSktVUmhqQWYzV2llOXJ5aVA5TU1ZT0tFNnBOMXlFWmJWZmRjTUNnZ0VBYXNTaWJhYlJGZS85UllZMAp1VUFVVUhoclpSdjR2dkpNV01ReExub0UyeEp2bXVpd0F1ckNjVnY1Q0oxa0R3Y0NHK011amw4U2l4MkRUamtyCkNIUDBHVDAwUTZSM2VLM3JZMWtYMWpmMWlQOWttMG1EUWdpNFVNY3RLdDFWV1M4Y3Y1b3RJOTJoMVFsR0tGZWcKV1NxTzRFZDAwZm9mV2xvN3NzL2VPSXlQazUyNVBZTWhvMVdmbWdvRjZLcVM2amJFSkRxTFVzc0k2aWl6N0daYQphWGVGNGVrUDl6MW9SVXgwSDlYTTRpczRzTXFEQ3lUU2VMYnFrNnFRU0dpUDkyOUtzb2FHRTdWUllOS2tNYnpECit1OWM3VDRrdUMybXdWSHhyenJhOUlRQnhMWUdoWFRtZkxkVnk3aUtTYks3SG5UeGlNWkpFejVMM051TVNGVUsKTVBxK3pRS0NBUUIrUE9FclMxc2d0YkFTWDlqZStVdlp2SzFDbUU1TmZsS1hSMFdOOTgrMGkyZW81UVEzVmRZdwpLcVRvT0d1OW5tZlJCZVVkcHUwUmtOdmY1Yktad051Zk01RzRvVGx2L1orWDQ5dTRtK3JMSzRiQjFRdU4vZG93CmlWNG9KaUpGMUJDSG9pam5lVUVGWmExR3R0dEs4a1g1MTkxSzZDTWtHVjhYbFlKOHFnREVyNFpCUmVNdUV3M2sKRUlGZVdoWThXSTVCS2RwVnpySzFFNU8ybXA5S0pnLzdZS1VqWHU0NGdJZXVlbW0vQ2JSLzFCVDRlc3VDdmlHWQppdDJkcStobzFYbzArTlNIYy9uWjhTM3RPblNnV2F1MzdUZ3JYRnhFRk1TYldWNjd1N2VsbWVCUVBrUm0rVjA5CjJucjZBcldUc3JwN1FJNkI2V2lkWFd6K0JTYW96RWFUQW9JQkFEY0FzNUQ4elg2dTN6dEhzSjg5amhFSzZ0L0wKOXo1UG5ocERocXA5Rk1BbC9jVUZvZkRQaTVOSGRrbjNPNDFQc1B1TFhHTWZ4SkFpOXFKZTdDVENmS3JvelVBbwprNzR3YXBQbEtIVjNsQm9hQUtDcUpQa1B0bjk3WG8rM215TUZJelRaVGtkRGNlY3ByZHZQeDZBbEZmME9BazBmCjlmZFJ0ZUNUSzhzUnV4WjJKTHJyRSs2ZG9HSm11NEpadmd2c0pheGliZm5yWFZjMVdvSkRCMEFvdmZ3SldxYVUKbStkaVhIcUhPQVE5NWVESlVJNWFsMXI2dW9xQWc4UStTYWY1T204dy9xckhZZWNiZ0pibG5tVDROelhDeFEvZgo1OTZzcFQ1dUllVjZCRFMzS2RBaDJtcWNNR1dwWGFnc3FHRjlNaFpyZlMvM0trMExFdk5xS0pqOVQraz0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K"

		referenceCSRBytes, err := base64.StdEncoding.DecodeString(referenceCSR64)
		if err != nil {
			t.Errorf("could not decode reference CSR")
		}

		referenceKeyBytes, err := base64.StdEncoding.DecodeString(referenceKey64)
		if err != nil {
			t.Errorf("could not decode reference key")
		}

		p, _ := pem.Decode(referenceCSRBytes)

		referenceCSR, err := x509.ParseCertificateRequest(p.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate request: %s", err)
		}

		// Public Key Algorithm: rsaEncryption
		// RSA Public-Key: (4096 bit)
		if referenceCSR.PublicKeyAlgorithm.String() != x509.RSA.String() {
			t.Errorf("public key algorithm of reference is not correct:\ngot: %s, want: %s", referenceCSR.PublicKeyAlgorithm.String(), x509.RSA.String())
		}
		if referenceCSR.PublicKey.(*rsa.PublicKey).N.BitLen() != 4096 {
			t.Errorf("public key length of reference is not correct:\ngot: %d, want: %d", referenceCSR.PublicKey.(*rsa.PublicKey).N.BitLen(), 4096)
		}

		// Checking whether the reference CSR is correct:
		// Subject: C=IS, ST=aurae, L=aurae, O=Aurae, OU=Runtime, CN=christoph.unsafe.aurae.io
		if referenceCSR.Subject.Country[0] != "IS" {
			t.Errorf("subject country of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.Country[0], "IS")
		}
		if referenceCSR.Subject.Province[0] != "aurae" {
			t.Errorf("subject state or province of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.Province[0], "aurae")
		}
		if referenceCSR.Subject.Locality[0] != "aurae" {
			t.Errorf("subject locality of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.Locality[0], "aurae")
		}
		if referenceCSR.Subject.Organization[0] != "Aurae" {
			t.Errorf("subject organization of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.Organization[0], "Aurae")
		}
		if referenceCSR.Subject.OrganizationalUnit[0] != "Runtime" {
			t.Errorf("subject organizational unit of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.OrganizationalUnit[0], "Runtime")
		}
		if referenceCSR.Subject.CommonName != "christoph.unsafe.aurae.io" {
			t.Errorf("subject common name of reference is not correct:\ngot: %s, want: %s", referenceCSR.Subject.CommonName, "christoph.unsafe.aurae.io")
		}

		// Genenerate a new CSR
		clientCsr, err := CreateClientCSR("", "unsafe.aurae.io", "christoph")
		if err != nil {
			t.Errorf("could create csr")
		}
		csrPem, _ := pem.Decode([]byte(clientCsr.CSR))
		generatedCSR, err := x509.ParseCertificateRequest(csrPem.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate request")
		}

		// Public Key Algorithm: rsaEncryption
		if generatedCSR.PublicKeyAlgorithm.String() != referenceCSR.PublicKeyAlgorithm.String() {
			t.Errorf("public key algorithm is not correct:\ngot: %s, want: %s", referenceCSR.PublicKeyAlgorithm.String(), x509.RSA.String())
		}
		// RSA Public-Key: (4096 bit)
		if generatedCSR.PublicKey.(*rsa.PublicKey).N.BitLen() != referenceCSR.PublicKey.(*rsa.PublicKey).N.BitLen() {
			t.Errorf("public key length is not correct:\ngot: %d, want: %d", referenceCSR.PublicKey.(*rsa.PublicKey).N.BitLen(), 4096)
		}

		// Checking whether the reference CSR is correct:
		// Subject: C=IS, ST=aurae, L=aurae, O=Aurae, OU=Runtime, CN=christoph.unsafe.aurae.io
		if generatedCSR.Subject.Country[0] != referenceCSR.Subject.Country[0] {
			t.Errorf("subject country is not correct:\ngot: %s, want: %s", generatedCSR.Subject.Country[0], referenceCSR.Subject.Country[0])
		}
		if generatedCSR.Subject.Province[0] != referenceCSR.Subject.Province[0] {
			t.Errorf("subject state or province is not correct:\ngot: %s, want: %s", generatedCSR.Subject.Province[0], referenceCSR.Subject.Province[0])
		}
		if generatedCSR.Subject.Locality[0] != referenceCSR.Subject.Locality[0] {
			t.Errorf("subject locality is not correct:\ngot: %s, want: %s", generatedCSR.Subject.Locality[0], referenceCSR.Subject.Locality[0])
		}
		if generatedCSR.Subject.Organization[0] != referenceCSR.Subject.Organization[0] {
			t.Errorf("subject organization is not correct:\ngot: %s, want: %s", generatedCSR.Subject.Organization[0], referenceCSR.Subject.Organization[0])
		}
		if generatedCSR.Subject.OrganizationalUnit[0] != referenceCSR.Subject.OrganizationalUnit[0] {
			t.Errorf("subject organizational unit is not correct:\ngot: %s, want: %s", generatedCSR.Subject.OrganizationalUnit[0], referenceCSR.Subject.OrganizationalUnit[0])
		}
		if generatedCSR.Subject.CommonName != referenceCSR.Subject.CommonName {
			t.Errorf("subject common name is not correct:\ngot: %s, want: %s", referenceCSR.Subject.CommonName, referenceCSR.Subject.CommonName)
		}

		// Testing attributes:
		//  Requested Extensions:
		//      X509v3 Subject Alternative Name:
		//          DNS:christoph.unsafe.aurae.io
		if generatedCSR.DNSNames[0] != referenceCSR.DNSNames[0] {
			t.Errorf("subject alternative name is not correct:\ngot: %s, want: %s", referenceCSR.DNSNames[0], referenceCSR.DNSNames[0])
		}

		keyPem, _ := pem.Decode(referenceKeyBytes)
		_, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes)
		if err != nil {
			t.Errorf("could not parse private key")
		}
	})
}
