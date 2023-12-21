package main

//little program to display a file-system tree and basic info
//once a file-path is selected, print it on stdout

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

var (
	statCache       = make(map[string]fs.FileInfo)
	dirCache        = make(map[string][]fs.FileInfo)
	showHiddenFiles = false
	selectedFile    string //full path of file selected in fileTable
	currentDir      string //full path of directory selected in dirTree
	startDir        string //starting directory arg or cwd()
)

const (
	timeFmt     = "02 Jan 06 15:04"
	nodeFlags   = imgui.TreeNodeFlagsOpenOnArrow | imgui.TreeNodeFlagsOpenOnDoubleClick | imgui.TreeNodeFlagsSpanFullWidth
	leafFlags   = imgui.TreeNodeFlagsLeaf
	tableFlags  = imgui.TableFlags_ScrollX | imgui.TableFlags_ScrollY | imgui.TableFlags_Resizable | imgui.TableFlags_SizingStretchProp
	selectFlags = imgui.SelectableFlagsAllowDoubleClick | imgui.SelectableFlagsSpanAllColumns
)

func mkSize(sz_ int64) string {
	sizes := []string{"KB", "MB", "GB", "TB"}
	sz := float64(sz_)
	add := ""
	for _, n := range sizes {
		if sz < 1024 {
			break
		}
		sz = sz / 1024
		add = n
	}
	if add == "" {
		return fmt.Sprint(sz_)
	} else {
		return fmt.Sprintf("%.2f %s", sz, add)
	}
}

//statFile follows symbolic links
func statFile(path string) (fs.FileInfo, error) {
	st, ok := statCache[path]
	if !ok {
		var err error
		st, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
		statCache[path] = st
	}
	return st, nil
}

//return a list of directory entries
func readDir(path string) ([]fs.FileInfo, error) {
	entry, ok := dirCache[path]
	if !ok {
		direntry, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		entry := make([]fs.FileInfo, len(direntry))
		for i, f := range direntry {
			childPath := filepath.Join(path, f.Name())

			if f.Type()&fs.ModeSymlink == 0 {
				entry[i], err = f.Info()
			} else {
				//default from ReadDir is lstat, which doesn't follow symlinks :(
				//stat it again using os.Stat
				entry[i], err = statFile(childPath)
			}
			if err != nil {
				return nil, err
			}
		}
		dirCache[path] = entry
	}
	return entry, nil
}

func getDirInfo(path string) (int, fs.FileInfo, []fs.FileInfo, bool) {
	info, err := statFile(path)
	if err != nil {
		return 0, nil, nil, false
	}

	entries, err := readDir(path)
	if err != nil {
		return 0, nil, nil, false
	}

	if info.Name()[0] == '.' && !showHiddenFiles {
		return 0, nil, nil, false
	}

	flags := leafFlags
	for _, e := range entries {
		if e.IsDir() {
			flags = nodeFlags
			break
		}
	}

	if path == currentDir {
		flags |= imgui.TreeNodeFlagsSelected
	}

	return flags, info, entries, true
}

func dirTree(path string) {
	flags, info, entries, ok := getDirInfo(path)
	if !ok {
		return
	}

	if strings.HasPrefix(startDir, path) {
		flags |= imgui.TreeNodeFlagsDefaultOpen
	}

	imgui.PushStyleVarFloat(imgui.StyleVarIndentSpacing, 7)
	defer imgui.PopStyleVar()

	open := imgui.TreeNodeV(info.Name(), flags)
	if imgui.IsItemClicked(int(giu.MouseButtonLeft)) {
		currentDir = path
	}
	if open {
		defer imgui.TreePop()
		for _, e := range entries {
			if e.IsDir() {
				name := filepath.Join(path, e.Name())
				dirTree(name)
			}
		}
	}
}

func isHidden(entry fs.FileInfo) bool {
	return entry.Name()[0] == '.'
}

func fileTable() {
	imgui.Text(currentDir)
	if imgui.BeginTable("FSTable", 3, tableFlags, imgui.ContentRegionAvail(), 0) {
		defer imgui.EndTable()
		imgui.TableSetupColumn("Name", 0, 10, 0)
		imgui.TableSetupColumn("Size", 0, 2, 0)
		imgui.TableSetupColumn("Time", 0, 4, 0)
		imgui.TableSetupScrollFreeze(1, 1)
		imgui.TableHeadersRow()
		//TODO: set up sorting

		entries, err := readDir(currentDir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.IsDir() || isHidden(e) {
				continue
			}
			path := filepath.Join(currentDir, e.Name())

			imgui.TableNextRow(0, 0)
			imgui.TableNextColumn()
			if imgui.SelectableV(e.Name(), path == selectedFile, selectFlags, imgui.Vec2{}) {
				selectedFile = path
				if imgui.IsMouseDoubleClicked(int(giu.MouseButtonLeft)) {
					selectFile()
				}
			}
			imgui.TableNextColumn()
			imgui.Text(mkSize(e.Size()))
			imgui.TableNextColumn()
			imgui.Text(e.ModTime().Format(timeFmt))
		}
	}
}

func selectFile() {
	if !filepath.IsAbs(selectedFile) {
		selectedFile = filepath.Join(currentDir, selectedFile)
	}
	fmt.Println(selectedFile)
	os.Exit(0)
}

func cancel() {
	os.Exit(1)
}

func mkNavBar() {
	width, _ := giu.GetAvailableRegion()
	giu.InputText(&selectedFile).Size(width).Build()
}

func loop() {
	giu.SingleWindow().Layout(
		giu.Custom(mkNavBar),
		giu.Custom(func() {
			//use a child frame to block the lists going off-screen
			//at the bottom of the screen is a row of buttons and such
			w, h := giu.GetAvailableRegion()
			//adjust for buttonrow height
			_, spacing := giu.GetItemSpacing()
			_, padding := giu.GetFramePadding()
			_, buttonH := giu.CalcTextSize("F")
			h -= buttonH + 2*padding + spacing
			giu.Child().Layout(
				giu.SplitLayout(giu.DirectionHorizontal, true, 200,
					giu.Custom(func() { dirTree(filepath.FromSlash("/")) }),
					giu.Custom(fileTable),
				),
			).Border(false).Size(w, h).Build()
		}),
		giu.Row(
			giu.Checkbox("Show Hidden", &showHiddenFiles),
			giu.Button("Cancel").OnClick(cancel),
			giu.Button("Select").OnClick(selectFile),
		),
	)
}

func main() {
	if len(os.Args) == 1 {
		startDir, _ = filepath.Abs(".")
	} else {
		startDir = os.Args[1]
		st, err := statFile(startDir)
		if err != nil {
			log.Fatalf("starting directory '%s': %v\n", startDir, err)
		}
		if !st.IsDir() {
			log.Fatalf("starting path '%s': not a directory\n", startDir)
		}
		startDir, _ = filepath.Abs(startDir)
	}
	currentDir = startDir
	giu.SetDefaultFont("DejavuSansMono.ttf", 12)
	w := giu.NewMasterWindow("FileTree", 800, 600, 0)
	w.Run(loop)
}
