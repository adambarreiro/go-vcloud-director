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

func (vcd *TestVCD) Test_RDE(check *C) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	skipOpenApiEndpointTest(vcd, check, endpoint)
	// TODO: Skip if not admin!

	dummyRDESchema := []byte(`
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
	err := json.Unmarshal(dummyRDESchema, &jsonSchema)
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

	allRDEs, err := vcd.client.GetAllRDETypes(nil)
	check.Assert(err, IsNil)
	alreadyPresentRDEs := len(allRDEs)

	newRDE, err := vcd.client.CreateRDEType(dummyRde)
	check.Assert(err, IsNil)
	check.Assert(newRDE, NotNil)
	check.Assert(newRDE.DefinedEntityType.Name, Equals, dummyRde.Name)
	check.Assert(newRDE.DefinedEntityType.Namespace, Equals, dummyRde.Namespace)
	check.Assert(newRDE.DefinedEntityType.Version, Equals, dummyRde.Version)
	check.Assert(newRDE.DefinedEntityType.Schema, NotNil)
	check.Assert(newRDE.DefinedEntityType.Schema["type"], Equals, "object")
	check.Assert(newRDE.DefinedEntityType.Schema["definitions"], NotNil)
	check.Assert(newRDE.DefinedEntityType.Schema["required"], NotNil)
	check.Assert(newRDE.DefinedEntityType.Schema["properties"], NotNil)

	AddToCleanupListOpenApi(newRDE.DefinedEntityType.ID, check.TestName(), types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEntityTypes+newRDE.DefinedEntityType.ID)

	allRDEs, err = vcd.client.GetAllRDETypes(nil)
	check.Assert(err, IsNil)
	check.Assert(len(allRDEs), Equals, alreadyPresentRDEs+1)

	obtainedRDE, err := vcd.client.GetRDETypeById(newRDE.DefinedEntityType.ID)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRDE.DefinedEntityType, DeepEquals, *newRDE.DefinedEntityType)

	obtainedRDE2, err := vcd.client.GetRDEType(obtainedRDE.DefinedEntityType.Vendor, obtainedRDE.DefinedEntityType.Namespace, obtainedRDE.DefinedEntityType.Version)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRDE2.DefinedEntityType, DeepEquals, *obtainedRDE.DefinedEntityType)

	deletedId := newRDE.DefinedEntityType.ID
	err = newRDE.Delete()
	check.Assert(err, IsNil)
	check.Assert(*newRDE.DefinedEntityType, DeepEquals, types.DefinedEntityType{})

	_, err = vcd.client.GetRDETypeById(deletedId)
	check.Assert(err, NotNil)
	check.Assert(strings.Contains(err.Error(), ErrorEntityNotFound.Error()), Equals, true)
}
