package main

import (
	G "github.com/AllenDang/giu"
	//I "github.com/AllenDang/imgui-go"
)

func draw() {
	G.MainMenuBar().Layout(mkMenu()).Build()

	//G.SingleWindowWithMenuBar().Layout(
	G.Window("Files").Pos(5, 30).Size(600, 600).Layout(
		G.PrepareMsgbox(),
		PrepareFileDialog(DialogOpen, actionOpen),
		PrepareFileDialog(DialogSaveAs, actionWriteFile),
		PrepareIntDialog(DialogGoto, actionGotoAddr),
		//G.MenuBar().Layout(mkMenu()),
		//makeToolBar(),
		mkTabWidget(),
	)
}

func main() {
	G.SetDefaultFont("DejavuSansMono.ttf", 12)
	w := G.NewMasterWindow("HexDunk", 800, 800, 0)
	w.Run(draw)
}
