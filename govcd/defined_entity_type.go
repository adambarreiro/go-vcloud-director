/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/url"
)

// DefinedEntityType is a type for handling Runtime Defined Entity (RDE) type definitions
type DefinedEntityType struct {
	DefinedEntityType *types.DefinedEntityType
	client            *Client
}

// CreateRDEType creates a Runtime Defined Entity type.
// Only System administrator can create RDE types.
func (vcdClient *VCDClient) CreateRDEType(rde *types.DefinedEntityType) (*DefinedEntityType, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("creating Runtime Defined Entity types requires System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	result := &DefinedEntityType{
		DefinedEntityType: &types.DefinedEntityType{},
		client:            &vcdClient.Client,
	}

	err = client.OpenApiPostItem(apiVersion, urlRef, nil, rde, result.DefinedEntityType, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllRDETypes retrieves all Runtime Defined Entity types. Query parameters can be supplied to perform additional filtering.
// Only System administrator can retrieve RDE types.
func (vcdClient *VCDClient) GetAllRDETypes(queryParameters url.Values) ([]*DefinedEntityType, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("getting Runtime Defined Entity types requires System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.DefinedEntityType{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParameters, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into DefinedEntityType types with client
	returnRDEs := make([]*DefinedEntityType, len(typeResponses))
	for sliceIndex := range typeResponses {
		returnRDEs[sliceIndex] = &DefinedEntityType{
			DefinedEntityType: typeResponses[sliceIndex],
			client:            &vcdClient.Client,
		}
	}

	return returnRDEs, nil
}

// GetRDEType gets a Runtime Defined Entity type by its unique combination of vendor, namespace and version.
// Only System administrator can retrieve RDE types.
func (vcdClient *VCDClient) GetRDEType(vendor, namespace, version string) (*DefinedEntityType, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("getting Runtime Defined Entity types requires System user")
	}

	queryParameters := url.Values{}
	queryParameters.Add("filter", fmt.Sprintf("vendor==%s;nss==%s;version==%s", vendor, namespace, version))
	rdeTypes, err := vcdClient.GetAllRDETypes(queryParameters)
	if err != nil {
		return nil, err
	}

	if len(rdeTypes) == 0 {
		return nil, fmt.Errorf("%s could not find the Runtime Defined Entity type with vendor %s, namespace %s and version %s", ErrorEntityNotFound, vendor, namespace, version)
	}

	if len(rdeTypes) > 1 {
		return nil, fmt.Errorf("found more than 1 Runtime Defined Entity type with vendor %s, namespace %s and version %s", vendor, namespace, version)
	}

	return rdeTypes[0], nil
}

// GetRDETypeById gets a Runtime Defined Entity type by its ID
// Only System administrator can retrieve RDEs.
func (vcdClient *VCDClient) GetRDETypeById(id string) (*DefinedEntityType, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("getting Runtime Defined Entity types requires System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	result := &DefinedEntityType{
		DefinedEntityType: &types.DefinedEntityType{},
		client:            &vcdClient.Client,
	}

	err = client.OpenApiGetItem(apiVersion, urlRef, nil, result.DefinedEntityType, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates the receiver Runtime Defined Entity type with the values given by the input.
// Only System administrator can update RDEs.
func (rdeType *DefinedEntityType) Update(rdeToUpdate types.DefinedEntityType) error {
	client := rdeType.client
	if !client.IsSysAdmin {
		return fmt.Errorf("updating Runtime Defined Entity types requires System user")
	}

	if rdeType.DefinedEntityType.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity type is empty")
	}

	if rdeToUpdate.ID != "" && rdeToUpdate.ID != rdeType.DefinedEntityType.ID {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity and the input ID don't match")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rdeType.DefinedEntityType.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiPutItem(apiVersion, urlRef, nil, rdeToUpdate, rdeType.DefinedEntityType, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the receiver Runtime Defined Entity type.
// Only System administrator can delete RDEs.
func (rdeType *DefinedEntityType) Delete() error {
	client := rdeType.client
	if !client.IsSysAdmin {
		return fmt.Errorf("deleting Runtime Defined Entity types requires System user")
	}

	if rdeType.DefinedEntityType.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity type is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rdeType.DefinedEntityType.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiDeleteItem(apiVersion, urlRef, nil, nil)
	if err != nil {
		return err
	}

	rdeType.DefinedEntityType = &types.DefinedEntityType{}
	return nil
}
