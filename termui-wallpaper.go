package main

import (
	//wdb "github.com/djent-/go-walldatabase"
	ui "github.com/gizak/termui"
)

type Pane struct {
	List ui.List // termui struct
	TotalItems []string // entire list of items
	CurrentIndex int // index of selected list item
	HasFocus bool
}

type Screen struct {
	Title ui.Par
	LeftPane Pane
	RightPane Pane
	StatusBar ui.Par
	HasFocus bool
}

func (s Screen) Draw() {
	ui.Render(&s.Title, &s.LeftPane.List, &s.RightPane.List)
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
	ui.Handle("sys/kbd/left", func(ui.Event) {
		// switch to left pane
	})
	ui.Handle("sys/kbd/right", func(ui.Event) {
		// switch to right pane
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
	ui.Handle("/timer/1s", func(e ui.Event) {
		// call draw
		// TODO: make the timer shorter to get 20+ FPS
		draw()
	})
	ui.Loop()
}

func CreateScreens() []Screen {	
	// Wallpaper screen
	wallpapers_t := ui.NewPar("Wallpapers")
	wallpapers_t.Height = 1
	wallpapers_t.Width = 10
	wallpapers_t.Y = 0
	wallpapers_t.X = 35
	wallpapers_t.Border = false
	wallpapers := &Screen{Title: *wallpapers_t, HasFocus: true}
	// filename pane (left)
	filename_l := ui.NewList()
	filename_l.Height = 15
	filename_l.Width = 45
	filename_l.BorderLabel = "Wallpapers"
	filename_l.X = 0
	filename_l.Y = 1
	filenames := &Pane{List: *filename_l, CurrentIndex: 1, HasFocus: true}
	wallpapers.LeftPane = *filenames
	// tag pane (right)
	tag_l := ui.NewList()
	tag_l.Height = 15
	tag_l.Width = 35
	tag_l.BorderLabel = "Tags"
	tag_l.X = 45
	tag_l.Y = 1
	tags := &Pane{List: *tag_l, CurrentIndex: 1, HasFocus: false}
	wallpapers.RightPane = *tags
	
	// Slideshow screen
	slideshows_t := ui.NewPar("Slideshows")
	slideshows_t.Height = 1
	slideshows_t.Width = 10
	slideshows_t.Y = 0
	slideshows_t.X = 35
	slideshows_t.Border = false
	slideshows := &Screen{Title: *slideshows_t, HasFocus: false}
	// slideshow pane (left)
	slideshow_l := ui.NewList()
	slideshow_l.Height = 15
	slideshow_l.Width = 35
	slideshow_l.Y = 1
	slideshow_l.X = 0
	slideshow_l.BorderLabel = "Slideshows"
	slideshow_p := &Pane{List: *slideshow_l, CurrentIndex: 1, HasFocus: false}
	slideshows.LeftPane = *slideshow_p
	// wallpaper pane (right)
	wallpaper_l := ui.NewList()
	wallpaper_l.Height = 15
	wallpaper_l.Width = 45
	wallpaper_l.BorderLabel = "Wallpapers"
	wallpaper_l.X = 35
	wallpaper_l.Y = 1
	wallpaper_p := &Pane{List: *wallpaper_l, CurrentIndex: 1, HasFocus: false}
	slideshows.RightPane = *wallpaper_p
	
	return []Screen{*wallpapers, *slideshows}
}