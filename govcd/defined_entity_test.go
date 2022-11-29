//go:build functional || openapi || rde || ALL
// +build functional openapi rde ALL

/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"encoding/json"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_RDE(check *C) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	skipOpenApiEndpointTest(vcd, check, endpoint)
	// TODO: Skip if not admin!

	schema := []byte(`
	{
		"definitions": {
			"foo": {
				"type": "object",
				"description": "Foo definition",
				"properties": {
					"key": {
						"type": "string",
						"description": "Key for foo"
					}
				}
			}
		},
		"type": "object",
		"required": [
			"foo"
		],
		"properties": {
			"bar": {
				"type": "string",
				"description": "Bar"
			},
			"prop2": {
				"type": "object",
				"properties": {
					"subprop1": {
						"type": "string"
					},
					"subprop2": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				}
			},
			"foo": {
				"$ref": "#/definitions/foo"
			}
		}
	}`)

	var jsonSchema map[string]any
	err := json.Unmarshal(schema, &jsonSchema)
	check.Assert(err, IsNil)

	rde := &types.DefinedEntity{
		Name:        check.TestName() + "_name",
		Namespace:   check.TestName() + "_nss",
		Version:     "1.2.3",
		Description: "Description of" + check.TestName(),
		Schema:      jsonSchema,
		Vendor:      "vmware",
		Interfaces:  []string{"urn:vcloud:interface:vmware:k8s:1.0.0"},
		IsReadOnly:  true,
	}

	newRde, err := vcd.client.CreateRDE(rde)
	check.Assert(err, IsNil)
	check.Assert(newRde, NotNil)
	check.Assert(newRde.DefinedEntity.Name, Equals, rde.Name)
	check.Assert(newRde.DefinedEntity.Namespace, Equals, rde.Namespace)
	check.Assert(newRde.DefinedEntity.Version, Equals, rde.Version)
	check.Assert(newRde.DefinedEntity.Schema, NotNil)
	check.Assert(newRde.DefinedEntity.Schema.(map[string]any)["type"], Equals, "object")
	check.Assert(newRde.DefinedEntity.Schema.(map[string]any)["definitions"], NotNil)
	check.Assert(newRde.DefinedEntity.Schema.(map[string]any)["required"], NotNil)
	check.Assert(newRde.DefinedEntity.Schema.(map[string]any)["properties"], NotNil)

	err = vcd.client.DeleteRDE(newRde.DefinedEntity.ID)
	check.Assert(err, IsNil)
}
