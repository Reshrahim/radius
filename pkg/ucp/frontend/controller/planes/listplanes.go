/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package planes

import (
	"context"
	"fmt"
	http "net/http"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	armrpc_controller "github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	armrpc_rest "github.com/radius-project/radius/pkg/armrpc/rest"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/ucp/datamodel"
	"github.com/radius-project/radius/pkg/ucp/datamodel/converter"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

var _ armrpc_controller.Controller = (*ListPlanes)(nil)

// ListPlanes is the controller implementation to get the list of all planes regardless of type.
type ListPlanes struct {
	armrpc_controller.Operation[*datamodel.GenericPlane, datamodel.GenericPlane]
}

// NewListPlanes creates a new controller for listing all planes regardless of type.
func NewListPlanes(opts armrpc_controller.Options) (armrpc_controller.Controller, error) {
	return &ListPlanes{
		Operation: armrpc_controller.NewOperation(opts,
			armrpc_controller.ResourceOptions[datamodel.GenericPlane]{
				RequestConverter:  converter.GenericPlaneDataModelFromVersioned,
				ResponseConverter: converter.GenericPlaneDataModelToVersioned,
			},
		),
	}, nil
}

// Run() queries the database client for planes in a given scope, creates a response with the results, and
// returns an OKResponse with the response. If an error occurs, it is returned.
func (e *ListPlanes) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (armrpc_rest.Response, error) {
	serviceCtx := v1.ARMRequestContextFromContext(ctx)
	logger := ucplog.FromContextOrDiscard(ctx)

	// The plane objects are all stored separately (by plane type). We need to query each type separately.
	planeTypes := []string{
		"aws",
		"azure",
		"radius",
	}

	objs := []database.Object{}
	for _, planeType := range planeTypes {
		query := database.Query{
			RootScope:    serviceCtx.ResourceID.String(),
			ResourceType: planeType,
			IsScopeQuery: true,
		}

		logger.Info(fmt.Sprintf("Listing planes of type %s in scope %s", query.ResourceType, query.RootScope))
		result, err := e.DatabaseClient().Query(ctx, query)
		if err != nil {
			return nil, err
		}

		objs = append(objs, result.Items...)
	}

	planes, err := e.createResponse(ctx, objs)
	if err != nil {
		return nil, err
	}
	var ok = armrpc_rest.NewOKResponse(planes)
	return ok, nil
}

func (p *ListPlanes) createResponse(ctx context.Context, objs []database.Object) (*v1.PaginatedList, error) {
	serviceCtx := v1.ARMRequestContextFromContext(ctx)
	items := v1.PaginatedList{}

	for _, item := range objs {
		var plane datamodel.GenericPlane
		err := item.As(&plane)
		if err != nil {
			return nil, err
		}

		versioned, err := p.ResponseConverter()(&plane, serviceCtx.APIVersion)
		if err != nil {
			return nil, err
		}

		items.Value = append(items.Value, versioned)
	}

	return &items, nil
}
