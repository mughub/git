// Package git contains an Endpoint implementation for the Git Git protocol.
package git

import (
	"context"
	"fmt"
	"github.com/mughub/mughub/bare"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
	"net"
	"os"
)

type endpoint struct {
	l server.Loader
}

func (e *endpoint) ListenAndServe(ctx context.Context) error {
	ls := getTCPListener(nil) // TODO: Retrieve listener from viper config
	conn, err := ls.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	refs := packp.NewReferenceUpdateRequest()
	err = refs.Decode(conn)
	fmt.Println(refs)
	return err
}

// NewEndpoint returns an Endpoint implementation which serves
// the Git Git protocol.
//
func NewEndpoint() bare.Endpoint {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return &endpoint{
		l: server.NewFilesystemLoader(osfs.New(wd)),
	}
}

func getTCPListener(cfg *viper.Viper) net.Listener {
	addr := fmt.Sprintf("%s:%d", cfg.GetString("addr"), cfg.GetInt("port"))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		l.Close()
		panic(err)
	}
	return l
}
