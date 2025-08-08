package role_ctx

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	KeyRoleCode = "role_code"
)

func ValueFromCtx(ctx context.Context) ([]string, error) {
	if roleCode, ok := ctx.Value("roleCode").(string); ok {
		return strings.Split(roleCode, ","), nil
	}

	if meta, ok := metadata.FromIncomingContext(ctx); ok {
		if value := meta.Get(KeyRoleCode); len(value) > 0 {
			return strings.Split(value[0], ","), nil
		}
	}

	return nil, errors.New("failed to get role code from context")
}
