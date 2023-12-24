package main

import (
	"fmt"
	"github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/ddkwork/golibrary/stream"
)

var (
	editor     *g.CodeEditorWidget
	errMarkers imgui.ErrorMarkers
)

func loop() {
	g.SingleWindow().Layout(
		g.Row(
			g.Button("Get Text").OnClick(func() {
				if editor.HasSelection() {
					fmt.Println(editor.GetSelectedText())
				} else {
					fmt.Println(editor.GetText())
				}

				column, line := editor.GetCursorPos()
				fmt.Println("Cursor pos:", column, line)

				column, line = editor.GetSelectionStart()
				fmt.Println("Selection start:", column, line)

				fmt.Println("Current line is", editor.GetCurrentLineText())
			}),
			g.Button("Set Text").OnClick(func() {
				editor.Text("Set text")
			}),
			g.Button("Set Error Marker").OnClick(func() {
				errMarkers.Clear()
				errMarkers.Insert(1, "Error message")
				fmt.Println("ErrMarkers Size:", errMarkers.Size())

				editor.ErrorMarkers(errMarkers)
			}),
		),
		editor,
	)
}

func main() {
	wnd := g.NewMasterWindow("Code Editor", 800, 600, 0)

	errMarkers = imgui.NewErrorMarkers()

	editor = g.CodeEditor().
		ShowWhitespaces(false).
		TabSize(2).
		Text("select * from greeting\nwhere date > current_timestamp\norder by date").
		LanguageDefinition(giu.LanguageDefinitionSQL).
		Border(true)

	wnd.SetDropCallback(func(files []string) {
		f := files[0]
		s := stream.NewReadFile(f)
		editor.Text(s.String())
		g.Update()
	})

	wnd.Run(loop)
}