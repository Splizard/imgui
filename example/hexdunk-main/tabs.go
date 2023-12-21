package main

import (
	"fmt"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
	//I "github.com/AllenDang/imgui-go"
)

func mkTabWidget() G.Widget {
	return G.Custom(func() {
		if len(HD.Tabs) != 0 && I.BeginTabBar("HexViewerTabs") {
			for i, tab := range HD.Tabs {

				hf, ok := HD.Files[tab.name]
				if !ok {
					panic("tab opened on closed file")
				}

				flags := G.TabItemFlagsNone
				if hf.dirty {
					flags |= G.TabItemFlagsUnsavedDocument
				}
				if I.BeginTabItem(fmt.Sprint(i) + ": " + hf.stats.Name()) {
					HD.ActiveTab = i
					h := HexView(fmt.Sprint(i, ".hexview##", hf.name), hf.buf, tab.view)
					h.Build()
					I.EndTabItem()
				}
			}
			I.EndTabBar()
		}
	})
}

func OpenTab(hf *HexFile) {
	HD.Tabs = append(HD.Tabs, HexTab{name: hf.name, view: new(ViewState)})
}

func CloseTab(t int) {
	if t < 0 || len(HD.Tabs) < t {
		panic("closeTab number doesn't exist (shouldn't happen)")
	}
	tab := HD.Tabs[HD.ActiveTab]
	copy(HD.Tabs[t:], HD.Tabs[t+1:])
	HD.Tabs = HD.Tabs[:len(HD.Tabs)-1]
	if HD.ActiveTab == t {
		HD.ActiveTab--
	}

	//if another tab has the same file open, don't close the file permanently
	for _, t := range HD.Tabs {
		if t.name == tab.name {
			return
		}
	}
	CloseHexFile(tab.name) // we closed last tab on this guy
}
