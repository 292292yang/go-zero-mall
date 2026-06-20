package jwtx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

func GetUserId(ctx context.Context) (int64, error) {
	val := ctx.Value("userId")
	if val == nil {
		return 0, fmt.Errorf("missing userId in context")
	}

	switch v := val.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("invalid userId type: %T", val)
	}
}
