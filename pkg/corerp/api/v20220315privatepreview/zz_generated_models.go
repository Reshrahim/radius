//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package v20220315privatepreview

import "time"

// ApplicationProperties - Application properties
type ApplicationProperties struct {
	// REQUIRED; The resource id of the environment linked to application.
	Environment *string `json:"environment,omitempty"`

	// READ-ONLY; Provisioning state of the application at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// ApplicationResource - Radius Application.
type ApplicationResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; Application properties
	Properties *ApplicationProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ApplicationResourceList - The list of applications.
type ApplicationResourceList struct {
	// The link used to get the next page of applications list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of applications.
	Value []*ApplicationResource `json:"value,omitempty"`
}

// ApplicationsClientCreateOrUpdateOptions contains the optional parameters for the ApplicationsClient.CreateOrUpdate method.
type ApplicationsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsClientDeleteOptions contains the optional parameters for the ApplicationsClient.Delete method.
type ApplicationsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsClientGetOptions contains the optional parameters for the ApplicationsClient.Get method.
type ApplicationsClientGetOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsClientListByScopeOptions contains the optional parameters for the ApplicationsClient.ListByScope method.
type ApplicationsClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsClientUpdateOptions contains the optional parameters for the ApplicationsClient.Update method.
type ApplicationsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

type AzureIdentity struct {
	// REQUIRED; The client ID for workload and user assigned managed identity
	ClientID *string `json:"clientId,omitempty"`

	// REQUIRED; Identity Kind
	Kind *AzureIdentityKind `json:"kind,omitempty"`

	// The tenant ID for workload identity.
	TenantID *string `json:"tenantId,omitempty"`
}

type AzureKeyVaultVolumeProperties struct {
	// REQUIRED; The Azure AD identity settings
	Identity *AzureIdentity `json:"identity,omitempty"`

	// REQUIRED; The volume kind
	Kind *string `json:"kind,omitempty"`

	// REQUIRED; The ID of the keyvault to use for this volume resource
	Resource *string `json:"resource,omitempty"`

	// Fully qualified resource ID for the application that the volume is connected to.
	Application *string `json:"application,omitempty"`

	// The KeyVault certificates that this volume exposes
	Certificates map[string]*CertificateObjectProperties `json:"certificates,omitempty"`

	// The KeyVault keys that this volume exposes
	Keys map[string]*KeyObjectProperties `json:"keys,omitempty"`

	// The KeyVault secrets that this volume exposes
	Secrets map[string]*SecretObjectProperties `json:"secrets,omitempty"`

	// READ-ONLY; Provisioning state of the Volume at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`
}

// GetVolumeProperties implements the VolumePropertiesClassification interface for type AzureKeyVaultVolumeProperties.
func (a *AzureKeyVaultVolumeProperties) GetVolumeProperties() *VolumeProperties {
	return &VolumeProperties{
		Kind: a.Kind,
		ProvisioningState: a.ProvisioningState,
		Application: a.Application,
		Status: a.Status,
	}
}

// BasicResourceProperties - Basic properties of a Radius resource.
type BasicResourceProperties struct {
	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`
}

type CertificateObjectProperties struct {
	// REQUIRED; The name of the certificate
	Name *string `json:"name,omitempty"`

	// File name when written to disk.
	Alias *string `json:"alias,omitempty"`

	// Certificate object type to be downloaded - the certificate itself, private key or public key of the certificate
	CertType *CertType `json:"certType,omitempty"`

	// Encoding format. Default utf-8
	Encoding *Encoding `json:"encoding,omitempty"`

	// Certificate format. Default pem
	Format *Format `json:"format,omitempty"`

	// Certificate version
	Version *string `json:"version,omitempty"`
}

type ConnectionProperties struct {
	// REQUIRED; The source of the connection
	Source *string `json:"source,omitempty"`
	DisableDefaultEnvVars *bool `json:"disableDefaultEnvVars,omitempty"`
	Iam *IamProperties `json:"iam,omitempty"`
}

// Container - Definition of a container.
type Container struct {
	// REQUIRED; The registry and image to download and run in your container
	Image *string `json:"image,omitempty"`

	// Dictionary of
	Env map[string]*string `json:"env,omitempty"`

	// Properties for readiness/liveness probe
	LivenessProbe HealthProbePropertiesClassification `json:"livenessProbe,omitempty"`

	// Dictionary of
	Ports map[string]*ContainerPort `json:"ports,omitempty"`

	// Properties for readiness/liveness probe
	ReadinessProbe HealthProbePropertiesClassification `json:"readinessProbe,omitempty"`

	// Dictionary of
	Volumes map[string]VolumeClassification `json:"volumes,omitempty"`
}

// ContainerPort - Specifies a listening port for the container
type ContainerPort struct {
	// REQUIRED; The listening port number
	ContainerPort *int32 `json:"containerPort,omitempty"`

	// Protocol in use by the port
	Protocol *Protocol `json:"protocol,omitempty"`

	// Specifies a route provided by this port
	Provides *string `json:"provides,omitempty"`
}

// ContainerProperties - Container properties
type ContainerProperties struct {
	// REQUIRED; Specifies resource id of the application
	Application *string `json:"application,omitempty"`

	// REQUIRED; Definition of a container.
	Container *Container `json:"container,omitempty"`

	// Dictionary of
	Connections map[string]*ConnectionProperties `json:"connections,omitempty"`

	// Extensions spec of the resource
	Extensions []ExtensionClassification `json:"extensions,omitempty"`

	// READ-ONLY; Gets the status of the container at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`
}

// ContainerResource - Container
type ContainerResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; Container properties
	Properties *ContainerProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ContainerResourceList - The list of containers.
type ContainerResourceList struct {
	// The link used to get the next page of containers list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of containers.
	Value []*ContainerResource `json:"value,omitempty"`
}

// ContainersClientCreateOrUpdateOptions contains the optional parameters for the ContainersClient.CreateOrUpdate method.
type ContainersClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ContainersClientDeleteOptions contains the optional parameters for the ContainersClient.Delete method.
type ContainersClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ContainersClientGetOptions contains the optional parameters for the ContainersClient.Get method.
type ContainersClientGetOptions struct {
	// placeholder for future optional parameters
}

// ContainersClientListByScopeOptions contains the optional parameters for the ContainersClient.ListByScope method.
type ContainersClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// ContainersClientUpdateOptions contains the optional parameters for the ContainersClient.Update method.
type ContainersClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// DaprSidecarExtension - Specifies the resource should have a Dapr sidecar injected
type DaprSidecarExtension struct {
	// REQUIRED; The Dapr appId. Specifies the identifier used by Dapr for service invocation.
	AppID *string `json:"appId,omitempty"`

	// REQUIRED; Specifies the extensions of a resource.
	Kind *string `json:"kind,omitempty"`

	// The Dapr appPort. Specifies the internal listening port for the application to handle requests from the Dapr sidecar.
	AppPort *int32 `json:"appPort,omitempty"`

	// Specifies the Dapr configuration to use for the resource.
	Config *string `json:"config,omitempty"`

	// Specifies the Dapr app-protocol to use for the resource.
	Protocol *Protocol `json:"protocol,omitempty"`

	// Specifies the resource id of a dapr.io.InvokeHttpRoute that can route traffic to this resource.
	Provides *string `json:"provides,omitempty"`
}

// GetExtension implements the ExtensionClassification interface for type DaprSidecarExtension.
func (d *DaprSidecarExtension) GetExtension() *Extension {
	return &Extension{
		Kind: d.Kind,
	}
}

// EnvironmentComputeClassification provides polymorphic access to related types.
// Call the interface's GetEnvironmentCompute() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EnvironmentCompute, *KubernetesCompute
type EnvironmentComputeClassification interface {
	// GetEnvironmentCompute returns the EnvironmentCompute content of the underlying type.
	GetEnvironmentCompute() *EnvironmentCompute
}

// EnvironmentCompute - Compute resource used by application environment resource.
type EnvironmentCompute struct {
	// REQUIRED; Type of compute resource.
	Kind *string `json:"kind,omitempty"`

	// The resource id of the compute resource for application environment.
	ResourceID *string `json:"resourceId,omitempty"`
}

// GetEnvironmentCompute implements the EnvironmentComputeClassification interface for type EnvironmentCompute.
func (e *EnvironmentCompute) GetEnvironmentCompute() *EnvironmentCompute { return e }

// EnvironmentProperties - Application environment properties
type EnvironmentProperties struct {
	// REQUIRED; Compute resource used by application environment resource.
	Compute EnvironmentComputeClassification `json:"compute,omitempty"`

	// Dictionary of
	Recipes map[string]*EnvironmentRecipeProperties `json:"recipes,omitempty"`

	// READ-ONLY; Provisioning state of the environment at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// EnvironmentRecipeProperties - Properties of a Recipe linked to an Environment.
type EnvironmentRecipeProperties struct {
	// REQUIRED; Type of the connector this recipe can be consumed by. For example: 'Applications.Connector/mongoDatabases'
	ConnectorType *string `json:"connectorType,omitempty"`

	// REQUIRED; Path to the template provided by the recipe. Currently only link to Azure Container Registry is supported.
	TemplatePath *string `json:"templatePath,omitempty"`
}

// EnvironmentResource - Application environment.
type EnvironmentResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; Application environment properties
	Properties *EnvironmentProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// EnvironmentResourceList - The list of environments.
type EnvironmentResourceList struct {
	// The link used to get the next page of environments list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of environments.
	Value []*EnvironmentResource `json:"value,omitempty"`
}

// EnvironmentsClientCreateOrUpdateOptions contains the optional parameters for the EnvironmentsClient.CreateOrUpdate method.
type EnvironmentsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsClientDeleteOptions contains the optional parameters for the EnvironmentsClient.Delete method.
type EnvironmentsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsClientGetOptions contains the optional parameters for the EnvironmentsClient.Get method.
type EnvironmentsClientGetOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsClientListByScopeOptions contains the optional parameters for the EnvironmentsClient.ListByScope method.
type EnvironmentsClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsClientUpdateOptions contains the optional parameters for the EnvironmentsClient.Update method.
type EnvironmentsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// EphemeralVolume - Specifies an ephemeral volume for a container
type EphemeralVolume struct {
	// REQUIRED; The Volume kind
	Kind *string `json:"kind,omitempty"`

	// REQUIRED; Backing store for the ephemeral volume
	ManagedStore *ManagedStore `json:"managedStore,omitempty"`

	// The path where the volume is mounted
	MountPath *string `json:"mountPath,omitempty"`
}

// GetVolume implements the VolumeClassification interface for type EphemeralVolume.
func (e *EphemeralVolume) GetVolume() *Volume {
	return &Volume{
		Kind: e.Kind,
		MountPath: e.MountPath,
	}
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorDetail `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail `json:"error,omitempty"`
}

// ExecHealthProbeProperties - Specifies the properties for readiness/liveness probe using an executable
type ExecHealthProbeProperties struct {
	// REQUIRED; Command to execute to probe readiness/liveness
	Command *string `json:"command,omitempty"`

	// REQUIRED; The HealthProbeProperties kind
	Kind *string `json:"kind,omitempty"`

	// Threshold number of times the probe fails after which a failure would be reported
	FailureThreshold *float32 `json:"failureThreshold,omitempty"`

	// Initial delay in seconds before probing for readiness/liveness
	InitialDelaySeconds *float32 `json:"initialDelaySeconds,omitempty"`

	// Interval for the readiness/liveness probe in seconds
	PeriodSeconds *float32 `json:"periodSeconds,omitempty"`
}

// GetHealthProbeProperties implements the HealthProbePropertiesClassification interface for type ExecHealthProbeProperties.
func (e *ExecHealthProbeProperties) GetHealthProbeProperties() *HealthProbeProperties {
	return &HealthProbeProperties{
		Kind: e.Kind,
		InitialDelaySeconds: e.InitialDelaySeconds,
		FailureThreshold: e.FailureThreshold,
		PeriodSeconds: e.PeriodSeconds,
	}
}

// ExtensionClassification provides polymorphic access to related types.
// Call the interface's GetExtension() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *DaprSidecarExtension, *Extension, *ManualScalingExtension
type ExtensionClassification interface {
	// GetExtension returns the Extension content of the underlying type.
	GetExtension() *Extension
}

// Extension of a resource.
type Extension struct {
	// REQUIRED; Specifies the extensions of a resource.
	Kind *string `json:"kind,omitempty"`
}

// GetExtension implements the ExtensionClassification interface for type Extension.
func (e *Extension) GetExtension() *Extension { return e }

// GatewayProperties - Gateway properties
type GatewayProperties struct {
	// REQUIRED; The resource id of the application linked to Gateway resource.
	Application *string `json:"application,omitempty"`

	// REQUIRED; Routes attached to this Gateway
	Routes []*GatewayRoute `json:"routes,omitempty"`

	// Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one: mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.
	Hostname *GatewayPropertiesHostname `json:"hostname,omitempty"`

	// Sets Gateway to not be exposed externally (no public IP address associated). Defaults to false (exposed to internet).
	Internal *bool `json:"internal,omitempty"`

	// READ-ONLY; Provisioning state of the Gateway at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; URL of the gateway resource. Readonly.
	URL *string `json:"url,omitempty" azure:"ro"`
}

// GatewayPropertiesHostname - Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one:
// mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.
type GatewayPropertiesHostname struct {
	// Specify a fully-qualified domain name: myapp.mydomain.com. Mutually exclusive with 'prefix' and will take priority if both
// are defined.
	FullyQualifiedHostname *string `json:"fullyQualifiedHostname,omitempty"`

	// Specify a prefix for the hostname: myhostname.myapp.PUBLICHOSTNAMEORIP.nip.io. Mutually exclusive with 'fullyQualifiedHostname'
// and will be overridden if both are defined.
	Prefix *string `json:"prefix,omitempty"`
}

// GatewayResource - Gateway Resource that specifies how traffic is exposed to the application.
type GatewayResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; Gateway properties
	Properties *GatewayProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// GatewayResourceList - The list of Gateways.
type GatewayResourceList struct {
	// The link used to get the next page of Gateways list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of Gateways.
	Value []*GatewayResource `json:"value,omitempty"`
}

type GatewayRoute struct {
	// The HttpRoute to route to. Ex - myserviceroute.id.
	Destination *string `json:"destination,omitempty"`

	// The path to match the incoming request path on. Ex - /myservice.
	Path *string `json:"path,omitempty"`

	// Optionally update the prefix when sending the request to the service. Ex - replacePrefix: '/' and path: '/myservice' will
// transform '/myservice/myroute' to '/myroute'
	ReplacePrefix *string `json:"replacePrefix,omitempty"`
}

// GatewaysClientCreateOrUpdateOptions contains the optional parameters for the GatewaysClient.CreateOrUpdate method.
type GatewaysClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// GatewaysClientDeleteOptions contains the optional parameters for the GatewaysClient.Delete method.
type GatewaysClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// GatewaysClientGetOptions contains the optional parameters for the GatewaysClient.Get method.
type GatewaysClientGetOptions struct {
	// placeholder for future optional parameters
}

// GatewaysClientListByScopeOptions contains the optional parameters for the GatewaysClient.ListByScope method.
type GatewaysClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// GatewaysClientUpdateOptions contains the optional parameters for the GatewaysClient.Update method.
type GatewaysClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// HTTPGetHealthProbeProperties - Specifies the properties for readiness/liveness probe using HTTP Get
type HTTPGetHealthProbeProperties struct {
	// REQUIRED; The listening port number
	ContainerPort *int32 `json:"containerPort,omitempty"`

	// REQUIRED; The HealthProbeProperties kind
	Kind *string `json:"kind,omitempty"`

	// REQUIRED; The route to make the HTTP request on
	Path *string `json:"path,omitempty"`

	// Threshold number of times the probe fails after which a failure would be reported
	FailureThreshold *float32 `json:"failureThreshold,omitempty"`

	// Custom HTTP headers to add to the get request
	Headers map[string]*string `json:"headers,omitempty"`

	// Initial delay in seconds before probing for readiness/liveness
	InitialDelaySeconds *float32 `json:"initialDelaySeconds,omitempty"`

	// Interval for the readiness/liveness probe in seconds
	PeriodSeconds *float32 `json:"periodSeconds,omitempty"`
}

// GetHealthProbeProperties implements the HealthProbePropertiesClassification interface for type HTTPGetHealthProbeProperties.
func (h *HTTPGetHealthProbeProperties) GetHealthProbeProperties() *HealthProbeProperties {
	return &HealthProbeProperties{
		Kind: h.Kind,
		InitialDelaySeconds: h.InitialDelaySeconds,
		FailureThreshold: h.FailureThreshold,
		PeriodSeconds: h.PeriodSeconds,
	}
}

// HTTPRouteProperties - HTTP Route properties
type HTTPRouteProperties struct {
	// REQUIRED; The resource id of the application linked to HTTP Route resource.
	Application *string `json:"application,omitempty"`

	// The internal hostname accepting traffic for the HTTP Route. Readonly.
	Hostname *string `json:"hostname,omitempty"`

	// The port number for the HTTP Route. Defaults to 80. Readonly.
	Port *int32 `json:"port,omitempty"`

	// The scheme used for traffic. Readonly.
	Scheme *string `json:"scheme,omitempty"`

	// A stable URL that that can be used to route traffic to a resource. Readonly.
	URL *string `json:"url,omitempty"`

	// READ-ONLY; Provisioning state of the HTTP Route at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`
}

// HTTPRouteResource - Radius HTTP Route Resource.
type HTTPRouteResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; HTTP Route properties
	Properties *HTTPRouteProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// HTTPRouteResourceList - The list of HTTP Routes.
type HTTPRouteResourceList struct {
	// The link used to get the next page of HTTP Routes list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of HTTP Route.
	Value []*HTTPRouteResource `json:"value,omitempty"`
}

// HTTPRoutesClientCreateOrUpdateOptions contains the optional parameters for the HTTPRoutesClient.CreateOrUpdate method.
type HTTPRoutesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesClientDeleteOptions contains the optional parameters for the HTTPRoutesClient.Delete method.
type HTTPRoutesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesClientGetOptions contains the optional parameters for the HTTPRoutesClient.Get method.
type HTTPRoutesClientGetOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesClientListByScopeOptions contains the optional parameters for the HTTPRoutesClient.ListByScope method.
type HTTPRoutesClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesClientUpdateOptions contains the optional parameters for the HTTPRoutesClient.Update method.
type HTTPRoutesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// HealthProbePropertiesClassification provides polymorphic access to related types.
// Call the interface's GetHealthProbeProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ExecHealthProbeProperties, *HTTPGetHealthProbeProperties, *HealthProbeProperties, *TCPHealthProbeProperties
type HealthProbePropertiesClassification interface {
	// GetHealthProbeProperties returns the HealthProbeProperties content of the underlying type.
	GetHealthProbeProperties() *HealthProbeProperties
}

// HealthProbeProperties - Properties for readiness/liveness probe
type HealthProbeProperties struct {
	// REQUIRED; The HealthProbeProperties kind
	Kind *string `json:"kind,omitempty"`

	// Threshold number of times the probe fails after which a failure would be reported
	FailureThreshold *float32 `json:"failureThreshold,omitempty"`

	// Initial delay in seconds before probing for readiness/liveness
	InitialDelaySeconds *float32 `json:"initialDelaySeconds,omitempty"`

	// Interval for the readiness/liveness probe in seconds
	PeriodSeconds *float32 `json:"periodSeconds,omitempty"`
}

// GetHealthProbeProperties implements the HealthProbePropertiesClassification interface for type HealthProbeProperties.
func (h *HealthProbeProperties) GetHealthProbeProperties() *HealthProbeProperties { return h }

type IamProperties struct {
	// REQUIRED; The kind of IAM provider to configure
	Kind *Kind `json:"kind,omitempty"`

	// RBAC permissions to be assigned on the source resource
	Roles []*string `json:"roles,omitempty"`
}

type KeyObjectProperties struct {
	// REQUIRED; The name of the key
	Name *string `json:"name,omitempty"`

	// File name when written to disk.
	Alias *string `json:"alias,omitempty"`

	// Key version
	Version *string `json:"version,omitempty"`
}

// KubernetesCompute - Specifies the properties for Kubernetes compute environment
type KubernetesCompute struct {
	// REQUIRED; Type of compute resource.
	Kind *string `json:"kind,omitempty"`

	// REQUIRED; The namespace to use for the environment.
	Namespace *string `json:"namespace,omitempty"`

	// The resource id of the compute resource for application environment.
	ResourceID *string `json:"resourceId,omitempty"`
}

// GetEnvironmentCompute implements the EnvironmentComputeClassification interface for type KubernetesCompute.
func (k *KubernetesCompute) GetEnvironmentCompute() *EnvironmentCompute {
	return &EnvironmentCompute{
		Kind: k.Kind,
		ResourceID: k.ResourceID,
	}
}

// ManualScalingExtension - ManualScaling Extension
type ManualScalingExtension struct {
	// REQUIRED; Specifies the extensions of a resource.
	Kind *string `json:"kind,omitempty"`

	// Replica count.
	Replicas *int32 `json:"replicas,omitempty"`
}

// GetExtension implements the ExtensionClassification interface for type ManualScalingExtension.
func (m *ManualScalingExtension) GetExtension() *Extension {
	return &Extension{
		Kind: m.Kind,
	}
}

// PersistentVolume - Specifies a persistent volume for a container
type PersistentVolume struct {
	// REQUIRED; The Volume kind
	Kind *string `json:"kind,omitempty"`

	// REQUIRED; The source of the volume
	Source *string `json:"source,omitempty"`

	// The path where the volume is mounted
	MountPath *string `json:"mountPath,omitempty"`

	// Container read/write access to the volume
	Rbac *VolumeRbac `json:"rbac,omitempty"`
}

// GetVolume implements the VolumeClassification interface for type PersistentVolume.
func (p *PersistentVolume) GetVolume() *Volume {
	return &Volume{
		Kind: p.Kind,
		MountPath: p.MountPath,
	}
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ResourceStatus - Status of a resource.
type ResourceStatus struct {
	OutputResources []map[string]interface{} `json:"outputResources,omitempty"`
}

type SecretObjectProperties struct {
	// REQUIRED; The name of the secret
	Name *string `json:"name,omitempty"`

	// File name when written to disk.
	Alias *string `json:"alias,omitempty"`

	// Encoding format. Default utf-8
	Encoding *Encoding `json:"encoding,omitempty"`

	// Secret version
	Version *string `json:"version,omitempty"`
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// TCPHealthProbeProperties - Specifies the properties for readiness/liveness probe using TCP
type TCPHealthProbeProperties struct {
	// REQUIRED; The listening port number
	ContainerPort *int32 `json:"containerPort,omitempty"`

	// REQUIRED; The HealthProbeProperties kind
	Kind *string `json:"kind,omitempty"`

	// Threshold number of times the probe fails after which a failure would be reported
	FailureThreshold *float32 `json:"failureThreshold,omitempty"`

	// Initial delay in seconds before probing for readiness/liveness
	InitialDelaySeconds *float32 `json:"initialDelaySeconds,omitempty"`

	// Interval for the readiness/liveness probe in seconds
	PeriodSeconds *float32 `json:"periodSeconds,omitempty"`
}

// GetHealthProbeProperties implements the HealthProbePropertiesClassification interface for type TCPHealthProbeProperties.
func (t *TCPHealthProbeProperties) GetHealthProbeProperties() *HealthProbeProperties {
	return &HealthProbeProperties{
		Kind: t.Kind,
		InitialDelaySeconds: t.InitialDelaySeconds,
		FailureThreshold: t.FailureThreshold,
		PeriodSeconds: t.PeriodSeconds,
	}
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VolumeClassification provides polymorphic access to related types.
// Call the interface's GetVolume() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EphemeralVolume, *PersistentVolume, *Volume
type VolumeClassification interface {
	// GetVolume returns the Volume content of the underlying type.
	GetVolume() *Volume
}

// Volume - Specifies a volume for a container
type Volume struct {
	// REQUIRED; The Volume kind
	Kind *string `json:"kind,omitempty"`

	// The path where the volume is mounted
	MountPath *string `json:"mountPath,omitempty"`
}

// GetVolume implements the VolumeClassification interface for type Volume.
func (v *Volume) GetVolume() *Volume { return v }

// VolumePropertiesClassification provides polymorphic access to related types.
// Call the interface's GetVolumeProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AzureKeyVaultVolumeProperties, *VolumeProperties
type VolumePropertiesClassification interface {
	// GetVolumeProperties returns the VolumeProperties content of the underlying type.
	GetVolumeProperties() *VolumeProperties
}

type VolumeProperties struct {
	// REQUIRED; The volume kind
	Kind *string `json:"kind,omitempty"`

	// Fully qualified resource ID for the application that the volume is connected to.
	Application *string `json:"application,omitempty"`

	// READ-ONLY; Provisioning state of the Volume at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; Status of the resource
	Status *ResourceStatus `json:"status,omitempty" azure:"ro"`
}

// GetVolumeProperties implements the VolumePropertiesClassification interface for type VolumeProperties.
func (v *VolumeProperties) GetVolumeProperties() *VolumeProperties { return v }

// VolumeResource - Radius Volume Resource.
type VolumeResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED
	Properties VolumePropertiesClassification `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VolumeResourceList - The list of Volumes.
type VolumeResourceList struct {
	// The link used to get the next page of Volumes list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of Volume.
	Value []*VolumeResource `json:"value,omitempty"`
}

// VolumesClientCreateOrUpdateOptions contains the optional parameters for the VolumesClient.CreateOrUpdate method.
type VolumesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientDeleteOptions contains the optional parameters for the VolumesClient.Delete method.
type VolumesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientGetOptions contains the optional parameters for the VolumesClient.Get method.
type VolumesClientGetOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientListByScopeOptions contains the optional parameters for the VolumesClient.ListByScope method.
type VolumesClientListByScopeOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientUpdateOptions contains the optional parameters for the VolumesClient.Update method.
type VolumesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

