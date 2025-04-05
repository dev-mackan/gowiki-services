package apiserver

import (
	"fmt"
	"net/http"
	"strconv"
)

func parseUintParam(r *http.Request, param string) (uint, error) {
	paramStr := r.PathValue(param)
	paramU64, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", param, err)
	}
	const maxUint = ^uint(0)
	if paramU64 > uint64(maxUint) {
		return 0, fmt.Errorf("%s exceeds uint max value", param)
	}
	return uint(paramU64), nil
}
