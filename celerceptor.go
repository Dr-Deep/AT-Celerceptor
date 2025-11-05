/*
	!!! PROOF OF CONCEPT - HAFTUNGSAUSSCHLUSS !!!

Dieser Code dient rein der Illustration.
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
	var c = &Celerceptor{
		logger: logger,
		cfg:    cfg,
	}

	return c
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

	c.tryLogin()

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
			c.Lock()
			c.tryUpdateDataVol()
			c.Unlock()

		case <-c.interuptSignals:
			c.logger.Error("catched SIGINT/SIGTERM")
			c.Shutdown()
			return nil
		}
	}
}

func (c *Celerceptor) handlePanic() {
	if r := recover(); r != nil {
		c.logger.Error("PANIC", fmt.Sprintf("%#v", r))
		c.Shutdown()
	}
}

/*
 * AT FUNCS
 */

/*
* ticker um datenvolumen einzusehen (smart ticker) => (nachbuchen oder nix?)
* volumen merken
 */

func (c *Celerceptor) tryLogin() {
	c.logger.Info("trying to login")

	if c.atframClient == nil {
		client, err := atfram.New(
			&atfram.AldiTalkConfig{
				BaseURI:  atfram.ALDITALK_BASE_URI,
				Username: c.cfg.Tel,
				Password: c.cfg.Password,
			},
			c.logger,
		)
		if err != nil {
			// in paar ticks wiederversuchen?
			panic(err)
		}

		c.atframClient = client
	}

	if err := c.atframClient.Login(); err != nil {
		panic(err)
	}
}

func (c *Celerceptor) tryLogout() {
	c.logger.Info("trying to logout")
}

func (c *Celerceptor) tryUpdateDataVol() {
	c.logger.Info("trying to update datavol")
	c.tryGetDataVol()
}
func (c *Celerceptor) tryGetDataVol() {
	c.logger.Info("trying to get datavol")
}

/*
	// Session
	s, err := alditalk.NewAldiTalkSession(c.logger, c.cfg.Tel, c.cfg.Password)
	if err != nil {
		return err
	}

	c.session = s
*/
