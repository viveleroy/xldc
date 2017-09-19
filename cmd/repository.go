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
	"os"
	"strings"

	"github.com/viveleroy/goxldeploy"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var id string
var ciType string

// verifyCmd represents the verify command
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "handle repository operations",
	Long:  `do inserts deletions, updated and deletions from and to the the xldeploy repository database`,
}

var repositoryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an already existing ci in the repository",
	Run:   UpdateCI,
}

var repositoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a ci from the repository",
	Run:   GetCI,
}

var repositoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a ci in the xld repository",
	Long:  "usage: create --id <id> --type <type> <comma seperated list prop=val",
	Run:   CreateCI,
}

func init() {
	// add flags to the various previously defined commands
	repositoryGetCmd.Flags().StringVarP(&outputFile, "out", "", "", "specify an output file")
	repositoryGetCmd.Flags().StringVarP(&inputFile, "in", "", "", "specify an input file containing")

	repositoryCmd.AddCommand(repositoryGetCmd)

	repositoryUpdateCmd.Flags().BoolVarP(&merge, "merge", "m", false, "merge the update with the existing ci")
	repositoryCmd.AddCommand(repositoryUpdateCmd)
	// add the commands to da mothership

	repositoryCreateCmd.Flags().StringVarP(&id, "id", "i", "", "specify ci id")
	repositoryCreateCmd.Flags().StringVarP(&ciType, "type", "t", "", "specify ci type")
	repositoryCmd.AddCommand(repositoryCreateCmd)

	RootCmd.AddCommand(repositoryCmd)

}

func GetCI(cmd *cobra.Command, args []string) {
	xld := GetClient()

	if len(args) == 0 {
		jww.FATAL.Printf("%s: requires on argument", cmd.CommandPath())
		os.Exit(1)
	}

	ci, err := xld.Repository.GetCI(args[0])

	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error in retrieving Configuration item %s: %s", cmd.CommandPath(), args[0], err)
		os.Exit(1)
	}

	if outputFile != "" {
		WriteJSONToFile(ci, outputFile)
		os.Exit(0)
	}

	RenderJSON(ci)
}

func CreateCI(cmd *cobra.Command, args []string) {

	var ci goxldeploy.Ci

	xld := GetClient()

	if len(args) != 0 {
		ci = goxldeploy.NewCI(id, ciType, splitPropertiesString(args[0]))
	}

	ci = goxldeploy.NewCI(id, ciType, nil)

	ci, err := xld.Repository.CreateCI(ci)
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error creating configuration Item %s: %s", cmd.CommandPath(), id, err)
		os.Exit(1)
	}

	RenderJSON(ci)

}

func UpdateCI(cmd *cobra.Command, args []string) {
	var ci goxldeploy.Ci

	xld := GetClient()

	if len(args) != 0 {
		ci = goxldeploy.NewCI(id, ciType, splitPropertiesString(args[0]))
	}

	ci = goxldeploy.NewCI(id, ciType, nil)

	ci, err := xld.Repository.CreateCI(ci)
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error creating configuration Item %s: %s", cmd.CommandPath(), id, err)
		os.Exit(1)
	}

	RenderJSON(ci)
}

func splitPropertiesString(s string) map[string]interface{} {
	var properties map[string]interface{}

	ps := strings.Split(s, ",")
	for _, p := range ps {
		kv := strings.Split(p, "=")
		properties[kv[0]] = kv[1]

	}

	return properties
}
