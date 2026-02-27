# Pathfinder Puzzle

Một trò chơi giải đố mê cung theo thời gian thực (Real-time Maze Puzzle) được xây dựng bằng Golang và WebSocket. Người chơi sẽ điều khiển một nhân vật di chuyển qua một mê cung được tạo ngẫu nhiên để tìm đường đến đích.

## Tính năng (Features)

*   **Thời gian thực (Real-time)**: Giao tiếp hai chiều tốc độ cao giữa Client và Server thông qua WebSocket. Mọi di chuyển đều được xử lý ở phía Server đảm bảo tính nhất quán và ngăn chặn gian lận.
*   **Mê cung tạo ngẫu nhiên (Procedural Maze Generation)**: Sử dụng thuật toán đục tường (Recursive Backtracker) để tạo ra các mê cung mới hoàn toàn mỗi khi bắt đầu, không có ván nào giống ván nào.
*   **Tùy chỉnh Kích thước**: Người chơi có thể tự do thay đổi kích thước mê cung từ quy mô nhỏ (5x5) đến quy mô lớn (51x51).
*   **Sương mù (Fog of War)**: Cơ chế che khuất tầm nhìn tự động ở Client, chỉ cho phép người chơi nhìn thấy một phạm vi nhất định xung quanh mình, làm tăng độ khó, bí ẩn và sự thú vị.
*   **Giao diện hiện đại**: Frontend được xây dựng bằng cấu trúc HTML5 đơn giản kết hợp CSS Grid để render với hiệu suất cao, thiết kế giao diện tối tối (Dark Mode) kết hợp với TailwindCSS.

## Công nghệ sử dụng (Tech Stack)

*   **Backend**: Go (Golang)
*   **Thư viện**: `github.com/gorilla/websocket` (Xử lý giao tiếp WebSocket)
*   **Frontend**: HTML, Vanilla JavaScript, CSS (TailwindCSS CDN)

## Cấu trúc dự án (Project Structure)

```text
pathfinder-puzzle/
├── main.go                 # Điểm khởi chạy của Server, chạy Web server và router WebSocket.
├── go.mod / go.sum         # Quản lý module và thư viện phụ thuộc của Go.
├── README.md               # Tài liệu dự án.
├── internal/
│   └── game/               # Chứa core logic của trò chơi chạy trên Server.
│       ├── engine.go       # Quản lý State của Game, vị trí người chơi và xử lý logic đi/thắng.
│       └── maze.go         # Chứa lõi Thuật toán tạo và định dạng Mê cung.
├── static/                 # Chứa các file tĩnh thiết yếu để phục vụ cho Frontend.
│   └── index.html          # Giao diện chính của trò chơi và mã JavaScript tương tác.
```

## Cách thức hoạt động (How It Works)

Trò chơi áp dụng mô hình kiến trúc **Client-Server truyền tải trạng thái qua WebSockets**.

1.  **Kết nối**: Khi người chơi mở URL trên trình duyệt, tệp `index.html` được tải về và mở kết nối WebSocket `ws://localhost:8989/ws` đến Go Server.
2.  **Khởi tạo**: Backend Server chấp nhận nâng cấp kết nối. Nó tạo ra một thể hiện `GameEngine` mới (mặc định với kích thước 21x21), tạo một mảng mê cung 2D, chọn điểm bắt đầu và điểm kết thúc, sau đó gửi toàn bộ cấu trúc dữ liệu JSON này lần đầu cho Client.
3.  **Hiển thị (Rendering)**: JavaScript ở Client phân tích file JSON, và lặp qua toàn bộ ma trận (mảng `Grid`) để tạo ra các khối `div`. Bằng việc sử dụng kỹ thuật CSS Grid, hệ thống sẽ render bức tường, đường đi, vị trí nhân vật (màu xanh) và ô mục tiêu (màu đỏ). Quá trình render bao gồm cả việc tính khoảng cách Euclid (distance) từ ô nhân vật đến các ô xung quanh để áp dụng tính năng **Fog of War** (phủ css màu đen).
4.  **Tương tác (Gameplay)**:
    *   Client lắng nghe sự kiện `keydown` (các phim mũi tên lên, xuống, trái, phải).
    *   Mỗi khi bấm phím, client gửi một payload JSON gọn nhẹ: `{"action": "move", "direction": "up"}`.
    *   Go Server nhận tín hiệu, xác minh tọa độ xem phía trước người chơi có phải là vách tường hay viền bản đồ hay không. Nếu hợp lệ, hệ thống tịnh tiến đồ thị người chơi, cộng thêm số bước đi (Steps), kiểm tra điều kiện Win.
    *   Server phát sóng ngược toàn bộ khối lệnh `GameEngine` mới cập nhật cho Client. Client nhận được và xóa DOM cũ, render lại DOM mới dưới 1 millisecond mang cảm giác tức thì (real-time realtime feedback loop).

## Thuật toán tạo mê cung (Maze Generation Algorithm)

Toàn bộ hệ thống bản đồ trong game được sinh ra linh hoạt (procedural generation) bằng thuật toán **Recursive Backtracker** (Quay lui đệ quy). Thuật toán này tạo ra những mê cung "hoàn hảo" (perfect maze) - nghĩa là không bao giờ có phòng kín, không có vòng lặp (vòng kết nối tuần hoàn), và luôn tồn tại chính xác 1 đường đi duy nhất giữa 2 điểm bất kỳ.

**Quy trình đục tường của thuật toán:**
1.  **Khởi tạo khối:** Toàn bộ bản đồ cấu tạo ban đầu là vách tường nguyên khối (grid mang giá trị `1`). Kích thước đồ họa thực tế bắt buộc phải là số lẻ.
2.  **Bắt đầu:** Chọn ngẫu nhiên một điểm lẻ $(x, y)$ và "đục" nó dọn thành lối đi (giá trị `0`).
3.  **Khám phá 4 hướng:** Ở vị trí hiện tại, xét 4 phía (Lên, Xuống, Trái, Phải). Các điểm lân cận cách điểm hiện tại 2 ô (`+/- 2`) được đánh giá xem có hướng nào chưa bị đục hay không.
4.  **Tuyển chọn ngẫu nhiên:** Thuật toán ngẫu nhiên chọn một hướng hợp lệ. Nó sẽ đục ô lân cận đó (`+/- 2`) *và* đục luôn ô bức tường nằm ở chính giữa để thông nối 2 khoảng sân. 
5.  **Áp dụng Đệ quy:** Tiếp tục gọi đệ quy (recursive) lặp lại bước 3 từ vị trí mới này. Cứ như thế, thuật toán giống như một con cào cào máy đào đường mù quáng đến chừng nào đâm vào rào hoặc đường đã bị đào.
6.  **Quay lui (Backtracking):** Khi đào tới một ngõ cụt (dead-end - nơi xung quanh đều là khoảng vách mà đằng sau là đường đã đục), thuật toán trượt ngược (quay lui) trở lại những ô ngã rẽ trước đó, nhằm tìm ra những nhánh chưa được đào khác và tiếp tục cho đến khi mọi ô vuông lẻ trên ma trận đều đã được đi qua.

Nhờ việc tráo đổi mảng thứ tự các hướng đi ngẫu nhiên bằng `rand.Shuffle` ở bước 4, mỗi một truy vấn _"Generate New Maze"_ đều bảo đảm kết tinh ra những con đường rắc rối và hỗn mang hoàn toàn khác biệt!

## Hướng dẫn chạy cục bộ (How to Run)

### Yêu cầu môi trường
*   Môi trường máy tính của bạn đã được cài đặt **Go** (phiên bản khuyến nghị là 1.21+).

### Các lệnh thực thi
1.  Clone hoặc tải thư mục mã nguồn về.
2.  Mở Terminal (hoặc Command Prompt) và truy cập vào gốc dự án `pathfinder-puzzle/`.
3.  Tải các dependency packages:
    ```bash
    go mod tidy
    ```
4.  Cấp quyền biên dịch và khởi chạy Game Server qua cổng 8989:
    ```bash
    go run main.go
    ```
5.  Mở bất kỳ trình duyệt web nào và truy cập:
    ```text
    http://localhost:8989
    ```
6.  Enjoy the game! Dùng các **Phím mũi tên** để điều khiển, sử dụng thanh điều khiển phía trên để tắt hiệu ứng Sương mù và thay đổi ngẫu nhiên sang các mê cung kích thước khác nhau.

### AI Support
- Gemini 3.1 Pro