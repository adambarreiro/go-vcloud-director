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
	"strings"
)

func (vcd *TestVCD) Test_Rde(check *C) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	skipOpenApiEndpointTest(vcd, check, endpoint)
	// TODO: Skip if not admin!

	dummyRdeSchema := []byte(`
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

	var jsonSchema map[string]interface{}
	err := json.Unmarshal(dummyRdeSchema, &jsonSchema)
	check.Assert(err, IsNil)

	dummyRde := &types.DefinedEntityType{
		Name:        check.TestName() + "_name",
		Namespace:   check.TestName() + "_nss",
		Version:     "1.2.3",
		Description: "Description of" + check.TestName(),
		Schema:      jsonSchema,
		Vendor:      "vmware",
		Interfaces:  []string{"urn:vcloud:interface:vmware:k8s:1.0.0"},
		IsReadOnly:  true,
	}

	allRdes, err := vcd.client.GetAllRdeTypes(nil)
	check.Assert(err, IsNil)
	alreadyPresentRdes := len(allRdes)

	newRde, err := vcd.client.CreateRdeType(dummyRde)
	check.Assert(err, IsNil)
	check.Assert(newRde, NotNil)
	check.Assert(newRde.DefinedEntityType.Name, Equals, dummyRde.Name)
	check.Assert(newRde.DefinedEntityType.Namespace, Equals, dummyRde.Namespace)
	check.Assert(newRde.DefinedEntityType.Version, Equals, dummyRde.Version)
	check.Assert(newRde.DefinedEntityType.Schema, NotNil)
	check.Assert(newRde.DefinedEntityType.Schema["type"], Equals, "object")
	check.Assert(newRde.DefinedEntityType.Schema["definitions"], NotNil)
	check.Assert(newRde.DefinedEntityType.Schema["required"], NotNil)
	check.Assert(newRde.DefinedEntityType.Schema["properties"], NotNil)

	AddToCleanupListOpenApi(newRde.DefinedEntityType.ID, check.TestName(), types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEntityTypes+newRde.DefinedEntityType.ID)

	allRdes, err = vcd.client.GetAllRdeTypes(nil)
	check.Assert(err, IsNil)
	check.Assert(len(allRdes), Equals, alreadyPresentRdes+1)

	obtainedRde, err := vcd.client.GetRdeTypeById(newRde.DefinedEntityType.ID)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRde.DefinedEntityType, DeepEquals, *newRde.DefinedEntityType)

	obtainedRde2, err := vcd.client.GetRdeType(obtainedRde.DefinedEntityType.Vendor, obtainedRde.DefinedEntityType.Namespace, obtainedRde.DefinedEntityType.Version)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRde2.DefinedEntityType, DeepEquals, *obtainedRde.DefinedEntityType)

	deletedId := newRde.DefinedEntityType.ID
	err = newRde.Delete()
	check.Assert(err, IsNil)
	check.Assert(*newRde.DefinedEntityType, DeepEquals, types.DefinedEntityType{})

	_, err = vcd.client.GetRdeTypeById(deletedId)
	check.Assert(err, NotNil)
	check.Assert(strings.Contains(err.Error(), ErrorEntityNotFound.Error()), Equals, true)
}
