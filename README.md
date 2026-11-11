# Booking App backend ? Go + Gin

Khung cho một backend Go chia module nghiệp vụ, PostgreSQL dùng chung, API và worker là hai entry point. Không gồm AI. Theo [bản thiết kế](docs/backend-go-blueprint.md) và [chỉ mục nguồn](docs/backend-go-inventory.md) bạn cung cấp.

Repository Booking_agent phát triển nền tảng đặt lịch cho Booking App. Giai đoạn hiện tại xây backend nghiệp vụ; giai đoạn sau bổ sung agent tư vấn lịch phù hợp và thực hiện đặt lịch theo xác nhận của người dùng. Xem [hướng phát triển](docs/roadmap.md).

**Auth đã được triển khai:** 8 route `/api/auth/*` có đăng ký, gửi mail/xác thực email, đăng nhập, refresh token và logout. 103 route nghiệp vụ còn lại vẫn trả HTTP 501. Swagger phân biệt rõ route đã làm và route scaffold.

## Chạy ngay trên Windows

Yêu cầu Go 1.27+. Mở terminal tại Backend-Golang:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\migrate.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\run.ps1
```

Setup giữ các giá trị `.env` đã có, tạo mật khẩu PostgreSQL khi tạo file lần đầu và bổ sung hai khóa JWT ngẫu nhiên nếu chưa có. Chạy migration khi PostgreSQL đã sẵn sàng; API không tự tạo bảng. API vẫn mở được Swagger khi chưa chạy database. Nếu cổng 8081 đang có server cũ, dừng server cũ trước; Go không tự reload code.

Swagger: http://localhost:8081/swagger/index.html
Alias tương tự backend Java: http://localhost:8081/api/swagger-ui.html
OpenAPI: http://localhost:8081/api/v3/api-docs

Health: /healthz trả 200 khi process chạy.
Readiness: /readyz trả 503 nếu chưa có kết nối PostgreSQL; trả 200 khi DB truy cập được, không có nghĩa nghiệp vụ đã hoàn thành.

## PostgreSQL và Docker

Bật Docker Desktop với Linux containers. Cổng PostgreSQL riêng trên máy là **55433**, trong mạng Docker là **5432**.

Compose project riêng là `booking-agent-abclinic`: container PostgreSQL và volume dữ liệu được tách khỏi các dự án khác. Volume database là `booking-agent-abclinic_postgres-data`; không dùng lại volume `abclinic_pgdata` của môi trường cũ.

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

Không có auto-migrate hoặc SQL tạo bảng khi startup. Migration identity chạy chủ động bằng `scripts/migrate.ps1`; xem [migrations](migrations/README.md).

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
api/openapi.json            # Contract cho 111 route, gồm 8 route auth đã triển khai
migrations/                 # Migration identity chạy chủ động
tests/integration/          # Test DB tùy chọn
docs/                       # Tài liệu nguồn, tiến độ
```

Trong các module còn là scaffold (auth đã có triển khai thật):

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

Service auth triển khai interface `middleware.TokenVerifier`: kiểm tra HS256, issuer, expiry, jti thu hồi và user còn hoạt động; role được đọc lại từ DB. Khi triển khai các route còn lại, dùng `middleware.Authenticate(application.Auth)`, kiểm tra role và quyền trên từng hồ sơ. ADMIN không tự có mọi quyền.

Không đưa AI/OCR/Task, patient AI summary/history/session vào router. Attachment đích gắn trực tiếp patient và examination; chưa migrate bảng cũ. Không có Redis/broker vì tài liệu hiện chưa yêu cầu. Solver OR-Tools độc lập với AI, được giữ ở mức ranh giới tích hợp.

## Auth

Giữ nguyên tiền tố `/api` và response `statusCode/message/data` của luồng hiện tại.

| Method | Route | Hành vi |
|---|---|---|
| POST | /api/auth/register | Tạo tài khoản PATIENT bằng email, password, fullName |
| POST | /api/auth/register-patient | Thêm phone; tạo hoặc liên kết hồ sơ bệnh nhân trong cùng transaction |
| POST | /api/auth/login | Trả accessToken, refreshToken và user |
| POST | /api/auth/refresh_token | Nhận accessToken/refreshToken cùng user; thu hồi cặp cũ và cấp cặp mới |
| POST | /api/auth/logout | Thu hồi cả hai token, có thể gọi lặp lại |
| GET | /api/auth/verify-email | Mở trang xác nhận từ link trong mail; không thay đổi trạng thái |
| POST | /api/auth/verify-email | Nhận token, xác thực email một lần |
| POST | /api/auth/resend-verification | Nhận email, gửi lại link theo thời gian chờ |

Email được trim và chuyển chữ thường. Mật khẩu khi đăng ký phải từ **8–32 ký tự**, có **chữ thường, chữ hoa, số và ký tự đặc biệt** (dấu câu hoặc ký hiệu; khoảng trắng không thay thế ký tự đặc biệt). Mật khẩu được giữ nguyên, không trim, tối đa 72 byte UTF-8 theo giới hạn BCrypt và lưu với cost 10. Ví dụ hợp lệ: `Password123!`. Mật khẩu không đạt trả HTTP 400 kèm thông báo quy tắc; chính sách tạo mật khẩu mới không chặn đăng nhập tài khoản cũ. Role đăng ký cố định PATIENT; client không được tự gán ADMIN. Đăng ký tài khoản nhân viên thuộc module user, chưa triển khai.

Access token mặc định 15 phút, refresh token 168 giờ. Khóa `JWT_ACCESS_SECRET` và `JWT_REFRESH_SECRET` phải khác nhau, ít nhất 32 byte; setup tạo khóa ngẫu nhiên, không có khóa fallback trong code. Thiếu khóa thì Swagger vẫn mở nhưng auth trả 503 cho request hợp lệ. DB dùng các tên cột đã đối chiếu với nguồn Java: `users.password_hash`, `roles.role_id`, `patients.user_id`, `invalidated_token.id_token`.

Đăng ký bệnh nhân nhận số điện thoại dạng `0` và 9 chữ số. Nếu có một hồ sơ chưa liên kết cùng số điện thoại thì gắn account mới vào hồ sơ đó; nhiều hồ sơ trùng số điện thoại hoặc hồ sơ đã có account sẽ bị từ chối. Lỗi giữa chừng rollback cả tài khoản và liên kết hồ sơ. Chưa xác minh số điện thoại bằng OTP.

Refresh token chỉ dùng một lần, kể cả hai request đến đồng thời. Access token có thể đã hết hạn khi gọi refresh; logout và khóa tài khoản có hiệu lực trên lần xác minh kế tiếp. Không tự chạy job dọn token hết hạn trong bước auth này.

Chạy riêng test auth:

```powershell
go test ./tests -v
```

Nếu muốn kiểm tra SQL thật, đặt `TEST_DATABASE_URL` trỏ PostgreSQL kiểm thử trước khi chạy. User DB cần quyền tạo schema. Các test mới nằm trong một file `tests/auth_test.go`, tách khỏi source `internal/auth/`.


## Xác thực email qua SMTP

Đăng ký luôn tạo tài khoản với `emailVerified=false`. Chỉ cấp/refresh/xác minh JWT khi email đã được xác thực và tài khoản còn hoạt động. Link xác thực có hiệu lực **15 phút**, chỉ dùng một lần; gửi lại cách nhau ít nhất **60 giây**. Database chỉ lưu SHA-256 của token ngẫu nhiên 32 byte.

Sau khi chỉnh code, chạy `scripts/migrate.ps1` để áp dụng migration 2. Tài khoản đã có cũng được đánh dấu chưa xác thực; không tự coi email cũ đã có bằng chứng sở hữu. Dùng API gửi lại link để xác thực tài khoản cũ.

Điền SMTP vào `.env` của chính dự án này. Ví dụ Gmail:

```dotenv
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-google-app-password
SMTP_FROM=your-email@gmail.com
SMTP_ENCRYPTION=starttls
SMTP_TIMEOUT=10s
EMAIL_VERIFICATION_URL=http://localhost:8081/api/auth/verify-email
EMAIL_VERIFICATION_TTL=15m
EMAIL_RESEND_COOLDOWN=60s
```

Gmail dùng [SMTP cổng 587 với STARTTLS hoặc 465 với TLS](https://developers.google.com/workspace/gmail/imap/imap-smtp). Với adapter username/password này, dùng [App Password của Google](https://support.google.com/accounts/answer/185833), không dùng mật khẩu đăng nhập Gmail thông thường; tài khoản cần bật xác minh hai bước và có quyền tạo App Password. Không đưa mật khẩu vào source hoặc gửi qua chat. Với nhà cung cấp khác, điền SMTP tương ứng. Không có fallback gửi mật khẩu qua kết nối không mã hóa.

Nếu chưa cấu hình SMTP, Swagger vẫn mở nhưng đăng ký/gửi lại mail trả **503**, không âm thầm bỏ bước xác thực. Chương trình không thể gửi thư Gmail thật khi chưa có thông tin tài khoản gửi. Khởi động lại API sau khi sửa `.env`. `EMAIL_VERIFICATION_URL` phải là URL người nhận mở được: localhost chỉ dùng khi mở mail trên cùng máy; khi triển khai dùng địa chỉ HTTPS công khai của API.

Luồng thử trên Swagger:

1. `POST /api/auth/register-patient` với email bạn nhận được thư, mật khẩu mạnh, họ tên và số điện thoại.
2. Đăng ký trả **201** với `emailVerified=false`; đăng nhập lúc này trả **403**, `statusCode=2005`.
3. Mở link trong hộp thư và bấm **Xác nhận email**; sau đó đăng nhập trả **200** và cặp JWT.
4. Nếu cần gửi lại, gọi `POST /api/auth/resend-verification` với `{ "email": "your-email@example.com" }`.
5. Có thể gọi trực tiếp `POST /api/auth/verify-email` với `{ "token": "token lấy từ link email" }`.

Link chứa token trong fragment `#token=...`, không phải query; token không đi vào access log của request GET. Trang xác nhận xóa fragment khỏi thanh địa chỉ, gửi token bằng POST, không tự kích hoạt bằng GET. API không trả token xác thực trong response và không in link vào log.

Gửi lại trả cùng response **202** cho email không tồn tại, đã xác thực hoặc đang trong thời gian chờ. Link mới thay thế link cũ. Nếu SMTP lỗi sau khi tài khoản đã lưu, API trả **503** với thông báo tài khoản đang chờ xác thực; sửa cấu hình gửi mail rồi dùng API gửi lại sau thời gian chờ. Không đăng ký lại để tránh lỗi email trùng.

Toàn bộ test auth/email/SMTP mới vẫn nằm trong **`tests/auth_test.go`**. SMTP được test bằng máy chủ TLS/STARTTLS cục bộ; không gửi thư ra ngoài khi chạy test. Test PostgreSQL dùng schema riêng, kiểm tra hash token, hết hạn, dùng một lần, gửi lại và transaction.

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

Test smoke kiểm tra Swagger, health, 103 route scaffold trả 501 và 7 route auth POST từ chối body rỗng bằng 400 và trang xác nhận email trả HTML 200, không ghi dữ liệu bệnh nhân. Test auth mới nằm riêng tại `tests/auth_test.go`. Test PostgreSQL chỉ chạy khi có `TEST_DATABASE_URL`; test auth tạo và xóa schema ngẫu nhiên riêng, không dùng bảng ứng dụng.

Thứ tự triển khai đề xuất: auth/user/doctor → patient/medicine/examination/attachment → roster/solver/publish → booking/queue → leave/closure/notification/jobs.
