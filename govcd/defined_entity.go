/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/url"
)

// DefinedEntityType is a type for handling Runtime Defined Entities (RDE)
type DefinedEntityType struct {
	DefinedEntityType *types.DefinedEntityType
	client            *Client
}

// DefinedEntity represents an instance of a Runtime Defined Entity (RDE) definition
type DefinedEntity struct {
	DefinedEntity *types.DefinedEntity
	client        *Client
}

// CreateRDEType creates a Runtime Defined Entity type.
// Only System administrator can create RDEs.
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
// Only System administrator can retrieve RDEs.
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
func (rde *DefinedEntityType) Update(rdeToUpdate types.DefinedEntityType) error {
	client := rde.client
	if !client.IsSysAdmin {
		return fmt.Errorf("updating Runtime Defined Entity types requires System user")
	}

	if rdeToUpdate.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity type is empty")
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

	err = client.OpenApiPutItem(apiVersion, urlRef, nil, rdeToUpdate, rde.DefinedEntityType, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the receiver Runtime Defined Entity type.
// Only System administrator can delete RDEs.
func (rde *DefinedEntityType) Delete() error {
	client := rde.client
	if !client.IsSysAdmin {
		return fmt.Errorf("deleting Runtime Defined Entity types requires System user")
	}

	if rde.DefinedEntityType.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity type is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rde.DefinedEntityType.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiDeleteItem(apiVersion, urlRef, nil, nil)
	if err != nil {
		return err
	}

	rde.DefinedEntityType = &types.DefinedEntityType{}
	return nil
}

// CreateRDE creates an entity of the type of the receiver Runtime Defined Entity (RDE).
// Only System administrator can create defined entities.
func (rde *DefinedEntityType) CreateRDE(entity types.DefinedEntity) (*DefinedEntity, error) {
	client := rde.client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("creating Runtime Defined Entities requires System user")
	}

	if rde.DefinedEntityType.ID == "" {
		return nil, fmt.Errorf("ID of the receiver Runtime Defined Entity is empty")
	}

	if entity.EntityType == "" {
		return nil, fmt.Errorf("ID of the Runtime Defined Entity type is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntities
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
		client:        rde.client,
	}

	err = client.OpenApiPostItem(apiVersion, urlRef, nil, rde, result.DefinedEntity, nil)
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

	if rde.DefinedEntity.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntities
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, rde.DefinedEntity.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiPutItem(apiVersion, urlRef, nil, rdeToUpdate, rde.DefinedEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the receiver Runtime Defined Entity.
// Only System administrator can delete RDEs.
func (rde *DefinedEntity) Delete() error {
	client := rde.client
	if !client.IsSysAdmin {
		return fmt.Errorf("deleting Runtime Defined Entity types requires System user")
	}

	if rde.DefinedEntity.ID == "" {
		return fmt.Errorf("ID of the receiver Runtime Defined Entity is empty")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntities
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
