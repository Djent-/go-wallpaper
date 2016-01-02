package main

import (
	"github.com/nsf/termbox-go"
)

type EditBox struct {
	defaultVal string
	currentVal string
	foreground termbox.Attribute
	background termbox.Attribute
	posx int
	posy int
	height int
	width int
	active bool
}

func (e EditBox) ToggleActive() {
	e.active = !e.active
}

func (e EditBox) Keystroke(key termbox.Key) {
	// If the edit box is at its default value, set the current
	// value to the entered keystroke.
	
	// Don't allow spaces. Would break compatibility with other
	// programs.
	
	
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	
	// When ESC sequence is in the buffer and it doesn't 
	// match any known sequence. ESC means KeyEsc.
	termbox.SetInputMode(termbox.InputEsc)
	
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	// ...
	
	// Synchronizes the internal back buffer with the terminal.
	termbox.Flush()
	
	// Loop
}

/*
What do I want this program to do?
Obviously I can't display the wallpaper unless I find a Go
image-to-ascii library.
Maybe just have a selector pane on one side with filenames,
and on the other have the list of tags and a blank space for a new one.

################################################
#				Term-Wallpaper.go              #
################################################
#catwallpaper1.png#cat                         #
#*catwp2.jpg      #cute                        #
#snowyroad.jpeg   #animal                      #
#thecutestkitte...#                            #
#                 #                            #
#                 #                            #
#                 #                            #
################################################

I don't want to have to make my own ranger-like filesystem browser,
so I don't know how I would add wallpapers to this.
I'd also want other functionality, like creating slideshow playlists
and choosing one to play.

################################################
#             Term-Wallpaper.go                #
################################################
#     *Wallpapers*    #     Slideshows         #
################################################
#catwallpaper1.png    #cat                     #
#*catwp2.jpg*         #cute                    #
#snowyroad.jpeg       #animal                  #
#thecutestkitteneve...#+                       #
#                     #                        #
################################################

From the Wallpapers screen, 'a' to add a wallpaper.
Will set the state of the program to be a text entry box,
which will validate the file after the user presses enter.
The program will then return to the Wallpapers screen,
and set the newly added wallpaper is the currently selected
wallpaper from the list of wallpapers.
To edit tags, the user will use the arrow keys to navigate
to the right pane and either start typing in the tag slot
labelled '+' or press delete while selecting a tag.

*/