# Money Management System (MMS)

Aplikasi **Money Management System** dibuat dengan **Golang 1.25**, menggunakan pendekatan **Domain Driven Design (DDD)** dan **Clean Architecture**.  
Teknologi yang digunakan:

- **Gin** → Web Framework
- **Zerolog** → Structured Logging
- **Paseto** → Token Authentication
- **MySQL (mms-db)** → Database

---

## Struktur Folder

```text
mms/
├── cmd/
│   └── mms/
│       └── main.go
├── config/
│   ├── config.go
│   └── config.yaml
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── transaction/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── service.go
│   ├── usecase/
│   │   ├── user_usecase.go
│   │   └── transaction_usecase.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── mysql/
│   │   │   │   ├── db.go
│   │   │   │   ├── user_repo.go
│   │   │   │   └── tx_repo.go
│   │   │   └── migration/schema.sql
│   │   ├── security/paseto.go
│   │   ├── logger/zerologger.go
│   │   └── http/middleware/
│   │       ├── auth.go
│   │       └── logger.go
│   ├── interface/
│   │   └── http/
│   │       ├── router.go
│   │       ├── handler/
│   │       │   ├── user_handler.go
│   │       │   ├── transaction_handler.go
│   │       │   └── response.go
│   │       └── request/
│   │           ├── user_request.go
│   │           └── transaction_request.go
│   └── test/
│       ├── user_service_test.go
│       └── transaction_usecase_test.go
├── pkg/
│   ├── errors/custom_error.go
│   ├── utils/
│   │   ├── timeutil.go
│   │   └── validation.go
│   └── constants/roles.go
└── README.md
```

---

## Cara Menjalankan

```bash
cd mms
go run ./cmd/mms
```

---

## Testing

```
go test ./internal/test/...
```

---

## Catatan

Struktur ini mengikuti prinsip Clean Architecture:

- domain → pure business logic

- usecase → application logic

- infrastructure → DB, logger, paseto, dll

- interface → HTTP adapter (Gin)

- pkg → helper reusable
