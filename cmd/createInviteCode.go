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
	Use:     "create-invite-code",
	Short:   "Create a new invite code",
	Example: "pdsadmin create-invite-code",
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
