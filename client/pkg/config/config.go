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

package config

import (
	"fmt"
	"log"
	"os/user"
)

type Configs struct {
	Auth   Auth
	System System
}

type Config interface {
	Set(p *Configs) error
}

func Default() (*Configs, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Cannot get current user", err)

		return nil, err
	}

	cfg := &Configs{
		Auth: Auth{
			CaCert:     fmt.Sprintf("%s/.aurae/pki/ca.crt", usr.HomeDir),
			ClientCert: fmt.Sprintf("%s/.aurae/pki/_signed.client.nova.crt", usr.HomeDir),
			ClientKey:  fmt.Sprintf("%s/.aurae/pki/client.nova.key", usr.HomeDir),
			ServerName: "server.unsafe.aurae.io",
		},
		System: System{
			Protocol: "unix",
			Socket:   "/var/run/aurae/aurae.sock",
		},
	}

	return cfg, nil
}

func From(cfg ...Config) (*Configs, error) {
	c, err := Default()
	if err != nil {
		return nil, err
	}

	for _, config := range cfg {
		err := config.Set(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
