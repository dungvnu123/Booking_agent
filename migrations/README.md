# Database migrations

Migration `000001_identity.up.sql` tạo các bảng identity cần cho auth: `roles`, `users`, `patients`, `invalidated_token`; seed tên role và thêm unique index email/điện thoại tài khoản. Không seed tài khoản hoặc mật khẩu mặc định.

Chạy khi PostgreSQL đã sẵn sàng, từ thư mục dự án:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\migrate.ps1
```

Hoặc `go run ./cmd/migrate` khi môi trường Go đã được cấu hình. Lệnh dùng `DATABASE_URL` trong `.env`, ghi version vào `schema_migrations` và bao toàn bộ thay đổi trong một transaction có khóa chống chạy trùng. Chạy lại không áp dụng version đã ghi nhận. API và worker không tự chạy migration khi startup.

Tên cột identity được đối chiếu với entity Java gốc. Nếu dùng DB đã có dữ liệu, kiểm tra schema và dữ liệu trùng trước: index email không phân biệt hoa/thường và phone có thể làm migration dừng. Migration không tự xóa hoặc sửa dữ liệu trùng. Dùng DB riêng cho phát triển nếu chưa chốt chuyển dữ liệu.

Test PostgreSQL auth trong `tests/auth_test.go` dùng schema ngẫu nhiên riêng và kiểm tra migration lặp lại, transaction đăng ký bệnh nhân, đăng ký đồng thời và thu hồi token. Không trỏ test vào database không cho phép tạo schema.

Các migration nghiệp vụ khác chưa được viết. Khi port booking cần transaction và UPDATE nhận/nhả chỗ có điều kiện; không áp constraint `booked_count <= capacity` vì lễ tân có `allowOverCapacity`.

## Migration 2 — xác thực email

`000002_email_verification.up.sql` thêm `users.email_verified`, `email_verified_at` và bảng `email_verifications`. Tất cả email cũ ban đầu chưa xác thực; migration không tự đặt thành verified. Tài khoản giữ nguyên trạng thái `is_active`, dữ liệu hồ sơ và mật khẩu.

Bảng token chỉ lưu hash, thời hạn và mốc gửi gần nhất. Xác nhận email khóa hàng user/token trong transaction, kiểm tra thời hạn rồi cập nhật trạng thái và xóa token; không dùng lại được link. Gửi lại thay thế token sau thời gian chờ. Chạy cùng lệnh migration ở trên; runner áp dụng tuần tự version chưa có và không tự chạy lúc API khởi động.
