package callback

import (
	"net/http"

	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/utils/rpc"
)

const moduleName = "callback"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	header := http.Header{}
	addr := ""
	if c != nil {
		conf := Config{}
		if err := c.Unmarshal(&conf); err != nil {
			return err
		}
		addr = conf.Addr
	}
	callService = &CallbackService{
		addr:   addr,
		client: rpc.NewClientHeader(header),
	}
	return nil
}
