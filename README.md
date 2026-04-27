# go-shorten 🚀

A minimal, high-performance, self-hosted URL shortener built with Go. It features a modern brutalist web interface and a simple RESTful API.

## ✨ Features

- **Lightning Fast**: Built with Go and Gin for high performance.
- **Brutalist UI**: Simple, clean, and functional web interface.
- **REST API**: Easy to integrate with other tools and services.
- **SQLite Backend**: Lightweight data persistence with WAL mode enabled.
- **Analytics**: Basic click tracking for your shortened URLs.
- **Self-Hosted**: Designed to be easily deployed on your own infrastructure.

## 🛠️ Technology Stack

- **Backend**: [Go](https://go.dev/) (Gin Gonic)
- **Database**: [SQLite](https://sqlite.org/)
- **Frontend**: HTML, CSS, JavaScript (Vanilla)
- **Task Runner**: [Just](https://github.com/casey/just)
- **Live Reload**: [Air](https://github.com/air-verse/air)

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- [Just](https://github.com/casey/just) (optional)
- [Docker](https://www.docker.com/) (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ppablomunoz/go-shorten.git
   cd go-shorten
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   just run
   # or
   go run cmd/main.go
   ```

### Running with Docker

You can also run the application using Docker Compose:
```bash
docker-compose up -d
```

The server will start at `http://localhost:8000`.

### Backup & Recovery (Docker)

Since the application uses SQLite, creating a backup is as simple as copying the database file.

**Backup from container to host:**
```bash
docker cp go-shorten-app-1:/root/db/go-shorten.db ./backup.db
```

*Note: If you are using the default `docker-compose.yml`, your database is already persisted in the `./db` folder on your host machine.*

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/url` | Create a new short URL |
| `GET` | `/api/url` | List all shortened URLs |
| `PUT` | `/api/url/:code` | Update a short URL |
| `DELETE` | `/api/url/:code` | Delete a short URL |
| `GET` | `/:code` | Redirect to original URL |

## 🛠️ Development

For development with live reloading:
```bash
just air
```

## 📄 License

MIT License. See [LICENSE](LICENSE) for more information.
