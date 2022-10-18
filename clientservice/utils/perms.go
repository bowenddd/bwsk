package utils

import "strings"

func generatePermMap(perms string) map[string]struct{}{
	permMap := make(map[string]struct{})
	permList := strings.Split(perms, ",")
	for _, perm := range permList {
		permMap[perm] = struct{}{}
	}
	return permMap
}

func CheckPerm(perms string, perm string) bool {
	permMap := generatePermMap(perms)
	_, ok := permMap[perm]
	return ok
}