package main

import (
	"fmt"
	"math/rand"
)

const (
	Road = iota
	Wall
	Visited
	Entrance
	Exit
	Trap
	Treasure
)

var cellChars = map[int]string{
	Road:     " ",
	Wall:     "#",
	Visited:  " ",
	Entrance: "E",
	Exit:     "X",
	Trap:     "T",
	Treasure: "$",
}

type Maze struct {
	Grid     [][]int
	Height   int
	Width    int
	Entrance Cell
	Exit     Cell
	Traps    int
}

type Cell struct {
	X int
	Y int
}

type BackTracker []Cell

func (b *BackTracker) IsEmpty() bool {
	return len(*b) == 0
}

// Push a new value onto the stack
func (b *BackTracker) Push(cell Cell) {
	*b = append(*b, cell) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (b *BackTracker) Pop() (Cell, bool) {
	if b.IsEmpty() {
		return Cell{-1, -1}, false
	} else {
		index := len(*b) - 1   // Get the index of the top most element.
		element := (*b)[index] // Index into the slice and obtain the element.
		*b = (*b)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func NewMaze(height int, width int) *Maze {
	if height%2 == 0 {
		height--
	}
	if width%2 == 0 {
		width--
	}

	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}

	return &Maze{
		Grid:     grid,
		Height:   height,
		Width:    width,
		Entrance: Cell{1, 1},
		Exit:     Cell{height - 2, width - 2},
		Traps:    5,
	}
}

func (m *Maze) GenerateMaze() {

	for i := 0; i < m.Height; i++ { // preparation to generate
		for j := 0; j < m.Width; j++ {
			if i%2 != 0 && j%2 != 0 {
				m.Grid[i][j] = Road
			} else {
				m.Grid[i][j] = Wall
			}
			m.Grid[m.Height-1][j] = Wall
		}
		m.Grid[i][m.Width-1] = Wall
	}
	currentCell := m.Entrance

	var backTracker BackTracker

	backTracker.Push(currentCell)
	m.Grid[currentCell.X][currentCell.Y] = Visited

	//firstIteration := true
	var path []Cell
	for len(backTracker) > 0 {
		neighbors := m.getNeighbors(currentCell)
		lenN := len(neighbors)
		if lenN > 0 {
			nextCell := neighbors[rand.Int()%lenN]
			m.removeWall(currentCell, nextCell)
			m.Grid[currentCell.X][currentCell.Y] = Visited
			backTracker.Push(nextCell)
			currentCell = nextCell
		} else {
			currentCell, _ = backTracker.Pop()
			//if m.isReached(currentCell) && firstIteration {
			//	path = append(path, backTracker...)
			//	fmt.Println(path)
			//	firstIteration = false
			//}
			//uncomment if you need more walls
			m.Grid[currentCell.X][currentCell.Y] = Visited
		}
	}
	//m.generatePath(path)
	m.generateEnvironment(path)

}

func (m *Maze) getNeighbors(cell Cell) (neighbors []Cell) {
	distance := 2
	up := Cell{cell.X, cell.Y - distance}
	dw := Cell{cell.X, cell.Y + distance}
	lt := Cell{cell.X - distance, cell.Y}
	rt := Cell{cell.X + distance, cell.Y}
	directions := []Cell{up, dw, lt, rt}
	for _, d := range directions {
		if d.X > 0 && d.X < m.Height && d.Y > 0 && d.Y < m.Width && m.Grid[d.X][d.Y] != Visited {
			neighbors = append(neighbors, d)
		}
	}
	return
}

func (m *Maze) removeWall(first Cell, second Cell) {
	dx := second.X - first.X
	dy := second.Y - first.Y
	if dx > 0 {
		dx++
	} else {
		dx--
	}
	if dy > 0 {
		dy++
	} else {
		dy--
	}
	dx /= 2
	dy /= 2
	fmt.Println("r1", dx, dy, first.X, first.Y)
	m.Grid[first.X+dx][first.Y+dy] = Visited
}

func (m *Maze) isReached(cell Cell) bool {
	return cell.X == m.Exit.X && cell.Y == m.Exit.Y
}

//func (m *Maze) generatePath(path []Cell) {
//	for _, i := range path {
//		m.Grid[i.X][i.Y] = Road
//	}
//	for i := 0; i < len(path)-1; i += 1 {
//		first, second := path[i], path[i+1]
//		dx := second.X - first.X
//		dy := second.Y - first.Y
//		if dx > 0 && dy == 0 {
//			dx++
//		} else if dy == 0 {
//			dx--
//		}
//		if dy > 0 && dx == 0 {
//			dy++
//		} else if dx == 0 {
//			dy--
//		}
//		dx /= 2
//		dy /= 2
//		fmt.Println(dx, dy, "   ", first.X, first.Y, "   ", second.X, second.Y, "    ", first.X+dx, first.Y+dy)
//		m.Grid[first.X+dx][first.Y+dy] = Road
//	}
//}

func (m *Maze) generateEnvironment(path []Cell) {
	m.generateEntranceExit()
	//m.generateTraps()
}

func (m *Maze) generateEntranceExit() {
	m.Grid[m.Entrance.X][m.Entrance.Y] = Entrance
	m.Grid[m.Exit.X][m.Exit.Y] = Exit
}

func (m *Maze) generateTraps() {
	trapsLeft := m.Traps
	trapsPath := 0
	for trapsLeft > 0 {
		randomX := 1 + rand.Intn(m.Height-2)
		randomY := 1 + rand.Intn(m.Width-2)
		if m.Grid[randomX][randomY] == Road && trapsPath < 2 {
			m.Grid[randomX][randomY] = Trap
			trapsPath++
			trapsLeft--
			continue
		} else if m.Grid[randomX][randomY] == Visited {
			m.Grid[randomX][randomY] = Trap
			trapsLeft--
		}
	}
}

func (m *Maze) generateTreasure() {

}

func (m *Maze) Print() {
	for _, row := range m.Grid {
		fmt.Println(row)
	}
	for _, row := range m.Grid {
		for _, cell := range row {
			fmt.Printf("%3s", cellChars[cell])
		}
		fmt.Println()
	}
}

func main() {
	maze := NewMaze(15, 15)
	maze.GenerateMaze()
	maze.Print()
}
