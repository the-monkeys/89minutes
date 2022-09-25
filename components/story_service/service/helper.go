package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Errorf(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Errorf(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}
