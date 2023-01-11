//go:build functional || openapi || rde || ALL
// +build functional openapi rde ALL

/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"encoding/json"
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (vcd *TestVCD) Test_Rde(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointEntityTypes
	skipOpenApiEndpointTest(vcd, check, endpoint)

	// Load the RDE type schema
	rdeFilePath := "../test-resources/rde_type.json"
	_, err := os.Stat(rdeFilePath)
	if os.IsNotExist(err) {
		check.Skip(fmt.Sprintf("unable to find RDE type file '%s': %s", rdeFilePath, err))
	}

	rdeFile, err := os.OpenFile(filepath.Clean(rdeFilePath), os.O_RDONLY, 0400)
	if err != nil {
		check.Fatalf("unable to open RDE type file '%s': %s", rdeFilePath, err)
	}
	defer safeClose(rdeFile)

	rdeSchema, err := io.ReadAll(rdeFile)
	if err != nil {
		check.Fatalf("error reading RDE type file %s: %s", rdeFilePath, err)
	}

	var unmarshaledJson map[string]interface{}
	err = json.Unmarshal(rdeSchema, &unmarshaledJson)
	check.Assert(err, IsNil)

	dummyRdeType := &types.DefinedEntityType{
		Name:        "name1",
		Namespace:   "namespace9", // Can't put check.TestName() due to a bug that causes RDEs to fail on GET once created.
		Version:     "1.2.3",
		Description: "Description of " + check.TestName(),
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

	unmarshaledJson = make(map[string]interface{})
	err = json.Unmarshal(dummyRdeEntity, &unmarshaledJson)
	check.Assert(err, IsNil)

	rde, err := obtainedRdeType.CreateRde(types.DefinedEntity{
		Name:   "dummyRdeType",
		Entity: unmarshaledJson,
	})
	check.Assert(err, IsNil)
	check.Assert(rde.DefinedEntity.Name, Equals, "dummyRdeType")
	check.Assert(*rde.DefinedEntity.State, Equals, "PRE_CREATED")

	err = rde.Resolve()
	check.Assert(err, IsNil)
	check.Assert(*rde.DefinedEntity.State, Equals, "RESOLVED")

	// The RDE can't be deleted until rde.Resolve() is called
	AddToCleanupListOpenApi(rde.DefinedEntity.ID, check.TestName(), types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEntities+rde.DefinedEntity.ID)

	// Delete the RDE instance prior to the RDE type deletion
	deletedId := rde.DefinedEntity.ID
	err = rde.Delete()
	check.Assert(err, IsNil)
	check.Assert(*rde.DefinedEntity, DeepEquals, types.DefinedEntity{})

	_, err = newRdeType.GetRdeById(deletedId)
	check.Assert(err, NotNil)
	check.Assert(strings.Contains(err.Error(), ErrorEntityNotFound.Error()), Equals, true)

	// Delete the RDE type last
	deletedId = newRdeType.DefinedEntityType.ID
	err = newRdeType.Delete()
	check.Assert(err, IsNil)
	check.Assert(*newRdeType.DefinedEntityType, DeepEquals, types.DefinedEntityType{})

	_, err = vcd.client.GetRdeTypeById(deletedId)
	check.Assert(err, NotNil)
	check.Assert(strings.Contains(err.Error(), ErrorEntityNotFound.Error()), Equals, true)
}
