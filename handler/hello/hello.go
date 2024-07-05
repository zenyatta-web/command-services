package handler

import (
	"context"
	pb "zenyatta-web/command-services/proto/hello"
)

type HelloWord struct {
}

func (h *HelloWord) HelloWord(ctx context.Context, in *pb.Empty, out *pb.HelloWordResponse) error {
	out.Hello = "Hello Word"
	return nil
}
