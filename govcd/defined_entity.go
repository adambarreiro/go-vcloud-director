/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// DefinedEntity represents an instance of a Runtime Defined Entity (RDE)
type DefinedEntity struct {
	DefinedEntity *types.DefinedEntity
	client        *Client
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
