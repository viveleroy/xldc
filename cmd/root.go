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
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/viveleroy/goxldeploy"
)

// vars for app
var cfgFile string
var verbose bool

// vars for flags
var username string
var password string
var host string
var port int
var context string
var ssl bool
var scheme string
var outputFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "xldc",
	Short:            "XL-Deploy CLI",
	Long:             `XL-Deploy CLI does some stuff with XL-Deploy`,
	PersistentPreRun: preVerifyConnection,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(setVerbose, initConfig, processConfig, checkRequiredFlags)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/xldc.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	RootCmd.PersistentFlags().StringVar(&username, "username", "", "Username for the connection")
	RootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for the connection")
	RootCmd.PersistentFlags().StringVar(&host, "host", "", "Hostname of the Xl-Deploy server")
	RootCmd.PersistentFlags().IntVar(&port, "port", 0, "Port where Xl-Deploy is running")
	RootCmd.PersistentFlags().StringVar(&context, "context", "", "Context-root where XL-Deploy is running")
	RootCmd.PersistentFlags().BoolVar(&ssl, "ssl", false, "Use SSL")
	RootCmd.PersistentFlags().StringVarP(&outputFile, "out", "", "", "specify an output file")
	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("context", RootCmd.PersistentFlags().Lookup("context"))
	viper.BindPFlag("ssl", RootCmd.PersistentFlags().Lookup("ssl"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in current directory with name "xldc" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("xldc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		jww.INFO.Println("Config", viper.ConfigFileUsed(), "not found")
	}
}

// setVerbose checks if verbose flag is set and adjusts logger accordingly
func setVerbose() {
	if verbose {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelInfo)
	}
}

// checkRequiredFlags will check if required flags have values without validation
func checkRequiredFlags() {
	if username == "" {
		jww.FATAL.Println("Username is required")
		os.Exit(1)
	}
	if password == "" {
		jww.FATAL.Println("Password is required")
		os.Exit(1)
	}
	if host == "" {
		jww.FATAL.Println("Host is required")
		os.Exit(1)
	}
	if port == 0 {
		jww.FATAL.Println("Port is required")
		os.Exit(1)
	}
	if context == "" {
		jww.WARN.Println("No context set, using / ")
		context = "/"
	}
	if ssl {
		scheme = "https"
	} else {
		scheme = "http"
	}
}

// processConfig will use viper config if flag is not set
func processConfig() {
	if username == "" && viper.IsSet("username") {
		username = viper.GetString("username")
	}
	if password == "" && viper.IsSet("password") {
		password = viper.GetString("password")
	}
	if host == "" && viper.IsSet("host") {
		host = viper.GetString("host")
	}
	if port == 0 && viper.IsSet("port") {
		port = viper.GetInt("port")
	}
	if viper.IsSet("context") {
		context = viper.GetString("context")
	}
	if viper.IsSet("ssl") {
		ssl = viper.GetBool("ssl")
	}
}

// preVerifyConnection will check if the connection can be established
func preVerifyConnection(cmd *cobra.Command, args []string) {
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
	} else {
		jww.FATAL.Println("Connection to XL-Deploy failed")
		os.Exit(1)
	}
}

//RenderJSON function to render output as json
// returns a string object with json formated output
func RenderJSON(l interface{}) {

	b, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		panic(err)
	}
	s := string(b)

	fmt.Println(s)
}

//GetClient returns a configured XLD client object
func GetClient() *goxldeploy.Client {

	cfg := goxldeploy.Config{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
		Context:  context,
		Scheme:   scheme,
	}

	return goxldeploy.New(&cfg)

}

//WriteToFile writes any string output to file
func WriteJSONToFile(l interface{}, f string) {

	var err error

	b, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		panic(err)
	}

	s := string(b)

	d1 := []byte(s + "\n")

	err = ioutil.WriteFile(f, d1, 0644)

	if err != nil {
		panic(err)
	}
}
