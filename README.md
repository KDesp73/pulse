# pulse

Pulse provides real-time data streaming and statistical insights through API 
endpoints and Server-Sent Events (SSE). Users can create their own custom 
frontend to visualize and interact with live data. The backend, built in Go, 
handles data storage, updates, and statistics. 
With flexible data feeds and API support, the system can be easily adapted for 
a wide range of use cases, offering real-time monitoring and analytics for any 
application.

![image](https://github.com/user-attachments/assets/3cb9db60-5b7e-48cb-be33-18c8c80faec4)

<!--toc:start-->
  - [Get Started](#get-started)
    - [1. Clone the repo](#1-clone-the-repo)
    - [2. Build the application](#2-build-the-application)
    - [3. Generate and modify the configuration file](#3-generate-and-modify-the-configuration-file)
    - [4. Run pulse to see live updates of your IoT device](#4-run-pulse-to-see-live-updates-of-your-iot-device)
  - [License](#license)
<!--toc:end-->


## Get Started

### 1. Clone the repo

```bash
git clone https://github.com/KDesp73/pulse
cd pulse
```

### 2. Build the application

```bash
go build ./cmd/pulse
```

### 3. Generate and modify the configuration file

Running

```bash
./pulse -generate-config
```

Will generate

```yml
mqtt:
  server: "broker.mqtt"
  port: 8883
  username: "your-username"
  password: "your-password"
  topic: "your/topic"
web:
  port: 8080
  page: your-dashboard.html
```

### 4. Run pulse to see live updates of your IoT device

```bash
./pulse
```

## License

[MIT](./LICENSE)
