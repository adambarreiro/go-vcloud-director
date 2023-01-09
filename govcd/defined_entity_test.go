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

	var unmarshaledJson map[string]interface{}
	err := json.Unmarshal(dummyRdeSchema, &unmarshaledJson)
	check.Assert(err, IsNil)

	dummyRdeType := &types.DefinedEntityType{
		Name:        check.TestName() + "_name",
		Namespace:   check.TestName() + "_nss",
		Version:     "1.2.3",
		Description: "Description of" + check.TestName(),
		Schema:      unmarshaledJson,
		Vendor:      "vmware",
		Interfaces:  []string{"urn:vcloud:interface:vmware:k8s:1.0.0"},
		IsReadOnly:  true,
	}

	allRdeTypes, err := vcd.client.GetAllRdeTypes(nil)
	check.Assert(err, IsNil)
	alreadyPresentRdes := len(allRdeTypes)

	newRdeType, err := vcd.client.CreateRdeType(dummyRdeType)
	check.Assert(err, IsNil)
	check.Assert(newRdeType, NotNil)
	check.Assert(newRdeType.DefinedEntityType.Name, Equals, dummyRdeType.Name)
	check.Assert(newRdeType.DefinedEntityType.Namespace, Equals, dummyRdeType.Namespace)
	check.Assert(newRdeType.DefinedEntityType.Version, Equals, dummyRdeType.Version)
	check.Assert(newRdeType.DefinedEntityType.Schema, NotNil)
	check.Assert(newRdeType.DefinedEntityType.Schema["type"], Equals, "object")
	check.Assert(newRdeType.DefinedEntityType.Schema["definitions"], NotNil)
	check.Assert(newRdeType.DefinedEntityType.Schema["required"], NotNil)
	check.Assert(newRdeType.DefinedEntityType.Schema["properties"], NotNil)

	AddToCleanupListOpenApi(newRdeType.DefinedEntityType.ID, check.TestName(), types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEntityTypes+newRdeType.DefinedEntityType.ID)

	allRdeTypes, err = vcd.client.GetAllRdeTypes(nil)
	check.Assert(err, IsNil)
	check.Assert(len(allRdeTypes), Equals, alreadyPresentRdes+1)

	obtainedRdeType, err := vcd.client.GetRdeTypeById(newRdeType.DefinedEntityType.ID)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRdeType.DefinedEntityType, DeepEquals, *newRdeType.DefinedEntityType)

	obtainedRdeType2, err := vcd.client.GetRdeType(obtainedRdeType.DefinedEntityType.Vendor, obtainedRdeType.DefinedEntityType.Namespace, obtainedRdeType.DefinedEntityType.Version)
	check.Assert(err, IsNil)
	check.Assert(*obtainedRdeType2.DefinedEntityType, DeepEquals, *obtainedRdeType.DefinedEntityType)

	dummyRdeEntity := []byte(`
	{
		"foo": {
			"key": "stringValue"
		},
		"bar": "stringValue2",
		"prop2": {
			"subprop1": "stringValue3",
			"subprop2": [
				"stringValue4",
				"stringValue5"
			]
		}
	}`)

	err = json.Unmarshal(dummyRdeEntity, &unmarshaledJson)
	check.Assert(err, IsNil)

	rde, err := obtainedRdeType.CreateRde(types.DefinedEntity{
		Name:       "dummyRdeType",
		ExternalId: "123",
		Entity:     unmarshaledJson,
	})
	check.Assert(err, IsNil)
	check.Assert(rde.DefinedEntity.Name, Equals, "dummyRdeType")

	AddToCleanupListOpenApi(rde.DefinedEntity.ID, check.TestName(), types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEntities+rde.DefinedEntity.ID)

	deletedId := newRdeType.DefinedEntityType.ID
	err = newRdeType.Delete()
	check.Assert(err, IsNil)
	check.Assert(*newRdeType.DefinedEntityType, DeepEquals, types.DefinedEntityType{})

	_, err = vcd.client.GetRdeTypeById(deletedId)
	check.Assert(err, NotNil)
	check.Assert(strings.Contains(err.Error(), ErrorEntityNotFound.Error()), Equals, true)
}
