![introduce](https://socialify.git.ci/gucooing/BaPs/image?description=1&font=Source+Code+Pro&forks=1&issues=1&language=1&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

# BaPs 🎮

#### [Chinese](README.md)

> ⚠️ This project is for learning purposes only, strictly prohibited for commercial use, please delete within 24 hours.

# For learning purposes only, strictly prohibited for commercial use, please delete within 24 hours!!!

> 🌟 Due to its stateless design, it may require slightly higher memory
  
> 📅 Currently supported version: Japan 1.54.327262

## 📍Discord

[![Discord](https://img.shields.io/badge/Join-Discord-blue?logo=discord&logoSize=auto)](https://discord.gg/222yVp6pUq)


---
## 🚀 Features Implemented
```
- Login  
- Newbie tutorial  
- Team management  
- Gacha  
- Story (pending testing)  
- Basic account management  
- MomoTalk  
- Mail global/personal management  
- Character growth management  
- Inventory management  
- Dungeons - Bounty / Commission / Scrimmage / Joint Firing Drill  
- Auto restoration of recoverable items  
- Cafe  
- Friend management  
- Lesson  
- Club  
- Battle Assistance  
- Total Assault  
- Daily login rewards  
- Final Restriction  
- Grand Assault
- Shop
```
---
## 🛠️ Usage

#### Prerequisites

1. Go to [Releases](https://github.com/gucooing/BaPs/releases/latest) and download the latest version and copy it to the run directory (please download according to your system)
2. Copy the data folder from the repository to the run directory
3. Download the Excel.bin file from [Releases](https://github.com/gucooing/BaPs/releases/latest) and replace it in the data folder
4. Run once to automatically generate the config.json file, open and edit the config.json file
5. Run the app again

> If Excel.bin cannot be found, please download it from the data folder in the source code
---

### 🐳 Docker Deployment
```bash
docker run -d \
  -p 5000:5000 \
  -v /data/baps/config.json:/usr/ba/config.json \
  -v /data/baps/sqlite/BaPs.db:/usr/ba/BaPs.db \
  -v /data/baps/sqlite/Rank.db:/usr/ba/Rank.db \
  ghcr.io/gucooing/baps:latest
``` 
<details>
You have expanded an available mirror acceleration, this mirror acceleration site comes from the network

```
docker run -d \
  -p 5000:5000 \
  -v /data/baps/config.json:/usr/ba/config.json \
  -v /data/baps/sqlite/BaPs.db:/usr/ba/BaPs.db \
  -v /data/baps/sqlite/Rank.db:/usr/ba/Rank.db \
  ghcr.nju.edu.cn/gucooing/baps:latest
```
</details>

---

## ⚙️ Configuration Instructions
> Note that comments cannot exist in the actual json file
```
{
  "LogLevel": "info",
  "ResourcesPath": "./resources", // Not used in release version
  "DataPath": "./data",
  "GucooingApiKey": "123456", // Key to authenticate when using API
  "AutoRegistration": true, // Auto registration
  "Tutorial": false, // Enable tutorial - incomplete
  "HttpNet": {
    "InnerAddr": "0.0.0.0", // Listening address
    "InnerPort": "5000", // Listening port
    "OuterAddr": "10.0.0.3", // External address
    "OuterPort": "5000", // External port
    "Tls": false, // Enable SSL
    "CertFile": "./data/cert.pem",
    "KeyFile":   "./data/key.pem"
  },
  "GateWay": {
    "MaxPlayerNum": 0, // Max online players
    "MaxCachePlayerTime": 720, // Max player cache time
    "BlackCmd": {}, // Not used in release version
    "IsLogMsgPlayer": true // Not used in release version
  },
  "DB": {
    "dbType": "sqlite", // Database type, supports sqlite and mysql
    "dsn": "BaPs.db" // Database address, if mysql please fill mysql url
  },
  "RaidRankDB": {
    "dbType": "sqlite", // Database type, supports sqlite and mysql
    "dsn": "Rank.db" // Database address, if mysql please fill mysql url
  },
  "Irc": { // Can use general IRC server address
    "HostAddress": "127.0.0.1", // Club chat server IRC address
    "Port": 16666, // Club chat server IRC port
    "Password": "mx123" // Club chat server IRC password
  }
}
```
---

## 🌐 Proxy Settings
Proxy the following addresses: where `http://127.0.0.1:5000` is the server address
```plaintext
https://ba-jp-sdk.bluearchive.jp  →  http://127.0.0.1:5000
https://yostar-serverinfo.bluearchiveyostar.com → http://127.0.0.1:5000
```

### ⛓️Proxy Solution

You can view the following docs
- [Android_MitmProxy Proxy Solution](Android_Mitmproxy_Readme_EN.md)

---

## ⌨️ GM Tool
Go to the following repository to download the GM tool for convenient use

- [BlueArchiveGM](https://github.com/PrimeStudentCouncil/BlueArchiveGM/releases/latest)

Online version of GM menu

- [BlueArchiveGM Web](https://gm.bluearchive.cc)

Default connection address: `http://127.0.0.1:5000` 
Default key: `123456` (can be changed in config.json)

**GM updates may not be timely, if you need to use locally, please use the local version**
**The online version of GM supports local use too~**

---
## 🤝 Contribute
We welcome everyone who wants to help us, you can help us in the following ways:
- 🐛 Submit an Issue to report problems
- 💡 Submit a Pull Request to improve the code
- 📖 Improve project documentation
- 🚀 Join the Discord channel to provide suggestions
---

## ⚠️ Notes
1. Due to copyright reasons, the resources used in dev will not be made public
2. Due to copyright reasons, some source codes will not be made public, but we can guarantee that the non-public parts of the code have no malicious content
3. Player data will not be saved to the database in real-time, if you need the latest data, you can access player data through the API

---
## 🤜 Acknowledgements

- Thanks to [zset](https://github.com/liyiheng/zset) for implementing the leaderboard