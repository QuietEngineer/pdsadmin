/*
Copyright © 2025 QuietEngineer <qtengineer@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ServerCreateInviteCode_Input struct {
	UseCount int `json:"useCount"`
}

type ServerCreateInviteCode_Output struct {
	Code string `json:"code"`
}

var number *int

// createInviteCodeCmd represents the createInviteCode command
var createInviteCodeCmd = &cobra.Command{
	Use:     "createInviteCode",
	Short:   "Create a new invite code",
	Example: "pdsadmin createInviteCode",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if *number < 1 {
			fmt.Println("number must be >=1")
			os.Exit(1)
		}
		body := ServerCreateInviteCode_Input{
			UseCount: *number,
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("could not create json body: %s\n", err)
			os.Exit(1)
		}
		bodyReader := bytes.NewReader(jsonBody)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/xrpc/com.atproto.server.createInviteCode", viper.GetString("hostname")), bodyReader)
		if err != nil {
			fmt.Printf("could not create request: %s\n", err)
			os.Exit(1)
		}
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth("admin", viper.GetString("admin_password"))

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("could not read response body: %s\n", err)
			os.Exit(1)
		}

		var inviteCode ServerCreateInviteCode_Output
		if err := json.Unmarshal(resBody, &inviteCode); err != nil {
			fmt.Printf("could not get invite code: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(inviteCode.Code)
	},
}

func init() {
	rootCmd.AddCommand(createInviteCodeCmd)
	number = createInviteCodeCmd.Flags().IntP("number", "n", 1, "number of times the code can be used")
}
