# BanPickSys

## 项目简介

BanPickSys 是一个为电子竞技比赛设计的禁用/选择（Ban/Pick）系统。该系统允许多队伍在比赛前进行角色的禁用和选择，支持多种比赛格式和规则。

## 主要功能

- 支持多种常见的禁用/选择模式
- 可自定义的角色库和队伍信息
- 计时器功能，确保选择过程按时进行
- 历史记录功能，可查看和分析之前的禁选结果

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
```

### 本地运行

```bash
# 进入项目根目录
git submodule update --remote	# 更新前端代码
go run .
```

## API接口

- `/game` - 启动一场新 bp (GET)
- `/game/{gameId}` - 进入指定 bp 场次 (GET)
- `/game/{gameId}/status` - 获取当前游戏状态
- `/game/{gameId}/entries` - 获取可选角色列表
- `/game/{gameId}/result` - 获取游戏结果
- `/game/{gameId}` - 选择角色 (POST)

## 项目结构

```bash
BanPickSys
├── api	# api接口
│   ├── init.go
│   └── router.go
├── config
│   └── config.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go	# 入口文件
├── model	# 结构体定义
│   ├── entry.go
│   ├── game.go
│   ├── player.go
│   └── stage.go
├── pkg
│   ├── init.go
│   └── uuid.go
├── README.md
├── service	# 业务编排
│   ├── entry.go
│   ├── game.go
│   └── stage.go
└── static	# 前端
```

