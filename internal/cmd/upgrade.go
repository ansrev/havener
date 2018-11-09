// Copyright © 2018 The Havener
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
	"io/ioutil"

	"github.com/homeport/havener/pkg/havener"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Kubernetes with new Helm Charts",
	Long:  `TODO please do this later`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" && viper.GetString("havenerconfig") == "" {
			havener.ExitWithError("please provide configuration via --config or environment variable HAVENERCONFIG", fmt.Errorf("no havener configuration file set"))
		}

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		cfgdata, err := ioutil.ReadFile(viper.GetString("havenerconfig"))
		if err != nil {
			havener.ExitWithError("unable to read file", err)
		}

		var config havener.Config
		if err := yaml.Unmarshal(cfgdata, &config); err != nil {
			havener.ExitWithError("failed to unmarshal config file", err)
		}

		for _, release := range config.Releases {
			overrides, err := havener.TraverseStructureAndProcessShellOperators(release.Overrides)
			if err != nil {
				havener.ExitWithError("failed to process overrides section", err)
			}

			overridesData, err := yaml.Marshal(overrides)
			if err != nil {
				havener.ExitWithError("failed to marshal overrides structure into bytes", err)
			}

			fmt.Printf("Going to upgrade existing %s chart\n", release.ChartName)
			if _, err := havener.UpdateHelmRelease(release.ChartName, release.ChartLocation, overridesData); err != nil {
				havener.ExitWithError("Error updating chart", err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
