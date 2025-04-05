package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func SanitizeTitle(title string) string {
	sanitized := strings.ReplaceAll(title, " ", "_")
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized = re.ReplaceAllString(sanitized, "")
	return sanitized
}

func ParseUintFromStr(uintStr string) (uint, error) {
	paramU64, err := strconv.ParseUint(uintStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", uintStr, err)
	}

	const maxUint = ^uint(0)
	if paramU64 > uint64(maxUint) {
		return 0, fmt.Errorf("%s exceeds uint max value", uintStr)
	}

	return uint(paramU64), nil
}
