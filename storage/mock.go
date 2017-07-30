package storage

import (
	"encoding/json"
)

var permissionItemsMock PermissionItems = nil
var permissionParentsMock PermissionItems = nil
var assignmentsMock Assignments = nil

// from News_Permissions_Cache_AuthData, warmup (_data['items'])
func GetPermissionItemsMock() PermissionItems {
	if permissionItemsMock == nil {
		// get from file
		data := []byte("")
		err := json.Unmarshal(data, permissionItemsMock)
		if err != nil {
			panic("Mock permissions data decoding error.")
		}
	}

	return permissionItemsMock
}

// from News_Permissions_Cache_AuthData, warmup (_data['parents'])
func GetPermissionParentsMock() PermissionItems {
	if permissionParentsMock == nil {
		// get form file
		data := []byte("")
		err := json.Unmarshal(data, permissionParentsMock)
		if err != nil {
			panic("Mock permission parents data decoding error.")
		}
	}

	return permissionParentsMock
}

// from News_Permissions_Cache_AuthData, warmup (_data['assignments'])
func GetAssignmentsMock() Assignments {
	if assignmentsMock == nil {
		// get from file
		data := []byte("")
		err := json.Unmarshal(data, assignmentsMock)
		if err != nil {
			panic("Mock permission assignments decoding error.")
		}
	}

	return assignmentsMock
}
