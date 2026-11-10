# Database migrations

Chưa có DDL nghiệp vụ được áp dụng. Tài liệu nguồn là khảo sát tĩnh, không phải schema PostgreSQL thực tế.

Trước khi viết migration đầu tiên: export schema thật, đối chiếu Java entity và resources/db/manual, chốt bảng AI bỏ/giữ và thiết kế attachment gắn patient/examination.

Đặt SQL có version tại đây (ví dụ 000001_identity.up.sql và .down.sql). Chọn một migration runner khi thêm migration đầu tiên; API/worker không tự tạo bảng, không dùng auto-migrate và không chạy SQL khi startup.

Khi port booking: dùng transaction, UPDATE nhận/nhả chỗ có điều kiện và kiểm tra affected rows. Không áp constraint booked_count <= capacity vì có allowOverCapacity của lễ tân. Thiết kế riêng khóa chống booking trùng lịch và nhiều CALLED trong cùng ca.
