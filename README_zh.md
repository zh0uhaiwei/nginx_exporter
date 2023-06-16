# <p align="center">nginx_exporter</p>

[English](README.md) | 简体中文
### <a href="https://github.com/nginxinc/nginx-prometheus-exporter">NGINX官方exporter</a>并不支持多实例抓取，通过nginx_exporter可以实现多实例抓取

## 开始
### 配置Nginx Stub Status Module
```sh
location = /stub_status {
    stub_status;
  }
```

### 使用Docker
```shell
docker pull zh0uhaiwei/nginx_exporter:latest
docker run -it --name nginx_exporter -v config.yml:/app/config.yml -p 9113:9113 zh0uhaiwei/nginx_exporter:latest
```

### 使用Shell
```shell
mkdir -p /opt/nginx_exporter
cd /opt/nginx_exporter
wget https://github.com/zh0uhaiwei/nginx_exporter/releases/download/v1.0/nginx_exporter-v1.0-linux-x86-64
mv nginx_exporter-v1.0-linux-x86-64 nginx_exporter
chmod +x nginx_exporter
#add nginx host in config.yml
./nginx_exporter
```

### 浏览指标http://localhost:9113/metrics

### 配置Prometheus
Add a block to the `scrape_configs` of your prometheus.yml config file:
```yaml
scrape_configs:
- job_name: nginx_exporter
  static_configs:
  - targets: ["localhost:9113"]
```

### 配置Grafana，浏览面板
Upload `grafana/nginx_exporter.json` and view the dashboard.

## 许可
This software is free to use under the MIT License [MIT license](/LICENSE).
