# BanPickSys

## 项目简介

BanPickSys 是一个为电子竞技比赛设计的禁用/选择（Ban/Pick）系统。该系统允许多队伍在比赛前进行角色的禁用和选择，支持多种比赛格式和规则。

## 主要功能

- 支持多种常见的禁用/选择模式
- 可自定义的角色库和队伍信息
- 计时器功能，确保选择过程按时进行
- 历史记录功能，可查看和分析之前的禁选结果
- WebSocket 实时通信，保证游戏状态同步

## 使用方法

> 服务启动后访问 `http://localhost:10088`

### Docker

#### 自行构建 or 使用现有镜像

```bash
# 自行构建
# 进入项目根目录
git submodule update --remote	# 更新前端代码
docker build -t bpsys .
docker run -d -p 10088:10088 bpsys

# 使用现有镜像
docker run -d -p 10088:10088 jiny14/bpsys

# 使用 docker-compose
docker-compose up -d
```

### 本地运行

```bash
# 进入项目根目录
git submodule update --remote	# 更新前端代码
go run .
```

## API接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/bp` | GET | 启动一场新 bp |
| `/bp/{bpID}` | GET | 进入指定 bp 场次 |
| `/bp/{bpID}/status` | GET | 获取当前bp状态 |
| `/bp/{bpID}/entries` | GET | 获取可选角色列表 |
| `/bp/{bpID}/result` | GET | 获取bp结果 |
| `/bp/{bpID}/submit` | POST | 选择角色 |
| `/bp/{bpID}/join` | POST | 加入指定 bp 场次 |
| `/bp/{bpID}/leave` | POST | 离开指定 bp 场次 |


## 项目结构

```bash
BanPickSys
├── api             # API接口定义
│   ├── init.go
│   └── router.go
├── config
│   └── config.yaml
├── docker-compose.yaml
├── Dockerfile
├── game            # 游戏逻辑处理
│   └── handler.go
├── go.mod
├── go.sum
├── logs            # 日志文件
│   └── default.log
├── main.go         # 程序入口
├── model           # 数据模型定义
│   ├── bp.go
│   ├── entry.go
│   ├── player.go
│   └── stage.go
├── pkg             # 通用功能包
│   ├── init.go
│   ├── uuid.go
│   └── websocket.go
├── README.md
├── service         # 业务逻辑层
│   ├── bp.go
│   ├── stage.go
│   └── ws.go
└── static          # 前端资源
```

## 前端说明

前端实现了以下功能：
- 实时显示当前BP阶段和计时
- 角色选择界面，支持禁用和选择操作
- 队伍结果展示，包括已禁用和已选择的角色
- WebSocket通信，保证游戏状态实时同步

详细信息请参考 [前端README](https://github.com/jiny3/BanPickSys-web/blob/main/README.md)

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。

