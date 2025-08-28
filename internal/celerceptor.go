package internal

import (
	"os"
	"sync"
)

type Celerceptor struct {
	sync.Mutex

	//logger
	//config
	//aldi sessions

	// Signals
	interuptSignals chan os.Signal
	reloadSignals   chan os.Signal
}

func NewCelerceptor() *Celerceptor {}

func (c *Celerceptor) Launch() error {
	c.Lock()

	// aldi talk new session (login)

	c.Unlock()
	return c.run()
}

func (c *Celerceptor) Reload() error {
	c.Lock()

	// config neu lesen und anwenden

	c.Unlock()
}

func (c *Celerceptor) Shutdown() error {
	c.Lock()

	//kein ticker, logout,session beenden

}

func (c *Celerceptor) run() error {
	/*
	* ticker um datenvolumen einzusehen (smart ticker) => (nachbuchen oder nix?)
	* volumen merken
	* und os Signals
	 */

	defer c.handlePanic()
	for {
		select {
			//!case ticker

		case <-c.interuptSignals:
			//!b.logger.Error("catched SIGINT/SIGTERM")
			c.Shutdown()
			return 

		case <-c.reloadSignals:
			//!log
			if err := c.Reload(); err != nil {
				//! log
			}

		default:
			continue
	}
}

func (c *Celerceptor) handlePanic() {
	if r := recover(); r != nil {
		//bot.Logger.Error("PANIC", fmt.Sprintf("%#v", r))
		//bot.Shutdown()
	}
}
