# Chỉ mục backend không AI để chuyển sang Go

Đọc cùng [bản thiết kế](backend-go-blueprint.md). Sinh từ source workspace ngày 10/11/2026; đây không phải OpenAPI hoặc schema DB runtime. Mọi route dưới đây có tiền tố mặc định `/api`. Quyền phải đọc cả service, không chỉ controller.

## 1. HTTP routes và handler

### Auth

Nguồn: [AuthController.java](../backend/src/main/java/com/abclinic/backend/controller/AuthController.java).

| Method | Path | Handler |
|---|---|---|
| POST | `/auth/register` | `register` |
| POST | `/auth/register-patient` | `registerPatient` |
| POST | `/auth/login` | `login` |
| POST | `/auth/refresh_token` | `refreshToken` |
| POST | `/auth/logout` | `logout` |

### AdminUser

Nguồn: [AdminUserController.java](../backend/src/main/java/com/abclinic/backend/controller/AdminUserController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/admin/users/roles` | `getRoles` |
| GET | `/admin/users` | `getUsers` |
| GET | `/admin/users/search` | `searchUsers` |
| GET | `/admin/users/{id}` | `getUser` |
| POST | `/admin/users` | `createUser` |
| PUT | `/admin/users/{id}` | `updateUser` |
| PATCH | `/admin/users/{id}/status` | `changeStatus` |
| PATCH | `/admin/users/{id}/password` | `resetPassword` |

### AdminDoctorProfile

Nguồn: [AdminDoctorProfileController.java](../backend/src/main/java/com/abclinic/backend/controller/AdminDoctorProfileController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/admin/doctor-profiles` | `listDoctorProfiles` |
| POST | `/admin/doctor-profiles` | `assignProfile` |

### Patient

Nguồn: [PatientController.java](../backend/src/main/java/com/abclinic/backend/controller/PatientController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/patients/{id}/summary` | `getPatientSummary` — tách/bỏ AI |
| GET | `/patients/{id}/summary-history/{type}` | `getPatientTaskHistory` — tách/bỏ AI |
| GET | `/patients/{id}/ai-sessions` | `getPatientAiSessions` — tách/bỏ AI |
| POST | `/patients` | `createPatient` |
| GET | `/patients/search` | `searchPatients` |
| GET | `/patients/stats` | `getPatientStats` |
| GET | `/patients` | `getAllPatients` |
| GET | `/patients/{id}` | `getPatientById` |
| PUT | `/patients/{id}` | `updatePatient` |
| POST | `/patients/import` | `importPatients` |
| POST | `/patients/import/batch` | `importPatientsBatch` |

### Medicine

Nguồn: [MedicineController.java](../backend/src/main/java/com/abclinic/backend/controller/MedicineController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/medicines` | `getMedicines` |
| GET | `/medicines/metadata` | `getMetadata` |
| GET | `/medicines/{id}` | `getMedicineById` |
| POST | `/medicines` | `createMedicine` |
| PUT | `/medicines/{id}` | `updateMedicine` |
| DELETE | `/medicines/{id}` | `deleteMedicine` |
| POST | `/medicines/bulk` | `bulkImport` |

### Examination

Nguồn: [ExaminationController.java](../backend/src/main/java/com/abclinic/backend/controller/ExaminationController.java).

| Method | Path | Handler |
|---|---|---|
| POST | `/examinations/{patientId}` | `create` |
| GET | `/examinations/patient/{patientId}` | `getByPatient` |
| GET | `/examinations/{id}` | `getById` |
| PUT | `/examinations/{id}` | `update` |

### File

Nguồn: [FileController.java](../backend/src/main/java/com/abclinic/backend/controller/FileController.java).

| Method | Path | Handler |
|---|---|---|
| POST | `/files/upload` | `uploadFile` |
| GET | `/files/list` | `listFiles` |

### DoctorMonthPick

Nguồn: [DoctorMonthPickController.java](../backend/src/main/java/com/abclinic/backend/controller/DoctorMonthPickController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/roster/monthly/{monthlyRosterId}/picks/me` | `getMyPick` |
| PUT | `/roster/monthly/{monthlyRosterId}/picks/me` | `submitPick` |

### MonthlyRoster

Nguồn: [MonthlyRosterController.java](../backend/src/main/java/com/abclinic/backend/controller/MonthlyRosterController.java).

| Method | Path | Handler |
|---|---|---|
| POST | `/roster/monthly` | `openMonth` |
| GET | `/roster/monthly` | `listMonths` |
| GET | `/roster/monthly/current` | `getCurrentOpenMonth` |
| GET | `/roster/monthly/my-preview` | `getMyPreview` |
| GET | `/roster/monthly/{id}` | `getMonth` |
| GET | `/roster/monthly/capacity-config` | `getAvgCapacity` |
| PUT | `/roster/monthly/solver-config` | `updateSolverConfig` |
| PUT | `/roster/monthly/{id}/override-config` | `overrideConfig` |
| DELETE | `/roster/monthly/{id}/override-config` | `clearOverrideConfig` |
| POST | `/roster/monthly/{id}/solve` | `runSolver` |
| GET | `/roster/monthly/{id}/result` | `getResult` |
| GET | `/roster/monthly/{id}/publish-preview` | `getPublishPreview` |
| POST | `/roster/monthly/{id}/approve` | `approve` |
| POST | `/roster/monthly/{id}/revert-to-draft` | `revertToDraft` |
| POST | `/roster/monthly/{id}/publish` | `publish` |
| POST | `/roster/monthly/{id}/publish-now` | `publishNow` |
| POST | `/roster/monthly/{id}/assignments` | `addAssignment` |
| DELETE | `/roster/monthly/{id}/assignments` | `removeAssignment` |
| POST | `/roster/monthly/{id}/assignments/lock` | `lockAssignment` |
| POST | `/roster/monthly/{id}/lock` | `lockCell` |
| GET | `/roster/monthly/{id}/cell-candidates` | `listCellCandidates` |
| GET | `/roster/monthly/{id}/doctor-statuses` | `listDoctorPickStatuses` |
| POST | `/roster/monthly/{id}/doctors/{doctorId}/reopen` | `reopenForDoctor` |
| POST | `/roster/monthly/{id}/doctors/{doctorId}/lock` | `lockForDoctor` |
| POST | `/roster/monthly/{id}/doctors/{doctorId}/unlock` | `unlockForDoctor` |
| GET | `/roster/monthly/{id}/published-swap/impact` | `getPublishedSwapImpact` |
| POST | `/roster/monthly/{id}/published-swap` | `swapPublishedDoctor` |
| GET | `/roster/monthly/{id}/audit-log` | `listAuditLog` |

### RosterWorkSchedule

Nguồn: [RosterWorkScheduleController.java](../backend/src/main/java/com/abclinic/backend/controller/RosterWorkScheduleController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/roster/work-schedules` | `listWorkSchedules` |
| GET | `/roster/work-schedules/me` | `listMyWorkSchedules` |

### LeaveRequest

Nguồn: [LeaveRequestController.java](../backend/src/main/java/com/abclinic/backend/controller/LeaveRequestController.java).

| Method | Path | Handler |
|---|---|---|
| POST | `/roster/leave-requests` | `create` |
| GET | `/roster/leave-requests/me` | `listMine` |
| GET | `/roster/leave-requests` | `listAll` |
| POST | `/roster/leave-requests/preview` | `preview` |
| GET | `/roster/leave-requests/{id}/impact` | `impact` |
| POST | `/roster/leave-requests/{id}/approve` | `approve` |
| POST | `/roster/leave-requests/{id}/reject` | `reject` |
| POST | `/roster/leave-requests/{id}/withdraw` | `withdraw` |

### ClinicClosure

Nguồn: [ClinicClosureController.java](../backend/src/main/java/com/abclinic/backend/controller/ClinicClosureController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/clinic-closures` | `listByMonth` |
| GET | `/clinic-closures/impact` | `getImpact` |
| GET | `/clinic-closures/impact-bulk` | `getImpactBulk` |
| POST | `/clinic-closures` | `add` |
| POST | `/clinic-closures/bulk` | `addBulk` |
| DELETE | `/clinic-closures/{id}` | `delete` |

### Booking

Nguồn: [BookingController.java](../backend/src/main/java/com/abclinic/backend/controller/BookingController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/booking/available-slots` | `getAvailableSlots` |
| GET | `/booking/slot-board` | `getSlotBoard` |
| POST | `/booking/hold` | `holdSlot` |
| POST | `/booking/release-hold` | `releaseHold` |
| POST | `/booking` | `createBooking` |
| POST | `/booking/staff` | `createStaffBooking` |
| GET | `/booking/me` | `listMyBookings` |
| POST | `/booking/{id}/cancel` | `cancelBooking` |
| GET | `/booking/work-schedules/{id}` | `listBookingsForWorkSchedule` |
| POST | `/booking/{id}/check-in` | `checkIn` |
| POST | `/booking/{id}/undo-check-in` | `undoCheckIn` |
| POST | `/booking/{id}/no-show` | `markNoShow` |
| POST | `/booking/{id}/call` | `callBooking` |
| POST | `/booking/{id}/finish` | `finishBooking` |
| POST | `/booking/{id}/call-failed` | `callBookingFailed` |
| POST | `/booking/{id}/staff-cancel` | `staffCancelBooking` |
| GET | `/booking/{id}/alternative-doctors` | `getAlternativeDoctors` |
| GET | `/booking/{id}/staff-alternative-doctors` | `getAlternativeDoctorsForStaff` |
| POST | `/booking/{id}/switch-doctor` | `switchDoctor` |
| POST | `/booking/{id}/staff-switch-doctor` | `staffSwitchDoctor` |
| POST | `/booking/{id}/reschedule` | `reschedule` |
| POST | `/booking/{id}/staff-reschedule` | `staffReschedule` |
| GET | `/booking/queue` | `getQueue` |
| POST | `/booking/{id}/swap` | `manualSwap` |

### Notification

Nguồn: [NotificationController.java](../backend/src/main/java/com/abclinic/backend/controller/NotificationController.java).

| Method | Path | Handler |
|---|---|---|
| GET | `/notifications/me` | `listMine` |
| POST | `/notifications/{id}/read` | `markRead` |

## 2. DTO đầu vào theo controller

Đây là chỉ mục tới DTO gốc để giữ nguyên field, validation và nested type. Một số request là query, multipart hoặc Map trực tiếp, cần đọc chữ ký handler trong controller.

### Auth

- [LoginRequest](../backend/src/main/java/com/abclinic/backend/dto/request/auth/LoginRequest.java)
- [LogoutRequest](../backend/src/main/java/com/abclinic/backend/dto/request/auth/LogoutRequest.java)
- [RefreshTokenRequest](../backend/src/main/java/com/abclinic/backend/dto/request/auth/RefreshTokenRequest.java)
- [RegisterRequest](../backend/src/main/java/com/abclinic/backend/dto/request/auth/RegisterRequest.java)
- [RegisterPatientRequest](../backend/src/main/java/com/abclinic/backend/dto/request/auth/RegisterPatientRequest.java)

### AdminUser

- [CreateUserRequest](../backend/src/main/java/com/abclinic/backend/dto/request/admin/CreateUserRequest.java)
- [ResetPasswordRequest](../backend/src/main/java/com/abclinic/backend/dto/request/admin/ResetPasswordRequest.java)
- [UpdateUserRequest](../backend/src/main/java/com/abclinic/backend/dto/request/admin/UpdateUserRequest.java)
- [UserStatusRequest](../backend/src/main/java/com/abclinic/backend/dto/request/admin/UserStatusRequest.java)

### AdminDoctorProfile

- [AssignDoctorProfileRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/AssignDoctorProfileRequest.java)

### Patient

- [CreatePatientRequest](../backend/src/main/java/com/abclinic/backend/dto/request/patient/CreatePatientRequest.java)
- [ImportPatientRowRequest](../backend/src/main/java/com/abclinic/backend/dto/request/patient/ImportPatientRowRequest.java)
- [UpdatePatientRequest](../backend/src/main/java/com/abclinic/backend/dto/request/patient/UpdatePatientRequest.java)

### Medicine

- [BulkImportMedicineRequest](../backend/src/main/java/com/abclinic/backend/dto/request/medicine/BulkImportMedicineRequest.java)
- [CreateMedicineRequest](../backend/src/main/java/com/abclinic/backend/dto/request/medicine/CreateMedicineRequest.java)
- [UpdateMedicineRequest](../backend/src/main/java/com/abclinic/backend/dto/request/medicine/UpdateMedicineRequest.java)

### Examination

- [ExaminationRequest](../backend/src/main/java/com/abclinic/backend/dto/request/ExaminationRequest.java)

### DoctorMonthPick

- [SubmitMonthPickRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/SubmitMonthPickRequest.java)

### MonthlyRoster

- [AdjustMonthlyAssignmentRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/AdjustMonthlyAssignmentRequest.java)
- [LockMonthlyAssignmentRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/LockMonthlyAssignmentRequest.java)
- [LockMonthlyCellRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/LockMonthlyCellRequest.java)
- [OpenMonthlyRosterRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/OpenMonthlyRosterRequest.java)
- [OverrideRosterConfigRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/OverrideRosterConfigRequest.java)
- [SwapPublishedDoctorRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/SwapPublishedDoctorRequest.java)
- [UpdateSolverConfigRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/UpdateSolverConfigRequest.java)

### LeaveRequest

- [ApproveLeaveRequestRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/ApproveLeaveRequestRequest.java)
- [CreateLeaveRequestRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/CreateLeaveRequestRequest.java)
- [PreviewLeaveRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/PreviewLeaveRequest.java)
- [RejectLeaveRequestRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/RejectLeaveRequestRequest.java)

### ClinicClosure

- [AddClinicClosureRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/AddClinicClosureRequest.java)
- [AddClinicClosuresBulkRequest](../backend/src/main/java/com/abclinic/backend/dto/request/roster/AddClinicClosuresBulkRequest.java)

### Booking

- [CallBookingRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/CallBookingRequest.java)
- [CreateBookingRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/CreateBookingRequest.java)
- [HoldSlotRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/HoldSlotRequest.java)
- [ManualSwapRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/ManualSwapRequest.java)
- [StaffCreateBookingRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/StaffCreateBookingRequest.java)
- [SwitchDoctorRequest](../backend/src/main/java/com/abclinic/backend/dto/request/booking/SwitchDoctorRequest.java)

## 3. Entity: field và annotation

Trích khai báo field và annotation từ entity; các quan hệ tới AI vẫn được hiển thị để phát hiện điểm cần tách. Không gồm các bảng chỉ phục vụ AI. @Column/@JoinColumn không chỉ rõ name thì cần kiểm tra naming strategy và DB thực tế. Không bao gồm thân lifecycle method; đọc nguồn để xem giá trị tự sinh.

### Role → `roles`

Nguồn: [Role.java](../backend/src/main/java/com/abclinic/backend/model/entity/Role.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
@Column(name = "role_id")
Long id;

@Column(name = "role_name", length = 50)
String roleName;

@OneToMany(mappedBy = "role")
List<User> users;
```

### User → `users`

Nguồn: [User.java](../backend/src/main/java/com/abclinic/backend/model/entity/User.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.UUID)
String id;

@ManyToOne(fetch = FetchType.EAGER)
@JoinColumn(name = "role_id")
Role role;

@Column(name = "username", length = 50)
String username;

@Column(name = "full_name", length = 200)
String fullName;

@Column(name = "email", length = 200)
String email;

@Column(name = "password_hash", length = 255)
String password;

@Column(name = "phone", length = 20)
String phone;

@Column(name = "is_active")
Boolean isActive;

@Column(name = "last_login_at")
LocalDateTime lastLoginAt;

@Column(name = "created_at")
LocalDateTime createdAt;

@Column(name = "updated_at")
LocalDateTime updatedAt;

@OneToMany(mappedBy = "doctor")
List<Examination> examinations;

@OneToMany(mappedBy = "user")
List<AiSession> aiSessions;

@OneToMany(mappedBy = "user")
List<FinalRecord> finalRecords;
```

### InvalidatedToken → `invalidated_token`

Nguồn: [InvalidatedToken.java](../backend/src/main/java/com/abclinic/backend/model/entity/InvalidatedToken.java).

```java
@Id
@Column(name = "id_token")
String idToken;

@Column(name = "expiry_time", nullable = false)
Date expiryTime;
```

### DoctorProfile → `doctor_profiles`

Nguồn: [DoctorProfile.java](../backend/src/main/java/com/abclinic/backend/model/entity/DoctorProfile.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@OneToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "user_id", nullable = false, unique = true)
User user;

@Column(name = "consultation_duration_minutes", nullable = false)
@Builder.Default
Integer consultationDurationMinutes = 15;

@Column(name = "patients_per_hour", nullable = false, columnDefinition = "integer not null default 5")
@Builder.Default
Integer patientsPerHour = 5;

@Column(name = "is_active", nullable = false)
@Builder.Default
Boolean isActive = true;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;

@UpdateTimestamp
@Column(name = "updated_at")
LocalDateTime updatedAt;
```

### Patient → `patients`

Nguồn: [Patient.java](../backend/src/main/java/com/abclinic/backend/model/entity/Patient.java).

```java
@Id
@Column(name = "id", length = 50)
String id;

@Column(name = "patient_code", length = 50)
String patientCode;

@Column(name = "full_name", length = 200)
String fullName;

@Column(name = "date_of_birth")
LocalDate dateOfBirth;

@Column(name = "gender", length = 10)
String gender;

@Column(name = "phone", length = 20)
String phone;

@Column(name = "email")
String email;

@Column(name = "job", length = 50)
String job;

@Column(name = "address", columnDefinition = "TEXT")
String address;

@Column(name = "allergies", columnDefinition = "TEXT")
String allergies;

@Column(name = "blood_type", length = 5)
String bloodType;

@Column(name = "insurance_id", length = 50)
String insuranceId;

@Column(name = "citizen_id", length = 12)
String citizenId;

@Column(name = "cooperation_ability", length = 50)
String cooperationAbility;

@Column(name = "created_at")
LocalDateTime createdAt;

@Column(name = "updated_at")
LocalDateTime updatedAt;

@Column(name = "last_visit_date")
LocalDate lastVisitDate;

@OneToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "user_id", unique = true)
User user;

@OneToMany(mappedBy = "patient")
List<PatientMedicine> patientMedicines;

@OneToMany(mappedBy = "patient")
List<FamilyHistory> familyHistories;

@OneToMany(mappedBy = "patient")
List<MedicalHistory> medicalHistories;

@OneToMany(mappedBy = "patient")
List<Examination> examinations;

@OneToMany(mappedBy = "patient")
List<AiSession> aiSessions;

@OneToOne(mappedBy = "patient", cascade = CascadeType.ALL)
PatientProfile patientProfile;
```

### PatientProfile → `patient_profiles`

Nguồn: [PatientProfile.java](../backend/src/main/java/com/abclinic/backend/model/entity/PatientProfile.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@JsonIgnore
@OneToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id")
Patient patient;

@Column(name = "sos", columnDefinition = "TEXT")
String sos;

@Column(name = "long_term_diagnosis", columnDefinition = "TEXT")
String longTermDiagnosis;

@Column(name = "perception_and_thought", columnDefinition = "TEXT")
String perceptionAndThought;

@Column(name = "needs", columnDefinition = "TEXT")
String needs;

@Column(name = "service_usage", columnDefinition = "TEXT")
String serviceUsage;

@Column(name = "service_behavior", columnDefinition = "TEXT")
String serviceBehavior;

@Column(name = "care_status", columnDefinition = "TEXT")
String careStatus;

@Column(name = "doctor_route", columnDefinition = "TEXT")
String doctorRoute;

@Column(name = "otc", columnDefinition = "TEXT")
String otc;

@Column(name = "medicine_characteristics", columnDefinition = "TEXT")
String medicineCharacteristics;

@Column(name = "non_medicine_route", columnDefinition = "TEXT")
String nonMedicineRoute;

@Column(name = "hotline_route", columnDefinition = "TEXT")
String hotlineRoute;

@Column(name = "prevention", columnDefinition = "TEXT")
String prevention;

@Column(name = "lifestyle", columnDefinition = "TEXT")
String lifestyle;

@Column(name = "lifestyle_change", columnDefinition = "TEXT")
String lifestyleChange;

@Column(name = "intervention", columnDefinition = "TEXT")
String intervention;
```

### MedicalHistory → `medical_histories`

Nguồn: [MedicalHistory.java](../backend/src/main/java/com/abclinic/backend/model/entity/MedicalHistory.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@JsonIgnore
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id")
Patient patient;

@Column(name = "condition", columnDefinition = "TEXT")
String condition;

@Column(name = "type", columnDefinition = "TEXT")
String type;

@Column(name = "diagnosed_at")
String diagnosedAt;

@Column(name = "notes", columnDefinition = "TEXT")
String notes;
```

### FamilyHistory → `family_histories`

Nguồn: [FamilyHistory.java](../backend/src/main/java/com/abclinic/backend/model/entity/FamilyHistory.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@JsonIgnore
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id")
Patient patient;

@Column(name = "relationship", length = 50)
String relationship;

@Column(name = "condition", columnDefinition = "TEXT")
String condition;

@Column(name = "notes", columnDefinition = "TEXT")
String notes;
```

### PatientMedicine → `patient_medicines`

Nguồn: [PatientMedicine.java](../backend/src/main/java/com/abclinic/backend/model/entity/PatientMedicine.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@JsonIgnore
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id")
Patient patient;

@Column(name = "medicine_name", columnDefinition = "TEXT")
String medicineName;

@Column(name = "dosage", columnDefinition = "TEXT")
String dosage;

@Enumerated(EnumType.STRING)
@Column(name = "category", length = 30)
PatientMedicineCategory category;

@Column(name = "can_use")
Boolean canUse;

@JsonIgnore
@Column(name = "noted_at")
LocalDateTime notedAt;
```

### Medicine → `medicines`

Nguồn: [Medicine.java](../backend/src/main/java/com/abclinic/backend/model/entity/Medicine.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Column(name = "medicine_code", length = 50, nullable = false, unique = true)
String medicineCode;

@Column(name = "medicine_name", length = 150, nullable = false)
String medicineName;

@Enumerated(EnumType.STRING)
@Column(name = "type", length = 30, nullable = false)
MedicineType type;

@Column(name = "category", length = 100)
String category;

@Column(name = "drug_group", length = 100)
String drugGroup;

@Column(name = "manufacturer", length = 200)
String manufacturer;

@Column(name = "country", length = 100)
String country;

@Column(name = "supplier_name", length = 200)
String supplierName;

@Column(name = "unit", length = 50)
String unit;

@Column(name = "pack_size")
Integer packSize;

@Column(name = "stock")
Integer stock;

@Column(name = "cabinet", length = 100)
String cabinet;

@Column(name = "shelf", length = 100)
String shelf;

@Column(name = "purchase_price", precision = 15, scale = 2)
BigDecimal purchasePrice;

@Column(name = "selling_price", precision = 15, scale = 2)
BigDecimal sellingPrice;

@Column(name = "discount_percentage", precision = 5, scale = 2)
BigDecimal discountPercentage;

@Column(name = "status", length = 50)
String status;

@Column(name = "expiry_date")
LocalDate expiryDate;

@Column(name = "image_url", columnDefinition = "TEXT")
String imageUrl;

@OneToOne(mappedBy = "medicine", cascade = CascadeType.ALL)
MedicineDescription description;

@BatchSize(size = 50)
@OneToMany(mappedBy = "medicine", cascade = CascadeType.ALL)
Set<Ingredient> ingredients;
```

### Ingredient → `ingredients`

Nguồn: [Ingredient.java](../backend/src/main/java/com/abclinic/backend/model/entity/Ingredient.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "medicine_id")
Medicine medicine;

@Column(name = "ingredient_name", length = 150)
String ingredientName;

@Column(name = "strength", length = 50)
String strength;

@Column(name = "description", columnDefinition = "TEXT")
String description;
```

### MedicineDescription → `medicine_descriptions`

Nguồn: [MedicineDescription.java](../backend/src/main/java/com/abclinic/backend/model/entity/MedicineDescription.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@OneToOne(fetch = FetchType.LAZY)
@JoinColumn(
name = "medicine_id",
nullable = false,
unique = true
)
Medicine medicine;

@Column(
name = "indications",
columnDefinition = "TEXT"
)
String indications;

@Column(
name = "contraindications",
columnDefinition = "TEXT"
)
String contraindications;

@Column(name = "dosage_per_day", length = 255)
String dosagePerDay;

@Column(name = "age_range", length = 100)
String ageRange;

@ElementCollection
@CollectionTable(
name = "medicine_usage_notes",
joinColumns = @JoinColumn(name = "medicine_description_id")
)
@Column(name = "usage_note", columnDefinition = "TEXT")
Set<String> usageNotes;

@Column(
name = "storage_instruction",
columnDefinition = "TEXT"
)
String storageInstruction;
```

### Examination → `examinations`

Nguồn: [Examination.java](../backend/src/main/java/com/abclinic/backend/model/entity/Examination.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Column(name = "external_id", length = 50)
String externalId;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id")
Patient patient;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_id")
User doctor;

@Column(name = "exam_date")
LocalDateTime examDate;

@Column(name = "symptom", columnDefinition = "TEXT")
String symptom;

@Column(name = "diagnosis", columnDefinition = "TEXT")
String diagnosis;

@Column(name = "doctor_advice", columnDefinition = "TEXT")
String doctorAdvice;

@Column(name = "patient_feedback", columnDefinition = "TEXT")
String patientFeedback;

@Column(name = "notes", columnDefinition = "TEXT")
String notes;

@Column(name = "status", length = 20)
String status;

@Column(name = "created_at")
LocalDateTime createdAt;

@Column(name = "updated_at")
LocalDateTime updatedAt;

@OneToMany(mappedBy = "examination")
List<TestRecord> tests;
```

### TestRecord → `test`

Nguồn: [TestRecord.java](../backend/src/main/java/com/abclinic/backend/model/entity/TestRecord.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "examination_id")
Examination examination;

@Column(name = "test_name", length = 200)
String testName;

@Column(name = "test_value", columnDefinition = "TEXT")
String testValue;

@Column(name = "ref_range", length = 50)
String refRange;

@Column(name = "note", columnDefinition = "TEXT")
String note;

@Column(name = "recorded_at")
LocalDateTime recordedAt;
```

### ClinicScheduleConfig → `clinic_schedule_configs`

Nguồn: [ClinicScheduleConfig.java](../backend/src/main/java/com/abclinic/backend/model/entity/ClinicScheduleConfig.java).

```java
@Id
Long id;

@Column(name = "min_per_session", nullable = false, columnDefinition = "integer not null default 1")
@Builder.Default
Integer minPerSession = 1;

@Column(name = "max_sessions_per_month", nullable = false, columnDefinition = "integer not null default 44")
@Builder.Default
Integer maxSessionsPerMonth = 44;

@Column(name = "pref_weight", nullable = false, columnDefinition = "integer not null default 5")
@Builder.Default
Integer prefWeight = 5;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "updated_by")
User updatedBy;

@Column(name = "updated_at")
LocalDateTime updatedAt;
```

### ClinicScheduleConfigLog → `clinic_schedule_config_logs`

Nguồn: [ClinicScheduleConfigLog.java](../backend/src/main/java/com/abclinic/backend/model/entity/ClinicScheduleConfigLog.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Column(name = "field", length = 50, nullable = false)
String field;

@Column(name = "old_value", length = 50)
String oldValue;

@Column(name = "new_value", length = 50, nullable = false)
String newValue;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "actor_id", nullable = false)
User actor;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### MonthlyRoster → `monthly_rosters`

Nguồn: [MonthlyRoster.java](../backend/src/main/java/com/abclinic/backend/model/entity/MonthlyRoster.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Column(name = "month", length = 7, nullable = false, unique = true)
String month;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
MonthlyRosterStatus status = MonthlyRosterStatus.PENDING;

@Column(name = "registration_opens_at")
LocalDateTime registrationOpensAt;

@Column(name = "pick_deadline", nullable = false)
LocalDateTime pickDeadline;

@Column(name = "publish_due_at")
LocalDateTime publishDueAt;

@Column(name = "registration_opened_notified_at")
LocalDateTime registrationOpenedNotifiedAt;

@Column(name = "registration_locked_at")
LocalDateTime registrationLockedAt;

@Column(name = "min_per_session", nullable = false)
@Builder.Default
Integer minPerSession = 1;

@Column(name = "max_sessions_per_month", nullable = false, columnDefinition = "integer not null default 44")
@Builder.Default
Integer maxSessionsPerMonth = 44;

@Column(name = "pref_weight", nullable = false)
@Builder.Default
Integer prefWeight = 5;

@Column(name = "override_min_per_session")
Integer overrideMinPerSession;

@Column(name = "override_max_sessions_per_month")
Integer overrideMaxSessionsPerMonth;

@Column(name = "override_pref_weight")
Integer overridePrefWeight;

@Column(name = "solver_run_at")
LocalDateTime solverRunAt;

@Column(name = "solver_status", length = 20)
String solverStatus;

@Column(name = "solve_ms")
Integer solveMs;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "solver_run_by")
User solverRunBy;

@Column(name = "publish_reminder_sent_on")
LocalDate publishReminderSentOn;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "approved_by")
User approvedBy;

@Column(name = "approved_at")
LocalDateTime approvedAt;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "published_by")
User publishedBy;

@Column(name = "published_at")
LocalDateTime publishedAt;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;

@UpdateTimestamp
@Column(name = "updated_at")
LocalDateTime updatedAt;
```

### DoctorMonthPick → `doctor_month_picks`

Nguồn: [DoctorMonthPick.java](../backend/src/main/java/com/abclinic/backend/model/entity/DoctorMonthPick.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "monthly_roster_id", nullable = false)
MonthlyRoster monthlyRoster;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_id", nullable = false)
User doctor;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
DoctorMonthPickStatus status = DoctorMonthPickStatus.DRAFT;

@Column(name = "reused_from_month", length = 7)
String reusedFromMonth;

@Column(name = "applied_scope_months")
Integer appliedScopeMonths;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "applied_from_pick_id")
DoctorMonthPick appliedFromPick;

@Column(name = "patients_per_hour")
Integer patientsPerHour;

@Column(name = "locked", nullable = false)
@Builder.Default
Boolean locked = false;

@Column(name = "submitted_at")
LocalDateTime submittedAt;

@Column(name = "locked_at")
LocalDateTime lockedAt;

@Column(name = "reopen_until")
LocalDateTime reopenUntil;

@Column(name = "reopen_expired_notified_at")
LocalDateTime reopenExpiredNotifiedAt;
```

### DoctorPickDay → `doctor_pick_days`

Nguồn: [DoctorPickDay.java](../backend/src/main/java/com/abclinic/backend/model/entity/DoctorPickDay.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_month_pick_id", nullable = false)
DoctorMonthPick doctorMonthPick;

@Column(name = "pick_date", nullable = false)
LocalDate pickDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClinicSession session;

@Enumerated(EnumType.STRING)
@Column(name = "level", length = 10, nullable = false)
AvailabilityLevel level;
```

### DoctorPlannedDay → `doctor_planned_days`

Nguồn: [DoctorPlannedDay.java](../backend/src/main/java/com/abclinic/backend/model/entity/DoctorPlannedDay.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "origin_pick_id", nullable = false)
DoctorMonthPick originPick;

@Column(name = "pick_date", nullable = false)
LocalDate pickDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClinicSession session;

@Enumerated(EnumType.STRING)
@Column(name = "level", length = 10, nullable = false)
AvailabilityLevel level;
```

### MonthlyRosterAssignmentCell → `monthly_roster_assignment_cells`

Nguồn: [MonthlyRosterAssignmentCell.java](../backend/src/main/java/com/abclinic/backend/model/entity/MonthlyRosterAssignmentCell.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "monthly_roster_id", nullable = false)
MonthlyRoster monthlyRoster;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_id", nullable = false)
User doctor;

@Column(name = "pick_date", nullable = false)
LocalDate pickDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClinicSession session;

@Enumerated(EnumType.STRING)
@Column(name = "source", length = 20, nullable = false)
AssignmentSource source;

@Column(name = "locked", nullable = false)
@Builder.Default
Boolean locked = false;

@Enumerated(EnumType.STRING)
@Column(name = "lock_type", length = 20)
AssignmentLockType lockType;
```

### MonthlyRosterAuditLog → `monthly_roster_audit_logs`

Nguồn: [MonthlyRosterAuditLog.java](../backend/src/main/java/com/abclinic/backend/model/entity/MonthlyRosterAuditLog.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "monthly_roster_id", nullable = false)
MonthlyRoster monthlyRoster;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "actor_id", nullable = false)
User actor;

@Column(name = "action", length = 50, nullable = false)
String action;

@Column(name = "pick_date")
LocalDate pickDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10)
ClinicSession session;

@Column(name = "doctor_id")
String doctorId;

@Column(name = "target_doctor_id")
String targetDoctorId;

@Column(name = "reason", length = 500)
String reason;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### WorkSchedule → `work_schedules`

Nguồn: [WorkSchedule.java](../backend/src/main/java/com/abclinic/backend/model/entity/WorkSchedule.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_id", nullable = false)
User doctor;

@Column(name = "work_date", nullable = false)
LocalDate workDate;

@Column(name = "start_time", nullable = false)
LocalTime startTime;

@Column(name = "end_time", nullable = false)
LocalTime endTime;

@Column(name = "slot_duration_minutes", nullable = false)
Integer slotDurationMinutes;

@Column(name = "buffer_minutes", nullable = false)
@Builder.Default
Integer bufferMinutes = 5;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
WorkScheduleStatus status = WorkScheduleStatus.PUBLISHED;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "approved_by")
User approvedBy;

@Column(name = "approved_at")
LocalDateTime approvedAt;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### Slot → `slots`

Nguồn: [Slot.java](../backend/src/main/java/com/abclinic/backend/model/entity/Slot.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "work_schedule_id", nullable = false)
WorkSchedule workSchedule;

@Column(name = "start_time", nullable = false)
LocalTime startTime;

@Column(name = "end_time", nullable = false)
LocalTime endTime;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
SlotStatus status = SlotStatus.AVAILABLE;

@Column(name = "held_until")
LocalDateTime heldUntil;

@Column(name = "held_by_patient_id", length = 50)
String heldByPatientId;

@Column(name = "capacity", nullable = false, columnDefinition = "integer not null default 1")
@Builder.Default
Integer capacity = 1;

@Column(name = "booked_count", nullable = false, columnDefinition = "integer not null default 0")
@Builder.Default
Integer bookedCount = 0;
```

### SlotHold → `slot_holds`

Nguồn: [SlotHold.java](../backend/src/main/java/com/abclinic/backend/model/entity/SlotHold.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "slot_id", nullable = false)
Slot slot;

@Column(name = "patient_id", length = 50, nullable = false)
String patientId;

@Column(name = "expires_at", nullable = false)
LocalDateTime expiresAt;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### Booking → `bookings`

Nguồn: [Booking.java](../backend/src/main/java/com/abclinic/backend/model/entity/Booking.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "patient_id", nullable = false)
Patient patient;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "slot_id", nullable = false)
Slot slot;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
BookingStatus status = BookingStatus.BOOKED;

@Column(name = "note", columnDefinition = "TEXT")
String note;

@Column(name = "cancel_reason", columnDefinition = "TEXT")
String cancelReason;

@Column(name = "checked_in_at")
LocalDateTime checkedInAt;

@Column(name = "checked_in_late")
Boolean checkedInLate;

@Column(name = "queue_order_at")
LocalDateTime queueOrderAt;

@Column(name = "queue_slot_start")
LocalTime queueSlotStart;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "original_slot_id")
Slot originalSlot;

@Column(name = "called_at")
LocalDateTime calledAt;

@Column(name = "finished_at")
LocalDateTime finishedAt;

@Column(name = "missed_call_count")
Integer missedCallCount;

@Column(name = "call_note", columnDefinition = "TEXT")
String callNote;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "created_by")
User createdBy;

@Column(name = "booking_for_name", length = 200)
String bookingForName;

@Column(name = "booking_for_relationship", length = 100)
String bookingForRelationship;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;

@UpdateTimestamp
@Column(name = "updated_at")
LocalDateTime updatedAt;
```

### BookingSwapLog → `booking_swap_logs`

Nguồn: [BookingSwapLog.java](../backend/src/main/java/com/abclinic/backend/model/entity/BookingSwapLog.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "booking_id", nullable = false)
Booking lateBooking;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "from_slot_id", nullable = false)
Slot fromSlot;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "to_slot_id", nullable = false)
Slot toSlot;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "swapped_in_booking_id", nullable = false)
Booking earlyBooking;

@CreationTimestamp
@Column(name = "swapped_at", nullable = false, updatable = false)
LocalDateTime swappedAt;

@Column(name = "note", length = 500)
String note;
```

### BookingCallLog → `booking_call_logs`

Nguồn: [BookingCallLog.java](../backend/src/main/java/com/abclinic/backend/model/entity/BookingCallLog.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "failed_booking_id", nullable = false)
Booking failedBooking;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "failed_slot_id", nullable = false)
Slot failedSlot;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "replacement_booking_id")
Booking replacementBooking;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "replacement_slot_id")
Slot replacementSlot;

@Column(name = "missed_call_count")
Integer missedCallCount;

@CreationTimestamp
@Column(name = "created_at", nullable = false, updatable = false)
LocalDateTime createdAt;

@Column(name = "note", length = 500)
String note;
```

### LeaveRequest → `leave_requests`

Nguồn: [LeaveRequest.java](../backend/src/main/java/com/abclinic/backend/model/entity/LeaveRequest.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "doctor_id", nullable = false)
User doctor;

@Column(name = "from_date", nullable = false)
LocalDate fromDate;

@Column(name = "to_date", nullable = false)
LocalDate toDate;

@Column(name = "reason", length = 500)
String reason;

@Enumerated(EnumType.STRING)
@Column(name = "status", length = 20, nullable = false)
@Builder.Default
LeaveRequestStatus status = LeaveRequestStatus.PENDING;

@Column(name = "reject_reason", length = 500)
String rejectReason;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "reviewed_by")
User reviewedBy;

@Column(name = "reviewed_at")
LocalDateTime reviewedAt;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;

@UpdateTimestamp
@Column(name = "updated_at")
LocalDateTime updatedAt;
```

### LeaveRequestItem → `leave_request_items`

Nguồn: [LeaveRequestItem.java](../backend/src/main/java/com/abclinic/backend/model/entity/LeaveRequestItem.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "leave_request_id", nullable = false)
LeaveRequest leaveRequest;

@Column(name = "leave_date", nullable = false)
LocalDate leaveDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClinicSession session;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "proposed_replacement_doctor_id")
User proposedReplacementDoctor;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "approved_replacement_doctor_id")
User approvedReplacementDoctor;
```

### ClinicClosure → `clinic_closures`

Nguồn: [ClinicClosure.java](../backend/src/main/java/com/abclinic/backend/model/entity/ClinicClosure.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Column(name = "closure_date", nullable = false)
LocalDate closureDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClosureSession session;

@Column(name = "reason", length = 255)
String reason;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "created_by", nullable = false)
User createdBy;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### ClinicClosureLog → `clinic_closure_logs`

Nguồn: [ClinicClosureLog.java](../backend/src/main/java/com/abclinic/backend/model/entity/ClinicClosureLog.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@Enumerated(EnumType.STRING)
@Column(name = "action", length = 10, nullable = false)
ClosureLogAction action;

@Column(name = "closure_date", nullable = false)
LocalDate closureDate;

@Enumerated(EnumType.STRING)
@Column(name = "session", length = 10, nullable = false)
ClosureSession session;

@Column(name = "reason", length = 255)
String reason;

@Column(name = "cancelled_booking_count", nullable = false)
@Builder.Default
Integer cancelledBookingCount = 0;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "actor_id", nullable = false)
User actor;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### Notification → `notifications`

Nguồn: [Notification.java](../backend/src/main/java/com/abclinic/backend/model/entity/Notification.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "recipient_id", nullable = false)
User recipient;

@Column(name = "type", length = 50, nullable = false)
String type;

@Column(name = "title", length = 200, nullable = false)
String title;

@Column(name = "message", columnDefinition = "TEXT", nullable = false)
String message;

@Column(name = "related_booking_id")
Long relatedBookingId;

@Column(name = "is_read", nullable = false)
@Builder.Default
Boolean isRead = false;

@CreationTimestamp
@Column(name = "created_at", updatable = false)
LocalDateTime createdAt;
```

### FileAttach → `file_attachments`

Nguồn: [FileAttach.java](../backend/src/main/java/com/abclinic/backend/model/entity/FileAttach.java).

```java
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
Long id;

@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "session_id")
AiSession session;

@Column(name = "file_name", columnDefinition = "TEXT")
String fileName;

@Column(name = "url", columnDefinition = "TEXT")
String url;

@Column(name = "file_size")
Long fileSize;

@Column(name = "file_type", columnDefinition = "TEXT")
String fileType;

@Column(name = "created_at")
LocalDateTime createdAt;
```

## 4. Enum và contract đặc biệt

- BookingStatus: BOOKED, CHECKED_IN, LATE, CALLED, COMPLETED, NO_SHOW, CANCELLED.
- MonthlyRosterStatus: PENDING, DRAFT, APPROVED, PUBLISHED.
- DoctorMonthPickStatus: DRAFT, SUBMITTED.
- AvailabilityLevel: OFF, OK, PREF.
- ClinicSession: SANG, CHIEU.
- ClosureSession: FULL_DAY, SANG, CHIEU.
- LeaveRequestStatus: PENDING, APPROVED, REJECTED, WITHDRAWN.
- WorkScheduleStatus: PUBLISHED, CANCELLED.
- SlotStatus: AVAILABLE, HELD, BOOKED (HELD còn phục vụ dữ liệu cũ).
- AssignmentSource: SOLVER, MANUAL_ADD, POST_PUBLISH_ADD.
- AssignmentLockType: MANUAL, MANUAL_CELL, LONG_TERM (legacy).
- PatientMedicineCategory: COMPATIBLE, INCOMPATIBLE, OPTIMAL.
- MedicineType: DB MEDICINE/SUPPLEMENT; JSON xuất Thuốc/TPCN.

## 5. Phạm vi bị loại khỏi chỉ mục API

Không liệt kê AiSessionController, AiAgentTestController, OcrController, TaskController, InternalPatientController, InternalMedicineController. Các route bệnh nhân gắn AI vẫn đánh dấu ở phần 1 vì nằm chung controller với hồ sơ bệnh nhân. File và import cần chỉnh model trước khi loại các bảng AI.
