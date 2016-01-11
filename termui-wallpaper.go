package main

import (
	//wdb "github.com/djent-/go-walldatabase"
	ui "github.com/gizak/termui"
	"time"
	//"strconv"
	"fmt"
)

type Pane struct {
	List ui.List // termui struct
	TotalItems []string // entire list of items
	CurrentIndex int // index of selected list item
	HasFocus bool
}

type Screen struct {
	Title ui.Par
	Panes []Pane
	StatusBar ui.Par
	HasFocus bool
	Active int
}

func (s Screen) Draw() {
	// s.Title.Text = time.Now().String() // debug (works)
	ui.Render(&s.Title, &s.Panes[0].List, &s.Panes[1].List)
}

func (s Screen) ToggleActivePane() {
	s.Panes[s.Active].HasFocus = false
	s.Panes[s.Active].List.BorderLabel = "Inactive" // debug
	s.Active = s.Active ^ 1
	s.Title.Text = "Test" // debug (does not work)
	s.Title.Text = fmt.Sprint("%d", s.Active) // debug (does not work)
	s.Panes[s.Active].HasFocus = true
	s.Panes[s.Active].List.BorderLabel = "Active" // debug
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	
	Screens := CreateScreens()
	active := 0
	
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
	ui.Handle("sys/kbd/<tab>", func(ui.Event) {
		// toggle active screen
		// ^ is Go's xor operator. Works only on ints
		active = active ^ 1
	})
	ui.Handle("sys/kbd", func(e ui.Event) {
		// send event to active screen
	})
	draw := func() {
		Screens[active].Draw()
	}
	// this is https://github.com/gizak/termui/issues/58
	tick := time.Second/24
	ui.Merge("timer/update", ui.NewTimerCh(tick))
	ui.Handle("/timer/"+tick.String(), func(e ui.Event) {
		// call draw
		draw()
	})
	ui.Loop()
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
	filenames := &Pane{List: *filename_l, CurrentIndex: 1, HasFocus: true}
	// tag pane (right)
	tag_l := ui.NewList()
	tag_l.Height = SCREENHEIGHT - TITLEHEIGHT
	tag_l.Width = 35
	tag_l.BorderLabel = "Tags"
	tag_l.X = 45
	tag_l.Y = 1
	tags := &Pane{List: *tag_l, CurrentIndex: 1, HasFocus: false}
	wallpapers.Panes = []Pane{*filenames, *tags}
	
	// Slideshow screen
	slideshows_t := ui.NewPar("Slideshows")
	slideshows_t.Height = TITLEHEIGHT
	slideshows_t.Width = 10
	slideshows_t.Y = 0
	slideshows_t.X = 35
	slideshows_t.Border = false
	slideshows := &Screen{Title: *slideshows_t, HasFocus: false}
	// slideshow pane (left)
	slideshow_l := ui.NewList()
	slideshow_l.Height = SCREENHEIGHT - TITLEHEIGHT
	slideshow_l.Width = 35
	slideshow_l.Y = 1
	slideshow_l.X = 0
	slideshow_l.BorderLabel = "Slideshows"
	slideshow_p := &Pane{List: *slideshow_l, CurrentIndex: 1, HasFocus: false}
	// wallpaper pane (right)
	wallpaper_l := ui.NewList()
	wallpaper_l.Height = SCREENHEIGHT - TITLEHEIGHT
	wallpaper_l.Width = 45
	wallpaper_l.BorderLabel = "Wallpapers"
	wallpaper_l.X = 35
	wallpaper_l.Y = 1
	wallpaper_p := &Pane{List: *wallpaper_l, CurrentIndex: 1, HasFocus: false}
	slideshows.Panes = []Pane{*slideshow_p, *wallpaper_p}
	
	return []Screen{*wallpapers, *slideshows}
}