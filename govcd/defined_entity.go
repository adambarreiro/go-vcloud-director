/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/url"
)

// DefinedEntity is a type for handling Runtime Defined Entities (RDE)
type DefinedEntity struct {
	DefinedEntity *types.DefinedEntity
	client        *Client
}

// CreateRDE creates a Runtime Defined Entity.
// Only System administrator can create RDEs.
func (vcdClient *VCDClient) CreateRDE(rde *types.DefinedEntity) (*DefinedEntity, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("creating Runtime Defined Entities requires System user")
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

	result := &DefinedEntity{
		DefinedEntity: &types.DefinedEntity{},
		client:        &vcdClient.Client,
	}

	err = client.OpenApiPostItem(apiVersion, urlRef, nil, rde, result.DefinedEntity, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllRDEs retrieves all Runtime Defined Entities. Query parameters can be supplied to perform additional filtering.
// Only System administrator can retrieve RDEs.
func (vcdClient *VCDClient) GetAllRDEs(queryParameters url.Values) ([]*DefinedEntity, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("getting Runtime Defined Entities requires System user")
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

	typeResponses := []*types.DefinedEntity{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParameters, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into DefinedEntity types with client
	returnRDEs := make([]*DefinedEntity, len(typeResponses))
	for sliceIndex := range typeResponses {
		returnRDEs[sliceIndex] = &DefinedEntity{
			DefinedEntity: typeResponses[sliceIndex],
			client:        &vcdClient.Client,
		}
	}

	return returnRDEs, nil
}

// GetRDEById gets a Runtime Defined Entity by its ID
// Only System administrator can retrieve RDEs.
func (vcdClient *VCDClient) GetRDEById(id string) (*DefinedEntity, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("getting Runtime Defined Entities requires System user")
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

	result := &DefinedEntity{
		DefinedEntity: &types.DefinedEntity{},
		client:        &vcdClient.Client,
	}

	err = client.OpenApiGetItem(apiVersion, urlRef, nil, result.DefinedEntity, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates the receiver Runtime Defined Entity with the values given by the input.
// Only System administrator can update RDEs.
func (rde *DefinedEntity) Update(rdeToUpdate types.DefinedEntity) error {
	client := rde.client
	if !client.IsSysAdmin {
		return fmt.Errorf("updating Runtime Defined Entities requires System user")
	}

	if rdeToUpdate.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rdeToUpdate.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiPutItem(apiVersion, urlRef, nil, rdeToUpdate, rde.DefinedEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the receiver runtime defined entity.
// Only System administrator can delete RDEs.
func (rde *DefinedEntity) Delete() error {
	client := rde.client
	if !client.IsSysAdmin {
		return fmt.Errorf("deleting Runtime Defined Entities requires System user")
	}

	if rde.DefinedEntity.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rde.DefinedEntity.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiDeleteItem(apiVersion, urlRef, nil, nil)
	if err != nil {
		return err
	}

	rde.DefinedEntity = &types.DefinedEntity{}
	return nil
}
