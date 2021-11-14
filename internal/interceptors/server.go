package interceptors

import (
	"context"
	todoapiv1 "github.com/itt-learning-groups/proto-contracts/todoapi/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func PrintRequestHeaders(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//log.Println("I'm the PrintRequestHeaders unary interceptor")

	// Question: This is *server* middleware. So should we be looking at "incoming" context or "outgoing" context?
	requestHeadersFromIncomingCtx, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("warning: server interceptor failed to read metadata from incoming context")
	}
	if requestHeadersFromIncomingCtx != nil {
		log.Printf("incoming ctx request headers: %+v", requestHeadersFromIncomingCtx)
	}

	// Let's see what's in the outgoing context...
	requestHeadersFromOutgoingCtx, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		log.Printf("warning: server interceptor failed to read metadata from outgoing context")
	}
	if requestHeadersFromOutgoingCtx != nil {
		log.Printf("outgoing ctx request headers: %+v", requestHeadersFromOutgoingCtx)
	}

	response, err := handler(ctx, request)

	// Question: Can we grab headers from response context on its way back "out" to the client? Why not? Does that make sense?
	if res, ok := response.(*todoapiv1.CreateTodoResponse); ok {
		log.Printf("response: %+v", res.GetTodo())
		log.Printf("error: %v", err)
	}

	return response, err
}

func AddCustomHeader(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//log.Println("I'm the AddCustomHeader unary interceptor")

	requestHeadersFromIncomingCtx, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("warning: server interceptor failed to read metadata from incoming context")
	}
	if requestHeadersFromIncomingCtx != nil {
		//ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("custom-header", "value")) // Why is this a bad idea?
		ctx = metadata.NewIncomingContext(ctx, metadata.Join(requestHeadersFromIncomingCtx, metadata.Pairs("custom-header", "value")))

		modifiedHeaders, _ := metadata.FromIncomingContext(ctx)
		log.Printf("modified headers: %+v", modifiedHeaders)
	}

	response, err := handler(ctx, request)
	return response, err
}
