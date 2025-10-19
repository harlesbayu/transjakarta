# 🚍 Transjakarta Fleet System (Dockerized)

Project ini terdiri dari beberapa komponen berbasis **Go** yang saling terhubung melalui **PostgreSQL**, **RabbitMQ**, dan **Mosquitto (MQTT)** — semuanya dijalankan menggunakan **Docker Compose**.

---

## 📦 Struktur Proyek

```
transjakarta/
├── fleet-backend/          # Aplikasi utama (HTTP + integrasi MQ & DB)
├── geofence-worker/        # Worker consumer event RabbitMQ
├── mock-publisher/         # Simulator/publisher MQTT
├── dockercompose/
│   ├── configs/
│   │   ├── mosquitto/mosquitto.conf
│   │   └── postgres/init.sql
│   ├── data/
│   │   ├── mosquitto/
│   │   └── postgres/
│   └── logs/
│       ├── mosquitto/
│       └── postgres/
└── docker-compose.yml
```

---

## 🚀 Menjalankan Project

Pastikan **Docker** dan **Docker Compose** sudah terinstal di sistem Anda.

### 1️⃣ Jalankan semua service

Dari direktori root project (`transjakarta`), jalankan:

```bash
docker compose up -d
```

Perintah ini akan:
- Membuat dan menjalankan seluruh container.
- Mengatur network antar service (`fleetnet`).
- Melakukan inisialisasi database PostgreSQL otomatis dari `init.sql`.

Jika butuh build ulang, jalankan:

```bash
docker compose up -d --build
```

---

## 🧰 Melihat Status dan Log

### 🔹 Melihat daftar container yang berjalan
```bash
docker ps
```

### 🔹 Melihat log semua service secara bersamaan
```bash
docker compose logs -f
```
> Tekan **Ctrl + C** untuk keluar dari tampilan log streaming.

### 🔹 Melihat log salah satu service
Contoh melihat log `geofence-worker`:
```bash
docker compose logs -f geofence-worker
```

---

## 🧹 Menghentikan dan Membersihkan Container

### 🔸 Hentikan semua service
```bash
docker compose down
```

### 🔸 Hentikan dan hapus volume data juga (bersih total)
```bash
docker compose down -v
```

---

## 🧩 Koneksi ke Service Terkait

| Service       | Port Lokal | Keterangan |
|----------------|-------------|-------------|
| PostgreSQL     | `5432` | Dapat diakses dari DBeaver / psql |
| RabbitMQ       | `15672` | Web UI: http://localhost:15672 (user: `guest`, pass: `guest`) |
| Mosquitto (MQTT) | `1883` | Endpoint MQTT Broker |
| Fleet Backend  | `8000` | HTTP API utama |

---

## 🧾 Catatan Tambahan

- Folder `dockercompose/data/` dan `dockercompose/logs/` **tidak dikomit ke Git**, agar data dan log lokal tidak bercampur dengan repository.
- Jika ingin menginisialisasi ulang database, hapus isi folder `dockercompose/data/postgres/` sebelum menjalankan ulang `docker compose up -d`.