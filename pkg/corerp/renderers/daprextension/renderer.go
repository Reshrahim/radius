// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package daprextension

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/project-radius/radius/pkg/armrpc/api/conv"
	connector "github.com/project-radius/radius/pkg/connectorrp/datamodel"
	"github.com/project-radius/radius/pkg/corerp/datamodel"
	"github.com/project-radius/radius/pkg/corerp/renderers"
	"github.com/project-radius/radius/pkg/providers"
	"github.com/project-radius/radius/pkg/ucp/resources"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Renderer is the renderers.Renderer implementation for the dapr sidecar extension.
type Renderer struct {
	Inner renderers.Renderer
}

// GetDependencyIDs returns dependencies for the container datamodel passed in
func (r Renderer) GetDependencyIDs(ctx context.Context, dm conv.DataModelInterface) ([]resources.ID, []resources.ID, error) {
	radiusDependencyIDs, azureDependencyIDs, err := r.Inner.GetDependencyIDs(ctx, dm)
	if err != nil {
		return nil, nil, err
	}

	extension, err := r.findExtension(dm)
	if err != nil {
		return nil, nil, err
	}

	if extension == nil {
		return radiusDependencyIDs, azureDependencyIDs, nil
	}

	if extension.Provides == "" {
		return radiusDependencyIDs, azureDependencyIDs, nil
	}

	parsed, err := resources.Parse(extension.Provides)
	if err != nil {
		return nil, nil, err
	}

	return append(radiusDependencyIDs, parsed), azureDependencyIDs, nil
}

// Render augments the container's kubernetes output resource with value for dapr sidecar extension.
func (r *Renderer) Render(ctx context.Context, dm conv.DataModelInterface, options renderers.RenderOptions) (renderers.RendererOutput, error) {
	resource, ok := dm.(datamodel.ContainerResource)
	if !ok {
		return renderers.RendererOutput{}, conv.ErrInvalidModelConversion
	}
	dependencies := options.Dependencies
	output, err := r.Inner.Render(ctx, resource, renderers.RenderOptions{Dependencies: dependencies})
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	extension, err := r.findExtension(resource)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	if extension == nil {
		return output, nil
	}

	// If we get here then we found a Dapr Sidecar extension. We need to update the Kubernetes deployment with
	// the desired annotations.

	// Resolve the AppID:
	// 1. If there's a DaprHttpRoute then it *must* specify an app id.
	// 2. The extension specifies an app id (must not conflict with 1)
	// 3. (none)

	appID, err := r.resolveAppId(extension, dependencies)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	for i := range output.Resources {
		if output.Resources[i].ResourceType.Provider != providers.ProviderKubernetes {
			// Not a Kubernetes resource
			continue
		}

		o, ok := output.Resources[i].Resource.(runtime.Object)
		if !ok {
			return renderers.RendererOutput{}, errors.New("found Kubernetes resource with non-Kubernetes payload")
		}

		annotations, ok := r.getAnnotations(o)
		if !ok {
			continue
		}

		annotations["dapr.io/enabled"] = "true"

		if appID != "" {
			annotations["dapr.io/app-id"] = appID
		}
		if appPort := extension.AppPort; appPort != 0 {
			annotations["dapr.io/app-port"] = fmt.Sprintf("%d", appPort)
		}
		if config := extension.Config; config != "" {
			annotations["dapr.io/config"] = config
		}
		if extension.Protocol != "" {
			annotations["dapr.io/protocol"] = string(extension.Protocol)
		}

		r.setAnnotations(o, annotations)
	}

	return output, nil
}

func (r *Renderer) findExtension(dm conv.DataModelInterface) (*datamodel.DaprSidecarExtension, error) {
	container, ok := dm.(datamodel.ContainerResource)
	if !ok {
		return nil, conv.ErrInvalidModelConversion
	}

	for _, t := range container.Properties.Extensions {
		switch t.Kind {
		case datamodel.DaprSidecar:
			return t.DaprSidecar, nil
		}
	}
	return nil, nil
}

func (r *Renderer) resolveAppId(extension *datamodel.DaprSidecarExtension, dependencies map[string]renderers.RendererDependency) (string, error) {
	// We're being extra pedantic here about reporting error cases. None of these
	// cases should be possible to trigger with user input, they would result from internal bugs.
	routeAppID := ""
	if extension.Provides != "" {
		routeDependency, ok := dependencies[extension.Provides]
		if !ok {
			return "", fmt.Errorf("failed to find depenendency with id %q", extension.Provides)
		}

		route := connector.DaprInvokeHttpRouteProperties{}
		err := convertDefinition(&route, routeDependency)
		if err != nil {
			return "", err
		}
		routeAppID = route.AppId
	}

	appID := extension.AppID
	if appID != "" && routeAppID != "" && appID != routeAppID {
		return "", fmt.Errorf("the appId specified on a daprInvokeHttpRoutes must match the appId specified on the extension. Route: %q, Extension: %q", routeAppID, appID)
	}

	if routeAppID != "" {
		return routeAppID, nil
	}

	return appID, nil
}

func (r *Renderer) getAnnotations(o runtime.Object) (map[string]string, bool) {
	dep, ok := o.(*appsv1.Deployment)
	if ok {
		if dep.Spec.Template.Annotations == nil {
			dep.Spec.Template.Annotations = map[string]string{}
		}

		return dep.Spec.Template.Annotations, true
	}

	un, ok := o.(*unstructured.Unstructured)
	if ok {
		if a := un.GetAnnotations(); a != nil {
			return a, true
		}

		return map[string]string{}, true
	}

	return nil, false
}

func (r *Renderer) setAnnotations(o runtime.Object, annotations map[string]string) {
	un, ok := o.(*unstructured.Unstructured)
	if ok {
		un.SetAnnotations(annotations)
	}
}

func convertDefinition(properties interface{}, r renderers.RendererDependency) error {
	b, err := json.Marshal(r.Definition)
	if err != nil {
		return fmt.Errorf("failed to marshal resource definition: %w", err)
	}

	err = json.Unmarshal(b, properties)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resource definition: %w", err)
	}

	return nil
}