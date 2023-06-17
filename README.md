# <p align="center">nginx_exporter</p>

English | [简体中文](README_zh.md)

### Inspired by <a href="https://github.com/nginxinc/nginx-prometheus-exporter">NGINX Prometheus Exporter</a>,this nginx_exporter is used to collect multiple nginx instance metrics through <a href="https://nginx.org/en/docs/http/ngx_http_stub_status_module.html">ngx_http_stub_status_module</a>.

## GetStarted
### Configure Nginx Stub Status Module
```sh
location = /stub_status {
    stub_status;
  }
```

### Use Docker
```shell
docker pull zh0uhaiwei/nginx_exporter:latest
docker run --name nginx_exporter -v $(pwd)/config.yml:/app/config.yml -p 9113:9113 zh0uhaiwei/nginx_exporter:latest
```

### Use Shell
```shell
mkdir nginx_exporter
cd nginx_exporter
wget https://github.com/zh0uhaiwei/nginx_exporter/releases/download/v1.0/nginx_exporter-v1.0-linux-x86-64
mv nginx_exporter-v1.0-linux-x86-64 nginx_exporter
chmod +x nginx_exporter
#add nginx host in config.yml
./nginx_exporter
```

### Browse http://localhost:9113/metrics

### Configure Prometheus
Add a block to the `scrape_configs` of your prometheus.yml config file:
```yaml
scrape_configs:
- job_name: nginx_exporter
  static_configs:
  - targets: ["localhost:9113"]
```

### Configure Grafana
Upload `grafana/nginx_exporter.json` and view the dashboard.

## License
This software is free to use under the MIT License [MIT license](/LICENSE).
