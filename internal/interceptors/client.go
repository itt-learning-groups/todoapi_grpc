package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func PrintHeaders(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//log.Println("client PrintHeaders interceptor works!")

	hmd := &metadata.MD{}
	opts = append(opts, grpc.Header(hmd))

	tmd := &metadata.MD{}
	opts = append(opts, grpc.Trailer(tmd))

	err := invoker(ctx, method, req, reply, cc, opts...)

	log.Printf("response headers: %+v", hmd)
	log.Printf("response trailers: %+v", tmd)

	return err
}