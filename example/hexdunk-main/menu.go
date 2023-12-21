package main

import (
	G "github.com/AllenDang/giu"
)

func ifActiveFile(w G.Widget) G.Widget {
	disabled := G.Style().SetDisabled(true).To(w)
	return G.Condition(ActiveFile() != nil, G.Layout{w}, G.Layout{disabled})
}

func ifClipboard(w G.Widget) G.Widget {
	disabled := G.Style().SetDisabled(true).To(w)
	return G.Condition(HD.ClipBoard != nil, G.Layout{w}, G.Layout{disabled})
}

func ifSelection(w G.Widget) G.Widget {
	tab := ActiveTab()
	disabled := G.Style().SetDisabled(true).To(w)
	return G.Condition(tab != nil && tab.view.selectionSize > 0, G.Layout{w}, G.Layout{disabled})
}

func ifUndo(w G.Widget) G.Widget {
	file := ActiveFile()
	disabled := G.Style().SetDisabled(true).To(w)
	return G.Condition(file != nil && len(file.undo) > 0, G.Layout{w}, G.Layout{disabled})
}

func ifRedo(w G.Widget) G.Widget {
	file := ActiveFile()
	disabled := G.Style().SetDisabled(true).To(w)
	return G.Condition(file != nil && len(file.redo) > 0, G.Layout{w}, G.Layout{disabled})
}

func menuFile() G.Widget {
	return G.Layout{
		G.MenuItem("New").OnClick(actionNewFile),
		G.MenuItem("Open").OnClick(actionOpenFile),
		G.Separator(),
		ifActiveFile(G.MenuItem("Save").OnClick(actionSaveFile)),
		ifActiveFile(G.MenuItem("Save As").OnClick(actionSaveAs)),
		ifActiveFile(G.MenuItem("Close Tab").OnClick(actionCloseTab)),
		G.Separator(),
		//G.MenuItem("Settings").OnClick(menuEditSettings),
		//G.Separator(),
		G.MenuItem("Quit").OnClick(actionQuit),
	}
}

func menuEdit() G.Widget {
	return G.Layout{
		ifSelection(G.MenuItem("Cut        x").OnClick(actionCut)),
		ifSelection(G.MenuItem("Copy       y").OnClick(actionCopy)),
		ifClipboard(G.MenuItem("Paste      p").OnClick(actionPaste)),
		G.Separator(),
		ifUndo(G.MenuItem("Undo       u").OnClick(actionUndo)),
		ifRedo(G.MenuItem("Redo       r").OnClick(actionRedo)),
	}
}

func menuPlugin() G.Widget {
	return G.Layout{
		G.MenuItem("Load"),
		G.Separator(),
		G.Menu("Plugin Foo").Layout(
			G.MenuItem("DoBar"),
			G.MenuItem("DoBaz"),
			G.MenuItem("Quux"),
			G.Separator(),
			G.MenuItem("Settings"),
		),
	}
}

func mkMenu() G.Widget {
	return G.Layout{
		G.Menu("File").Layout(menuFile()),
		G.Menu("Edit").Layout(menuEdit()),
		//G.Menu("Plugin").Layout(menuPlugin()),
	}
}
