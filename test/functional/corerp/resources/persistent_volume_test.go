// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package resource_test

import (
	"testing"

	"github.com/project-radius/radius/test/functional"
	"github.com/project-radius/radius/test/functional/corerp"
	"github.com/project-radius/radius/test/step"
	"github.com/project-radius/radius/test/validation"
)

func Test_PersistentVolume(t *testing.T) {
	t.Skipf("Skip until workload identity is supported")

	template := "testdata/corerp-resources-volume-azure-keyvault.bicep"
	name := "corerp-resources-volume-azkv"

	requiredSecrets := map[string]map[string]string{}

	test := corerp.NewCoreRPTest(t, name, []corerp.TestStep{
		{
			Executor: step.NewDeployExecutor(template, functional.GetMagpieImage()),
			CoreRPResources: &validation.CoreRPResourceSet{
				Resources: []validation.CoreRPResource{
					{
						Name: name,
						Type: validation.ApplicationsResource,
					},
					{
						Name: "volume-azkv-ctnr",
						Type: validation.ContainersResource,
						App:  name,
					},
					{
						Name: "volume-azkv",
						Type: validation.VolumesResource,
						App:  name,
					},
				},
			},
			K8sObjects: &validation.K8sObjectSet{
				Namespaces: map[string][]validation.K8sObject{
					"default": {
						validation.NewK8sPodForResource(name, "volume-azkv-ctnr"),
					},
				},
			},
		},
	}, requiredSecrets)

	test.Test(t)
}