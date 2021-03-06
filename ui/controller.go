package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/edouardparis/lntop/ui/views"
)

type controller struct {
	logger logging.Logger
	models *models.Models
	views  *views.Views
}

func (c *controller) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	return c.views.Layout(g, maxX, maxY)
}

func (c *controller) cursorDown(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return view.CursorDown()
	}
	return nil
}

func (c *controller) cursorUp(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return view.CursorUp()
	}
	return nil
}

func (c *controller) cursorRight(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return view.CursorRight()
	}
	return nil
}

func (c *controller) cursorLeft(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return view.CursorLeft()
	}
	return nil
}

func (c *controller) SetModels(ctx context.Context) error {
	err := c.models.RefreshInfo(ctx)
	if err != nil {
		return err
	}

	err = c.models.RefreshWalletBalance(ctx)
	if err != nil {
		return err
	}

	err = c.models.RefreshChannelsBalance(ctx)
	if err != nil {
		return err
	}

	return c.models.RefreshChannels(ctx)
}

func (c *controller) Listen(ctx context.Context, g *gocui.Gui, sub chan *events.Event) {
	c.logger.Debug("Listening...")
	refresh := func(fn ...func(context.Context) error) {
		for i := range fn {
			err := fn[i](ctx)
			if err != nil {
				c.logger.Error("failed", logging.Error(err))
			}
		}
		g.Update(func(*gocui.Gui) error { return nil })
	}

	for event := range sub {
		c.logger.Debug("event received", logging.String("type", event.Type))
		switch event.Type {
		case events.BlockReceived:
			refresh(c.models.RefreshInfo)
		case events.ChannelPending:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.ChannelActive:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.ChannelInactive:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.InvoiceSettled:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.PeerUpdated:
			refresh(c.models.RefreshInfo)
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *controller) Help(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	view := c.views.Get(g.CurrentView())
	if view == nil {
		return nil
	}

	if view.Name() != views.HELP {
		c.views.SetPrevious(view)
		return c.views.Help.Set(g, 0, -1, maxX, maxY)
	}

	err := g.DeleteView(views.HELP)
	if err != nil {
		return err
	}

	if c.views.Previous != nil {
		_, err := g.SetCurrentView(c.views.Previous.Name())
		return err
	}

	return nil
}

func (c *controller) OnEnter(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	view := c.views.Get(v)
	if view == nil {
		return nil
	}

	switch view.Name() {
	case views.CHANNELS:
		c.views.SetPrevious(view)
		_, cy := v.Cursor()
		err := c.models.SetCurrentChannel(context.Background(), cy)
		if err != nil {
			return err
		}

		err = c.views.Channel.Set(g, 0, 6, maxX-1, maxY)
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView(c.views.Channel.Name())
		return err

	case views.CHANNEL:
		err := c.views.Channel.Delete(g)
		if err != nil {
			return err
		}

		if c.views.Previous != nil {
			_, err := g.SetCurrentView(c.views.Previous.Name())
			return err
		}

		err = c.views.Channels.Set(g, 0, 6, maxX-1, maxY)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) setKeyBinding(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyF10, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, c.cursorLeft)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, c.cursorRight)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, c.OnEnter)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, c.Help)
	if err != nil {
		return err
	}

	return nil
}

func newController(app *app.App) *controller {
	m := models.New(app)
	return &controller{
		logger: app.Logger.With(logging.String("logger", "controller")),
		models: m,
		views:  views.New(m),
	}
}
