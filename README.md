# ABClinic backend — Go + Gin

Khung cho một backend Go chia module nghiệp vụ, PostgreSQL dùng chung, API và worker là hai entry point. Không gồm AI. Theo [bản thiết kế](docs/backend-go-blueprint.md) và [chỉ mục nguồn](docs/backend-go-inventory.md) bạn cung cấp.

Repository Booking_agent phát triển nền tảng đặt lịch cho ABClinic. Giai đoạn hiện tại xây backend nghiệp vụ; giai đoạn sau bổ sung agent tư vấn lịch phù hợp và thực hiện đặt lịch theo xác nhận của người dùng. Xem [hướng phát triển](docs/roadmap.md).

**Hiện chỉ là khung:** 108 route nghiệp vụ được đăng ký dưới /api và trả HTTP 501. Chưa có đăng nhập/JWT verifier, truy vấn nghiệp vụ, CRUD thật, xếp lịch hay job chạy tự động. Swagger mô tả rõ trạng thái này. Không dùng dữ liệu demo để giả làm chức năng hoàn chỉnh.

## Chạy ngay trên Windows

Yêu cầu Go 1.27+. Mở terminal tại Backend-Golang:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\run.ps1
```

Setup tạo .env nếu chưa có, với mật khẩu PostgreSQL ngẫu nhiên; không ghi đè .env đã có. API vẫn mở được Swagger khi chưa chạy database. Nếu cổng 8081 đang có server cũ, dừng server cũ trước; Go không tự reload code.

Swagger: http://localhost:8081/swagger/index.html
Alias tương tự backend Java: http://localhost:8081/api/swagger-ui.html
OpenAPI: http://localhost:8081/api/v3/api-docs

Health: /healthz trả 200 khi process chạy.
Readiness: /readyz trả 503 nếu chưa có kết nối PostgreSQL; trả 200 khi DB truy cập được, không có nghĩa nghiệp vụ đã hoàn thành.

## PostgreSQL và Docker

Bật Docker Desktop với Linux containers. Cổng PostgreSQL trên máy là **5433**, trong mạng Docker là **5432**. Không dùng cổng 5435 từ cấu hình Java cũ.

Chỉ chạy DB để code Go trên máy:

```powershell
docker compose up -d postgres
```

Chạy cả API và DB bằng Docker (dừng API local nếu cùng cổng):

```powershell
docker compose up --build -d
docker compose logs -f api
```

Solver Python chưa được khởi động bởi Compose này. ROSTERING_URL dành cho Go chạy local; ROSTERING_DOCKER_URL dành cho container (mặc định host.docker.internal:8000). Client solver hiện chưa gửi request.

Nếu dùng tên user/database/cổng khác, sửa DATABASE_URL trong .env tương ứng. Mật khẩu ngẫu nhiên ban đầu chỉ dành cho môi trường phát triển; .env bị gitignore. Đổi mật khẩu .env sau khi đã có volume PostgreSQL không tự đổi mật khẩu trong database.

Không có auto-migrate hoặc SQL tạo bảng khi startup. Xem [migrations](migrations/README.md) để bắt đầu sau khi xác minh schema thật.

## Cấu trúc đọc code

```text
cmd/
  api/main.go
  worker/main.go
internal/
  app/                      # Khởi tạo và ghép dependency, đăng ký router
  config/                   # .env và biến môi trường, validation
  platform/
    postgres/               # pgx pool + transaction helper
    httpx/                  # Response chung và ánh xạ lỗi scaffold
    middleware/             # Request ID, CORS, ranh giới auth/role
    clock/                  # Clock thật/fixed; Asia/Ho_Chi_Minh
    storage/                # Interface lưu file, chưa có adapter
    swagger/                # Swagger UI nhúng, chạy không cần CDN
  auth/
  user/
  doctor/
  patient/
  medicine/
  examination/
  roster/
  leave/
  closure/
  booking/
  notification/
  attachment/
  integration/rostering/    # Ranh giới gọi OR-Tools, chưa gửi request
  jobs/                     # Danh mục job; tất cả đang tắt
api/openapi.json            # Contract scaffold cho 108 route
migrations/                 # Chưa có DDL nghiệp vụ
tests/integration/          # Test DB tùy chọn
docs/                       # Tài liệu nguồn, tiến độ
```

Trong mỗi module:

- handler.go: đăng ký route và chuyển request sang service.
- service.go: tên use case tương ứng inventory; hiện trả ErrNotImplemented.
- dto.go: chỗ bổ sung request/response đã xác minh, chưa đoán field từ tên DTO.
- model.go: mô hình domain ban đầu, chưa đầy đủ cột/constraint của DB.
- repository.go: chỗ định nghĩa interface theo use case thật, chưa tạo empty interface.
- postgres.go: chỗ viết SQL sau khi xác minh schema, hiện không chạy truy vấn.

Đọc main.go → app.go → module/handler.go → service.go. Khi có nghiệp vụ thật, service gọi repository; transaction bao trùm toàn use case.

## Contract và phạm vi

Tiền tố /api, không phải /api/v1 của demo trước. User ID là chuỗi UUID; patient ID là chuỗi nghiệp vụ; ID khác theo inventory. pgtype.Date giữ kiểu date và pgtype.Numeric giữ tiền thập phân; không dùng float64. DTO JSON phải camelCase và tách khỏi model.

Response scaffold dùng statusCode=1501 và HTTP 501. Đây là mã tạm của khung, chưa xác nhận mã lỗi Java; không coi OpenAPI này là contract tương thích frontend hoàn chỉnh. Examination/file và pagination cần kiểm tra riêng khi port.

Middleware auth chỉ nhận TokenVerifier qua interface và từ chối khi thiếu verifier. Chưa có JWT verifier thật, chưa gắn quyền cho các stub vì chúng chỉ trả 501. Trước khi triển khai từng API cần gắn kiểm tra chữ ký/issuer/expiry/jti thu hồi/user active, role và quyền trên từng hồ sơ; không mặc định ADMIN có mọi quyền.

Không đưa AI/OCR/Task, patient AI summary/history/session vào router. Attachment đích gắn trực tiếp patient và examination; chưa migrate bảng cũ. Không có Redis/broker vì tài liệu hiện chưa yêu cầu. Solver OR-Tools độc lập với AI, được giữ ở mức ranh giới tích hợp.

## Worker và kiểm tra

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\run.ps1 worker
powershell -ExecutionPolicy Bypass -File .\scripts\check.ps1
```

Worker hiện chỉ chờ tín hiệu dừng; 4 job nghiệp vụ đều chưa bật. Khi triển khai nhiều worker cần DB lock hoặc một chủ sở hữu job.

Khi API đang chạy, terminal khác:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\test-api.ps1
```

Test smoke kiểm tra Swagger, health và 108 route trả đúng 501, không ghi dữ liệu bệnh nhân. Test PostgreSQL integration chỉ chạy khi có TEST_DATABASE_URL trỏ database kiểm thử; nếu không sẽ skip.

Thứ tự triển khai đề xuất: auth/user/doctor → patient/medicine/examination/attachment → roster/solver/publish → booking/queue → leave/closure/notification/jobs.
