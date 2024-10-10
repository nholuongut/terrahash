/*
Copyright © 2024 the terrahash authors

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.

*/
package cmd

import (
	"fmt"
	"os"
	"log/slog"

	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/jedib0t/go-pretty/v6/table"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a mod lock file if one doesn't already exist.",
	Long: `Init scans the current Terraform configuration and produces a mod lock
	file if one doesn't already exist. This command will warn if a mod lock file
	is found or the Terraform configuration hasn't been initialized yet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("init command called")
		
		path, err := setPath(Source)
		if err != nil {
			return err
		}

		slog.Debug("check to see if terraform has been initialized")
		msg, init := terraformInitialized(path)

		if !init {
			slog.Warn(msg)
			slog.Warn("has terraform init been run?")
			return nil
		}

		// If the modFile already exists, run a check to see if any changes
		// are required, otherwise create the modFile
		slog.Debug("check to see if the " + modFileName + "file exists")
		if _, err := os.Stat(path + modFileName); err == nil {
			fmt.Println("Existing mod lock file found")
			checkErr := check(path)
			if checkErr != nil {
				slog.Error("Changes detected between configuration and lock file")
				slog.Error(checkErr.Error())
				return fmt.Errorf("run terrhash upgrade to update the mod lock file")
			}
			fmt.Println("No changes to mod lock file required")
			return nil
		}

		slog.Debug("get the modules used by the configuration")
		sourcedMods, err := processModules(path, IgnoreLocal)
		if err != nil {
			return fmt.Errorf("error processing modules %v", err)
		}

		if len(sourcedMods.Modules) == 0 {
			slog.Warn("no external modules found, exiting")
			return nil
		}

		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"Name","Version","Source"})
		for _,v := range sourcedMods.Modules {
			tw.AppendRow(table.Row{v.Key,v.Version,v.Source})
		}
		fmt.Println("The following modules are being added to the mod lock file:")
		fmt.Println(tw.Render())


		//Prepare the json to look nice
		bytes, _ := json.MarshalIndent(sourcedMods, "", "  ")

		// Create the mod lock file
		slog.Debug("writing modules out to file")
		os.WriteFile(path + modFileName, bytes, os.ModePerm)

		fmt.Printf("\n\nSummary: %v modules added to mod lock file.\n\n", len(sourcedMods.Modules))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
