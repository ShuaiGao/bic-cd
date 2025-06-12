package util

import (
	"crypto/md5"
	"encoding/hex"
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

// MD5 md5 encryption， 长度32的md5字符串
func MD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
