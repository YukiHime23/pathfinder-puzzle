package main

import (
	"log"
	"net/http"

	"github.com/YukiHime23/pathfinder-puzzle/internal/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Cho phép mọi nguồn kết nối
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Lỗi nâng cấp WebSocket:", err)
		return
	}
	defer ws.Close()

	log.Println("Người chơi đã kết nối thành công!")

	// Tạo game engine
	gameEngine := game.NewGame(21, 21)

	// Gửi trạng thái ban đầu để client hiển thị
	ws.WriteJSON(gameEngine)

	// Trong vòng lặp WebSocket
	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		if action, ok := msg["action"].(string); ok {
			switch action {
			case "move":
				if dir, ok := msg["direction"].(string); ok {
					gameEngine.Move(dir) // Gọi logic xử lý
				}
			case "restart":
				w := 21
				h := 21
				if widthF, ok := msg["width"].(float64); ok {
					w = int(widthF)
				}
				if heightF, ok := msg["height"].(float64); ok {
					h = int(heightF)
				}

				// Đảm bảo kích thước tối thiểu, tối đa và là số lẻ để thuật toán tạo mê cung hoạt động đúng
				if w < 5 {
					w = 5
				}
				if w > 51 {
					w = 51
				}
				if w%2 == 0 {
					w++
				}

				if h < 5 {
					h = 5
				}
				if h > 51 {
					h = 51
				}
				if h%2 == 0 {
					h++
				}

				gameEngine = game.NewGame(w, h)
			}
		}

		// Gửi lại trạng thái cập nhật cho người chơi
		ws.WriteJSON(gameEngine)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	log.Println("Server đang chạy tại http://localhost:8989")
	http.ListenAndServe(":8989", nil)
}
