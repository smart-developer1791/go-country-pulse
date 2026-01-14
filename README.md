# 🌍 Country Pulse

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![SSE](https://img.shields.io/badge/SSE-Real--time-FF6B6B?style=for-the-badge&logo=lightning&logoColor=white)
![Tailwind](https://img.shields.io/badge/Tailwind-CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)
![API](https://img.shields.io/badge/API-REST%20Countries-blue?style=for-the-badge)
![Render](https://img.shields.io/badge/Render-Deployed-46E3B7?style=for-the-badge&logo=render&logoColor=white)

> Real-time country discovery stream powered by Server-Sent Events

## ✨ Features

- 🌐 **Live Streaming** — Discover random countries every 4 seconds via SSE
- 🏳️ **Rich Data** — Flags, capitals, population, languages, currencies
- 🗺️ **Google Maps Integration** — Explore countries on the map
- 📜 **Discovery History** — Track recently discovered countries
- 📊 **Live Statistics** — Total countries and discovery counter
- 🎨 **Modern UI** — Beautiful glassmorphism design with Tailwind CSS

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.21+ with Chi router |
| Streaming | Server-Sent Events (SSE) |
| Frontend | Tailwind CSS |
| API | REST Countries API |

## 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/smart-developer1791/go-country-pulse
cd go-country-pulse
```

Initialize dependencies and run:

```bash
go mod tidy
go run .
```

Open in browser:

```text
http://localhost:8080
```

## 📡 API Endpoints

| Endpoint | Description |
|----------|-------------|
| GET / | Main dashboard UI |
| GET /api/stream | SSE stream of random countries |
| GET /api/random | Single random country JSON |
| GET /api/stats | Countries statistics |

## 🌐 Data Source

Data provided by [REST Countries API](https://restcountries.com/) — a free public API with comprehensive country information.

### Country Information Includes:

- 🏳️ Flag (SVG/PNG)
- 🏛️ Capital city
- 👥 Population
- 📐 Area (km²)
- 🗣️ Languages
- 💰 Currency
- 🕐 Timezone
- 🌍 Region & Continent
- 🗺️ Bordering countries

## 📁 Project Structure

```text
go-country-pulse/
├── main.go          # Application entry point
├── render.yaml      # Render deployment config
├── .gitignore       # Git ignore rules
└── README.md        # Documentation
```

## ⚙️ Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |

## 🎯 How It Works

```text
┌─────────────────────────────────────────────────────────┐
│                    Country Pulse                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│   ┌──────────┐      ┌───────────┐      ┌─────────────┐  │
│   │  Client  │◄────►│  Go/Chi   │◄────►│ REST        │  │
│   │ Browser  │ SSE  │  Server   │ HTTP │ Countries   │  │
│   └──────────┘      └───────────┘      │ API         │  │
│                                         └─────────────┘  │
│                                                          │
│   Stream Flow:                                           │
│   1. Server fetches all countries on startup             │
│   2. Client connects to /api/stream                      │
│   3. Random country sent every 4 seconds                 │
│   4. UI updates with animations                          │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## 🔧 Dependencies

```text
github.com/go-chi/chi/v5 — Lightweight HTTP router
```

## 🌟 Features Showcase

### Real-time Discovery
Countries stream automatically every 4 seconds with smooth animations.

### Rich Country Cards
Each country displays comprehensive information including flag, capital, population, languages, and more.

### Interactive History
Track your discoveries with a visual history grid showing flags and timestamps.

### One-Click Maps
Instantly explore any country on Google Maps.

## 📈 Performance

- ⚡ Lightweight — Single binary ~10MB
- 🚀 Fast startup — Countries cached on boot
- 📡 Efficient SSE — Minimal bandwidth usage
- 🔄 Auto-reconnect — Resilient connections

---

## Deploy in 10 seconds

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)
