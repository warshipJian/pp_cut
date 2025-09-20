# 平平剪切板

平平剪切板是一个基于Web的剪切板，支持跨端获取剪切板内容

成品展示：
http://c.xiaoxiaoguo.cn/

# 项目结构
```
.
├── backend (后端代码)
│   ├── main.go
├── frontend (前端代码)
│   ├── index.html
```

# frontend
静态页，直接部署

# backend

使用gin框架

## 初始化
go mod tidy

## 构建
go build .

## 使用 
./pp_cut

## 使用systemd守护
```
cat /lib/systemd/system/pp_cut.service
[Unit]
Description=pp cut

[Service]
ExecStart=/opt/pp_cut/pp_cut
StandardOutput=append:/var/log/pp_cut/service.log
StandardError=append:/var/log/pp_cut/service_error.log

[Install]
WantedBy=multi-user.target
```

## 使用NGINX代理
```
server {

	listen 443 ssl http2;

	ssl_certificate /path/to/xxx.pem; 
        ssl_certificate_key /path/to/xxx.key;

        ssl_session_cache shared:SSL:50m;
        ssl_session_timeout 1h;	

	add_header Access-Control-Allow-Origin '*';
        add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
        add_header Access-Control-Allow-Headers 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization';

        if ($request_method = 'OPTIONS') {
           return 204;
        }

	location / {
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forward-For $remote_addr;
	    proxy_pass http://127.0.0.1:8080;
	    keepalive_timeout 0;
	}
}
```
