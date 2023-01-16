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
	"context"
	"testing"

	"github.com/aurae-runtime/ae/pkg/cli/testsuite"
	"github.com/prometheus/common/version"
)

func TestVersionCMD(t *testing.T) {
	version.Version = "v0.1.0"
	version.BuildDate = "2023-01-07"
	version.Revision = "a7c46aa017bc447ece506629196bd0548cbbc469"
	tests := []testsuite.Test{
		{
			Title: "empty args",
			Cmd:   NewCMD(context.Background()),
			Args:  []string{},
			ExpectedStdout: "{\n" +
				"    \"buildTime\": \"2023-01-07\",\n" +
				"    \"version\": \"v0.1.0\",\n" +
				"    \"commit\": \"a7c46aa017bc447ece506629196bd0548cbbc469\"\n" +
				"}\n",
		},
		{
			Title: "print version in json",
			Cmd:   NewCMD(context.Background()),
			Args:  []string{"--output", "json"},
			ExpectedStdout: "{\n" +
				"    \"buildTime\": \"2023-01-07\",\n" +
				"    \"version\": \"v0.1.0\",\n" +
				"    \"commit\": \"a7c46aa017bc447ece506629196bd0548cbbc469\"\n" +
				"}\n",
		},
		{
			Title: "print short version",
			Cmd:   NewCMD(context.Background()),
			Args:  []string{"--short"},
			ExpectedStdout: "{\n" +
				"    \"version\": \"v0.1.0\"\n" +
				"}\n",
		},
	}
	testsuite.ExecuteSuite(t, tests)
}
