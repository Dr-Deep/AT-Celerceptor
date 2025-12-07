/*
	!!! PROOF OF CONCEPT - HAFTUNGSAUSSCHLUSS !!!

Der Code in dieser Repo  dient rein der Illustration.
Nutzung auf eigene Gefahr. Keine Garantie für Funktionalität oder Richtigkeit.
Der Entwickler übernimmt keine Haftung für Schäden oder Folgenutzung.
Vor produktivem Einsatz rechtliche und technische Risiken prüfen.
*/

package celerceptor

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dr-Deep/AT-Celerceptor/atfram"
	"github.com/Dr-Deep/AT-Celerceptor/config"
	"github.com/Dr-Deep/logging-go"
)

type Celerceptor struct {
	sync.Mutex

	atframClient *atfram.Client
	logger       *logging.Logger
	cfg          *config.Configuration

	// Signals
	ticker          *time.Ticker
	interuptSignals chan os.Signal
}

func NewCelerceptor(logger *logging.Logger, cfg *config.Configuration) *Celerceptor {
	client, err := atfram.New(
		&atfram.AldiTalkConfig{
			BaseURI:  atfram.ALDITALK_BASE_URI,
			Username: cfg.Tel,
			Password: cfg.Password,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return &Celerceptor{
		atframClient: client,
		logger:       logger,
		cfg:          cfg,
	}
}

func (c *Celerceptor) Launch() error {
	c.Lock()

	c.AdjustTicker(
		time.Duration(c.cfg.CheckInterval) * time.Second,
	)

	c.interuptSignals = make(chan os.Signal, 1)
	signal.Notify(
		c.interuptSignals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go c.initAldiTalk()

	c.Unlock()

	return c.run()
}

func (c *Celerceptor) AdjustTicker(interval time.Duration) {
	c.logger.Info("interval set to", interval.Abs().String())
	c.ticker = time.NewTicker(interval)
}

func (c *Celerceptor) Shutdown() error {
	c.Lock()
	c.logger.Info("shutting down...")

	// Signals
	c.ticker.Stop()
	signal.Stop(c.interuptSignals)
	close(c.interuptSignals)

	// AT Logout
	c.tryLogout()

	return nil
}

func (c *Celerceptor) run() error {
	defer c.handlePanic()

	for {
		select {
		case <-c.ticker.C:
			go func() {
				c.logger.Error("TICKER BEGIN")
				c.Lock()
				c.logger.Error("TICKER END")
				c.Unlock()
			}()

		case <-c.interuptSignals:
			c.logger.Error("catched SIGINT/SIGTERM")
			return c.Shutdown()
		}
	}
}

func (c *Celerceptor) handlePanic() {
	if r := recover(); r != nil {
		c.logger.Error("PANIC", fmt.Sprintf("%#v", r))
		c.Shutdown()
	}
}
