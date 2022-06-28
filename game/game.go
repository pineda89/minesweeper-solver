package game

import (
	"github.com/pineda89/minesweeper-solver/providers"
	"github.com/pineda89/minesweeper-solver/util"
	"github.com/pineda89/minesweeper-solver/winapi"
	"log"
	"time"
)

type Game struct {
	Movements           map[XY]bool
	Provider            providers.Provider
	UnsafeCtr           int
	numRows, numColumns int
	Board               [][]int
}

type XY struct {
	X int
	Y int
}

func (game *Game) MinesweeperUnSafeIteration() bool {
	var bestChoise, bestRow, bestCol int
	var bestNeighborn []XY
	for row := 0; row < len(game.Board); row++ {
		for column := 0; column < len(game.Board[row]); column++ {
			switch {
			case util.InBetween(game.Board[row][column], 1, 6):
				if settedBombs, v, va := game.getVecinos(row, column); (v - game.Board[row][column] + settedBombs) > bestChoise {
					bestChoise = v - game.Board[row][column] + settedBombs
					bestRow, bestCol = row, column
					bestNeighborn = va
				}
			}
		}
	}

	if bestChoise > 0 {
		game.Provider.MoveToCeld(bestNeighborn[0].X, bestNeighborn[0].Y)
		winapi.ClickMouse(2*time.Millisecond, false)

		log.Println("random click...", bestRow, bestCol, bestChoise, bestNeighborn)

		return true
	}
	return false
}

func (game *Game) safeCell(row int, column int, bombsNear int) bool {
	settedBombs, unknownNeighbornsCtr, unknownNeighborns := game.getVecinos(row, column)
	if settedBombs < bombsNear {
		if settedBombs+unknownNeighbornsCtr == bombsNear {
			for i := range unknownNeighborns {
				if !game.Movements[XY{X: unknownNeighborns[i].X, Y: unknownNeighborns[i].Y}] {
					game.Movements[XY{X: unknownNeighborns[i].X, Y: unknownNeighborns[i].Y}] = true
					game.Provider.MoveToCeld(unknownNeighborns[i].X, unknownNeighborns[i].Y)
					winapi.ClickMouse(2*time.Millisecond, true)
				}
			}

			log.Println("tab", bombsNear, "can be resolved", row, column, "#", unknownNeighbornsCtr, unknownNeighborns)
			return true
		}
	} else if unknownNeighbornsCtr > 0 {
		for i := range unknownNeighborns {
			if !game.Movements[XY{X: unknownNeighborns[i].X, Y: unknownNeighborns[i].Y}] {
				game.Movements[XY{X: unknownNeighborns[i].X, Y: unknownNeighborns[i].Y}] = true
				game.Provider.MoveToCeld(unknownNeighborns[i].X, unknownNeighborns[i].Y)
				winapi.ClickMouse(2*time.Millisecond, false)
			}
		}
		return true
	}
	return false
}

func (game *Game) MinesweeperSafeIteration() (res bool) {
	for row := 0; row < len(game.Board); row++ {
		for column := 0; column < len(game.Board[row]); column++ {
			switch {
			case util.InBetween(game.Board[row][column], 1, 6):
				if game.safeCell(row, column, game.Board[row][column]) {
					res = true
				}
			}
		}
	}
	return res
}

func (game *Game) getVecinos(row int, column int) (settedBombs int, uncoveredNeighborsCtr int, uncoveredNeighbors []XY) {
	for i := row - 1; i <= row+1; i++ {
		for j := column - 1; j <= column+1; j++ {

			if util.InBetween(i, 0, len(game.Board)-1) && util.InBetween(j, 0, len(game.Board[i])-1) {
				if !(i == row && j == column) {
					switch game.Board[i][j] {
					case providers.CELL_FLAG:
						settedBombs++
					case providers.CELL_UNKNOWN:
						uncoveredNeighborsCtr++
						uncoveredNeighbors = append(uncoveredNeighbors, XY{X: j, Y: i})
					}
				}
			}
		}
	}
	return settedBombs, uncoveredNeighborsCtr, uncoveredNeighbors
}
