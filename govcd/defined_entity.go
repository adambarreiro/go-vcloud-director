/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// DefinedEntity is a type for handling Runtime Defined Entities (RDE)
type DefinedEntity struct {
	DefinedEntity *types.DefinedEntity
	client        *Client
}

// CreateRDE creates a Runtime Defined Entity
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

// GetRDEById gets a Runtime Defined Entity
func (vcdClient *VCDClient) GetRDEById(rdeId string) error {
	return nil
}

// GetRDEByName gets a Runtime Defined Entity
func (vcdClient *VCDClient) GetRDEByName(rdeId string) error {
	return nil
}

// GetRDEByIdOrName gets a Runtime Defined Entity
func (vcdClient *VCDClient) GetRDEByIdOrName(rdeId string) error {
	return nil
}

// UpdateRDE updates a Runtime Defined Entity
func (vcdClient *VCDClient) UpdateRDE(rde types.DefinedEntity) (*DefinedEntity, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, fmt.Errorf("updating Runtime Defined Entities requires System user")
	}
	return nil, nil
}

// DeleteRDE deletes a Runtime Defined Entity
func (vcdClient *VCDClient) DeleteRDE(id string) error {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return fmt.Errorf("deleting Runtime Defined Entities requires System user")
	}

	if id == "" {
		return fmt.Errorf("empty Runtime Defined Entity identifier")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	apiVersion, err := client.getOpenApiHighestElevatedVersion(endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return err
	}

	err = client.OpenApiDeleteItem(apiVersion, urlRef, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
