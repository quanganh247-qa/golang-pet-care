package perms

import (
	"fmt"
)



func CheckPermission(method []Permission, typeApi TypeApi, permissions *map[string]bool) bool {
	if method == nil || len(method) == 0 {
		return true
	}
	if len(method) > 1 {
		for _, m := range method {
			if CheckPermission([]Permission{m}, typeApi, permissions) {
				return true
			}
		}
		return false
	}
	permissionAPI := fmt.Sprintf("%s-%s", method[0], typeApi)
	allPermission := fmt.Sprintf("ALL-%s", typeApi)
	if  (typeApi) == "" {
		return true
	}
	if permissions == nil {
		return false
	}
	if (*permissions)[permissionAPI] || (*permissions)[allPermission] {
		return true
	}
	fullPerms := "ALL-ALL"
	halfAllPerms := fmt.Sprintf("%s-ALL", method[0])
	if (*permissions)[fullPerms] || (*permissions)[halfAllPerms] {
		return true
	}
	return false
}
