package game

type Player struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type GameEngine struct {
	Maze   *Maze  `json:"maze"`
	Player Player `json:"player"`
	Steps  int    `json:"steps"`
	Won    bool   `json:"won"`
}

// Khởi tạo Game mới
func NewGame(width, height int) *GameEngine {
	m := NewMaze(width, height)
	m.Generate(1, 1) // Tạo mê cung

	// Đặt đích đến ở góc xa nhất
	m.Grid[height-2][width-2] = 2

	// Khởi tạo GameEngine
	engine := &GameEngine{
		Maze:   m,
		Player: Player{X: 1, Y: 1}, // Bắt đầu tại (1,1)
		Steps:  0,
		Won:    false,
	}

	return engine
}

func (g *GameEngine) Move(direction string) {
	if g.Won {
		return
	} // Game đã kết thúc thì không đi nữa

	newX, newY := g.Player.X, g.Player.Y

	switch direction {
	case "up":
		newY--
	case "down":
		newY++
	case "left":
		newX--
	case "right":
		newX++
	}

	// 1. Kiểm tra biên và va chạm tường
	if newY >= 0 && newY < g.Maze.Height && newX >= 0 && newX < g.Maze.Width {
		if g.Maze.Grid[newY][newX] != 1 { // Nếu không phải tường
			g.Player.X = newX
			g.Player.Y = newY
			g.Steps++

			// 2. Kiểm tra nếu chạm đích (ô giá trị là 2)
			if g.Maze.Grid[newY][newX] == 2 {
				g.Won = true
			}
		}
	}
}
