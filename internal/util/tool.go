package util

import (
	"regexp"
	"time"
)

var LocationShanghai *time.Location

func init() {
	LocationShanghai, _ = time.LoadLocation("Asia/Shanghai")
}

func IsValidName(name string) bool {
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

var reTag = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func IsValidTag(tag string) bool {
	return reTag.Match([]byte(tag))
}
