package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

const DefaultTimeDuration = 2 * time.Second

func UnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	if timeout.Seconds() == 0 {
		timeout = DefaultTimeDuration
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		done := make(chan struct{}, 0)

		deadline, ok := ctx.Deadline()
		if !ok || deadline.After(time.Now().Add(timeout)) {
			deadline = time.Now().Add(timeout)
		}

		var (
			cancel context.CancelFunc
		)
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()

		go func() {
			resp, err = handler(ctx, req)
			done <- struct{}{}
		}()

		select {
		case <-done:

		case <-ctx.Done():
			return nil, ctx.Err()
		}

		return resp, err
	}
}
