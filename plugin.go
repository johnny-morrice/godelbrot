package godelbrot

import (
	"image"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type RenderResult struct {
	Image *image.NRGBA
	Error error
}

type RemoteRender interface {
	Render(*WireInfo) RenderResult
}

type GodelPlugin struct {
	R RemoteRender
}

var _ plugin.Plugin = (*GodelPlugin)(nil)

func (p *GodelPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	s := &GodelRPCServer{}
	s.R = p.R
	return s, nil
}

func (p *GodelPlugin) Client(b *plugin.MuxBroker, c *rpc.Client)  (interface {}, error) {
	cli := &GodelRPC{}
	cli.Client = c

	return cli, nil
}

type GodelRPCServer struct {
	R RemoteRender
}

func (s *GodelRPCServer) Render(args *WireInfo, resp *RenderResult) error {
	*resp = s.R.Render(args)
	return nil
}

type GodelRPC struct {
	Client *rpc.Client
}

func (c *GodelRPC) Render(info *WireInfo) RenderResult {
	var resp RenderResult
	err := c.Client.Call("Plugin.Render", info, &resp)

	if err != nil {
		resp.Error = err
		return resp
	}

	return resp
}