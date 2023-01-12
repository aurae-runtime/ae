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

package version

import (
	"testing"

	cmdTest "github.com/aurae-runtime/ae/cmd/test"
	"github.com/prometheus/common/version"
)

func TestVersionCMD(t *testing.T) {
	version.Version = "v0.1.0"
	version.BuildDate = "2023-01-07"
	version.Revision = "a7c46aa017bc447ece506629196bd0548cbbc469"
	testSuite := []cmdTest.Suite{
		{
			Title:           "empty args",
			Args:            []string{},
			IsErrorExpected: false,
			ExpectedMessage: `{
	"buildTime": "2023-01-07",
	"version": "v0.1.0",
	"commit": "a7c46aa017bc447ece506629196bd0548cbbc469"
}`,
		},
		{
			Title:           "print short version",
			Args:            []string{"--short"},
			IsErrorExpected: false,
			ExpectedMessage: `{
	"version": "v0.1.0"
}`,
		},
	}
	cmdTest.ExecuteSuiteTest(t, NewCMD, testSuite)
}
