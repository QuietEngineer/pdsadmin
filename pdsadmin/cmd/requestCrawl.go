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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SyncRequestCrawl_Input struct {
	// hostname: Hostname of the current service (eg, PDS) that is requesting to be crawled.
	Hostname string `json:"hostname"`
}

// requestCrawlCmd represents the requestCrawl command
var requestCrawlCmd = &cobra.Command{
	Use:   "requestCrawl",
	Short: "Request a crawl from a relay host",
	Example: `pdsadmin requestCrawl bsky.network
pdsadmin requestCrawl bsky.network,second-relay.example.com`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var relayHosts []string
		if len(args) == 1 {
			relayHosts = strings.Split(args[0], ",")
		} else {
			relayHosts = viper.GetStringSlice("crawlers")
		}
		if len(relayHosts) == 0 {
			fmt.Println("ERROR: missing RELAY HOST parameter")
			os.Exit(1)
		}

		client := &http.Client{}
		for _, host := range relayHosts {
			if host == "" {
				continue
			}
			fmt.Printf("Requesting crawl from %s\n", host)
			if !strings.HasPrefix(host, "https:") && !strings.HasPrefix(host, "http:") {
				host = fmt.Sprintf("https://%s", host)
			}
			body := SyncRequestCrawl_Input{
				Hostname: viper.GetString("hostname"),
			}
			jsonBody, err := json.Marshal(body)
			if err != nil {
				fmt.Printf("could not create json body: %s\n", err)
				continue
			}
			bodyReader := bytes.NewReader(jsonBody)

			res, err := client.Post(fmt.Sprintf("%s/xrpc/com.atproto.sync.requestCrawl", host), "application/json", bodyReader)
			if err != nil {
				fmt.Printf("error making http request: %s\n", err)
				continue
			}

			if _, err := io.ReadAll(res.Body); err != nil {
				fmt.Printf("could not read response body: %s\n", err)
				continue
			}
		}
		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(requestCrawlCmd)
}
