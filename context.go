package miru

import (
	"io/ioutil"

	rdb "github.com/dancannon/gorethink"
)

type Context struct {
	Db     *rdb.Session
	Config *Config
	Queues *Queues
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
		data = DefaultConfig
	}

	conf, err := LoadConfig(data)
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
	c.Queues = NewQueues()
}
