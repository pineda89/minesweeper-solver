package main

import (
	"github.com/pineda89/minesweeper-solver/game"
	"github.com/pineda89/minesweeper-solver/providers/googleminesweeper"
	"github.com/pineda89/minesweeper-solver/winapi"
	"image"
	"log"
	"os"
	"time"
)

// run over: https://www.google.com/fbx?fbx=minesweeper

func main() {
	g := &game.Game{
		Movements: make(map[game.XY]bool),
		Provider:  &googleminesweeper.GoogleProvider{},
	}

	for {
		img, _ := winapi.MyCaptureScreen()

		myImg := image.Image(img)

		g.Board = g.Provider.ProcessImage(myImg)

		if g.MinesweeperSafeIteration() {
			time.Sleep(150 * time.Millisecond)
			g.UnsafeCtr = 0
		} else {
			g.UnsafeCtr++

			if g.UnsafeCtr > 10 {
				if !g.MinesweeperUnSafeIteration() {
					log.Println("game seems ended!")
					os.Exit(0)
				}
				time.Sleep(150 * time.Millisecond)
				g.UnsafeCtr = 0
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
