package main

import (
	"container/list"
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
	Treasure bool
}

type Cell struct {
	X int
	Y int
}

type BackTracker []Cell

func (b *BackTracker) IsEmpty() bool {
	return len(*b) == 0
}

func (b *BackTracker) Push(cell Cell) {
	*b = append(*b, cell)
}

func (b *BackTracker) Pop() (Cell, bool) {
	if b.IsEmpty() {
		return Cell{-1, -1}, false
	} else {
		index := len(*b) - 1
		element := (*b)[index]
		*b = (*b)[:index]
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
		Treasure: true,
	}
}

func (m *Maze) GenerateMaze() {
	for i := 0; i < m.Height; i++ {
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
			m.Grid[currentCell.X][currentCell.Y] = Visited
		}
	}
	m.generateEnvironment()
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
	m.Grid[first.X+dx][first.Y+dy] = Visited
}

func (m *Maze) isReached(cell Cell) bool {
	return cell.X == m.Exit.X && cell.Y == m.Exit.Y
}

func (m *Maze) generateEnvironment() {
	m.bfsShortestPath()
	m.generateEntranceExit()
	m.generateTraps()
	m.generateTreasure()
}

func (m *Maze) bfsShortestPath() {
	visited := make([][]bool, m.Height)
	for i := range visited {
		visited[i] = make([]bool, m.Width)
	}

	queue := list.New()
	queue.PushBack(m.Entrance)
	visited[m.Entrance.X][m.Entrance.Y] = true

	parent := make(map[Cell]Cell)

	movements := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for queue.Len() > 0 {
		curr := queue.Remove(queue.Front()).(Cell)

		if curr.X == m.Exit.X && curr.Y == m.Exit.Y {
			path := []Cell{curr}
			for curr != m.Entrance {
				curr = parent[curr]
				path = append([]Cell{curr}, path...)
			}
			for _, i := range path {
				m.Grid[i.X][i.Y] = Road
			}
			break
		}

		for _, move := range movements {
			newX := curr.X + move[0]
			newY := curr.Y + move[1]

			if newX > 0 && newX < m.Height && newY > 0 && newY < m.Width && !visited[newX][newY] && m.Grid[newX][newY] == Visited {
				nextCell := Cell{newX, newY}
				queue.PushBack(nextCell)
				visited[newX][newY] = true
				parent[nextCell] = curr
			}
		}
	}
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
	treasure := m.Treasure
	for treasure {
		randomX := 1 + rand.Intn(m.Height-2)
		randomY := 1 + rand.Intn(m.Width-2)
		if m.Grid[randomX][randomY] == Road || m.Grid[randomX][randomY] == Visited {
			m.Grid[randomX][randomY] = Treasure
			break
		}
	}
}

func (m *Maze) Print() {
	//for _, cell := range m.Grid {
	//	fmt.Println(cell)
	//}
	for _, row := range m.Grid {
		for _, cell := range row {
			fmt.Printf("%3s", cellChars[cell])
		}
		fmt.Println()
	}
}

func main() {
	var height, width int
	fmt.Print("Enter height and width of maze:\n")
	_, err := fmt.Scanf("%d %d", &height, &width)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	maze := NewMaze(height, width)
	maze.GenerateMaze()
	maze.Print()
}
