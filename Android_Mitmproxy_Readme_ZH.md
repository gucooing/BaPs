## Mitmproxy方案使用教程
（By 北野桜奈）

## 前置要求
 1. [下载 mitmproxy](https://mitmproxy.org/) 并安装。
 2. 具备 WireGuard 和 Python 脚本的基本知识。
 3. 一台客户端设备（例如 Android 模拟器或手机）以及运行 `mitmproxy` 的主机。
 >本次将以模拟器的方式进行教程
### 安装步骤
- Linux/Mac
```markdown
  # Ubuntu/Debian

  sudo apt update
  sudo apt install mitmproxy

  # macOS
  brew install mitmproxy
  ```
- **Windows**: 从 [mitmproxy.org](https://mitmproxy.org/) 下载 `.exe` 安装程序，并按说明完成安装。

### 验证安装
运行以下命令验证安装是否成功：
```bash
mitmproxy --version
```

---

## 第二步：在客户端以及服务端安装 CA 证书

为了解密 HTTPS 流量，客户端需要信任 `mitmproxy` 的 CA 证书。

### 操作步骤
 1. 启动 `mitmproxy`生成证书：
   ```bash
   mitmdump
   ```
 2. 在电脑端的C:\Users\用户\ .mitmproxy安装电脑证书（mitmproxy-ca.p12）
 3. mitmproxy的目录下会有mitmproxy-ca-cert.cer
 4. 将 mitmproxy-ca-cert.cer 重命名为 c8750f0d.0
 5. 将证书安装为系统 CA
---
 1. 将证书移动到系统 CA 目录：
   ```bash
   adb root
   adb remount
   adb shell mv /sdcard/c8750f0d.0 /system/etc/security/cacerts/
   ```
 2. 设置正确的权限：
   ```bash
   adb shell chmod 644 /system/etc/security/cacerts/c8750f0d.0
   ```
 3. 重启设备：
   ```bash
   adb reboot
   ```

---

## 第三步：下载脚本

```python
# KitanoSakura
# 脚本还没完善，请使用WireGuard进行代理

from mitmproxy import http

# 定义重定向规则
redirects = {
    "https://ba-jp-sdk.bluearchive.jp": "http://127.0.0.1:5000",
    "https://prod-gateway.bluearchiveyostar.com:5100/api/gateway": "http://127.0.0.1:5000/getEnterTicket/gateway",
    "https://prod-game.bluearchiveyostar.com:5000/api/gateway": "http://127.0.0.1:5000/api/gateway",
    "https://prod-logcollector.bluearchiveyostar.com:5300": "http://127.0.0.1:5000/game/log",
}

def request(flow: http.HTTPFlow) -> None:
    # 判断请求的URL是否在重定向规则中
    for original_url, redirected_url in redirects.items():
        if flow.request.pretty_url.startswith(original_url):
            # 如果匹配，修改请求的URL为本地地址
            flow.request.url = flow.request.pretty_url.replace(original_url, redirected_url)
            print(f"Redirecting {original_url} to {redirected_url}")
            break
```

---

## 第四步：启动 mitmproxy 并加载脚本

运行以下命令以使用重定向脚本启动 `mitmproxy`：
```bash
mitmweb -m wireguard --no-http2 -s redirect_server.py --set termlog_verbosity=warn --ignore 这里输入你的IP地址
```

### 参数说明：
- `-m wireguard`: 使用 WireGuard 作为网络层。
- `--no-http2`: 禁用 HTTP/2 以提高兼容性。
- `-s redirect_server.py`: 加载重定向脚本。
- `--set termlog_verbosity=warn`: 设置日志级别为警告。

你可以通过 `http://localhost:8081` 访问 `mitmweb` 界面监控流量。

---

## 第五步：安装并配置 WireGuard

使用 WireGuard 将客户端流量路由到 `mitmproxy`。

### 安装步骤
- **Android**: [下载 WireGuard](https://play.google.com/store/apps/details?id=com.wireguard.android)。
- **其他平台**: 参考 [WireGuard 官方安装指南](https://www.wireguard.com/install/)。

### 配置步骤
 1. 打开 WireGuard 客户端，点击左下角＋号，选择扫描二维码
 2. 选择后模拟器会弹出扫一扫窗口，选择实时截屏
 3. 选择截屏后，会有获取图像窗口，移动到Mitmproxy浏览器页面上的二维码（没有的话在设置里面）
 4. 启用该配置。

---

## 故障排查

### Client TLS handshake failed. The client does not trust the proxy's certificate for xxx.com (OpenSSL Error([('SSL routines', '', 'ssl/tls alert certificate unknown')]))
- 确保电脑端以及客户端证书为内容一样的
- 确保双端安装了Mitmproxy证书

### 安卓端安装后证书消失？
- 可以使用MT管理器授予SU权限
- 然后前往/system/etc/security/cacerts/
- 找到 c8750f0d.0 给予 664 权限。用户组为 root

### （手机端）无权限修改系统目录？
- 安装证书到用户证书
- 然后安装模块自动将用户证书转到系统证书
- 之后重启手机查看有没有相关证书
---