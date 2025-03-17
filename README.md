# Is It Live - Self-Hosted Uptime Monitoring Platform

A self-hosted uptime monitoring platform where users can create public dashboards with lifetime statistics, custom domains, and optional email alerts.

## Features

- **User System**
  - Sign up without email (username/password) with JWT authentication
  - Optional email field for downtime alerts
  - Public profile URLs: `/user/{username}` displaying all monitored sites

- **Uptime Monitoring**
  - Add sites via URL with backend checks every 60 seconds
  - Track response time, status codes, and uptime percentages (lifetime + 7/30/90-day stats)
  - Store granular data for historical graphs (1-minute intervals, aggregated daily for long-term)

- **Public Dashboards**
  - Grid layout showing all monitored sites with live status indicators
  - Interactive charts (response times, uptime trends) using Chart.js
  - Shareable URLs: `yourserver.com/site/{sitename}` or `yourserver.com/user/{username}`

- **Custom Domains**
  - Configure `status.theirdomain.com` to point to user dashboards
  - Backend validates DNS CNAME record (e.g., to `yourserver.com`)

## Tech Stack

- **Frontend**: Vue.js 3 (Composition API) + Pinia, Tailwind CSS (dark mode + teal theme)
- **Backend**: Go (Gin framework) with structured logging
- **Database**: PostgreSQL for users/sites + TimescaleDB for time-series metrics
- **Deployment**: Docker-compose for easy self-hosting

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Domain name (for production deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/is-site-live.git
   cd is-site-live
   ```

2. Configure environment variables:
   ```bash
   cp .env.example .env
   # Edit .env file with your settings
   ```

3. Start the application:
   ```bash
   docker-compose up -d
   ```

4. Access the application at `http://localhost:3000`

## Development Setup

### Backend (Go)

1. Install Go 1.21 or later
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the development server:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend (Vue.js)

1. Install Node.js 18 or later
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
3. Install dependencies:
   ```bash
   npm install
   ```
4. Run the development server:
   ```bash
   npm run dev
   ```

## Configuration

Configuration is done through environment variables:

- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: PostgreSQL username
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: PostgreSQL database name
- `JWT_SECRET`: Secret key for JWT tokens
- `SMTP_HOST`: SMTP server host (optional)
- `SMTP_PORT`: SMTP server port (optional)
- `SMTP_USER`: SMTP server username (optional)
- `SMTP_PASSWORD`: SMTP server password (optional)
- `SMTP_FROM`: Email sender address (optional)

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
