* Added support for Defined Interfaces with client methods `VCDClient.CreateDefinedInterface`, `VCDClient.GetAllDefinedInterfaces`, 
  `VCDClient.GetDefinedInterface`, `VCDClient.GetDefinedInterfaceById` and methods to manipulate them `DefinedInterface.Update`, 
  `DefinedInterface.Delete` [GH-527]
* Added support for Runtime Defined Entity types with client methods `VCDClient.CreateRdeType`, `VCDClient.GetAllRdeTypes`,
  `VCDClient.GetRdeType`, `VCDClient.GetRdeTypeById` and methods to manipulate them `DefinedEntityType.Update`,
  `DefinedEntityType.Delete` [GH-527]
* Added support for Runtime Defined Entity instances with methods `DefinedEntityType.GetAllRdes`, `DefinedEntityType.GetRdeByName`,
  `DefinedEntityType.GetRdeById`, `DefinedEntityType.CreateRde` and methods to manipulate them `DefinedEntity.Resolve`,
  `DefinedEntity.Update`, `DefinedEntity.Delete` [GH-527]
