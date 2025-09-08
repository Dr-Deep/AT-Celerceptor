package internal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	alditalk "github.com/Dr-Deep/AT-Celerceptor/internal/aldi-talk"
	"github.com/Dr-Deep/AT-Celerceptor/internal/config"
	"github.com/Dr-Deep/logging-go"
)

type Celerceptor struct {
	sync.Mutex

	logger *logging.Logger
	cfg    *config.Configuration

	//aldi sessions
	session *alditalk.AldiTalkSession

	// Signals
	ticker          *time.Ticker
	interuptSignals chan os.Signal
	reloadSignals   chan os.Signal
}

func NewCelerceptor(logger *logging.Logger, cfg *config.Configuration) *Celerceptor {
	var c = &Celerceptor{
		logger: logger,
		cfg:    cfg,
	}

	return c
}

func (c *Celerceptor) Launch() error {
	c.Lock()

	c.AdjustTicker(
		time.Duration(c.cfg.CheckInterval),
	)

	c.interuptSignals = make(chan os.Signal, 1)
	signal.Notify(
		c.interuptSignals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	// reload sig

	// Session
	s, err := alditalk.NewAldiTalkSession(c.logger, c.cfg.Tel, c.cfg.Password)
	if err != nil {
		return err
	}

	c.session = s

	c.Unlock()
	return c.run()
}

func (c *Celerceptor) AdjustTicker(interval time.Duration) {
	c.logger.Info("interval set to", interval.Abs().String())
	c.ticker = time.NewTicker(interval)
}

func (c *Celerceptor) Reload() error {
	c.Lock()

	// config neu lesen und anwenden

	c.Unlock()

	return nil
}

func (c *Celerceptor) Shutdown() error {
	c.Lock()
	c.logger.Info("shutting down...")

	// Signals
	c.ticker.Stop()
	signal.Stop(c.interuptSignals)
	close(c.interuptSignals)
	//c.reloadSignals

	// aldi logout?

	return nil
}

func (c *Celerceptor) run() error {
	/*
	* ticker um datenvolumen einzusehen (smart ticker) => (nachbuchen oder nix?)
	* volumen merken
	 */

	defer c.handlePanic()
	for {
		select {
		case <-c.ticker.C:
			c.Lock()
			//! aldi talk check
			c.Unlock()

		case <-c.interuptSignals:
			c.logger.Error("catched SIGINT/SIGTERM")
			c.Shutdown()
			return nil

		case <-c.reloadSignals:
			c.logger.Info("catched SIGUSR1/SIGUSR2")
			if err := c.Reload(); err != nil {
				c.logger.Error("reload failed", err.Error())
			}

		default:
			continue
		}
	}
}

func (c *Celerceptor) handlePanic() {
	if r := recover(); r != nil {
		c.logger.Error("PANIC", fmt.Sprintf("%#v", r))
		c.Shutdown()
	}
}
