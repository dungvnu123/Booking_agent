# Trạng thái khung Booking App

Đã có: Go module chung, Gin router theo nghiệp vụ, API/worker, config .env, timezone, PostgreSQL pool/transaction helper, CORS/request ID, ranh giới auth/role fail-closed, Swagger nhúng, Docker/Compose và script Windows.

Có **108** route nghiệp vụ. Tất cả trả **501**; không route nào giả trả thành công. Xem api/openapi.json và handler/service của từng module.

Chưa có: DTO Java thực tế, full DB mapping/schema snapshot, migration nghiệp vụ, repository SQL, JWT/BCrypt, phân quyền endpoint/ownership, business errors cuối cùng, CRUD, booking concurrency, solver DTO/call, storage adapter, job chạy thật.

Các file dto.go/repository.go/postgres.go chỉ là vị trí bổ sung code, có TODO cụ thể; chưa có implementation giả hoặc interface rỗng.

Bản thiết kế và inventory được giữ nguyên nội dung do người dùng cung cấp. Link ../backend/... trong hai tài liệu trỏ source Java gốc, source đó không có trong workspace Go này.

Khác biệt có chủ ý: chuẩn hóa tên path parameter monthlyRosterId thành id ở nhánh /roster/monthly, vì Gin yêu cầu tên parameter thống nhất trong cùng cây route; URL thực tế không thay đổi. Không suy luận loại ID từ cùng tên parameter, cần đối chiếu từng DTO khi triển khai.

Migrations chưa có SQL vì chưa có schema thật đã xác minh. Không tự áp DDL lên DB cũ. Các bảng có liên kết AI cũ chưa được migrate hoặc xóa.

Bản demo user/catalog và các tài nguyên không thuộc Booking App được bỏ theo yêu cầu; không tiếp tục đăng ký API /api/v1/users.
