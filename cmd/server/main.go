package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	todoapiv1 "github.com/itt-learning-groups/proto-contracts/todoapi/gen/go/v1"
	"github.com/itt-learning-groups/todoapi-grpc/cmd/server/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {
	grpcServerPort := 50051
	gatewayServerPort := 8080
	grpcServerAddress := fmt.Sprintf("localhost:%d", grpcServerPort)
	gatewayServerAddress := fmt.Sprintf("localhost:%d", gatewayServerPort)

	// grpc
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatalf("cannot listen to address %s: %v", grpcServerAddress, err)
	}

	grpcServer := grpc.NewServer()

	todoService := &services.TodoapiServer{}

	todoapiv1.RegisterTodoapiServiceServer(grpcServer, todoService)

	reflection.Register(grpcServer)

	log.Printf("serving TodoapiServer (grpc) on %s", grpcServerAddress)

	// Question: What happens if we don't serve the (underlying) grpc server as a child goroutine?
	//   Answer: The `Serve` call will block. The parent goroutine will never proceed and try to run the grpc-gateway
	//   server (unless the `Serve` call returns, which would mean the server shut down, which wouldn't work for creating
	//   the grpc-gateway server anyway).
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// grpc-gateway
	// Question: Do you see the (literal) connection between the underlying grpc server and the gateway server?
	//   Answer: The `cxn` below is the literal connection: The gateway (REST) server has to dial the underlying grpc
	//   server because it proxies (i.e. relays) all the calls it receives to that underlying server.
	// Question: What happens if we don't use the `WithInsecure` dial option here? Don't people say that part of the
	// advantage of grpc is that it runs on HTTP2, which requires TLS? What gives?
	//   Answer: We get a dial failure because we haven't configured a TLS certificate to support an https handshake when
	//   connecting to the grpc server. The TLS requirement can be bypassed. The TLS requirement for HTTP2 really only
	//   comes into play when connecting remotely across internet routers that may or may not attempt to inspect and
	//   validate HTTP headers if they aren't encrypted. And if you are connecting a client and server within the same
	//   local system, ignoring TLS is not as terrible idea as it sounds. But, more importantly, sophisticated Kubernetes
	//   implementations tend to outsource the TLS config to a service-mesh networking layer (e.g. Istio) so that
	//   individual servers don't have to worry about it.
	cxn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial grpc server at address %s", grpcServerAddress)
	}

	mux := runtime.NewServeMux()
	if err = todoapiv1.RegisterTodoapiServiceHandler(context.Background(), mux, cxn); err != nil {
		log.Fatalf("cannot register gateway handler: %v", err)
	}

	restServer := &http.Server{
		Addr:    gatewayServerAddress,
		Handler: mux,
	}

	log.Printf("serving TodoapiServer (rest) on %s", gatewayServerAddress)

	// Question: Why do Kount Go servers tend to run this `ListenAndServe` call in a child goroutine, too (like the
	// underlying grpc server `Serve` call)?
	//   Answer: This allows the parent goroutine to create an independent blocking `select` block that can execute a
	//   graceful server shutdown if the server encounters a fatal error. That means we can allow time for things like
	//   in-progress concurrent handlers to complete DB calls before the server abruptly cuts them off, decreasing the
	//   likelihood that a fatal error will result in unpredictable/non-deterministic app state.
	if err = restServer.ListenAndServe(); err != nil {
		log.Fatalf("cannot serve gateway: %v", err)
	}

	// Question: Can we "hit" our CreateTodo endpoint now on either port 50051 (as a gRPC call) or port 8080 (as a REST call)?
	//   Answer: Yep. Try it both ways, e.g. a grpc call from a grpc client like grpcurl; a REST call from curl or Postman.
	// Question: How does the way a Go server is exposed to a typical client when running on your local machine differ
	// from the way it is exposed when running in a Pod in a Kubernetes cluster? What implications does this have for
	// the way we set up grpc + grpc-gateway servers that will run in Kubernetes?
	//   Answer: The difference is that when we run a client call using e.g. grpcurl or Postman against our locally-running
	//   server, we're running it as a local client. In other words, the call is coming from inside the house. When we
	//   run the server as a Pod in K8s, the client will be a different pod in the same K8s cluster or an actual external
	//   network client. So we'll need to set up a K8s `Service` resource to expose the server to those clients. And you
	//   have one `Service` for each port, so if we want to expose both the underlying grpc server and the gateway server,
	//   we need two `Service`s. Sometimes we may want to do that, but more commonly we'll only want to expose the gateway
	//   server; in that case, we'll have only a single `Service` for the REST port. Bottom line, we have more than one
	//   potential pattern for setting up a grpc + grpc-gateway server depending on the access pattern we want to support.
}
