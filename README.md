# 📱 TZone - Platform So Sánh Cấu Hình Điện Thoại

<div align="center">

**API mạnh mẽ để so sánh cấu hình chi tiết của các điện thoại thông minh**

[![Go Version](https://img.shields.io/badge/Go-1.25.4%2B-blue?logo=go)](https://golang.org)
[![Supabase](https://img.shields.io/badge/Supabase-Latest-blue?logo=supabase)](https://supabase.com)
[![MongoDB](https://img.shields.io/badge/MongoDB-Latest-green?logo=mongodb)](https://www.mongodb.com)
[![License](https://img.shields.io/badge/License-MIT-purple)](#-giấy-phép)

</div>

---

## 📑 Mục Lục

- [🎯 Giới Thiệu](#-giới-thiệu)
- [🚀 Công Nghệ Sử Dụng](#-công-nghệ-sử-dụng)
- [📂 Cấu Trúc Dự Án](#-cấu-trúc-dự-án)
- [🛠️ Yêu Cầu Trước Tiên](#️-yêu-cầu-trước-tiên)
- [⚙️ Cài Đặt & Thiết Lập](#️-cài-đặt--thiết-lập)
- [▶️ Chạy Ứng Dụng](#️-chạy-ứng-dụng)
- [📝 Các Endpoint API](#-các-endpoint-api)
- [📱 Model Specifications Chi Tiết](#-model-specifications-chi-tiết)
- [📖 Hướng Dẫn Logging](#-hướng-dẫn-logging--best-practices)
- [🚀 Development Guide](#-development-guide)
- [🔧 Troubleshooting](#-troubleshooting)
- [🔐 Security Best Practices](#-security-best-practices)
- [🤝 Đóng Góp](#-đóng-góp)
- [📄 Giấy Phép](#-giấy-phép)

---

## 🏗️ Kiến Trúc Toàn Diện

### Sơ Đồ Tổng Quan Hệ Thống

```mermaid
graph TB
    subgraph Client["🖥️ Client Layer"]
        Web["React/Vite Frontend<br/>3000"]
        UI["UI Components<br/>Theme + Dark/Light"]
        Ctx["Contexts<br/>Auth/Favorites/Theme"]
    end

    subgraph API["🌐 API Gateway"]
        Gin["Gin HTTP Server<br/>8080"]
        Auth["JWT Auth<br/>Rate Limit"]
        RBAC["RBAC Middleware"]
    end

    subgraph Handler["📨 Delivery Layer - Handlers"]
        BrandH["Brand Handler"]
        DeviceH["Device Handler"]
        AuthH["Auth Handler"]
        FavH["Favorite Handler"]
        CommonH["Common Handler"]
    end

    subgraph Service["🧠 Service Layer - Business Logic"]
        AuthSvc["Auth Service<br/>JWT + OTP"]
        BrandSvc["Brand Service"]
        DeviceSvc["Device Service<br/>Finder"]
        FavSvc["Favorite Service<br/>Sync Logic"]
        PermSvc["Permission Service"]
    end

    subgraph Repository["🗄️ Repository Layer - Data Access"]
        UserRepo["User Repository"]
        BrandRepo["Brand Repository<br/>MongoDB"]
        DeviceRepo["Device Repository<br/>MongoDB"]
        FavRepo["Favorite Repository<br/>PostgreSQL"]
        TokenRepo["Token Repository<br/>PostgreSQL"]
    end

    subgraph DB["💾 Database Layer"]
        PG["PostgreSQL<br/>Users, Auth,<br/>Favorites, Tokens"]
        Mongo["MongoDB<br/>Brands, Devices,<br/>Specifications"]
    end

    subgraph External["🔗 External Services"]
        MinIO["MinIO<br/>Object Storage<br/>Device Images"]
        SMTP["SMTP<br/>Email Service<br/>OTP Codes"]
        JWT["JWT<br/>Token Generation"]
    end

    Client -->|HTTP| API
    API --> Handler
    Handler --> Service
    Service --> Repository
    Repository -->|SQL| DB
    Repository -->|Query| DB
    Service -->|Upload| External
    Service -->|Email OTP| External
    AuthSvc -->|Gen Token| External
    
    style Client fill:#e1f5ff
    style API fill:#fff3e0
    style Handler fill:#f3e5f5
    style Service fill:#e8f5e9
    style Repository fill:#fce4ec
    style DB fill:#ede7f6
    style External fill:#fff9c4
```

### Luồng Request Chi Tiết

```mermaid
sequenceDiagram
    participant Browser as 🌐 Browser
    participant FE as React Frontend
    participant API as 🚀 Gin API
    participant Auth as 🔐 Auth Handler
    participant AuthSvc as Auth Service
    participant UserRepo as User Repo
    participant PG as PostgreSQL
    participant JWT as JWT Module

    Browser->>FE: User clicks Login
    FE->>API: POST /auth/login {email, password}
    
    API->>Auth: Validate request
    Auth->>AuthSvc: Login(email, password)
    
    AuthSvc->>UserRepo: FindByEmailWithRole(email)
    UserRepo->>PG: Query users table
    PG-->>UserRepo: Return user + role
    
    AuthSvc->>AuthSvc: Compare password hash
    AuthSvc->>JWT: GenerateTokenPair(userID)
    JWT-->>AuthSvc: accessToken, refreshToken
    
    AuthSvc->>UserRepo: SaveRefreshToken()
    UserRepo->>PG: INSERT into tokens
    
    AuthSvc-->>Auth: Return tokens + user
    Auth-->>API: Success response
    
    API-->>FE: {accessToken, user, role}
    FE->>FE: Store in localStorage
    FE->>FE: Load Favorites + Theme
    Browser-->>Browser: Redirect to home
```

### Kiến Trúc Layer-wise

```mermaid
graph LR
    subgraph Delivery
        Handler["Handler<br/>Request Parsing<br/>Validation"]
    end

    subgraph Business["Service<br/>Logic"]
        Logic["Business Rules<br/>OTP Generation<br/>Favorite Sync"]
    end

    subgraph Data["Repository<br/>Data Access"]
        Query["Query Building<br/>Error Handling"]
    end

    subgraph DB["Database"]
        SQLite["SQL/NoSQL"]
    end

    HTTP["HTTP Request"]
    HTTP -->|1. Parse| Delivery
    Delivery -->|2. Delegate| Business
    Business -->|3. Access| Data
    Data -->|4. Query| DB
    DB -->|5. Return| Data
    Data -->|6. Aggregate| Business
    Business -->|7. Format| Delivery
    Delivery -->|8. Respond| HTTP

    style Delivery fill:#ffe0b2
    style Business fill:#c8e6c9
    style Data fill:#f8bbd0
    style DB fill:#b3e5fc
```

### Entity Relationships (ERD)

```mermaid
erDiagram
    USERS ||--o{ REFRESH_TOKENS : has
    USERS ||--o{ FAVORITES : has
    USERS ||--o{ USER_ROLES : has

    USERS {
        uuid id PK
        string email UK
        string password_hash
        timestamp created_at
    }

    REFRESH_TOKENS {
        uuid id PK
        uuid user_id FK
        timestamp expires_at
        timestamp created_at
    }

    FAVORITES {
        int id PK
        uuid user_id FK
        string device_id FK
        timestamp created_at
    }

    BRANDS {
        objectId _id PK
        string brand_name
        array devices
    }

    DEVICES {
        objectId _id PK
        string model_name
        string imageUrl
        object specifications
    }

    ROLES {
        int id PK
        string name UK
    }

    USER_ROLES {
        uuid user_id FK
        int role_id FK
    }
```

### Data Flow: Guest → Login → Favorites Sync

```mermaid
flowchart TD
    Start["🌐 User Visits Site"]
    GuestFav["💾 Add to Favorites<br/>localStorage:'favorites'"]
    LocalStore["📱 Store in<br/>Browser Memory"]
    
    Login["🔐 Click Login Button"]
    Input["📝 Enter Email/Password"]
    Submit["✅ Submit Form"]
    
    APIAuth["🚀 POST /auth/login"]
    DBCheck["🔍 Query PostgreSQL<br/>Check user exists"]
    PassMatch["🔑 Verify password hash"]
    
    Decision1{Match?}
    ErrorLogin["❌ Invalid credentials"]
    
    GenToken["🎟️ Generate JWT<br/>accessToken, refreshToken"]
    SaveToken["💾 Save refreshToken<br/>to PostgreSQL"]
    
    FELoad["🎨 FE Detects<br/>isAuthenticated"]
    ReadLocal["📖 Read localStorage<br/>favorites"]
    
    Decision2{Has Guest<br/>Favorites?}
    GetServer["📡 GET /favorites<br/>from server"]
    
    SyncAPI["📤 POST /api/v1/favorites/sync<br/>{ device_ids: [...] }"]
    ValidateDevice["✅ Check each device<br/>exists in MongoDB"]
    DeDupe["🔄 Remove duplicates<br/>via DB unique constraint"]
    Insert["💾 INSERT new favorites<br/>to PostgreSQL"]
    
    ClearLocal["🗑️ localStorage.removeItem<br/>('favorites')"]
    FinalList["📋 Return synced list<br/>device_ids: [...]"]
    
    UIUpdate["🎨 FE Updates UI<br/>Hearts filled on devices"]
    Success["✅ Sync Complete"]
    
    Start --> GuestFav --> LocalStore
    LocalStore -.->|Later| Login
    
    Login --> Input --> Submit
    Submit --> APIAuth
    
    APIAuth --> DBCheck --> PassMatch
    PassMatch --> Decision1
    Decision1 -->|No| ErrorLogin
    Decision1 -->|Yes| GenToken
    
    GenToken --> SaveToken
    SaveToken -->|Return tokens| FELoad
    
    FELoad --> ReadLocal
    ReadLocal --> Decision2
    Decision2 -->|No| GetServer
    Decision2 -->|Yes| SyncAPI
    
    SyncAPI --> ValidateDevice
    ValidateDevice --> DeDupe
    DeDupe --> Insert
    
    Insert --> ClearLocal
    ClearLocal --> FinalList
    FinalList --> UIUpdate
    UIUpdate --> Success
    
    style Start fill:#e1f5ff
    style GuestFav fill:#fff9c4
    style Login fill:#ffe0b2
    style APIAuth fill:#f8bbd0
    style GenToken fill:#c8e6c9
    style SyncAPI fill:#f3e5f5
    style Success fill:#c8e6c9
```

### Frontend State Management

```mermaid
graph TB
    App["App Root<br/>ThemeProvider<br/>AuthProvider<br/>FavoritesProvider"]
    
    Theme["🎨 ThemeContext<br/>theme: dark/light<br/>toggleTheme()"]
    
    Auth["🔐 AuthContext<br/>user, token<br/>isAuthenticated<br/>login, logout"]
    
    Fav["❤️ FavoritesContext<br/>favoriteIds: string array<br/>isFavorite(id)<br/>toggleFavorite(id)"]
    
    Pages["📄 Pages"]
    Device["Device Detail"]
    Brand["Brand List"]
    Finder["Device Finder"]
    Favorites["Favorites Page"]
    
    LocalStorage["📱 localStorage<br/>theme<br/>favorites<br/>access_token<br/>user"]
    
    API["🚀 API Calls<br/>GET /favorites<br/>POST /favorites/sync<br/>DELETE /favorites/:id"]
    
    App -->|wrap| Theme
    App -->|wrap| Auth
    App -->|wrap| Fav
    
    Theme -->|useTheme| Pages
    Auth -->|useAuth| Pages
    Fav -->|useFavorites| Pages
    
    Pages --> Device
    Pages --> Brand
    Pages --> Finder
    Pages --> Favorites
    
    Fav -->|persist| LocalStorage
    Auth -->|persist| LocalStorage
    Theme -->|persist| LocalStorage
    
    Fav -->|sync on login| API
    API -->|update| Fav
    
    style App fill:#e1f5ff
    style Theme fill:#fff9c4
    style Auth fill:#ffe0b2
    style Fav fill:#f8bbd0
    style LocalStorage fill:#c8e6c9
    style API fill:#f3e5f5
```

### Deployment Architecture

```mermaid
graph LR
    Client["🖥️ Client<br/>Browser"]
    
    CDN["📦 CDN<br/>dist/ files<br/>Static assets"]
    
    LB["⚖️ Load Balancer<br/>nginx/reverse proxy"]
    
    subgraph Container["🐳 Docker Container"]
        API["Gin API<br/>Port 8080"]
        Worker["Background Jobs<br/>OTP Cleanup"]
    end
    
    PG["🐘 PostgreSQL<br/>Port 5432"]
    Mongo["🍃 MongoDB<br/>Port 27017"]
    MinIO["📦 MinIO<br/>Port 9000"]
    SMTP["📧 SMTP<br/>Email Service"]
    
    Client -->|HTTP/HTTPS| CDN
    Client -->|API calls| LB
    LB -->|Route| Container
    
    Container -->|Query| PG
    Container -->|Query| Mongo
    Container -->|Upload/Get| MinIO
    Container -->|Send OTP| SMTP
    
    style Client fill:#e1f5ff
    style CDN fill:#fff9c4
    style LB fill:#ffe0b2
    style Container fill:#f8bbd0
    style PG fill:#e0f2f1
    style Mongo fill:#c8e6c9
    style MinIO fill:#f3e5f5
```

---

## 🎯 Giới Thiệu

TZone là một nền tảng **REST API** hiện đại được xây dựng bằng **Go**, sử dụng **framework Gin** và hỗ trợ nhiều loại cơ sở dữ liệu (**MongoDB** và **Supabase**). Ứng dụng tuân theo **Clean Architecture** để đảm bảo khả năng mở rộng, dễ bảo trì và hiệu suất cao.

### ✨ Tính Năng Chính

- 🔍 **So sánh chi tiết** cấu hình giữa các điện thoại thông minh
- 📊 **Quản lý danh mục** nhãn hiệu (Brands) và thiết bị (Devices)
- 🗄️ **Hỗ trợ đa cơ sở dữ liệu** - MongoDB và Supabase (PostgreSQL)
- 📝 **Logging chi tiết** với mã emoji tiêu chuẩn
- ⚙️ **Cấu hình linh hoạt** qua biến môi trường
- 🛡️ **Xử lý lỗi toàn diện** không sử dụng panic
- 🚀 **Khởi động nhanh** và khả năng mở rộng tốt

---

## 🐳 Chạy Local bằng Docker

Để khởi động toàn bộ stack local gồm API, frontend, PostgreSQL, MongoDB và MinIO (đã seed từ `phoneExample.json`), chạy:

```bash
docker compose up --build
```

### Dịch vụ được tạo

- **Frontend:** http://localhost:3000
- **API:** http://localhost:8080
- **PostgreSQL:** `localhost:5432`
- **MongoDB:** `localhost:27017`
- **MinIO API:** http://localhost:9000
- **MinIO Console:** http://localhost:9001

### Dữ liệu khởi tạo

- PostgreSQL được khởi tạo với extension `pgcrypto` để GORM có thể tạo UUID mặc định.
- MongoDB sẽ được import lại toàn bộ collection `brands` trong database `Cluster0` từ file `phoneExample.json`.
- API container dùng biến môi trường `POSTGRES_DSN` như chuỗi kết nối PostgreSQL local, còn `MONGODB_URL` trỏ tới MongoDB container.
- Thư mục `media/` sẽ được đồng bộ lên bucket MinIO và `imageUrl` trong seed sẽ được chuẩn hóa thành URL MinIO public.
- Device tạo mới/cập nhật từ UI sẽ upload ảnh trực tiếp vào MinIO thay vì lưu local.

Mật khẩu mặc định cho PostgreSQL local trong compose là:

```text
user: postgres
password: postgres
database: tzone
```

---

## 🚀 Công Nghệ Sử Dụng

- **Ngôn Ngữ:** [Go](https://golang.org) 1.25.4+
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin) (HTTP router)
- **Cơ Sở Dữ Liệu:**
  - [MongoDB Atlas](https://www.mongodb.com/products/platform?msockid=34596b69b607646d3daf7d27b725651b) (Go Driver v2)
  - [Supabase](https://supabase.com) (PostgreSQL)
- **ORM:** [GORM](https://gorm.io) (cho Supabase)
- **Quản Lý Cấu Hình:** [godotenv](https://github.com/joho/godotenv)
- **Tạo Giao Diện:** [Gomponents](https://github.com/maragu/gomponents)
- **Kiến Trúc:** Clean Architecture / Layered Architecture

## 📂 Cấu Trúc Dự Án

TZone tuân theo **Clean Architecture** với sự phân tách rõ ràng giữa các lớp:

```
tzone/
├── cmd/
│   └── main.go                           # 🎯 Điểm vào ứng dụng
│
├── infrastructure/                       # 🏗️ Lớp cơ sở hạ tầng
│   ├── configuration/
│   │   └── configuration.go              # ⚙️ Tải & xác thực biến môi trường
│   └── database/
│       ├── mongodb.go                    # 🍃 Kết nối MongoDB
│       └── supabase.go                   # 🐘 Kết nối Supabase/PostgreSQL
│
├── internal/                             # 🔒 Lõi ứng dụng (private)
│   ├── delivery/                         # 📨 HTTP API layer
│   │   ├── handler/
│   │   │   ├── device_handler.go         # 📱 Xử lý request thiết bị
│   │   │   ├── brand_handler.go          # 🏷️ Xử lý request nhãn hiệu
│   │   │   └── common_handler.go         # 🔄 Handler dùng chung
│   │   └── route/
│   │       ├── device_route.go           # 📍 Route cho device API
│   │       ├── brand_route.go            # 📍 Route cho brand API
│   │       └── common_route.go           # 📍 Route chung
│   │
│   ├── dto/                              # 📦 Data Transfer Objects
│   │   ├── request.go                    # 📥 Định dạng request từ client
│   │   └── response.go                   # 📤 Định dạng response gửi về
│   │
│   ├── model/                            # 📋 Domain Models
│   │   ├── device.go                     # 📱 Model Device (với 15+ thông số)
│   │   └── brand.go                      # 🏷️ Model Brand
│   │
│   ├── repository/                       # 🗄️ Truy cập dữ liệu
│   │   ├── mongodb_repository.go         # 🍃 Query MongoDB
│   │   ├── postgre_repository.go         # 🐘 Query PostgreSQL
│   │   └── supabase_repository.go        # 🌐 Query Supabase API
│   │
│   ├── service/                          # 🧠 Lớp logic kinh doanh
│   │   ├── device_service.go             # 💼 Logic xử lý thiết bị
│   │   └── brand_service.go              # 💼 Logic xử lý nhãn hiệu
│   │
│   └── server/                           # 🌐 HTTP Server
│       ├── server.go                     # 🚀 Thiết lập server Gin
│       └── handler.go                    # 🔌 Khởi tạo handlers
│
├── util/                                 # 🛠️ Tiện ích
│   └── response/
│       └── response.go                   # 📋 Utility cho response
│
├── web/                                  # 🎨 Tài nguyên web
│   └── page/
│       ├── index.go                      # 🏠 Trang chủ
│       └── shared/
│           ├── header.go                 # 📝 Phần header
│           ├── footer.go                 # 📝 Phần footer
│           └── layout.go                 # 📐 Layout
│
├── docs/                                 # 📚 Tài liệu
│   └── LOGGING_GUIDE.md                  # 📋 Hướng dẫn logging chi tiết
│
├── .env                                  # 🔐 Biến môi trường (không commit)
├── .env.example                          # 📝 Mẫu biến môi trường
├── phone.json                            # 📱 Dữ liệu điện thoại
├── go.mod                                # 📌 Go module definition
├── go.sum                                # 🔒 Hash phụ thuộc
└── README.md                             # 📖 Tài liệu dự án
```

### 🏗️ Kiến Trúc Chi Tiết

```
HTTP Request
    ↓
[Delivery Layer - Handler] (📨 handler/)
    ↓ (Parse & Validate)
[Service Layer] (🧠 service/)
    ↓ (Business Logic)
[Repository Layer] (🗄️ repository/)
    ↓ (Data Access)
[Database] (🗄️ MongoDB / PostgreSQL)
```

## 🛠️ Yêu Cầu Trước Tiên

Trước khi chạy ứng dụng, hãy đảm bảo bạn đã cài đặt:

- [Go](https://go.dev/dl/) (phiên bản 1.25.4 hoặc cao hơn)
- **Cơ Sở Dữ Liệu** (chọn ít nhất 1):
  - [MongoDB](https://www.mongodb.com/try/download/community) - cục bộ hoặc [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
  - [Supabase](https://supabase.com) - PostgreSQL dựa trên cloud

## ⚙️ Cài Đặt & Thiết Lập

### 1️⃣ Clone Repository

```bash
git clone https://github.com/LuuDinhTheTai/tzone.git
cd tzone
```

### 2️⃣ Cài Đặt Dependencies

```bash
go mod download
```

Hoặc nếu bạn cập nhật dependencies:

```bash
go mod tidy
```

### 3️⃣ Cấu Hình Biến Môi Trường

Sao chép `.env.example` để tạo file `.env`:

```bash
# Windows
copy .env.example .env

# Linux/Mac
cp .env.example .env
```

Chỉnh sửa file `.env` với cấu hình của bạn:

```env
# 🌐 Server Configuration
SERVER_PORT=8080

# 🍃 MongoDB (local hoặc cloud)
MONGODB_URL=mongodb://localhost:27017/?directConnection=true

# 🐘 PostgreSQL local
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/tzone?sslmode=disable

# 🌐 Supabase hosted API (tuỳ chọn)
SUPABASE_API_URL=https://your-project.supabase.co
SUPABASE_KEY=your-api-key-here

# 🛡️ API rate limit (tuỳ chọn)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_API_RPM=120
RATE_LIMIT_API_BURST=30
RATE_LIMIT_AUTH_RPM=20
RATE_LIMIT_AUTH_BURST=5
```

**Lưu ý:**
- Ít nhất một cơ sở dữ liệu phải được cấu hình
- Nếu không cấu hình, ứng dụng sẽ hiển thị cảnh báo nhưng vẫn khởi động
- Cả hai cơ sở dữ liệu có thể hoạt động đồng thời

### 4️⃣ (Tuỳ Chọn) Cài Đặt MongoDB Cục Bộ

Nếu bạn muốn chạy MongoDB cục bộ:

```bash
# Tải MongoDB Community Edition
# https://www.mongodb.com/try/download/community

# Sau khi cài đặt, khởi động MongoDB
# Windows
mongod

# Linux/Mac
brew services start mongodb-community
```

Cập nhật `.env`:
```env
MONGODB_URL=mongodb://localhost:27017
```

## ▶️ Chạy Ứng Dụng

### Chế Độ Development

Để khởi động máy chủ ở chế độ phát triển:

```bash
go run cmd/main.go
```

Bạn sẽ thấy output tương tự:

```
🚀 Starting TZone Application...
📅 Date: 2026-01-26
🔄 Loading environment configuration...
✅ Loaded configuration from .env file
✅ Configuration validated successfully
🔧 Connecting to MongoDB...
✅ MongoDB connected and ready
🔧 Initializing HTTP server...
🌐 Starting HTTP server on port 8080...
```

### Xây Dựng Binary

Để tạo executable cho môi trường production:

```bash
# Windows
go build -o tzone.exe ./cmd

# Linux/Mac
go build -o tzone ./cmd
```

Chạy binary:

```bash
# Windows
./tzone.exe

# Linux/Mac
./tzone
```

### Chế Độ Monitoring

Ứng dụng hỗ trợ structured logging với `slog`. Để xem chi tiết logging:

```bash
go run cmd/main.go 2>&1 | tee app.log
```

Máy chủ sẽ khởi động trên cổng được chỉ định trong `.env` (mặc định: `8080`)

Truy cập trang chủ: [http://localhost:8080](http://localhost:8080)

## 📝 Các Endpoint API

### Base URL
```
http://localhost:8080/api/v1
```

### Thiết Bị (Devices)

#### Lấy Danh Sách Tất Cả Thiết Bị
```http
GET /devices
```
**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Devices retrieved successfully",
  "data": [
    {
      "model_name": "iPhone 15 Pro",
      "imageUrl": "https://example.com/image.jpg",
      "specifications": {
        "network": {...},
        "launch": {...},
        "body": {...},
        "display": {...},
        "platform": {...},
        "memory": {...},
        "mainCamera": {...},
        "selfieCamera": {...},
        "sound": {...},
        "comms": {...},
        "features": {...},
        "battery": {...},
        "misc": {...}
      }
    }
  ]
}
```

#### Tìm Thiết Bị Theo Tên
```http
GET /devices/search?name=iphone&page=1&limit=10
```
**Response:** `200 OK`

Trả về danh sách device có `model_name` khớp một phần với tham số `name` theo kiểu không phân biệt hoa/thường.

#### Tạo Thiết Bị Mới
```http
POST /devices
Content-Type: multipart/form-data

{
  "model_name": "Samsung Galaxy S24",
  "brand_id": "...",
  "image": "<file upload>",
  "specifications": { /* full specifications object */ }
}
```
**Response:** `201 Created`

> Ghi chú: device seed vẫn có thể dùng ảnh URL từ nguồn ngoài, nhưng các device tạo mới/cập nhật sẽ lưu ảnh local để frontend load qua `/uploads`.

#### Lấy Chi Tiết Thiết Bị
```http
GET /devices/:id
```
**Response:** `200 OK`

#### Cập Nhật Thiết Bị
```http
PUT /devices/:id
Content-Type: multipart/form-data

{
  "model_name": "Updated Name",
  "brand_id": "...",
  "image": "<file upload optional>",
  "specifications": { /* updated specs */ }
}
```
**Response:** `200 OK`

#### Xóa Thiết Bị
```http
DELETE /devices/:id
```
**Response:** `200 OK`

### Nhãn Hiệu (Brands)

#### Lấy Danh Sách Tất Cả Nhãn Hiệu
```http
GET /brands
```
**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Brands retrieved successfully",
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "brand_name": "Apple",
      "devices": [
        { /* device objects */ }
      ]
    }
  ]
}
```

#### Tìm Nhãn Hiệu Theo Tên
```http
GET /brands/search?name=apple&page=1&limit=10
```
**Response:** `200 OK`

Trả về danh sách brand có `brand_name` khớp một phần với tham số `name` theo kiểu không phân biệt hoa/thường.

#### Tạo Nhãn Hiệu Mới
```http
POST /brands
Content-Type: application/json

{
  "brand_name": "OnePlus",
  "devices": []
}
```
**Response:** `201 Created`

#### Lấy Chi Tiết Nhãn Hiệu
```http
GET /brands/:id
```
**Response:** `200 OK`

#### Cập Nhật Nhãn Hiệu
```http
PUT /brands/:id
Content-Type: application/json

{
  "brand_name": "Updated Brand Name"
}
```
**Response:** `200 OK`

#### Xóa Nhãn Hiệu
```http
DELETE /brands/:id
```
**Response:** `200 OK`

### Status Codes

| Code | Ý Nghĩa |
|------|---------|
| 200 | OK - Thành công |
| 201 | Created - Tạo mới thành công |
| 400 | Bad Request - Dữ liệu không hợp lệ |
| 404 | Not Found - Tài nguyên không tìm thấy |
| 500 | Server Error - Lỗi máy chủ |

## 📖 Hướng Dẫn Logging & Best Practices

TZone sử dụng **logging emoji-based** tiêu chuẩn để dễ dàng theo dõi hoạt động:

### Emoji Reference

| Emoji | Ý Nghĩa | Ví Dụ |
|-------|---------|--------|
| 🚀 | Khởi động app | Starting application |
| 🔄 | Đang xử lý | Connecting to database |
| ✅ | Thành công | Connection successful |
| ❌ | Lỗi | Operation failed |
| ⚠️ | Cảnh báo | Config missing |
| 🔧 | Setup/Config | Initializing server |
| 🔑 | Auth/Credentials | Loading credentials |
| 🌐 | Network | Server listening |
| 📝 | Info | Using defaults |
| 🛑 | Fatal | Application terminated |
| 👋 | Shutdown | Graceful shutdown |

### Error Handling Rules

✅ **Đúng cách:** Trả về error, không dùng panic
```go
if err != nil {
    log.Printf("❌ Operation failed: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
```

❌ **Sai:** Dùng panic
```go
if err != nil {
    panic(err)  // ❌ Không được phép!
}
```

Để biết thêm chi tiết, vui lòng tham khảo [LOGGING_GUIDE.md](docs/LOGGING_GUIDE.md).

## 📱 Model Specifications Chi Tiết

Device model bao gồm các thông số kỹ thuật chi tiết:

### Network (Mạng)
- Technology: 2G, 3G, 4G, 5G
- Bands: Các loại dải tần 2G, 3G, 4G, 5G
- Speed: Tốc độ mạng tối đa

### Launch (Ra mắt)
- Announced: Ngày công bố
- Status: Trạng thái (Available, Discontinued, etc.)

### Body (Thân thiết bị)
- Dimensions: Kích thước
- Weight: Trọng lượng
- Build: Chất liệu xây dựng
- SIM: Loại SIM
- IP Rating: Chứng chỉ chống nước/bụi

### Display (Màn hình)
- Type: AMOLED, IPS, LCD, etc.
- Size: Kích thước màn hình
- Resolution: Độ phân giải

### Platform (Nền tảng)
- OS: Hệ điều hành
- Chipset: Chip xử lý
- CPU: Bộ vi xử lý
- GPU: Card đồ họa

### Memory (Bộ nhớ)
- Card Slot: Hỗ trợ thẻ nhớ
- Internal: Bộ nhớ nội bộ

### Cameras (Camera)
- **Main Camera:** Triple/Quad, Features, Video
- **Selfie Camera:** Độ phân giải, Video

### Sound (Âm thanh)
- Loudspeaker: Loa ngoài
- 3.5mm Jack: Jack tai nghe

### Comms (Kết nối)
- WLAN: Chuẩn Wi-Fi
- Bluetooth: Phiên bản Bluetooth
- Positioning: GPS, GLONASS, etc.
- NFC: Hỗ trợ NFC
- Radio: Đài FM
- USB: Loại cổng USB

### Features (Tính năng)
- Sensors: Danh sách các cảm biến

### Battery (Pin)
- Type: Loại pin
- Charging: Hỗ trợ sạc nhanh

### Misc (Khác)
- Các thông tin khác

## 📖 Hướng Dẫn Logging & Best Practices

TZone sử dụng **logging emoji-based** tiêu chuẩn để dễ dàng theo dõi hoạt động:

### Emoji Reference

| Emoji | Ý Nghĩa | Ví Dụ |
|-------|---------|--------|
| 🚀 | Khởi động app | Starting application |
| 🔄 | Đang xử lý | Connecting to database |
| ✅ | Thành công | Connection successful |
| ❌ | Lỗi | Operation failed |
| ⚠️ | Cảnh báo | Config missing |
| 🔧 | Setup/Config | Initializing server |
| 🔑 | Auth/Credentials | Loading credentials |
| 🌐 | Network | Server listening |
| 📝 | Info | Using defaults |
| 🛑 | Fatal | Application terminated |
| 👋 | Shutdown | Graceful shutdown |

### Error Handling Rules

✅ **Đúng cách:** Trả về error, không dùng panic
```go
if err != nil {
    log.Printf("❌ Operation failed: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
```

❌ **Sai:** Dùng panic
```go
if err != nil {
    panic(err)  // ❌ Không được phép!
}
```

Để biết thêm chi tiết, vui lòng tham khảo [LOGGING_GUIDE.md](docs/LOGGING_GUIDE.md).

---

## 🚀 Development Guide

### Cấu Trúc Dự Án - Chi Tiết Thêm

#### Thêm Một Repository Mới

1. Tạo file `internal/repository/xyz_repository.go`
2. Kế thừa từ cấu trúc phổ biến:
```go
type XYZRepository struct {
    mongoClient *mongo.Client
}

func NewXYZRepository(mongoClient *mongo.Client) *XYZRepository {
    return &XYZRepository{mongoClient: mongoClient}
}

func (r *XYZRepository) Create(ctx context.Context, data XYZ) error {
    // Implementation
    return nil
}
```

#### Thêm Một Service Mới

1. Tạo file `internal/service/xyz_service.go`
2. Inject repository:
```go
type XYZService struct {
    repo *repository.XYZRepository
}

func NewXYZService(repo *repository.XYZRepository) *XYZService {
    return &XYZService{repo: repo}
}
```

#### Thêm Một Handler Mới

1. Tạo file `internal/delivery/handler/xyz_handler.go`
2. Sử dụng service:
```go
type XYZHandler struct {
    service *service.XYZService
}

func (h *XYZHandler) GetAll(c *gin.Context) {
    // Implementation
}
```

#### Thêm Route Mới

1. Tạo file `internal/delivery/route/xyz_route.go`
2. Map routes:
```go
func MapXYZRoutes(r *gin.Engine, handler *handler.XYZHandler) {
    group := r.Group("/api/v1")
    group.GET("/xyz", handler.GetAll)
}
```

3. Đăng ký route trong `internal/server/server.go`

### Testing

Chạy các test:
```bash
go test ./...
```

Chạy test với coverage:
```bash
go test -cover ./...
```

### Code Quality

Kiểm tra linting:
```bash
go vet ./...
```

Format code:
```bash
go fmt ./...
```

## 🔧 Troubleshooting

### MongoDB Connection Error

**Vấn đề:** `connection refused`

**Giải pháp:**
```bash
# Kiểm tra MongoDB đang chạy không
# Windows (đã cài MongoDB)
mongod

# Hoặc sử dụng MongoDB Atlas - cập nhật MONGODB_URL trong .env
```

### Port Already in Use

**Vấn đề:** `Address already in use`

**Giải pháp:**
```bash
# Thay đổi port trong .env
SERVER_PORT=8081

# Hoặc tìm và kill process đang dùng port 8080
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

### Module Not Found

**Vấn đề:** `cannot find module`

**Giải pháp:**
```bash
go mod tidy
go mod download
```

### Environment Variables Not Loaded

**Vấn đề:** App không nhận biến môi trường từ `.env`

**Giải pháp:**
1. Kiểm tra file `.env` tồn tại trong thư mục gốc
2. Kiểm tra định dạng biến:
```env
SERVER_PORT=8080  # Đúng
SERVER_PORT = 8080  # Sai (có dấu cách)
SERVER_PORT="8080"  # Có thể dùng quotes
```

## 🔐 Security Best Practices

- ✅ Không commit file `.env` (đã trong `.gitignore`)
- ✅ Sử dụng environment variables cho credentials
- ✅ Validate input từ client
- ✅ Không log sensitive data (passwords, keys)
- ✅ Sử dụng HTTPS trong production

## 📊 Performance Tips

- 💨 Sử dụng MongoDB indexing để query nhanh hơn
- 💨 Cache results nếu có thể
- 💨 Batch insert/update khi xử lý lượng lớn data
- 💨 Sử dụng pagination cho large datasets

## 📚 Tài Liệu Bổ Sung

- [LOGGING_GUIDE.md](docs/LOGGING_GUIDE.md) - Chi tiết về logging system
- [Gin Documentation](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/)
- [GORM Documentation](https://gorm.io/)
- [Supabase Documentation](https://supabase.com/docs)

## 🤝 Đóng Góp

Chúng tôi hoan nghênh các đóng góp! Vui lòng tự do gửi Pull Request.

### Quy Trình Đóng Góp

1. Fork dự án
2. Tạo branch feature (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

### Hướng Dẫn Code

- Tuân theo style Go chuẩn
- Viết test cho các tính năng mới
- Cập nhật tài liệu nếu cần
- Sử dụng emoji-based logging

## 📄 Giấy Phép

[MIT](https://choosealicense.com/licenses/mit/)
