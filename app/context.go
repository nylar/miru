package app

import (
	"io/ioutil"

	"github.com/nylar/miru/config"
	"github.com/nylar/miru/queue"

	rdb "github.com/dancannon/gorethink"
)

type Context struct {
	Db     *rdb.Session
	Config *config.Config
	Queues *queue.Queues
}

func NewContext() *Context {
	ctx := new(Context)
	ctx.InitQueues()
	return ctx
}

func (c *Context) LoadConfig(f string) error {
	file, err := ioutil.ReadFile(f)
	data := string(file)
	if err != nil {
		data = config.DefaultConfig
	}

	conf, err := config.LoadConfig(data)
	if err != nil {
		return err
	}

	c.Config = conf
	return nil
}

func (c *Context) Connect(host string) error {
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address:  host,
		Database: c.Config.Database.Name,
	})

	if err != nil {
		return err
	}

	c.Db = session
	return nil
}

func (c *Context) InitQueues() {
	c.Queues = queue.NewQueues()
}
