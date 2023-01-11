package utils

import (
	"strings"
)

type Path []string

// NewPathWithString Create a Path instance using path as source.
// path is the string that represent the path of the http request.
func NewPathWithString(path string) Path {
	output := []string{}
	if path == "" || path == "/" {
		return Path(output)
	}
	if path[0] != '/' {
		return Path(output)
	}
	output = strings.Split(path, "/")[1:]
	return Path(output)
}

// RemovePrefixes Clean up the path of prefixes, this
// is without "[/beta]/api/content-sources/v.../"
// Returns a new Path without prefixes.
func (p Path) RemovePrefixes() Path {
	output := []string(p)
	lenOutput := len(output)
	idx := 0
	if lenOutput < 4 {
		return []string{}
	}
	if output[idx] == "beta" {
		if lenOutput < 5 {
			return []string{}
		}
		idx++
	}
	if output[idx] != "api" {
		return []string{}
	}
	idx++
	if output[idx] != "content-sources" {
		return []string{}
	}
	idx++
	if output[idx][0] != 'v' {
		return []string{}
	}
	idx++
	return output[idx:]
}

// StartWithResources check if the indicated resources match
// the path. It can be combined with a previous call to RemovePrefixes
// to get a path list of items without prefixes.
// resources is a variadic argument and each item is a slice with the part of the path to match.
// Return true if some of the resources match with the starting items.
func (p Path) StartWithResources(resources ...[]string) bool {
	lenComponents := len(p)
	for _, r := range resources {
		lenResource := len(r)
		if lenComponents < lenResource {
			continue
		}
		flag := true
		for i := 0; i < lenResource; i++ {
			if p[i] != r[i] {
				flag = false
				break
			}
		}
		if flag {
			return true
		}
	}
	return false
}
