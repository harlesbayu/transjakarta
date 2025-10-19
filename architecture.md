# 🏗️ Arsitektur Sistem — Transjakarta Fleet System

Sistem ini terdiri dari beberapa komponen **berbasis microservice** yang saling berinteraksi melalui **Docker network internal**.  
Semua komponen berjalan dalam container terisolasi namun saling terhubung di dalam satu jaringan (`fleetnet`).

---

## ⚙️ 1. Fleet Backend (REST API + Integrator)
Service utama yang menangani komunikasi antar sistem.  
Fungsi utamanya:
- Menyimpan data posisi kendaraan ke **PostgreSQL**.
- Mengirim event *geofence* ke **RabbitMQ**.
- Menerima update lokasi dari **MQTT broker (Mosquitto)**.
- Menyediakan endpoint HTTP untuk keperluan monitoring atau kontrol.

**Interaksi:**
- 🔁 Menulis data ke **PostgreSQL**
- 📤 Publish event ke **RabbitMQ (Exchange: fleet.events)**
- 📥 Subscribe ke **MQTT topic** dari `mock-publisher`

---

## 🐇 2. RabbitMQ (Event Broker)
Berfungsi sebagai *message queue* antar service.  
**Fleet Backend** mengirim event ke *exchange* `fleet.events`, sementara **Geofence Worker** mengonsumsi event dari *queue* `geofence_alerts`.

**Flow:**
```
Fleet Backend  →  RabbitMQ (fleet.events)  →  Geofence Worker
```

---

## 📡 3. Mosquitto (MQTT Broker)
Menangani pesan *publish/subscribe* antar perangkat.  
**Mock Publisher** mengirimkan lokasi kendaraan via MQTT topic (misalnya `/fleet/vehicle/location`),  
dan **Fleet Backend** menerima data tersebut untuk diproses lebih lanjut.

**Flow:**
```
Mock Publisher  →  Mosquitto  →  Fleet Backend
```

---

## 🗄️ 4. PostgreSQL (Database)
- Menyimpan data posisi kendaraan (`vehicle_location`)
- Diinisialisasi otomatis menggunakan file:
  ```
  dockercompose/configs/postgres/init.sql
  ```
- Volume data dan log disimpan di:
  ```
  dockercompose/data/postgres/
  dockercompose/logs/postgres/
  ```

---

## ⚙️ 5. Geofence Worker
Service *consumer* yang mendengarkan event dari **RabbitMQ**.  
Saat menerima event `geofence_entry`, worker akan memproses atau mencatatnya.  
Service ini dapat dikembangkan lebih lanjut untuk:
- Mengirim notifikasi real-time.
- Melakukan analitik data pergerakan kendaraan.

---

## 🧪 6. Mock Publisher
Simulator kendaraan yang mengirim data lokasi secara periodik ke **Mosquitto (MQTT)**.  
Berfungsi untuk melakukan pengujian dan simulasi sistem tanpa perangkat fisik.

**Contoh Alur:**
1. Mock Publisher mengirim data posisi kendaraan.
2. Fleet Backend menerima pesan via MQTT.
3. Data disimpan ke PostgreSQL.
4. Event dikirim ke RabbitMQ.
5. Geofence Worker memproses event geofence.

---

## 🌐 Diagram Alur Interaksi

```text
              ┌──────────────────────────┐
              │      Mock Publisher      │
              │ (Simulate Vehicle MQTT)  │
              └────────────┬─────────────┘
                           │ MQTT (1883)
                           ▼
                  ┌──────────────────┐
                  │   Mosquitto MQ   │
                  │   (MQTT Broker)  │
                  └────────┬─────────┘
                           │
                           ▼
                  ┌──────────────────┐
                  │  Fleet Backend   │
                  │  (Core Service)  │
                  ├──────────────────┤
                  │ - Save to DB     │
                  │ - Publish event  │
                  └────────┬─────────┘
                           │ AMQP (5672)
                           ▼
                  ┌──────────────────┐
                  │   RabbitMQ MQ    │
                  │   (Event Bus)    │
                  └────────┬─────────┘
                           │
                           ▼
                  ┌──────────────────┐
                  │ Geofence Worker  │
                  │ (Event Consumer) │
                  └──────────────────┘
```

---

## 🧩 Integrasi Antar Komponen

| Komponen         | Port  | Peran Utama                        | Komunikasi Dengan        |
|------------------|-------|------------------------------------|--------------------------|
| Fleet Backend    | 8000  | Core service (HTTP + MQ + DB)      | PostgreSQL, RabbitMQ, MQTT |
| RabbitMQ         | 5672 / 15672 | Event bus & message broker       | Fleet Backend, Geofence Worker |
| Mosquitto        | 1883 / 9001   | MQTT broker                     | Mock Publisher, Fleet Backend |
| PostgreSQL       | 5432  | Database utama                     | Fleet Backend |
| Geofence Worker  | —     | Consumer RabbitMQ                  | RabbitMQ |
| Mock Publisher   | —     | Simulasi kendaraan (MQTT publish)  | Mosquitto |

---

## 📦 Docker Network

Semua container dijalankan di jaringan internal Docker bernama `fleetnet`,  
sehingga tiap service bisa diakses melalui **service name** (bukan `localhost`).

Contoh koneksi:
```bash
host: rabbitmq
port: 5672
user: guest
pass: guest
```

---

## 🔄 Alur Singkat Sistem

1. **Mock Publisher** mengirim data lokasi via MQTT.
2. **Fleet Backend** menerima data dan menyimpannya di **PostgreSQL**.
3. Jika kendaraan masuk radius geofence, **Fleet Backend** mengirim event ke **RabbitMQ**.
4. **Geofence Worker** menerima event tersebut dan memprosesnya (misalnya mencatat atau mengirim notifikasi).

---

## 🧱 Teknologi yang Digunakan

- **Golang (Echo Framework)** → Backend dan Worker
- **PostgreSQL 16** → Database utama
- **RabbitMQ 3-management** → Message Broker (AMQP)
- **Eclipse Mosquitto 2** → MQTT Broker
- **Docker Compose** → Orkestrasi multi-container
- **MQTT + AMQP** → Komunikasi antar service asinkron
