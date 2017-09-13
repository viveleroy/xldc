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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/viveleroy/goxldeploy"
)

// verifyCmd represents the verify command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "meta displays metadata",
	Long:  `retrieves and displays metadata form xldeploy`,
}

var metaTypeCommand = &cobra.Command{
	Use:   "type",
	Short: "Display metadata for types",
	Long:  "fetches metadata from xldeploy for a single type",
	Run:   getTypeMetadata,
}
var metaOrchestratorCommand = &cobra.Command{
	Use:   "orchestrators",
	Short: "Display metadata for orchestrators",
	Long:  "fetches a list of available orchestrators from xldeploy ",
	Run:   getOrchestratorMetadata,
}
var metaPermissionsCommand = &cobra.Command{
	Use:   "permissions",
	Short: "Display metadata for Permissions",
	Long:  "fetches metadata from xldeploy concerning permissions",
	Run:   getPermissionMetadata,
}

func init() {
	metaCmd.AddCommand(metaTypeCommand)
	metaCmd.AddCommand(metaOrchestratorCommand)
	metaCmd.AddCommand(metaPermissionsCommand)
	RootCmd.AddCommand(metaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTypeMetadata(cmd *cobra.Command, args []string) {

	//Lets declare an interface for output
	var o interface{}
	var err error

	cfg := goxldeploy.Config{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
		Context:  context,
		Scheme:   scheme,
	}

	xld := goxldeploy.New(&cfg)

	if len(args) == 0 {
		o, err = xld.Metadata.GetTypeList()
		if err != nil {
			jww.FATAL.Printf("%s: encounterd a fatal error in retrieving metadata: %s", cmd.CommandPath(), err)
			os.Exit(1)
		}
	} else {
		o, err = xld.Metadata.GetType(args[0])
		if err != nil {
			jww.FATAL.Printf("%s: encounterd a fatal error in retrieving metadata for %s: %s", cmd.CommandPath(), args[0], err)
			os.Exit(1)
		}
	}

	fmt.Println(RenderJSON(o))

}

func getOrchestratorMetadata(cmd *cobra.Command, args []string) {

	cfg := goxldeploy.Config{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
		Context:  context,
		Scheme:   scheme,
	}

	xld := goxldeploy.New(&cfg)

	o, err := xld.Metadata.GetOrchestrators()
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error while retrieving metadata: %s", cmd.CommandPath(), err)
		os.Exit(1)
	}

	fmt.Println(RenderJSON(o))
}

func getPermissionMetadata(cmd *cobra.Command, args []string) {

	cfg := goxldeploy.Config{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
		Context:  context,
		Scheme:   scheme,
	}

	xld := goxldeploy.New(&cfg)

	o, err := xld.Metadata.GetPermissions()
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error while retrieving metadata: %s", cmd.CommandPath(), err)
		os.Exit(1)
	}

	fmt.Println(RenderJSON(o))
}

//RenderJSON function to render output as json
// returns a string object with json formated output
func RenderJSON(l interface{}) string {

	b, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		panic(err)
	}
	s := string(b)

	return s
}
