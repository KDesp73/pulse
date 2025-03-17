# pulse

Pulse is an MQTT-based system designed to collect real-time data and provide 
live feeds and statistics through API endpoints. It can be easily integrated 
into custom frontends, offering users a flexible way to display and analyze 
incoming MQTT messages for various applications.

![image](https://github.com/user-attachments/assets/3cb9db60-5b7e-48cb-be33-18c8c80faec4)

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

Will generate `pulse.yml`

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
