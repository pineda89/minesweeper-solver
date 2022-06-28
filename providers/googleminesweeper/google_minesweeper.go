package googleminesweeper

import (
	"github.com/pineda89/minesweeper-solver/providers"
	"github.com/pineda89/minesweeper-solver/util"
	"github.com/pineda89/minesweeper-solver/winapi"
	"image"
	"log"
	"time"
)

type GoogleProvider struct {
	boardInitialized                         bool
	initialX, initialY, endpointX, endpointY int
	rowSize, colSize                         float64
	numColumns, numRows                      int
	board                                    [][]int
}

type RGB struct {
	R         uint32
	G         uint32
	B         uint32
	Tolerance uint32
}

func (inst *RGB) Match(r, g, b uint32) bool {
	return util.InBetween(int(r), int(inst.R)-int(inst.Tolerance), int(inst.R)+int(inst.Tolerance)) &&
		util.InBetween(int(g), int(inst.G)-int(inst.Tolerance), int(inst.G)+int(inst.Tolerance)) &&
		util.InBetween(int(b), int(inst.B)-int(inst.Tolerance), int(inst.B)+int(inst.Tolerance))
}

var (
	COLOR_BROWN_DARK = RGB{R: 216, G: 185, B: 154, Tolerance: 5}
	COLOR_BROWN      = RGB{R: 230, G: 195, B: 160, Tolerance: 5}
	COLOR_GREEN_DARK = RGB{R: 163, G: 210, B: 73, Tolerance: 5}
	COLOR_GREEN      = RGB{R: 171, G: 216, B: 81, Tolerance: 5}

	COLOR_BLUE         = RGB{R: 25, G: 118, B: 210, Tolerance: 10}
	COLOR_GREEN_NUMBER = RGB{R: 56, G: 142, B: 60, Tolerance: 10}
	COLOR_RED_NUMBER   = RGB{R: 211, G: 50, B: 49, Tolerance: 10}
	COLOR_FUCSIA       = RGB{R: 123, G: 31, B: 162, Tolerance: 10}
	COLOR_ORANGE       = RGB{R: 255, G: 143, Tolerance: 10}
	COLOR_CYAN         = RGB{G: 151, B: 167, Tolerance: 10}
	COLOR_RED_FLAG     = RGB{R: 230, G: 51, B: 7, Tolerance: 10}
)

func (prov *GoogleProvider) ProcessImage(img image.Image) [][]int {
	if !prov.boardInitialized {

		prov.initialX, prov.initialY = prov.findInitialPoint(img)
		prov.endpointX, prov.endpointY, prov.numColumns, prov.numRows = prov.findEndPoint(img, prov.initialX, prov.initialY)
		log.Println(prov.initialX, prov.initialY, prov.endpointX, prov.endpointY)

		winapi.MoveMouse((prov.endpointX+prov.initialX)/2, (prov.endpointY+prov.initialY)/2)
		log.Println("moving to", (prov.endpointX+prov.initialX)/2, (prov.endpointY+prov.initialY)/2)
		winapi.ClickMouse(30*time.Millisecond, false)

		prov.rowSize, prov.colSize = float64(prov.endpointX-prov.initialX)/float64(prov.numColumns), float64(prov.endpointY-prov.initialY)/float64(prov.numRows)

		prov.board = make([][]int, prov.numRows)
		for row := 0; row < prov.numRows; row++ {
			prov.board[row] = make([]int, prov.numColumns)
			for col := 0; col < prov.numColumns; col++ {
				prov.board[row][col] = -1
			}
		}

		prov.boardInitialized = true

		log.Println(len(prov.board), prov.numRows, prov.numColumns, prov.rowSize, prov.colSize)
	}

	for row := 0; row < prov.numRows; row++ {
		for col := 0; col < prov.numColumns; col++ {
			var counters = make([]int, 10)

			for x := prov.initialX + int(prov.colSize*float64(col)) + int(prov.colSize*5/100); x < prov.initialX+int(prov.colSize*float64(col))+int(prov.colSize*95/100); x++ {
				for y := prov.initialY + int(prov.rowSize*float64(row)) + int(prov.rowSize*5/100); y < prov.initialY+int(prov.rowSize*float64(row))+int(prov.rowSize*95/100); y++ {
					r, g, b := providers.GetColor(img, x, y)
					switch {
					case COLOR_BROWN.Match(r, g, b), COLOR_BROWN_DARK.Match(r, g, b):
						counters[providers.CELL_DISCOVERED]++
					case COLOR_GREEN.Match(r, g, b), COLOR_GREEN_DARK.Match(r, g, b):
						counters[providers.CELL_UNKNOWN]++
					case COLOR_BLUE.Match(r, g, b):
						counters[providers.CELL_1_NEIGHBORN]++
					case COLOR_GREEN_NUMBER.Match(r, g, b):
						counters[providers.CELL_2_NEIGHBORN]++
					case COLOR_RED_NUMBER.Match(r, g, b):
						counters[providers.CELL_3_NEIGHBORN]++
					case COLOR_FUCSIA.Match(r, g, b):
						counters[providers.CELL_4_NEIGHBORN]++
					case COLOR_ORANGE.Match(r, g, b):
						counters[providers.CELL_5_NEIGHBORN]++
					case COLOR_CYAN.Match(r, g, b):
						counters[providers.CELL_6_NEIGHBORN]++
					case COLOR_RED_FLAG.Match(r, g, b):
						counters[providers.CELL_FLAG]++
					}
				}
			}

			prov.board[row][col] = providers.GetTag(counters)
		}
	}

	return prov.board
}

func (prov *GoogleProvider) findEndPoint(img image.Image, initialX int, initialY int) (int, int, int, int) {
	var endpointX, endpointY int
	var numColumns, numRows = 1, 1

	fr, fg, fb := providers.GetColor(img, initialX, initialY)
	for x := initialX; x < img.Bounds().Dx(); x++ {
		r, g, b := providers.GetColor(img, x, initialY)
		if !util.InBetween(r, fr-5, fr+5) || !util.InBetween(g, fg-5, fg+5) || !util.InBetween(b, fb-5, fb+5) {
			if !util.InBetween(r, 160, 175) || !util.InBetween(g, 200, 225) || !util.InBetween(b, 65, 85) {
				endpointX = x
				break
			}
			fr, fg, fb = r, g, b
			numColumns++
		}
	}

	fr, fg, fb = providers.GetColor(img, initialX, initialY)
	for y := initialY; y < img.Bounds().Dy(); y++ {
		r, g, b := providers.GetColor(img, initialX, y)
		if !util.InBetween(r, fr-5, fr+5) || !util.InBetween(g, fg-5, fg+5) || !util.InBetween(b, fb-5, fb+5) {
			if !util.InBetween(r, 160, 175) || !util.InBetween(g, 200, 225) || !util.InBetween(b, 65, 85) {
				endpointY = y
				break
			}
			fr, fg, fb = r, g, b
			numRows++
		}
	}

	return endpointX, endpointY, numColumns, numRows
}

func (prov *GoogleProvider) findInitialPoint(img image.Image) (int, int) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b := providers.GetColor(img, x, y)

			switch {
			case r == 171 && g == 216 && b == 81:
				return x, y
			case r == 163 && g == 210 && b == 73:
				return x, y
			}
		}
	}
	return -1, -1
}

func (prov *GoogleProvider) MoveToCeld(x, y int) {
	winapi.MoveMouse(prov.initialX+int(prov.colSize*float64(x))+int(prov.colSize/2), prov.initialY+int(prov.rowSize*float64(y))+int(prov.rowSize/2))
	time.Sleep(10 * time.Millisecond)
}
