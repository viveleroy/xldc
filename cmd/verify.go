// Copyright © 2017 Roy Kliment <roy.kliment@cinqict.nl>
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
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/viveleroy/goxldeploy"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies connection to XL-Deploy",
	Long:  `Verifies connection to XL-Deploy using the goxldeploy library`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyConnection()
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func verifyConnection() {
	cfg := goxldeploy.Config{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
		Context:  context,
		Scheme:   scheme,
	}
	xld := goxldeploy.New(&cfg)
	c := xld.Connected()
	if c {
		jww.INFO.Println("Connection to XL-Deploy verified")
		if verbose == false {
			fmt.Println("Connection to XL-Deploy verified")
		}
	} else {
		jww.FATAL.Println("Connection to XL-Deploy failed")
		os.Exit(1)
	}
}
