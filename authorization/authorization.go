package authorization

import (
	"strings"
	"regexp"
	"fmt"
)

const (
	Admin         = "admin"
	Developer     = "developer"
)

var Mapping = map[string][]string{
	"users-list": { Admin },
	"resource-get": { Admin, Developer },
	"resource-delete": { Admin },
	"resource-list": { Admin, Developer },
	"resource-discretion": { Admin, Developer },
}

// Compares the requested endpoint with roles associated with a user.
func ConfirmRoleMapping(role string, endpoint string) bool {

	fmt.Println(endpoint)

	rolesWithAccess := Mapping[endpoint]

	if rolesWithAccess != nil {

		for _, roleWithAccess := range rolesWithAccess {
			if role == roleWithAccess {
				return true
			}
		}
	}

	return false
}

// Decodes a request path to identify what resource is being requested.
func DecodeRolePath(path string, params map[string]string) string {

	// Remove query parameters from path
	for _, param := range params {
		path = strings.Split(path, "/"+param)[0]
	}

	re := regexp.MustCompile("{(.*?)}")
	rm := re.FindStringSubmatch(path)

	for _, field := range rm {
		path = strings.Split(path, "/"+field)[0]
	}

	path = strings.Replace(strings.TrimPrefix(path, "/"), "/", "-", -1)

	// Create role names matching endpoints
	return path
}