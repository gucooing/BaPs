## Mitmproxy Scheme Usage Tutorial
(By KitanoSakura)

## Prerequisites
1. [Download mitmproxy](https://mitmproxy.org/) and install it.
2. Basic knowledge of WireGuard and Python scripts.
3. A client device (e.g., Android emulator or phone) and a host running `mitmproxy`.
> This tutorial will use an emulator.

### Installation Steps
- Linux/Mac
```markdown
  # Ubuntu/Debian

  sudo apt update
  sudo apt install mitmproxy

  # macOS
  brew install mitmproxy
  ```
- **Windows**: Download the `.exe` installer from [mitmproxy.org](https://mitmproxy.org/) and follow the instructions to complete the installation.

### Verify Installation
Run the following command to verify the installation:
```bash
mitmproxy --version
```

---

## Step 2: Install CA Certificate on Client and Server

To decrypt HTTPS traffic, the client needs to trust the `mitmproxy` CA certificate.

### Steps
1. Start `mitmproxy` to generate the certificate:
   ```bash
   mitmdump
   ```
2. Install the certificate (mitmproxy-ca.p12) on the computer at C:\Users\username\.mitmproxy
3. There will be mitmproxy-ca-cert.cer in the mitmproxy directory
4. Rename mitmproxy-ca-cert.cer to c8750f0d.0
5. Install the certificate as a system CA
---
1. Move the certificate to the system CA directory:
   ```bash
   adb root
   adb remount
   adb shell mv /sdcard/c8750f0d.0 /system/etc/security/cacerts/
   ```
2. Set the permissions:
   ```bash
   adb shell chmod 644 /system/etc/security/cacerts/c8750f0d.0
   ```
3. Reboot the device:
   ```bash
   adb reboot
   ```

---

## Step 3: Download the Script

```python
# KitanoSakura
# The script is not yet complete, please use WireGuard for proxy

from mitmproxy import http

# Define redirection rules
redirects = {
    "https://ba-jp-sdk.bluearchive.jp": "http://127.0.0.1:5000",
    "https://prod-gateway.bluearchiveyostar.com:5100/api/gateway": "http://127.0.0.1:5000/getEnterTicket/gateway",
    "https://prod-game.bluearchiveyostar.com:5000/api/gateway": "http://127.0.0.1:5000/api/gateway",
    "https://prod-logcollector.bluearchiveyostar.com:5300": "http://127.0.0.1:5000/game/log",
}

def request(flow: http.HTTPFlow) -> None:
    # Check if the request URL is in the redirection rules
    for original_url, redirected_url in redirects.items():
        if flow.request.pretty_url.startswith(original_url):
            # If it matches, modify the request URL to the local address
            flow.request.url = flow.request.pretty_url.replace(original_url, redirected_url)
            print(f"Redirecting {original_url} to {redirected_url}")
            return
```

---

## Step 4: Start mitmproxy and Load the Script

Run the following command to start `mitmproxy` with the redirection script:
```bash
mitmweb -m wireguard --no-http2 -s redirect_server.py --set termlog_verbosity=warn --ignore [your IP address here]
```

### Parameter Explanation:
- `-m wireguard`: Use WireGuard as the network layer.
- `--no-http2`: Disable HTTP/2 for better compatibility.
- `-s redirect_server.py`: Load the redirection script.
- `--set termlog_verbosity=warn`: Set log level to warning.

You can monitor traffic through the mitmweb interface at `http://localhost:8081`.

---

## Step 5: Install and Configure WireGuard

Use WireGuard to route client traffic to `mitmproxy`.

### Installation Steps
- **Android**: [Download WireGuard](https://play.google.com/store/apps/details?id=com.wireguard.android).
- **Other Platforms**: Refer to the [WireGuard official installation guide](https://www.wireguard.com/install/).

### Configuration Steps
1. Open the WireGuard client, click the + button in the lower-left corner, and select Scan QR code.
2. The emulator will pop up a scan window, select Real-time screenshot.
3. After selecting the screenshot, a capture image window will appear, move it to the QR code on the Mitmproxy browser page (if not, find it in settings).
4. Enable the proxy.

---

## Troubleshooting

### Client TLS handshake failed. The client does not trust the proxy's certificate for xxx.com (OpenSSL Error([('SSL routines', '', 'ssl/tls alert certificate unknown')]))
- Ensure the certificates on both the computer and client are the same.
- Ensure the mitmproxy certificate installed as system certificate.

### Certificate disappears after installation on Android?
- You can use MT Manager to grant SU permissions.
- Then go to /system/etc/security/cacerts/
- Find c8750f0d.0 and set it to 664 permissions. The user group should be root.

### (Mobile) No permission to modify the system directory?
- Install the certificate as a user certificate.
- Then install a module to automatically move the user certificate to the system certificate.
- Then restart the phone to see if the relevant certificate is present.