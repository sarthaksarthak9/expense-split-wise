# 💰 Simplified expense-splite-wise – Expense Split Backend

A production-ready microservices backend for splitting expenses among groups, built with **Go**, **MongoDB**, **Redis**, and **RabbitMQ**.  
Perfect for learning backend architecture, async processing, and distributed systems.


## 🎯 Overview

This project is a backend system for splitting expenses among friends, similar to **Splitwise**. Users can:

- Create groups for shared expenses (e.g., trips, roommates, events)
- Add expenses with details on who paid and how to split
- View balances to see who owes whom
- Track expense history for each group

The system demonstrates real-world backend patterns including **REST APIs**, **async job processing**, **caching**, and **message queues**.


## ✨ Features

### Core Functionality

- ✅ Create and manage expense groups  
- ✅ Add members to groups dynamically  
- ✅ Record expenses with flexible splitting (equal split among members)  
- ✅ Automatic balance calculation  
- ✅ Real-time expense tracking  
- ✅ Complete expense history  

## 🚀 Quick Start

Get the entire backend running locally in just a few minutes using Docker Compose.

---

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/yourusername/expense-spliwise.git
cd expense-backend
```

### 2️⃣ Start All Services
```bash
# Build images and start all containers
docker-compose up --build
```

### 3️⃣ Verify Running Services
```bash
docker-compose ps
```

#### Expected output:
```
NAME                 STATUS                   PORTS
splitwise-api        Up (healthy)             0.0.0.0:8080->8080/tcp
splitwise-worker     Up (healthy)
splitwise-mongo      Up (healthy)             0.0.0.0:27017->27017/tcp
splitwise-redis      Up (healthy)             0.0.0.0:6379->6379/tcp
splitwise-rabbitmq   Up (healthy)             0.0.0.0:5672->5672/tcp, 15672/tcp
```

### 4️⃣ Test the API
```
curl http://localhost:8080/health
```

#### Expected response:
```
{
  "status": "ok"
}
```

🎉 Success! Your backend is up and running.
