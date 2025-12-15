# IPTV M3U 转换工具

一个功能完整的 IPTV 频道列表获取和转换工具，支持自动获取组播源、频道汇总、定时任务、日志记录和消息推送等功能。

## 功能特性

- ✅ **自动获取组播源**：从 tonkiang 自动获取最新的组播源 IP 列表
- ✅ **批量处理**：支持从配置文件批量读取 URL，自动汇总所有频道
- ✅ **自动去重**：基于 URL 自动去重，避免重复频道
- ✅ **多格式输出**：同时生成 M3U 和 CSV 两种格式
- ✅ **定时任务**：支持 Cron 表达式配置定时执行
- ✅ **日志系统**：完整的日志记录，支持 INFO、WARN、ERROR、DEBUG 级别
- ✅ **消息推送**：ERROR 日志自动推送到 Bark
- ✅ **文件重定向**：支持将输出文件自动拷贝到指定位置
- ✅ **服务管理**：提供启动、停止、重启脚本

## 项目结构

```
iptv/
├── bin/                    # 服务管理脚本
│   ├── start.sh           # 启动服务
│   ├── stop.sh            # 停止服务
│   └── restart.sh         # 重启服务
├── config/                 # 配置文件目录
│   ├── app.yml            # 主配置文件
│   ├── app.yml.example    # 配置示例文件
│   └── source.txt         # URL 列表文件
├── dto/                    # 数据传输对象
│   └── channel.go         # 频道数据结构
├── pkg/                    # 内部包
│   ├── bark/              # Bark 推送功能
│   ├── config/            # 配置管理
│   ├── cron/              # 定时任务
│   ├── html/              # HTML 解析工具
│   └── log/               # 日志系统
├── logs/                   # 日志文件目录
├── output/                 # 输出文件目录
│   ├── iptv.m3u          # M3U 格式输出
│   ├── local.txt         # CSV 格式输出
│   └── debug.html        # Debug HTML（如果启用）
├── main.go                # 主程序入口
├── task.go                # 任务执行逻辑
├── channel.go             # 频道解析功能
├── multicast.go           # 组播源获取功能
└── README.md              # 本文档
```

## 快速开始

### 1. 编译程序

```bash
go build -o iptv .
```

### 2. 配置应用

复制配置示例文件并编辑：

```bash
cp config/app.yml.example config/app.yml
vim config/app.yml
```

### 3. 运行程序

**方式一：直接运行**

```bash
./iptv
```

**方式二：使用服务管理脚本**

```bash
# 启动服务
./bin/start.sh

# 停止服务
./bin/stop.sh

# 重启服务
./bin/restart.sh
```

## 配置文件说明

配置文件位于 `config/app.yml`，主要配置项如下：

### 应用配置

```yaml
app:
  debug: true  # 是否启用 debug 模式
```

### 组播源配置

```yaml
multicastIP:
  limit: 5  # 从组播源列表获取的 IP 数量
```

### Cookie 配置

```yaml
cookie:
  data: "你的cookies值"  # 从浏览器获取的 cookies
```

**获取 Cookies 方法：**

1. 在浏览器中打开 `https://tonkiang.us/channellist.html`
2. 打开开发者工具（F12）
3. 切换到 **Network** 标签
4. 刷新页面或执行搜索操作
5. 找到 `getall.php` 请求
6. 查看 **Request Headers** 中的 **Cookie** 值
7. 复制完整的 Cookie 字符串到配置文件中

**注意**：`cf_clearance` cookie 会定期过期，需要定期更新。

### 定时任务配置

```yaml
crontab:
  enable: true           # 是否启用定时任务
  job: "0 1 * * *"       # Cron 表达式（每天凌晨1点）
```

Cron 表达式格式：`分 时 日 月 星期`

示例：
- `0 1 * * *` - 每天凌晨1点
- `0 */6 * * *` - 每6小时
- `0 0 * * 0` - 每周日凌晨

### 输出配置

```yaml
output:
  m3u: output/iptv.m3u        # M3U 格式输出文件
  local: output/local.txt      # CSV 格式输出文件
  debug: output/debug.html     # Debug HTML 文件
```

### 日志配置

```yaml
log:
  path: logs  # 日志文件目录
```

日志文件命名格式：`logs/app-YYYY-MM-DD.log`

### 推送配置

```yaml
push:
  bark:
    host: https://bark.abcd.xyz  # Bark 服务器地址
    key: "你的Bark密钥"            # Bark 推送密钥
```

**获取 Bark 密钥：**

1. 下载 Bark 应用
2. 在应用中获取设备密钥
3. 配置服务器地址和密钥

### 文件重定向配置

```yaml
redirectOutput:
  enable: true                              # 是否启用文件重定向
  move: output/local.txt                    # 源文件路径
  to: /path/to/destination/local.txt       # 目标文件路径
```

## 工作流程

程序执行流程：

1. **更新组播源列表**
   - 从 `iptvmulticast.php` 获取最新的组播源 IP
   - 自动更新 `config/source.txt` 文件

2. **读取 URL 列表**
   - 从 `config/source.txt` 读取所有 URL

3. **获取频道数据**
   - 遍历每个 URL，获取频道列表
   - 自动去重（基于 URL）

4. **输出结果**
   - 生成 M3U 格式文件
   - 生成 CSV 格式文件

5. **文件重定向**（如果启用）
   - 将输出文件拷贝到指定位置

## 输出格式

### M3U 格式 (`output/iptv.m3u`)

标准的 M3U 播放列表格式：

```
#EXTM3U
#EXTINF:-1,频道名称1
http://stream.url1
#EXTINF:-1,频道名称2
http://stream.url2
...
```

### CSV 格式 (`output/local.txt`)

简单的文本格式，每行一个源：

```
频道名称1,http://stream.url1
频道名称2,http://stream.url2
...
```

## 日志系统

程序使用完整的日志系统，所有输出都会记录到日志文件中：

- **日志位置**：`logs/app-YYYY-MM-DD.log`
- **日志级别**：INFO、WARN、ERROR、DEBUG
- **日志格式**：`[时间戳] [级别] 消息内容`
- **ERROR 推送**：ERROR 级别的日志会自动推送到 Bark

## 服务管理

### 启动服务

```bash
./bin/start.sh
```

功能：
- 检查程序是否已在运行
- 后台启动程序
- 保存进程 ID 到 `iptv.pid`

### 停止服务

```bash
./bin/stop.sh
```

功能：
- 优雅停止进程
- 如果 10 秒内未停止，强制终止
- 清理 PID 文件

### 重启服务

```bash
./bin/restart.sh
```

功能：
- 先停止服务
- 再启动服务

## 开发说明

### 项目依赖

- `github.com/PuerkitoBio/goquery` - HTML 解析
- `gopkg.in/yaml.v3` - YAML 配置解析

### 编译

```bash
go build -o iptv .
```

### 运行测试

```bash
go test ./...
```

## 注意事项

1. **Cookies 必需**：网站使用 Cloudflare 保护，必须配置有效的 cookies
2. **Cookies 会过期**：`cf_clearance` cookie 通常每 1-2 小时过期，需要定期更新
3. **网络连接**：确保能够访问目标网站
4. **文件权限**：确保程序有读写配置文件和输出目录的权限
5. **日志文件**：日志文件会持续增长，建议定期清理

## 故障排查

### 问题：获取不到频道数据

**可能原因：**
- Cookies 已过期
- 网络连接问题
- URL 参数不正确

**解决方法：**
1. 检查并更新 `config/app.yml` 中的 cookies
2. 启用 debug 模式，查看 `output/debug.html`
3. 检查日志文件中的错误信息

### 问题：定时任务不执行

**可能原因：**
- Cron 表达式格式错误
- 定时任务未启用

**解决方法：**
1. 检查 `config/app.yml` 中的 `crontab.enable` 是否为 `true`
2. 验证 Cron 表达式格式是否正确
3. 查看日志文件确认定时任务是否已添加

### 问题：Bark 推送失败

**可能原因：**
- Bark 配置不正确
- 网络连接问题

**解决方法：**
1. 检查 `config/app.yml` 中的 Bark 配置
2. 验证 Bark 服务器地址和密钥是否正确
3. 检查网络连接

## 许可证

本项目仅供学习和研究使用。

## 贡献

欢迎提交 Issue 和 Pull Request。
