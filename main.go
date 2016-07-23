/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"log"
	"reflect"
	"runtime"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	runtime.LockOSThread()
}

func idle(data reflect.Value) bool {
	renderer := data.Interface().(*sdl.Renderer)

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	var rect sdl.Rect
	renderer.GetViewport(&rect)
	rect.W /= 2
	rect.H /= 2

	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&rect)

	renderer.Present()

	return true
}

func main() {
	gtk.Init(nil)
	sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("SDL2 in GTK3")
	window.SetSizeRequest(640, 360)
	window.ShowAll()

	sdlWindow, err := sdl.CreateWindowFrom(window.Widget.GetWindow().GetNativeWindow().ID)
	if err != nil {
		log.Fatalln(err)
	}

	renderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalln(err)
	}

	window.Connect("configure_event", func() {
		log.Println("GTK configure_event!")
	})

	window.Connect("delete_event", func() {
		renderer.Destroy()
		sdlWindow.Destroy()
		sdl.Quit()
		gtk.MainQuit()
	})

	glib.IdleAdd(idle, renderer)
	gtk.Main()

	log.Println("clean shutdown")
}
