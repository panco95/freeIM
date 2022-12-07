package chat

import (
	"context"
	"im/models"

	"github.com/smallnest/rpcx/client"
)

type RPC struct {
	chatSvc     *Service
	ProjectName string
}

func NewRPC(projectName string, chatSvc *Service) *RPC {
	return &RPC{
		ProjectName: projectName,
		chatSvc:     chatSvc,
	}
}

type SendMessageArgs struct {
	Message *models.Message
}

type SendMessageReply struct {
	Successed bool
}

func (r *RPC) SendMessageListen(ctx context.Context, args *SendMessageArgs, reply *SendMessageReply) error {
	r.chatSvc.Send(ctx, args.Message)
	return nil
}

func (r *RPC) SendMessageCall(ctx context.Context, message *models.Message) error {
	addrs, err := r.chatSvc.manager.GetAllServices()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
		if err != nil {
			return err
		}
		xClient := client.NewXClient(r.ProjectName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		defer xClient.Close()

		err = xClient.Call(ctx, "SendMessageListen", &SendMessageArgs{Message: message}, &SendMessageReply{})
		if err != nil {
			return err
		}
	}

	return nil
}
