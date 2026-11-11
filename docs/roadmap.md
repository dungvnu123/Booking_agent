# Hướng phát triển Booking_agent

Mục tiêu: xây backend đặt lịch cho Booking App, sau đó bổ sung agent có thể tư vấn lịch phù hợp và hỗ trợ người dùng đặt lịch tại phòng khám.

## Giai đoạn hiện tại: nền tảng backend

Khung Go + Gin chia module nghiệp vụ, PostgreSQL, API/worker, Swagger và Docker Compose. Các route nghiệp vụ hiện trả 501, chưa có chức năng đặt lịch hoặc agent hoạt động.

Thứ tự triển khai:

1. Auth, user và doctor: tài khoản, JWT, role và hồ sơ bác sĩ.
2. Patient, medicine, examination và attachment: hồ sơ và dữ liệu phòng khám.
3. Roster: đăng ký, xếp lịch, công bố lịch, sinh WorkSchedule/Slot.
4. Booking: tìm lịch, giữ chỗ, đặt/hủy/đổi lịch; bảo đảm transaction và xử lý request đồng thời.
5. Queue, leave, closure, notification và job nền.

## Giai đoạn sau: agent tư vấn và đặt lịch

Agent nhận nhu cầu về ngày, giờ hoặc bác sĩ; gọi API để tìm lịch còn trống, trình bày lựa chọn và hỗ trợ đặt lịch sau khi người dùng xác nhận.

Luồng dự kiến:

```text
Người dùng nêu nhu cầu
  → Agent hỏi thêm thông tin còn thiếu
  → Backend cung cấp lịch còn trống
  → Agent đề xuất lựa chọn và người dùng chọn lịch
  → Giữ chỗ theo chính sách của backend
  → Người dùng xác nhận
  → Backend kiểm tra lại và tạo booking
  → Trả mã đặt lịch, thời gian và bác sĩ
```

Ranh giới triển khai:

- Backend là nguồn quyết định lịch trống, sức chứa, quyền truy cập và trạng thái booking.
- Agent sử dụng API đã xác thực, không ghi trực tiếp vào database hay tự bỏ qua luật nghiệp vụ.
- Lịch có thể thay đổi sau khi đề xuất: xử lý giữ chỗ hết hạn, chỗ vừa bị đặt và đưa ra lựa chọn mới.
- Bổ sung cơ chế chống tạo booking lặp khi agent retry request.
- Agent đặt lịch là thành phần phát triển sau, chưa thêm framework LLM hoặc dependency agent vào khung hiện tại.
- Solver Python OR-Tools phục vụ xếp lịch bác sĩ là tích hợp riêng, không phải agent hội thoại đặt lịch.

Các API ứng viên từ inventory: GET /api/booking/available-slots, POST /api/booking/hold, POST /api/booking/release-hold, POST /api/booking và GET /api/booking/me. Contract đầu vào/đầu ra phải được hoàn thiện khi triển khai backend trước khi dùng làm công cụ cho agent.
