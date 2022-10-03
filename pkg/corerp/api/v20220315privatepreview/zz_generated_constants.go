//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package v20220315privatepreview

const (
	moduleName = "v20220315privatepreview"
	moduleVersion = "v0.0.1"
)

// AzureIdentityKind - Identity Kind
type AzureIdentityKind string

const (
	AzureIdentityKindWorkload AzureIdentityKind = "Workload"
	AzureIdentityKindUserAssigned AzureIdentityKind = "UserAssigned"
)

// PossibleAzureIdentityKindValues returns the possible values for the AzureIdentityKind const type.
func PossibleAzureIdentityKindValues() []AzureIdentityKind {
	return []AzureIdentityKind{	
		AzureIdentityKindWorkload,
		AzureIdentityKindUserAssigned,
	}
}

// CertType - Certificate object type to be downloaded - the certificate itself, private key or public key of the certificate
type CertType string

const (
	CertTypeCertificate CertType = "certificate"
	CertTypePublickey CertType = "publickey"
	CertTypePrivatekey CertType = "privatekey"
)

// PossibleCertTypeValues returns the possible values for the CertType const type.
func PossibleCertTypeValues() []CertType {
	return []CertType{	
		CertTypeCertificate,
		CertTypePublickey,
		CertTypePrivatekey,
	}
}

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication CreatedByType = "Application"
	CreatedByTypeKey CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns the possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{	
		CreatedByTypeApplication,
		CreatedByTypeKey,
		CreatedByTypeManagedIdentity,
		CreatedByTypeUser,
	}
}

// Encoding - Encoding format. Default utf-8
type Encoding string

const (
	EncodingUTF8 Encoding = "utf-8"
	EncodingHex Encoding = "hex"
	EncodingBase64 Encoding = "base64"
)

// PossibleEncodingValues returns the possible values for the Encoding const type.
func PossibleEncodingValues() []Encoding {
	return []Encoding{	
		EncodingUTF8,
		EncodingHex,
		EncodingBase64,
	}
}

// Format - Certificate format. Default pem
type Format string

const (
	FormatPem Format = "pem"
	FormatPfx Format = "pfx"
)

// PossibleFormatValues returns the possible values for the Format const type.
func PossibleFormatValues() []Format {
	return []Format{	
		FormatPem,
		FormatPfx,
	}
}

// Kind - The kind of IAM provider to configure
type Kind string

const (
	KindAzure Kind = "azure"
)

// PossibleKindValues returns the possible values for the Kind const type.
func PossibleKindValues() []Kind {
	return []Kind{	
		KindAzure,
	}
}

// ManagedStore - Backing store for the ephemeral volume
type ManagedStore string

const (
	ManagedStoreDisk ManagedStore = "disk"
	ManagedStoreMemory ManagedStore = "memory"
)

// PossibleManagedStoreValues returns the possible values for the ManagedStore const type.
func PossibleManagedStoreValues() []ManagedStore {
	return []ManagedStore{	
		ManagedStoreDisk,
		ManagedStoreMemory,
	}
}

// Protocol - Protocol in use by the port
type Protocol string

const (
	ProtocolGrpc Protocol = "grpc"
	ProtocolHTTP Protocol = "http"
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

// PossibleProtocolValues returns the possible values for the Protocol const type.
func PossibleProtocolValues() []Protocol {
	return []Protocol{	
		ProtocolGrpc,
		ProtocolHTTP,
		ProtocolTCP,
		ProtocolUDP,
	}
}

// ProvisioningState - Provisioning state of the resource at the time the operation was called.
type ProvisioningState string

const (
	ProvisioningStateAccepted ProvisioningState = "Accepted"
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	ProvisioningStateFailed ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{	
		ProvisioningStateAccepted,
		ProvisioningStateCanceled,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateProvisioning,
		ProvisioningStateSucceeded,
		ProvisioningStateUpdating,
	}
}

// VolumeRbac - Container read/write access to the volume
type VolumeRbac string

const (
	VolumeRbacRead VolumeRbac = "read"
	VolumeRbacWrite VolumeRbac = "write"
)

// PossibleVolumeRbacValues returns the possible values for the VolumeRbac const type.
func PossibleVolumeRbacValues() []VolumeRbac {
	return []VolumeRbac{	
		VolumeRbacRead,
		VolumeRbacWrite,
	}
}

