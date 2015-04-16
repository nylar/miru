package miru

import (
	"io/ioutil"

	rdb "github.com/dancannon/gorethink"
)

// Context holds database, configuration and queue data.
type Context struct {
	Db     *rdb.Session
	Config *Config
	Queues *Queues
}

// NewContext instantiates a new context and initialises a queue.
func NewContext() *Context {
	ctx := new(Context)
	ctx.InitQueues()
	return ctx
}

// LoadConfig reads a given file from the filesystem, if not found uses the
// default config.
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

// Connect creates a connection to the database.
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

// InitQueues initialises a new queue list.
func (c *Context) InitQueues() {
	c.Queues = NewQueues()
}
