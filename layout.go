package layout

import (
	"context"

	"honnef.co/go/js/dom"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

const (
	SmallScreenOnly = "mdl-layout--small-screen-only"
	LargeScreenOnly = "mdl-layout--large-screen-only"
)

var layoutTabUpdates chan string
var doc dom.Document

func init() {
	layoutTabUpdates = make(chan string)
	doc = dom.GetWindow().Document()
}

func switchState(id string, active bool) {
	toActive(id+"-bar", active)
	toActive(id, active)
}

func toActive(id string, state bool) {
	e := doc.GetElementByID(id)
	a := "is-active"
	if state {
		e.Class().Add(a)
	} else {
		e.Class().Remove(a)
	}
}

type Layout struct {
	vecty.Core
	IsJS              bool
	Header            *Header
	Children          vecty.MarkupOrComponentOrHTML
	Content           vecty.MarkupOrComponentOrHTML
	FixedDrawer       bool
	FixedHeader       bool
	NoDrawerBtn       bool
	NoDestopDrawerBtn bool
	LargeScreenOnly   bool
	TabView           bool
	Panes             []*Panel
	cancel            func()
	activePane        string
	Drawer            *Drawer
}

func New(tabview bool) *Layout {
	l := &Layout{TabView: tabview}
	if tabview {
		ctx, cancel := context.WithCancel(context.Background())
		l.cancel = cancel
		return l.watch(ctx)
	}
	return l
}

func (l *Layout) watch(ctx context.Context) *Layout {
	go func() {
	done:
		for {
			select {
			case id := <-layoutTabUpdates:
				if id != l.activePane {
					switchState(l.activePane, false)
					switchState(id, true)
					l.activePane = id
				}
			case <-ctx.Done():
				break done
			}
		}
	}()
	return l
}

func (l *Layout) Unmount() {
	if l.cancel != nil {
		l.cancel()
	}
}

func (l *Layout) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-layout"] = true
	if l.IsJS {
		c["mdl-js-layout"] = true
	}
	if l.FixedDrawer {
		c["mdl-layout--fixed-drawer"] = true
	}
	if l.FixedHeader {
		c["mdl-layout--fixed-header"] = true
	}
	if l.NoDrawerBtn {
		c["mdl-layout--no-drawer-button"] = true
	}
	if l.NoDestopDrawerBtn {
		c["mdl-layout--no-desktop-drawer-button"] = true
	}
	var p vecty.List

	if l.LargeScreenOnly {
		c["mdl-layout--large-screen-only"] = true
	}
	if l.Header != nil {
		p = append(p, l.Header)
	}
	if l.Drawer != nil {
		p = append(p, l.Drawer)
	}

	if l.TabView {
		b := &TabBar{}
		var pt vecty.List
		for i := 0; i < len(l.Panes); i++ {
			b.Links = append(b.Links, &TabLink{
				ID:       l.Panes[i].ID,
				Text:     l.Panes[i].Name,
				IsActive: l.Panes[i].IsActive,
			})
			pt = append(pt, l.Panes[i])
		}
		if l.Header == nil {
			l.Header = &Header{}
			p = append(p, l.Header)
		}
		l.Header.TabView = true
		l.Header.TabNav = b
		p = append(p, elem.Main(
			pt,
		))
	} else {
		if l.Content != nil {
			p = append(p, elem.Main(prop.Class("mdl-layout__content"), l.Content))
		}
	}

	return elem.Div(
		l.Children,
		c,
		p,
	)
}

type Header struct {
	vecty.Core
	Icon             *Icon
	Row              *HeaderRow
	Children         vecty.MarkupOrComponentOrHTML
	Scroll           bool
	Waterfall        bool
	WaterfallHideTop bool
	Transparent      bool
	Seamed           bool
	TabView          bool
	TabNav           *TabBar
	ManualTabSwitch  bool
}

func (h *Header) Render() *vecty.HTML {
	var l vecty.List
	if h.Icon != nil {
		l = append(l, h.Icon)
	}
	c := make(vecty.ClassMap)
	c["mdl-layout__header"] = true
	if h.Scroll {
		c["mdl-layout__header--scroll"] = true
	}
	if h.Waterfall {
		c["mdl-layout__header--waterfall"] = true
	}
	if h.WaterfallHideTop {
		c["mdl-layout__header--waterfall-hide-top"] = true
	}
	if h.Transparent {
		c["mdl-layout__header--transparent"] = true
	}
	if h.Seamed {
		c["mdl-layout__header--seamed"] = true
	}
	if h.TabView {
		return elem.Header(
			h.Children,
			c,
			h.Row,
			h.TabNav,
		)
	}
	return elem.Header(
		h.Children,
		c,
		h.Row,
	)

}

type Icon struct {
	vecty.Core
	Children vecty.MarkupOrComponentOrHTML
}

func (i *Icon) Render() *vecty.HTML {
	return elem.Div(i.Children, prop.Class("mdl-layout-icon"))
}

type HeaderRow struct {
	vecty.Core
	Children  vecty.MarkupOrComponentOrHTML
	Title     *Title
	AddSpacer bool
	Nav       *Nav
}

func (h *HeaderRow) Render() *vecty.HTML {
	var l vecty.List
	if h.Title != nil {
		l = append(l, h.Title)
	}
	if h.AddSpacer {
		l = append(l, elem.Div(prop.Class("mdl-layout-spacer")))
	}
	if h.Nav != nil {
		l = append(l, h.Nav)
	}
	return elem.Div(
		h.Children,
		prop.Class("mdl-layout__header-row"),
		l,
	)

}

type Nav struct {
	vecty.Core
	Links    []*NavLink
	Children vecty.MarkupOrComponentOrHTML
}

func (n *Nav) Render() *vecty.HTML {
	var l vecty.List
	for i := 0; i < len(n.Links); i++ {
		l = append(l, n.Links[i])
	}
	return elem.Navigation(
		n.Children,
		prop.Class("mdl-navigation"),
		l,
	)
}

type NavLink struct {
	vecty.Core
	Href     string
	Text     string
	Children vecty.MarkupOrComponentOrHTML
}

func (n *NavLink) Render() *vecty.HTML {
	return elem.Anchor(
		n.Children,
		prop.Class("mdl-navigation__link"),
		vecty.Text(n.Text),
	)
}

type Title struct {
	vecty.Core
	Text     string
	Children vecty.MarkupOrComponentOrHTML
}

func (t *Title) Render() *vecty.HTML {
	return elem.Span(
		t.Children,
		prop.Class("mdl-layout-title"),
		vecty.Text(t.Text),
	)
}

type TabBar struct {
	vecty.Core
	Links  []*TabLink
	Manual bool
}

type TabLink struct {
	vecty.Core
	ID       string
	IsActive bool
	Text     string
	Children vecty.MarkupOrComponentOrHTML
}

func (t *TabLink) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-layout__tab"] = true
	if t.IsActive {
		c["is-active"] = true
	}
	return elem.Anchor(
		t.Children,
		prop.ID("#"+t.ID),
		prop.ID(t.ID+"-bar"), c,
		vecty.Text(t.Text),
		event.Click(func(e *vecty.Event) {
			go func() {
				layoutTabUpdates <- t.ID
			}()
		}),
	)
}

type TabNav struct {
	vecty.Core
	TabBar   *TabBar
	Tabs     []*Panel
	Children vecty.MarkupOrComponentOrHTML
}

type Panel struct {
	vecty.Core
	ID       string
	Name     string
	IsActive bool
	Children vecty.MarkupOrComponentOrHTML
}

func (t *Panel) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-layout__tab-panel"] = true
	if t.IsActive {
		c["is-active"] = true
	}
	return elem.Section(
		c,
		prop.ID(t.ID),
		t.Children,
	)
}

type Drawer struct {
	vecty.Core
	Title *Title
	Nav   *Nav
}

func (d *Drawer) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-layout__drawer"] = true
	return elem.Div(
		c, d.Title, d.Nav,
	)
}
