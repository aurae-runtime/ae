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

package create

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"

	"github.com/aurae-runtime/ae/pkg/pki"
)

func TestPKICreateCMD(t *testing.T) {
	t.Run("ae pki create my.domain.com", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		cmd := NewCMD(context.Background())
		cmd.SetOut(buffer)
		cmd.SetErr(buffer)
		cmd.SetArgs([]string{"my.domain.com"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("ae pki create my.domain.com")
		}

		var ca pki.AuraeCA
		err = json.Unmarshal(buffer.Bytes(), &ca)
		if err != nil {
			t.Errorf("could not marshall certificate")
		}

		// load ca.cert
		cert, _ := pem.Decode([]byte(ca.Certificate))

		crt, err := x509.ParseCertificate(cert.Bytes)
		if err != nil {
			t.Errorf("could not parse certificate")
		}
		if crt.Subject.CommonName != "my.domain.com" {
			t.Errorf("certificate does not contain common name")
		}
	})
}
