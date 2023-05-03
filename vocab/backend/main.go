package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jonas27/mono/vocab/backend/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

const (
	exitFail = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitFail)
	}
}

type Server struct {
	proto.VocabServiceServer
	vocabs *proto.VocabListResponse
}

// List all vocab
func (s *Server) ListVocabs(csRequest *proto.VocabListRequest, stream proto.VocabService_ListVocabsServer) error {
	log.Println("ListVocabs")
	_ = csRequest.PageSize
	// vocabs := &proto.VocabListResponse{}

	err := stream.Send(s.vocabs)
	if err != nil {
		return fmt.Errorf("error sending metric message %s", err)
	}

	return nil
}

func run() error {
	apiServer := grpc.NewServer(
	// grpc.StreamInterceptor(streamInterceptor),
	)
	s := &Server{}

	s.vocabs = &proto.VocabListResponse{
		TotalCount: 2,
		Vocab: []*proto.Vocab{
			{Word: "1",
				Description: "1",
				Translation: "1",
				Info:        []string{"1"},
			},
			{Word: "2",
				Description: "2",
				Translation: "2",
				Info:        []string{"2"},
			},
		}}

	proto.RegisterVocabServiceServer(apiServer, s)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		klog.Fatalln(err)
	}

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return fmt.Errorf("failed to serve: %v", apiServer.Serve(lis))
	})

	// gRPC web code
	grpcWebServer := grpcweb.WrapServer(
		apiServer,
		// Enable CORS
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    "localhost:8082",
	}

	klog.Infof("http server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
