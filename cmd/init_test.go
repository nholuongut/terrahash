/*
Copyright © 2024 the terrahash authors

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/
package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestInitCmd(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test case 1: Missing folder
	t.Run("MissingFolder", func(t *testing.T) {
		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"init", "--source", "missing"})
		errMissing := rootCmd.Execute()

		assert.NotNil(t, errMissing)

	})

	// Test case 2: Terraform not initialized
	t.Run("TerraformNotInitialized", func(t *testing.T) {
		// Mock the setPath function to return the temporary directory
		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"init", "--source", tempDir})
		rootCmd.Execute()

		expected := ""

		assert.Equal(t, expected, actual.String())
	})

	// Test case 2: Mod lock file already exists
	t.Run("ModLockFileExists", func(t *testing.T) {
		// Create a mod lock file in the temporary directory
		modLockFile := tempDir + "/" + modFileName
		file, _ := os.Create(modLockFile)

		defer os.Remove(modLockFile)
		// TODO: Add assertions for the expected output
		file.WriteString(`{"Modules": {}}`)

		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"init", "--source", tempDir})
		rootCmd.Execute()

		expected := ""

		//Assert no error
		assert.Equal(t, expected, actual.String())
	})

	// Test case 3: No external modules found
	t.Run("NoExternalModulesFound", func(t *testing.T) {
		//Create an empty Terraform configuration
		terraformConfig := tempDir + "/main.tf"
		file, _ := os.Create(terraformConfig)
		defer os.Remove(terraformConfig)

		file.WriteString(`locals { testing = "test" }`)

		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: tempDir,})

		terraform.Init(t, terraformOptions)

		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"init", "--source", tempDir})
		rootCmd.Execute()

		expected := ""

		// Assert no errors
		assert.Equal(t, expected, actual.String())
		// Assert mod lock file is not created
		_, err := os.Stat(tempDir + "/" + modFileName)
		assert.NotNil(t, err)
	})

	// Test case 4: External modules found
	t.Run("ExternalModulesFound", func(t *testing.T) {
		
		//Create Terraform configuration with external module
		terraformConfig := tempDir + "/main.tf"
		file, _ := os.Create(terraformConfig)
		defer os.Remove(terraformConfig)

		file.WriteString(testConfig)

		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: tempDir,})

		terraform.Init(t, terraformOptions)

		actual := new(bytes.Buffer)
		rootCmd.SetOut(actual)
		rootCmd.SetErr(actual)
		rootCmd.SetArgs([]string{"init", "--source", tempDir})
		rootCmd.Execute()

		expected := ""

		// Assert no errors
		assert.Equal(t, expected, actual.String())

		//Assert mod lock file is created
		_, err := os.Stat(tempDir + "/" + modFileName)
		assert.Nil(t, err)
	})
}