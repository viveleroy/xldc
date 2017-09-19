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

//typeShort is used in provicing a non verbose display of a type or a typelist
type typeShort struct {
	Name        string `json:"type"`
	Description string `json:"description"`
}

// verifyCmd represents the verify command
var metaCmd = &cobra.Command{
	Use:   "metadata",
	Short: "meta displays metadata",
	Long:  `retrieves and displays metadata form xldeploy`,
}

var metaTypeCommand = &cobra.Command{
	Use:   "type",
	Short: "Display metadata for types",
	Long:  "fetches metadata from xldeploy for a single type",
	Run:   getTypeMetadata,
}
var metaTemplateCommand = &cobra.Command{
	Use:   "template",
	Short: "Display's template for type creation",
	Long:  "fetches metadata from xldeploy for a single type and returns a json template that can be used for the creation of said type",
	Run:   getTypeTemplate,
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

var longBool bool
var optionalBool bool

func init() {

	// add flags to the various previously defined commands
	metaTypeCommand.Flags().BoolVarP(&longBool, "long", "l", false, "print long listing instead of condensed output")
	metaTemplateCommand.Flags().BoolVarP(&optionalBool, "optional", "o", false, "include optional parameters in template")

	// add the commands to da mothership
	metaCmd.AddCommand(metaTypeCommand)

	metaCmd.AddCommand(metaTemplateCommand)

	metaCmd.AddCommand(metaOrchestratorCommand)

	metaCmd.AddCommand(metaPermissionsCommand)

	RootCmd.AddCommand(metaCmd)

}

//Gets a list of
func getTypeTemplate(cmd *cobra.Command, args []string) {
	var tl goxldeploy.TypeList
	var tmpl []map[string]interface{}

	xld := GetClient()

	if len(args) == 0 {
		jww.FATAL.Printf("%s:requires at least one argument", cmd.CommandPath())
		os.Exit(1)
	} else {
		for _, t := range args {
			tt, err := xld.Metadata.GetType(t)
			if err != nil {
				jww.FATAL.Printf("%s: encounterd a fatal error in retrieving metadata for %s: %s", cmd.CommandPath(), t, err)
				os.Exit(1)
			}
			tl = append(tl, tt)
		}
	}

	// loop over the typelist (tl) and create a template type map[string]interface{} for each element
	for _, t := range tl {
		templ := make(map[string]interface{})

		templ["name"] = ""
		templ["type"] = t.Type

		for _, p := range t.Properties {
			if p.Required == true {
				if p.Default == nil {
					templ[p.Name] = "required"
				} else {
					templ[p.Name] = p.Default
				}
			} else {
				if optionalBool == true {
					templ[p.Name] = "very optional"
				}
			}
		}
		tmpl = append(tmpl, templ)
	}

	// handle the output
	// if only one template was given then we do not want to bother our esteemd users with a slice representation
	var o interface{}

	if len(tmpl) == 1 {
		o = tmpl[0]
	} else {
		o = tmpl
	}

	// Yes ... yes we do file .. thank you ..
	if outputFile != "" {
		WriteJSONToFile(o, outputFile)
		os.Exit(0)
	}

	RenderJSON(o)
	os.Exit(0)
}

func getTypeMetadata(cmd *cobra.Command, args []string) {

	//Lets declare an interface for output
	var o interface{}
	// var out interface{}
	var err error

	xld := GetClient()

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

	//if the templateBool is set we will need to scan the returned metadata for required properties and
	// drop those in a map[string]Interface{}
	// then render that as valid output

	//if the longBool is false we want the short rundown of the types
	// so we have to determine if we're dealing with a list or a single type
	// when dealing with a single type .. get name and description from that
	// when dealing with a list .. compose an alternate list with type description pairs
	if longBool == false {

		//Figure out the type we got handed
		switch o.(type) {
		case goxldeploy.Type:
			//o is a interface{} lets assert it to goxldeploy.Type
			oT := o.(goxldeploy.Type)
			RenderJSON(typeShort{Name: oT.Type, Description: oT.Description})
			os.Exit(0)
		case goxldeploy.TypeList:
			var localOut []typeShort
			for _, t := range o.(goxldeploy.TypeList) {
				ts := typeShort{Name: t.Type, Description: t.Description}
				// fmt.Printf("%+v\n", t)
				localOut = append(localOut, ts)
			}
			RenderJSON(localOut)
			os.Exit(0)

		default:
			fmt.Printf("I don't know, ask stackoverflow.")
		}
	} else {
		RenderJSON(o)
	}

	// if not just
}

func getOrchestratorMetadata(cmd *cobra.Command, args []string) {

	xld := GetClient()

	o, err := xld.Metadata.GetOrchestrators()
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error while retrieving metadata: %s", cmd.CommandPath(), err)
		os.Exit(1)
	}

	RenderJSON(o)
}

func getPermissionMetadata(cmd *cobra.Command, args []string) {

	xld := GetClient()

	o, err := xld.Metadata.GetPermissions()
	if err != nil {
		jww.FATAL.Printf("%s: encounterd a fatal error while retrieving metadata: %s", cmd.CommandPath(), err)
		os.Exit(1)
	}

	RenderJSON(o)
}
