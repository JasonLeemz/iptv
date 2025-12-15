# IPTV M3U 转换工具

这个工具可以从 tonkiang.us 网站获取频道列表数据，并转换为标准的 M3U 格式。

## 功能

- 从 tonkiang.us 获取频道列表
- 自动解析 HTML 数据
- 转换为标准 M3U 格式
- 支持调试模式

## 使用方法

### 基本使用

程序会自动从 `config/source.txt` 文件中读取URL列表，汇总所有结果。

1. **准备配置文件**

   在 `config/source.txt` 文件中添加要爬取的URL，每行一个：

   ```
   https://tonkiang.us/channellist.html?ip=123.117.238.43&tk=c188c009&p=2
   https://tonkiang.us/channellist.html?ip=221.220.131.129&tk=c188c009&p=2
   https://tonkiang.us/channellist.html?ip=171.213.132.105&tk=c188c009&p=2
   ```

   支持注释行（以 `#` 开头）和空行。

2. **运行程序**

   ```bash
   go run main.go
   ```

   或使用编译后的二进制文件：

   ```bash
   ./iptv
   ```

3. **查看结果**

   程序会自动：
   - 遍历配置文件中的所有URL
   - 获取每个URL的频道数据
   - 自动去重（基于URL）
   - 汇总所有结果到 `output/output.m3u` 和 `output/output.txt`

### 设置 Cookies（重要）

由于网站使用 Cloudflare 保护，需要设置有效的 cookies 才能获取数据。

#### 方法1：从浏览器获取 Cookies

1. 在浏览器中打开 `https://tonkiang.us/channellist.html`
2. 打开开发者工具（F12 或 Cmd+Option+I）
3. 切换到 **Network** 标签
4. 刷新页面或执行搜索操作
5. 找到 `getall.php` 请求
6. 查看 **Request Headers** 中的 **Cookie** 值
7. 复制完整的 Cookie 字符串

然后设置环境变量：

```bash
export COOKIES='HstCfa4853344=...; cf_clearance=...; ...'
go run main.go "https://tonkiang.us/channellist.html?ip=123.117.238.43&tk=c188c009&p=2"
```

#### 方法2：从 curl 命令提取

如果你有完整的 curl 命令（包含 `-b` 参数），可以直接提取 cookies：

```bash
# 从 curl 命令的 -b 参数中提取 cookies
export COOKIES='HstCfa4853344=1748659924355; HstCnv4853344=1; ...'
go run main.go "https://tonkiang.us/channellist.html?ip=123.117.238.43&tk=c188c009&p=2"
```

**注意**：`cf_clearance` cookie 会定期过期，如果获取失败，需要重新从浏览器获取最新的 cookies。

### 调试模式

如果遇到问题，可以启用调试模式查看获取的 HTML 内容：

```bash
DEBUG=1 go run main.go
```

或

```bash
DEBUG=1 ./iptv
```

调试模式下，程序会将每次请求的 HTML 保存到 `debug.html` 文件中（会覆盖），方便检查数据结构。

## 功能特性

- ✅ 支持从配置文件批量读取URL
- ✅ 自动汇总多个源的所有频道
- ✅ 自动去重（基于URL）
- ✅ 同时生成M3U和CSV两种格式
- ✅ 输出到独立的output文件夹
- ✅ 显示处理进度和统计信息

## 输出

程序会在 `output` 文件夹中生成两个文件：

### 1. M3U 格式 (`output/output.m3u`)

标准的 M3U 播放列表格式，适用于大多数 IPTV 播放器：

```
#EXTM3U
#EXTINF:-1,频道名称1
http://stream.url1
#EXTINF:-1,频道名称2
http://stream.url2
...
```

### 2. CSV 格式 (`output/output.txt`)

简单的文本格式，每行一个源，格式为：`频道名称,接口地址`

```
CCTV-1 HD,http://123.117.238.43:5222/rtp/239.3.1.129:8008
CCTV-2 HD,http://123.117.238.43:5222/rtp/239.3.1.60:8084
CCTV-3 HD,http://123.117.238.43:5222/rtp/239.3.1.172:8001
...
```

这种格式便于：
- 导入到其他系统
- 批量处理
- 简单的文本编辑

## 注意事项

1. **Cookies 是必需的**：网站使用 Cloudflare 保护，必须设置有效的 cookies（特别是 `cf_clearance`）
2. **Cookies 会过期**：`cf_clearance` cookie 通常每 1-2 小时过期，需要定期更新
3. **Referer 必须正确**：程序会自动设置正确的 Referer header
4. **User-Agent**：程序使用完整的浏览器 User-Agent 来模拟真实请求
5. 如果获取不到数据：
   - 检查 cookies 是否有效（从浏览器重新获取）
   - 使用调试模式查看返回的 HTML
   - 确认 URL 参数是否正确

## 参数说明

URL 参数：
- `ip`: IP 地址或域名（必需）
- `tk`: 验证令牌（可选）
- `p`: 页码（可选）
- `c`: 其他参数（可选）

