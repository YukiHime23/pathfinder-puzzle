package game

import (
	"math/rand"
	"time"
)

type Maze struct {
	Width  int
	Height int
	Grid   [][]int
}

// Khởi tạo mê cung toàn tường (1)
func NewMaze(width, height int) *Maze {
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
		for j := range grid[i] {
			grid[i][j] = 1
		}
	}
	return &Maze{Width: width, Height: height, Grid: grid}
}

func (m *Maze) Generate(x, y int) {
	m.Grid[y][x] = 0 // Đánh dấu ô hiện tại là đường đi

	// Định nghĩa 4 hướng di chuyển (Lên, Xuống, Trái, Phải)
	// Chúng ta nhảy 2 bước để giữ lại các bức tường giữa các đường đi
	dirs := [][2]int{{0, 2}, {0, -2}, {2, 0}, {-2, 0}}

	// Xáo trộn hướng đi để mê cung ngẫu nhiên mỗi lần chạy
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]

		// Kiểm tra xem ô mới có nằm trong bản đồ và chưa được đào không
		if nx > 0 && nx < m.Width-1 && ny > 0 && ny < m.Height-1 && m.Grid[ny][nx] == 1 {
			// Đào xuyên qua bức tường ở giữa ô hiện tại và ô mới
			m.Grid[y+d[1]/2][x+d[0]/2] = 0
			// Đệ quy để tiếp tục đào từ ô mới
			m.Generate(nx, ny)
		}
	}
}
