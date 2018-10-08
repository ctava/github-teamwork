// Copyright © 2018 Chris Tava <chris1tava@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// VERSION of app
const VERSION = "0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version information, check for update and notify if new version available",
	Long:  `print version information, check for update and notify if new version available`,
	Run: func(cmd *cobra.Command, args []string) {
		app := "github-teamwork"
		fmt.Printf("%s v%s\n", app, VERSION)
		fmt.Println("\nChecking for new version...")

		resp, err := http.Get(fmt.Sprintf("https://github.com/ctava/%s/releases/latest", app))
		if err != nil {
			checkError(fmt.Errorf("Network error"))
		}
		items := strings.Split(resp.Request.URL.String(), "/")
		releasedVersion := ""
		if items[len(items)-1] == "" {
			releasedVersion = items[len(items)-2]
		} else {
			releasedVersion = items[len(items)-1]
		}
		if releasedVersion == VERSION {
			fmt.Printf("You are using the latest version of %s\n", app)
		}
		if releasedVersion < VERSION {
			fmt.Printf("New version available: %s %s at %s\n", app, releasedVersion, resp.Request.URL.String())
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
