# BrandTrekin 部署指南

## 快速开始

### 前提条件
- Go 1.23+ ✅ (已安装 1.24.5)
- Node.js > v18.16.0 ✅ (已安装 v22.19.0)
- MySQL 8.0+ ✅ (阿里云RDS已配置)
- **服务器内存至少2GB** ❌ (当前仅904MB，需升级)

---

## 方案一：推荐部署流程（在2GB+内存服务器上）

### 1. 后端部署

```bash
# 进入后端目录
cd /home/ec2-user/gin-vue-admin/server

# 下载依赖
go mod download

# 编译服务器
go build -o server .

# 运行服务器
./server

# 或者使用后台运行
nohup ./server > /tmp/gva-server.log 2>&1 &

# 查看日志
tail -f /tmp/gva-server.log
```

**服务器启动成功标志**:
```
[GIN-debug] Listening and serving HTTP on :8888
```

### 2. 前端部署

```bash
# 进入前端目录
cd /home/ec2-user/gin-vue-admin/web

# 安装依赖（如果还没安装）
npm install

# 开发模式运行
npm run serve

# 或者构建生产版本
npm run build
```

**开发服务器访问**:
```
http://localhost:8080
```

### 3. 访问页面

```bash
# 市场列表页
http://localhost:8080/#/markets

# 市场详情页
http://localhost:8080/#/markets/CNCRouter

# 品牌详情页
http://localhost:8080/#/markets/CNCRouter/brands/STEPPERONLINE
```

### 4. 测试数据导入

```bash
# 使用Postman或curl导入测试数据
# 首先需要登录获取token

# 登录
curl -X POST http://localhost:8888/base/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'

# 获取token后，执行批量导入
# 使用 /trekin-main/data/CNCRouter/ 下的文件
```

---

## 方案二：当前服务器部署（交叉编译）

由于当前服务器内存不足，无法直接编译Go程序。可以在本地机器编译后上传：

### 在本地机器（Mac/Linux/Windows）

```bash
# 克隆代码到本地
git clone <repository>
cd gin-vue-admin/server

# 交叉编译为Linux AMD64二进制
GOOS=linux GOARCH=amd64 go build -o server-linux .

# 上传到服务器
scp server-linux ec2-user@<server-ip>:/home/ec2-user/gin-vue-admin/server/server

# 上传配置文件（如果有修改）
scp config.yaml ec2-user@<server-ip>:/home/ec2-user/gin-vue-admin/server/
```

### 在服务器上运行

```bash
cd /home/ec2-user/gin-vue-admin/server

# 赋予执行权限
chmod +x server

# 运行
./server
```

---

## 方案三：Docker部署（推荐生产环境）

### 1. 后端Docker部署

```bash
cd /home/ec2-user/gin-vue-admin/server

# 构建镜像
docker build -t brandtrekin-backend:latest .

# 运行容器
docker run -d \
  --name brandtrekin-backend \
  -p 8888:8888 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  brandtrekin-backend:latest

# 查看日志
docker logs -f brandtrekin-backend
```

### 2. 前端Docker部署

```bash
cd /home/ec2-user/gin-vue-admin/web

# 构建镜像
docker build -t brandtrekin-frontend:latest .

# 运行容器
docker run -d \
  --name brandtrekin-frontend \
  -p 8080:80 \
  brandtrekin-frontend:latest
```

### 3. 使用Docker Compose（最简单）

创建 `docker-compose.yml`:
```yaml
version: '3.8'

services:
  backend:
    build: ./server
    ports:
      - "8888:8888"
    volumes:
      - ./server/config.yaml:/app/config.yaml
    environment:
      - GIN_MODE=release
    restart: always

  frontend:
    build: ./web
    ports:
      - "8080:80"
    depends_on:
      - backend
    restart: always
```

运行：
```bash
docker-compose up -d
```

---

## 数据库配置检查

当前配置（已在 `server/config.yaml` 中更新）:
```yaml
mysql:
  prefix: ""
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: brandtrekin
  username: brandtrekin
  password: "cl@2025@!"
  path: rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com
  engine: ""
  log-mode: error
  max-idle-conns: 10
  max-open-conns: 100
  singular: false
  log-zap: false
```

### 验证数据库连接

```bash
# 使用mysql客户端测试
mysql -h rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com \
      -u brandtrekin \
      -p'cl@2025@!' \
      -e "SELECT VERSION();"
```

---

## 生成Swagger文档

```bash
cd /home/ec2-user/gin-vue-admin/server

# 安装swag（如果还没安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init

# 启动服务器后访问
# http://localhost:8888/swagger/index.html
```

---

## 测试数据导入流程

### 1. 准备测试文件
```
测试文件位置：/home/ec2-user/gin-vue-admin/trekin-main/data/

可用市场：
- CNCRouter/
- LaserEngraver/
- ThermalCamera/

每个市场包含5个文件：
- Brand-Social.xlsx
- GKW.csv
- KeywordHistory.xlsx
- Product-US.xlsx
- product-US-sales.xlsx
```

### 2. 创建市场（如果不存在）
```sql
-- 连接到数据库
mysql -h rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com \
      -u brandtrekin \
      -p'cl@2025@!' \
      brandtrekin

-- 插入市场
INSERT INTO bt_markets (market_slug, market_name, status)
VALUES ('CNCRouter', 'CNC Router', 'active');
```

### 3. 导入数据（使用Postman）
```
URL: POST http://localhost:8888/btImport/batchImport
Headers:
  - Authorization: Bearer {token}
Content-Type: multipart/form-data

Form Data:
  - marketId: 1
  - replaceMode: true
  - brandSocial: [选择文件] Brand-Social.xlsx
  - gkw: [选择文件] GKW.csv
  - keywordHistory: [选择文件] KeywordHistory.xlsx
  - productUS: [选择文件] Product-US.xlsx
  - productSales: [选择文件] product-US-sales.xlsx
```

### 4. 验证数据
```sql
-- 检查品牌数量
SELECT COUNT(*) FROM bt_brands WHERE market_id = 1;

-- 检查产品数量
SELECT COUNT(*) FROM bt_products WHERE market_id = 1;

-- 检查CAGR计算
SELECT brand_name, cagr FROM bt_brands
WHERE market_id = 1
ORDER BY cagr DESC LIMIT 10;

-- 检查月度趋势
SELECT * FROM bt_brand_monthly_trends
WHERE brand_id = 1
ORDER BY date;
```

---

## 生产环境配置建议

### 1. Nginx反向代理配置

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # 前端静态文件
    location / {
        root /var/www/brandtrekin/web/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 后端其他路由
    location ~ ^/(btImport|base|user|fileUploadAndDownload|system|sysDictionary|swagger)/ {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 2. SSL证书配置（HTTPS）

```bash
# 使用Let's Encrypt免费证书
sudo certbot --nginx -d yourdomain.com
```

### 3. 系统服务配置（systemd）

创建 `/etc/systemd/system/brandtrekin-backend.service`:
```ini
[Unit]
Description=BrandTrekin Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user/gin-vue-admin/server
ExecStart=/home/ec2-user/gin-vue-admin/server/server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable brandtrekin-backend
sudo systemctl start brandtrekin-backend
sudo systemctl status brandtrekin-backend
```

---

## 常见问题排查

### 1. 后端无法启动
```bash
# 检查端口占用
lsof -i :8888

# 检查数据库连接
mysql -h rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com \
      -u brandtrekin -p'cl@2025@!'

# 查看日志
tail -f /tmp/gva-server.log
```

### 2. 前端无法连接后端
```bash
# 检查代理配置
cat web/.env.development

# 应该包含：
VITE_BASE_URL=http://127.0.0.1
VITE_BASE_PATH=/
VITE_SERVER_PORT=8888
```

### 3. 数据导入失败
```bash
# 检查文件权限
ls -l /home/ec2-user/gin-vue-admin/trekin-main/data/CNCRouter/

# 检查数据库连接
# 查看后端日志中的错误信息
```

### 4. CAGR未计算
```sql
-- 检查月度趋势数据
SELECT brand_id, COUNT(*) as months
FROM bt_brand_monthly_trends
GROUP BY brand_id;

-- CAGR需要至少12个月的数据
-- 如果少于12个月，CAGR会是NULL
```

---

## 性能优化建议

### 1. 数据库索引
```sql
-- 创建必要的索引
CREATE INDEX idx_market_slug ON bt_markets(market_slug);
CREATE INDEX idx_brand_name ON bt_brands(brand_name, market_id);
CREATE INDEX idx_product_asin ON bt_products(asin);
CREATE INDEX idx_brand_trend_date ON bt_brand_monthly_trends(brand_id, date);
CREATE INDEX idx_market_trend_date ON bt_market_monthly_trends(market_id, date);
```

### 2. Redis缓存（可选）
```bash
# 安装Redis
sudo yum install redis -y

# 启动Redis
sudo systemctl start redis
sudo systemctl enable redis

# 在config.yaml中配置Redis
redis:
  db: 0
  addr: 127.0.0.1:6379
  password: ""
```

### 3. 日志级别调整
```yaml
# config.yaml
zap:
  level: 'info'  # 生产环境建议使用info，开发用debug
  format: 'json'
  log-in-console: false
```

---

## 监控和维护

### 1. 服务器资源监控
```bash
# 内存使用
free -h

# CPU使用
top

# 磁盘使用
df -h
```

### 2. 应用日志监控
```bash
# 后端日志
tail -f /tmp/gva-server.log

# Nginx访问日志
tail -f /var/log/nginx/access.log

# Nginx错误日志
tail -f /var/log/nginx/error.log
```

### 3. 数据库备份
```bash
# 定时备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -h rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com \
          -u brandtrekin \
          -p'cl@2025@!' \
          brandtrekin > /backup/brandtrekin_$DATE.sql

# 保留最近7天的备份
find /backup/ -name "brandtrekin_*.sql" -mtime +7 -delete
```

添加到crontab:
```bash
# 每天凌晨2点备份
0 2 * * * /path/to/backup-script.sh
```

---

## 下一步开发建议

### 短期（1-2周）
1. [ ] 在2GB+服务器上完成部署和测试
2. [ ] 导入所有3个市场的测试数据
3. [ ] 验证前后端集成
4. [ ] 完成Swagger文档

### 中期（1个月）
1. [ ] 添加API分页支持
2. [ ] 添加Redis缓存
3. [ ] 优化数据库索引
4. [ ] 添加数据导入进度通知

### 长期（2-3个月）
1. [ ] 用户权限管理
2. [ ] 数据导出功能
3. [ ] 多时间段对比
4. [ ] Prometheus监控集成

---

## 联系和支持

- 项目仓库：[GitHub链接]
- 问题反馈：[Issues链接]
- 文档：`BRANDTREKIN_测试报告.md`

---

最后更新：2025-10-31
