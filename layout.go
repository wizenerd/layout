package layout

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Layout struct {
	vecty.Core
	IsJS     bool
	Header   *Header
	Children vecty.MarkupOrComponentOrHTML
	Content  vecty.MarkupOrComponentOrHTML
}

func (l *Layout) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-layout"] = true
	if l.IsJS {
		c["mdl-js-layout"] = true
	}
	var p vecty.List
	if l.Header != nil {
		p = append(p, l.Header)
	}
	if l.Content != nil {
		p = append(p, elem.Main(prop.Class("mdl-layout__content"), l.Content))
	}
	return elem.Div(
		l.Children,
		c,
		p,
	)
}

type Header struct {
	vecty.Core
	Icon     *Icon
	Rows     []*HeaderRow
	Children vecty.MarkupOrComponentOrHTML
}

func (h *Header) Render() *vecty.HTML {
	var l vecty.List
	if h.Icon != nil {
		l = append(l, h.Icon)
	}
	for i := 0; i < len(h.Rows); i++ {
		l = append(l, h.Rows[i])
	}
	return elem.Header(
		h.Children,
		prop.Class("mdl-layout__header"),
		l,
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
