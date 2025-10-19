# 🚍 Fleet Backend Service

**Fleet Backend** adalah komponen utama dalam sistem *Transjakarta Fleet*, yang berfungsi sebagai API backend dan pusat integrasi antara database PostgreSQL, broker RabbitMQ, dan MQTT broker (Mosquitto).  
Aplikasi ini ditulis menggunakan **Golang** dengan arsitektur modular dan clean architecture pattern.

---

## 🧩 Struktur Folder

```
fleet-backend/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point aplikasi
│
├── configs/
│   ├── config.json                  # Konfigurasi server, DB, dan message broker
│   └── geofence_locations.json      # Titik geofence untuk pemantauan kendaraan
│
├── go.mod, go.sum                   # Dependency management Go
│
├── internal/
│   ├── adapter/                     # Implementasi adapter eksternal
│   │   ├── http_adapter/
│   │   │   └── server.go            # HTTP server menggunakan Gin
│   │   ├── postgres_adapter/
│   │   │   ├── gorm.go              # Inisialisasi koneksi GORM PostgreSQL
│   │   │   └── vehicle_repository.go # Repository kendaraan
│   │   ├── mqtt_adapter/
│   │   │   ├── mqtt_client.go       # MQTT client (subscriber)
│   │   │   └── mqtt_subscriber.go   # Subscriber untuk data kendaraan
│   │   ├── rabbitmq_adapter/
│   │   │   └── rabbitmq.go          # Koneksi dan pengiriman event ke RabbitMQ
│   │
│   ├── config/
│   │   └── config.go                # Struktur konfigurasi dan loader
│   │
│   ├── container/
│   │   └── container.go             # Dependency injection container
│   │
│   ├── domain/
│   │   ├── model/
│   │   │   └── vehicle_location.go  # Model domain lokasi kendaraan
│   │   └── repository/
│   │       └── vehicle_repository.go# Interface repository domain
│   │
│   ├── handler/
│   │   ├── ping_handler.go          # Endpoint healthcheck
│   │   └── vehicle_handler.go       # Endpoint REST API kendaraan
│   │
│   └── usecase/
│       └── vehicle_usecase.go       # Logika bisnis untuk kendaraan
│
└── README.md
```

---

## ⚙️ Konfigurasi

File `configs/config.json` berisi konfigurasi seperti berikut:

```json
{
  "server": { "port": 8000 },
  "postgres": {
    "host": "postgres",
    "port": 5432,
    "user": "fleet",
    "password": "fleetpass",
    "dbname": "fleetdb",
    "sslmode": "disable"
  },
  "mqtt": { "host": "mosquitto", "port": 1883 },
  "rabbitmq": { "url": "amqp://guest:guest@rabbitmq:5672/" }
}
```

---

## 🚀 Menjalankan Aplikasi (Tanpa Docker)

Jika ingin menjalankan langsung dari Go:

```bash
cd cmd/server
go run main.go
```

Pastikan service berikut sudah berjalan:
- PostgreSQL
- RabbitMQ
- Mosquitto (MQTT broker)

---

## 🧠 Arsitektur Singkat

Aplikasi ini menggunakan pola **Clean Architecture**:
- **Handler** → menerima HTTP request
- **Usecase** → berisi logika bisnis
- **Repository (interface)** → mendefinisikan kontrak akses data
- **Adapter** → mengimplementasikan koneksi ke sumber eksternal (DB, MQ, MQTT)
- **Container** → menyatukan semua dependensi

---

## 🧪 Endpoint Dasar

| Endpoint | Method | Deskripsi |
|-----------|---------|-----------|
| `/ping` | GET | Mengecek apakah server aktif |
| `/vehicles/:vehicle_id/location` | GET | Mendapatkan lokasi terakhir kendaraan |
| `/vehicles/:vehicle_id/history` | GET | Mendapatkan histori lokasi kendaraan |
