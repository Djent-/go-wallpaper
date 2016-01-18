package main

import (
	wdb "github.com/djent-/go-walldatabase"
	ui "github.com/gizak/termui"
	"time"
	"fmt"
	//"log"
	//"strings"
)

type PaneType int

const (
	FILELIST PaneType = 0
	TAGLIST PaneType = 1
	DBFILE = "C:\\Users\\Patrick\\Documents\\wall.db"
)

type Pane struct {
	List ui.List // termui struct
	TotalItems []string // entire list of items
	CurrentIndex int // index of selected list item
	ListOffset int
	HasFocus bool
	Type PaneType
}

type Screen struct {
	Title ui.Par
	Panes []Pane
	StatusBar ui.Par
	HasFocus bool
	Active int
	KbdHandler func(ui.Event)
}

func (s *Screen) Draw() {
	ui.Render(&s.Title, &s.Panes[0].List, &s.Panes[1].List)
}

func (s *Screen) ToggleActivePane() {
	s.Panes[s.Active].HasFocus = false
	s.Active = s.Active ^ 1
	s.Panes[s.Active].HasFocus = true
	ui.Render(&s.Title, &s.Panes[0].List, &s.Panes[1].List)
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	
	Screens := CreateScreens()
	active := 0
	
	// Open the wallpaper database
	WallDB := wdb.OpenDB(DBFILE)
	Screens[0].Panes[0].PopulateWallpaperFilelistPane(WallDB)
	
	ui.Handle("sys/kbd/<escape>", func(ui.Event) {
		// press esc to quit
		ui.StopLoop()
	})
	ui.Handle("sys/kbd/<left>", func(ui.Event) {
		// switch to left pane
		Screens[active].ToggleActivePane()
	})
	ui.Handle("sys/kbd/<right>", func(ui.Event) {
		// switch to right pane
		Screens[active].ToggleActivePane()
	})
	ui.Handle("sys/kbd/<up>", func(ui.Event) {
		// decrement Pane List selected index
		if Screens[active].Panes[Screens[active].Active].CurrentIndex == 0 {
			return
		}
		Screens[active].Panes[Screens[active].Active].CurrentIndex -= 1
	})
	ui.Handle("sys/kbd/<down>", func(ui.Event) {
		// increment Pane List selected index
		// TODO: check for overflow
		Screens[active].Panes[Screens[active].Active].CurrentIndex += 1
	})
	ui.Handle("sys/kbd/<tab>", func(ui.Event) {
		// toggle active screen
		Screens[active].HasFocus = false
		// ^ is Go's xor operator. Works only on ints
		active = active ^ 1
		Screens[active].HasFocus = true
	})
	ui.Handle("sys/kbd", HandleKeyboardEvent)
	ui.Handle("termui-wallpaper/index/update", func(ui.Event) {
		// handle updates to CurrentIndex
	})
	
	// this is https://github.com/gizak/termui/issues/58
	tick := time.Second/24
	ui.Merge("timer/update", ui.NewTimerCh(tick))
	ui.Handle("/timer/"+tick.String(), func(e ui.Event) {
		// update pane lists
		Screens[0].Panes[0].UpdateWallpaperFilelistPane(WallDB)
		// call draw
		Screens[active].Draw()
	})
	ui.Loop()
}

func HandleKeyboardEvent(e ui.Event) {
	// Send event to active Screen -> Pane
}

func (p *Pane) PopulateWallpaperFilelistPane(w wdb.WallDatabase) error {
	// Get list of wallpapers from the database
	wallpapers := w.FetchAllWallpapers()
	// clear current list from pane
	p.TotalItems = []string{}
	// go through the wallpapers and add the filename to p.TotalItems
	for _, wallpaper := range(wallpapers) {
		p.TotalItems = append(p.TotalItems, wallpaper.Filename)
	}
	err := p.UpdateWallpaperFilelistPane(w)
	// FetchAllWallpapers() doesn't return an error yet, but it will
	return err
}

func (p *Pane) UpdateWallpaperFilelistPane(w wdb.WallDatabase) error {
	// clear p.List.Items
	p.List.Items = []string{}
	for index, filename := range(p.TotalItems) {
		// break if index is out of view bounds
		if index > p.ListOffset + 17 { // outside of visible range
			break
		}
		if index < p.ListOffset {
			continue
		}
		var filename_f string
		var filename_f1 string
		var filename_f2 string
		if len(filename) > 43 {
			for ind, char := range(filename) {
				if (ind < 20) { 
					filename_f1 = fmt.Sprintf("%s%c", filename_f1, char)
				} else if (ind > len(filename) - 22) {
					filename_f2 = fmt.Sprintf("%s%c", filename_f2, char)
				}
			}
			// add ellipsis to center of truncated string
			filename_f = fmt.Sprintf("%s…%s", filename_f1, filename_f2)
		}
		if index + p.ListOffset == p.CurrentIndex {
			filename_f = fmt.Sprintf("[%s](fg-white,bg-green)", filename_f)
		}
		p.List.Items = append(p.List.Items, filename_f)
	}
	return nil
}

func CreateScreens() []Screen {
	SCREENHEIGHT := 18
	TITLEHEIGHT := 1
	// Wallpaper screen
	wallpapers_t := ui.NewPar("          Wallpapers          ")
	wallpapers_t.Height = TITLEHEIGHT
	wallpapers_t.Width = 30
	wallpapers_t.Y = 0
	wallpapers_t.X = 25
	wallpapers_t.Border = false
	wallpapers := &Screen{Title: *wallpapers_t, HasFocus: true}
	// filename pane (left)
	filename_l := ui.NewList()
	filename_l.Height = SCREENHEIGHT - TITLEHEIGHT
	filename_l.Width = 45
	filename_l.BorderLabel = "Wallpapers"
	filename_l.X = 0
	filename_l.Y = 1
	filenames := &Pane{List: *filename_l, 
		CurrentIndex: 0, 
		HasFocus: true,
		Type: FILELIST,
		ListOffset: 0}
	// tag pane (right)
	tag_l := ui.NewList()
	tag_l.Height = SCREENHEIGHT - TITLEHEIGHT
	tag_l.Width = 36
	tag_l.BorderLabel = "Tags"
	tag_l.X = 44
	tag_l.Y = 1
	tags := &Pane{List: *tag_l, 
		CurrentIndex: 0, 
		HasFocus: false,
		Type: TAGLIST,
		ListOffset: 0}
	wallpapers.Panes = []Pane{*filenames, *tags}
	
	// Slideshow screen
	slideshows_t := ui.NewPar("          Slideshows          ")
	slideshows_t.Height = TITLEHEIGHT
	slideshows_t.Width = 30
	slideshows_t.Y = 0
	slideshows_t.X = 25
	slideshows_t.Border = false
	slideshows := &Screen{Title: *slideshows_t, HasFocus: false}
	// slideshow pane (left)
	slideshow_l := ui.NewList()
	slideshow_l.Height = SCREENHEIGHT - TITLEHEIGHT
	slideshow_l.Width = 36
	slideshow_l.Y = 1
	slideshow_l.X = 0
	slideshow_l.BorderLabel = "Slideshows"
	slideshow_p := &Pane{List: *slideshow_l, 
		CurrentIndex: 0, 
		HasFocus: false,
		Type: FILELIST,
		ListOffset: 0}
	// wallpaper pane (right)
	wallpaper_l := ui.NewList()
	wallpaper_l.Height = SCREENHEIGHT - TITLEHEIGHT
	wallpaper_l.Width = 45
	wallpaper_l.BorderLabel = "Wallpapers"
	wallpaper_l.X = 35
	wallpaper_l.Y = 1
	wallpaper_p := &Pane{List: *wallpaper_l, 
		CurrentIndex: 0, 
		HasFocus: false,
		Type: FILELIST,
		ListOffset: 0}
	slideshows.Panes = []Pane{*slideshow_p, *wallpaper_p}
	
	return []Screen{*wallpapers, *slideshows}
}