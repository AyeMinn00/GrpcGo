package main

import (
	"context"
	proto "github.com/gosample/pb"
	"log"
)

type sampleRouteGuideServer struct {
	proto.UnimplementedRouteGuideServer
}

func (server *sampleRouteGuideServer) GetFeature(ctx context.Context, point *proto.Point) (*proto.Feature, error) {
	return &proto.Feature{
		Name:     "Aye Min Oo",
		Location: point,
	}, nil
}

func (server *sampleRouteGuideServer) RouteChat(in *proto.RouteNote, stream proto.RouteGuide_RouteChatServer) error {
	log.Printf("received %s", in.Message)
	response := proto.RouteNote{
		Location: in.Location,
		Message:  "Hello " + in.Message,
	}
	if err := stream.Send(&response); err != nil {
		return err
	}
	return nil
}
