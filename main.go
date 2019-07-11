package main

import (
	"machine"

	"image/color"

	"time"

	"math/rand"

	"github.com/conejoninja/demoreel/tetris"
	"github.com/conejoninja/tinydraw"
	"github.com/conejoninja/tinyfont"
	"github.com/conejoninja/tinyfont/demoreel"
	"tinygo.org/x/drivers/waveshare-epd/epd2in13"
)

var display epd2in13.Device

var black = color.RGBA{1, 1, 1, 255}
var white = color.RGBA{0, 0, 0, 255}
var w int16 = 250
var h int16 = 122

func main() {
	machine.PowerSupplyActive(true)
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		Mode:      0,
		SCK:       machine.EPD_SCK_PIN,
		MOSI:      machine.EPD_MOSI_PIN,
	})

	display = epd2in13.New(machine.SPI0, machine.EPD_CS_PIN, machine.EPD_DC_PIN, machine.EPD_RESET_PIN, machine.EPD_BUSY_PIN)
	display.Configure(epd2in13.Config{Rotation: 1})

	// display.SetLUT(true) is the slow-but-no-ghosting update
	// this 4 lines will flash the screen black-white several times to "clean" it
	display.SetLUT(true)
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	display.SetLUT(false) // a bit dirty , but way faster

	// different effects, comment to disable
	for {
		blinky("LOOK", "AT ME")
		myNameIs("@conejo")
		talkDate("Thursday, 5PM", "at room 123 build XYZ")
		sunrays()
		loading()
		loadingInverted()
		tetrisfx()
		dvd("TinyGo")
	}
}

func blinky(topline, bottomline string) {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	// calculate the width of the text so we could center them later
	w32top, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte(topline))
	w32bottom, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte(bottomline))
	for i := int16(0); i < 10; i++ {
		// fill the screen with with
		tinydraw.FilledRectangle(&display, 0, 0, w, h, white)
		// show black text
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32top))/2, 50, []byte(topline), black)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32bottom))/2, 100, []byte(bottomline), black)

		// display
		display.Display()
		display.WaitUntilIdle()

		// repeat the other way around
		tinydraw.FilledRectangle(&display, 0, 0, w, h, black)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32top))/2, 50, []byte(topline), white)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32bottom))/2, 100, []byte(bottomline), white)

		display.Display()
		display.WaitUntilIdle()
	}
}

func myNameIs(name string) {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	var r int16 = 8

	// round corners
	tinydraw.FilledCircle(&display, r, r, r, black)
	tinydraw.FilledCircle(&display, w-r-1, r, r, black)
	tinydraw.FilledCircle(&display, r, h-r-1, r, black)
	tinydraw.FilledCircle(&display, w-r-1, h-r-1, r, black)

	// top band
	tinydraw.FilledRectangle(&display, r, 0, w-2*r-1, r, black)
	tinydraw.FilledRectangle(&display, 0, r, w, 26, black)

	// bottom band
	tinydraw.FilledRectangle(&display, r, h-r, w-2*r-1, r, black)
	tinydraw.FilledRectangle(&display, 0, h-2*r-1, w, r, black)

	// top text : my NAME is
	w32, _ := tinyfont.LineWidth(&demoreel.Regular12pt7b, []byte("my NAME is"))
	tinyfont.WriteLine(&display, &demoreel.Regular12pt7b, (w-int16(w32))/2, 24, []byte("my NAME is"), white)

	// middle text
	w32, _ = tinyfont.LineWidth(&demoreel.Bold12pt7b, []byte(name))
	tinyfont.WriteLine(&display, &demoreel.Bold12pt7b, (w-int16(w32))/2, 74, []byte(name), black)

	display.Display()
	display.WaitUntilIdle()
	time.Sleep(5 * time.Second)
}

func talkDate(dateString, roomString string) {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	// top text : Come see my talk on
	_, w32 := tinyfont.LineWidth(&demoreel.Bold12pt7b, []byte("Come see my talk on"))
	tinyfont.WriteLine(&display, &demoreel.Regular12pt7b, (w-int16(w32))/2, 28, []byte("Come see my talk on"), black)

	// middle text
	w32, _ = tinyfont.LineWidth(&demoreel.Bold12pt7b, []byte(dateString))
	tinyfont.WriteLine(&display, &demoreel.Bold12pt7b, (w-int16(w32))/2, 70, []byte(dateString), black)

	if roomString != "" {
		// bottom text : at room XYZ
		w32, _ = tinyfont.LineWidth(&demoreel.Regular12pt7b, []byte(roomString))
		tinyfont.WriteLine(&display, &demoreel.Regular12pt7b, (w-int16(w32))/2, 110, []byte(roomString), black)
	}

	display.Display()
	display.WaitUntilIdle()
	time.Sleep(5 * time.Second)
}

func sunrays() {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	colors := [7][]color.RGBA{
		{white, white, white, white, white, white},
		{black, white, white, white, white, white},
		{black, white, white, black, white, white},
		{black, white, white, black, black, white},
		{black, black, white, black, black, white},
		{black, black, white, black, black, black},
		{black, black, black, black, black, black},
	}

	for i := 0; i < 21; i++ {
		if i%2 == 0 {
			tinydraw.FilledRectangle(&display, 0, 0, w, h, white)
			rays(black)
		} else {
			tinydraw.FilledRectangle(&display, 0, 0, w, h, black)
			rays(white)
		}
		tinydraw.FilledRectangle(&display, 20, 20, 210, 102, white)
		w32, _ := tinyfont.LineWidth(&demoreel.Regular12pt7b, []byte("Badge coded with"))
		tinyfont.WriteLine(&display, &demoreel.Regular12pt7b, (w-int16(w32))/2, 50, []byte("Badge coded with"), black)

		w32, _ = tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte("TinyGo"))
		tinyfont.WriteLineColors(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2, 100, []byte("TinyGo"), colors[i%7])

		display.Display()
		display.WaitUntilIdle()
	}

}

func rays(color color.RGBA) {
	// center point at the bottom
	var cx int16 = 124
	var cy int16 = 122

	// left side rays
	tinydraw.FilledTriangle(&display, cx, cy, 0, 0, 0, 20, color)
	tinydraw.FilledTriangle(&display, cx, cy, 0, 40, 0, 60, color)
	tinydraw.FilledTriangle(&display, cx, cy, 0, 80, 0, 100, color)

	// top rays
	tinydraw.FilledTriangle(&display, cx, cy, 20, 0, 40, 0, color)
	tinydraw.FilledTriangle(&display, cx, cy, 60, 0, 80, 0, color)
	tinydraw.FilledTriangle(&display, cx, cy, 100, 0, 120, 0, color)
	tinydraw.FilledTriangle(&display, cx, cy, 140, 0, 160, 0, color)
	tinydraw.FilledTriangle(&display, cx, cy, 180, 0, 200, 0, color)
	tinydraw.FilledTriangle(&display, cx, cy, 220, 0, 240, 0, color)

	// right rays
	tinydraw.FilledTriangle(&display, cx, cy, 250, 10, 250, 30, color)
	tinydraw.FilledTriangle(&display, cx, cy, 250, 50, 250, 70, color)
	tinydraw.FilledTriangle(&display, cx, cy, 250, 90, 250, 110, color)
}

func loading() {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	for i := int16(0); i < 25; i++ {
		// draw a rectangle bigger each time
		tinydraw.FilledRectangle(&display, i*10, 0, 10, h, black)

		// draw text again since a part of it was behind the rectangle
		w32, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte("TinyGo"))
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2, 70, []byte("TinyGo"), white)

		display.Display()
		display.WaitUntilIdle()
	}
}

func loadingInverted() {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	for i := int16(0); i < 25; i++ {
		// this is the opposite, we draw the text and draw a rectangle of the same color as the background
		// to make the ilusion the text is revealing
		w32, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte("TinyGo"))
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2, 70, []byte("TinyGo"), black)

		tinydraw.FilledRectangle(&display, i*10, 0, 250-i*10, h, white)

		display.Display()
		display.WaitUntilIdle()
	}
}

// tetrisfx
// will create a new ramdom piece ramdomly rotated each time the previous one stopped
func tetrisfx() {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	tetris.NewBoard()

	tetris.NewPiece()
	failed := 0
	k := 0
	for {
		display.ClearBuffer()
		if tetris.MovePiece() {
			failed = 0
		} else {
			failed++
			tetris.NewPiece()
		}
		tetris.DrawBoard(&display)
		tetris.DrawPiece(&display)

		w32, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte("TinyGo"))
		// add a white broder around the text
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2-2, 70-2, []byte("TinyGo"), white)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2-2, 70+2, []byte("TinyGo"), white)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2+2, 70-2, []byte("TinyGo"), white)
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2+2, 70+2, []byte("TinyGo"), white)
		// add the text
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, (w-int16(w32))/2, 70, []byte("TinyGo"), black)

		display.Display()
		display.WaitUntilIdle()
		// stop after 5 pieces in a row that can not move (screen is kind of filled)
		if failed >= 5 {
			return
		}
		k++
	}
}

func dvd(text string) {
	display.ClearBuffer()
	display.Display()
	display.WaitUntilIdle()

	w32, _ := tinyfont.LineWidth(&demoreel.Bold24pt7b, []byte(text))
	maxW := w - int16(w32)
	maxH := h - 36 //assume line height is 36

	// random start point
	x := int16(rand.Int31n(int32(maxW)))
	y := int16(rand.Int31n(int32(maxH)))
	d := int16(4)
	dx := d
	dy := d

	for i := 0; i < 80; i++ { //duration 80 frames
		// paint white the previous text to "erase" it
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, x, y+36, []byte("TinyGo"), white)

		// move text
		x += dx
		y += dy

		// paint and show text
		tinyfont.WriteLine(&display, &demoreel.Bold24pt7b, x, y+36, []byte("TinyGo"), black)
		display.Display()
		display.WaitUntilIdle()

		// change direction if needed
		if x >= maxW {
			dx = -d
		}
		if x <= 0 {
			dx = d
		}
		if y >= maxH {
			dy = -d
		}
		if y <= 0 {
			dy = d
		}
	}
}
