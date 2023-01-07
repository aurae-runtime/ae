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

package test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

type Suite struct {
	Title           string
	Args            []string
	ExpectedMessage string
	IsErrorExpected bool
}

func ExecuteSuiteTest(t *testing.T, newCMD func() *cobra.Command, suites []Suite) {
	for _, test := range suites {
		t.Run(test.Title, func(t *testing.T) {
			buffer := bytes.NewBufferString("")
			cmd := newCMD()
			cmd.SetOut(buffer)
			cmd.SetErr(buffer)
			cmd.SetArgs(test.Args)

			err := cmd.Execute()
			if test.IsErrorExpected {
				if !IsNil(err) {
					if !IsEqualString(test.ExpectedMessage, err.Error()) {
						t.Errorf("\nError message of command \"%s\" was incorrect, got: \n%s\nwant: \n%s.", test.Title, err.Error(), test.ExpectedMessage)
					}
				} else {
					t.Errorf("\nError in command \"%s\" expected.", test.Title)
				}
			} else if IsNil(err) {
				if !IsEqualString(test.ExpectedMessage, buffer.String()) {
					t.Errorf("\nTest of command \"%s\" was incorrect, got: \n%swant: \n%s.", test.Title, buffer.String(), test.ExpectedMessage)
				}
			}
		})
	}
}

func IsEqualString(s1, s2 string) bool {
	return s1 == s2
}

func IsNil(object interface{}) bool {
	return object == nil
}
