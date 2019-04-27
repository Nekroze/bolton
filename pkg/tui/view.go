package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Nekroze/bolton/pkg/boltons"
	"github.com/Nekroze/bolton/pkg/hardpoints"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func Run(hps []*hardpoints.Point, bl boltons.Library) {
	app := tview.NewApplication()
	st := selectionTree(hps, bl)
	dv := diffView(st, hps, bl, app)
	flex := tview.NewFlex().
		AddItem(st, 0, 1, false).
		AddItem(dv, 0, 2, false)
	app.SetRoot(flex, true)
	app.SetFocus(st)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func hardpointLookupTable(hps []*hardpoints.Point) map[string]*hardpoints.Point {
	out := map[string]*hardpoints.Point{}
	for _, hp := range hps {
		out[hp.String()] = hp
	}
	return out
}

func selectionTree(hps []*hardpoints.Point, bl boltons.Library) *tview.TreeView {
	root := tview.NewTreeNode("Hardpoints / Boltons").SetColor(tcell.ColorRed).SetSelectable(false)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	for _, hp := range hps {
		hardpoint := hp.String()
		node := tview.NewTreeNode(hardpoint).
			SetReference(hardpoint).
			SetSelectable(true)
		root.AddChild(node)

		for _, tag := range hp.Tags {
			nodetag := tview.NewTreeNode(tag).
				SetReference(hardpoint).
				SetSelectable(true)
			node.AddChild(nodetag)
			for _, bolt := range bl[tag] {
				nodebolton := tview.NewTreeNode(bolt.Name).
					SetReference(hardpoint + ":" + tag + ":" + bolt.Name).
					SetSelectable(true)
				nodetag.AddChild(nodebolton)
			}
		}
	}
	return tree
}

func applySelector(app *tview.Application, hps []*hardpoints.Point, bl boltons.Library) func(*tview.TreeNode) {
	hpl := hardpointLookupTable(hps)
	return func(node *tview.TreeNode) {
		ref, ok := node.GetReference().(string)
		if !ok {
			return
		}
		parts := strings.Split(ref, ":")
		if len(parts) < 4 {
			return
		}

		filename := parts[0]
		lineno := parts[1]
		tag := parts[2]
		boltname := parts[3]

		hardpoint := hpl[fmt.Sprintf("%s:%s", filename, lineno)]
		if hardpoint == nil { // may not be needed
			return
		}

		b, err := ioutil.ReadFile(hardpoint.Path)
		if err != nil {
			panic(err)
		}
		text := string(b)

		bolton, err := bl.Get(tag + ":" + boltname)
		if err != nil {
			panic(err)
		}
		text = apply(text, bolton, hardpoint)
		f, err := os.Create(hardpoint.Path)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(text)
		if err != nil {
			f.Close()
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
		app.Stop()
	}
}

func changeSelector(dv *tview.TextView, hps []*hardpoints.Point, bl boltons.Library) func(*tview.TreeNode) {
	hpl := hardpointLookupTable(hps)
	return func(node *tview.TreeNode) {
		ref, ok := node.GetReference().(string)
		if !ok {
			return
		}
		parts := strings.Split(ref, ":")
		if len(parts) < 2 {
			return
		}

		filename := parts[0]
		lineno := parts[1]

		hardpoint := hpl[fmt.Sprintf("%s:%s", filename, lineno)]
		if hardpoint == nil { // may not be needed
			return
		}

		b, err := ioutil.ReadFile(hardpoint.Path)
		if err != nil {
			panic(err)
		}
		text := string(b)

		if len(parts) == 4 {
			tag := parts[2]
			boltname := parts[3]
			bolton, err := bl.Get(tag + ":" + boltname)
			if err != nil {
				panic(err)
			}
			text = apply(text, bolton, hardpoint)
		}

		dv.SetText(text)
		dv.ScrollToBeginning()
		if hardpoint.Line > 1 {
			dv.ScrollTo(hardpoint.Line-1, 0)
		}
	}
}

func apply(input string, b boltons.Bolton, h *hardpoints.Point) string {
	lines := strings.Split(input, "\n")

	newstring, err := b.Contents()
	if err != nil {
		panic(err)
	}
	newsection := strings.Split(newstring, "\n")

	output := append(append(lines[:h.Line], newsection...), lines[h.Line:]...)

	return strings.Join(output, "\n")
}

func diffView(st *tview.TreeView, hps []*hardpoints.Point, bl boltons.Library, app *tview.Application) *tview.TextView {
	dv := tview.NewTextView()
	dv.SetScrollable(true)
	st.SetChangedFunc(changeSelector(dv, hps, bl))
	st.SetSelectedFunc(applySelector(app, hps, bl))
	return dv
}
