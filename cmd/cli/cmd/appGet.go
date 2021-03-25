// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/radius/pkg/radclient"
	"github.com/spf13/cobra"
)

// getCmd command to get properties of an application
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get RAD application details",
	Long:  "Get RAD application details",
	RunE:  getApplication,
}

func init() {
	applicationCmd.AddCommand(getCmd)

	getCmd.Flags().String("name", "", "The application name")
	getCmd.MarkFlagRequired("name")
}

func getApplication(cmd *cobra.Command, args []string) error {
	applicationName, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	env, err := validateEnvironment()
	if err != nil {
		return err
	}

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		return err
	}

	radc := radclient.NewClient(env.SubscriptionID)
	radc.Authorizer = authorizer
	app, err := radc.GetApplication(cmd.Context(), env.ResourceGroup, applicationName)
	var applicationDetails []byte
	applicationDetails, err = json.MarshalIndent(app, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", applicationDetails)

	return err
}
