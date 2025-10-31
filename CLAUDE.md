# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在本仓库中工作时提供指导。

## 项目概述

**gin-vue-admin** 是一个基于 Vue 3 和 Gin 框架的全栈前后端分离管理平台。提供 JWT 鉴权、动态路由、Casbin 权限管理、表单生成器、代码生成器以及完整的 CRUD 功能。

- **在线演示**: http://demo.gin-vue-admin.com (admin/123456)
- **项目文档**: https://www.gin-vue-admin.com
- **视频教程**: https://www.bilibili.com/video/BV1Rg411u7xH

## 技术栈

### 后端 (Go 1.23+)
- **Gin** 1.10.0 - Web 框架
- **GORM** 1.25.12 - ORM 框架，支持自动迁移
- **Casbin** 2.103.0 - RBAC 权限管理
- **JWT** 5.2.2 - 身份认证
- **Viper** 1.19.0 - YAML 配置管理
- **Zap** 1.27.0 - 结构化日志
- **Redis** 9.7.0 - 缓存和会话管理
- **Swagger** - API 文档生成（通过 swaggo）
- **MCP** mark3labs/mcp-go v0.41.1 - AI 开发支持

### 前端 (Node > v18.16.0)
- **Vue** 3.5.7 - 使用 Composition API
- **Vite** 6.2.3 - 构建工具
- **Element Plus** 2.10.2 - UI 组件库
- **Pinia** 2.2.2 - 状态管理（取代 Vuex）
- **UnoCSS** 66.4.2 - 原子化 CSS
- **Axios** 1.8.2 - HTTP 客户端（通过统一的 request.js）
- **ECharts** 5.5.1 - 数据可视化

### 支持的数据库
MySQL（默认）、PostgreSQL、SQLite、SQL Server、MongoDB、Oracle

### 支持的云存储
阿里云 OSS、AWS S3、MinIO、七牛云、腾讯云 COS、Cloudflare R2、华为云 OBS

## 开发命令

### 后端 (server/)
```bash
# 启动开发服务器
cd server
go generate              # 安装依赖
go run .                # 在 8888 端口运行服务器

# 生成 Swagger 文档
swag init               # 生成 docs/swagger.json, docs/swagger.yaml, docs/docs.go
# 访问地址: http://localhost:8888/swagger/index.html

# 测试
go test ./...           # 运行所有测试
```

### 前端 (web/)
```bash
# 开发环境
cd web
npm install             # 安装依赖
npm run serve          # 在 8080 端口启动开发服务器
npm run dev            # 替代的开发命令

# 生产环境
npm run build          # 构建生产版本
npm run preview        # 预览生产构建
npm run limit-build    # 增加内存限制进行构建
```

### 使用 VSCode 工作区
打开 `gin-vue-admin.code-workspace`，使用 "Both (Backend & Frontend)" 启动配置可同时运行前后端。

### 使用 Makefile 构建
```bash
make build-local              # 同时构建前端和后端
make build-server-local       # 仅构建后端
make build-web-local         # 仅构建前端
make build-image-server      # 构建后端 Docker 镜像
make build-image-web         # 构建前端 Docker 镜像
make plugin PLUGIN="name"    # 将插件打包为 zip
```

## 架构概览

### 后端：严格的四层架构

**关键要点**：后端采用严格的分层架构，**绝不允许跨层调用**：

```
路由层 (Router) → API 层 → 服务层 (Service) → 模型层 (Model)
```

**各层职责：**

1. **模型层 (Model Layer)** (`server/model/`)
   - 数据库实体（GORM 模型）
   - 请求 DTO (`model/*/request/`)
   - 响应 DTO (`model/*/response/`)
   - 必须继承 `global.GVA_MODEL` 以包含基础字段（ID、CreatedAt、UpdatedAt、DeletedAt）

2. **服务层 (Service Layer)** (`server/service/`)
   - 业务逻辑实现
   - 通过 GORM 进行数据库 CRUD 操作
   - **禁止**使用 `gin.Context` - 该层与 HTTP 无关
   - 返回结果和 `error` 对象
   - 入口点：`service/enter.go` 暴露 `ServiceGroupApp`

3. **API 层 (API Layer)** (`server/api/v1/`)
   - HTTP 请求处理器
   - 参数绑定和验证
   - 通过 `service.ServiceGroupApp` 调用服务层
   - 通过 `response.OkWithDetailed()` 或 `response.FailWithMessage()` 返回格式化的 JSON
   - **强制要求**：每个端点必须有完整的 Swagger 注释
   - 入口点：`api/v1/enter.go` 暴露 `ApiGroupApp`

4. **路由层 (Router Layer)** (`server/router/`)
   - 路由定义和中间件配置
   - 通过 `api.ApiGroupApp` 将 URL 映射到 API 处理器
   - 配置中间件（JWT、Casbin、CORS、日志等）
   - 入口点：`router/enter.go` 暴露 `RouterGroupApp`

5. **初始化层 (Initialize Layer)** (`server/initialize/`)
   - 系统引导（DB、Redis、Router、插件）
   - 插件注册
   - 配置加载

### 前端：基于组件的 Vue 3 架构

```
视图组件 (View) → API 服务 → Axios (通过 request.js) → 后端
        ↓
   Pinia 状态管理
```

**目录结构：**
- `web/src/api/` - API 服务层（所有 HTTP 调用）
- `web/src/view/` - 页面组件（仅使用 Composition API）
- `web/src/components/` - 可复用的 UI 组件
- `web/src/pinia/modules/` - Pinia 状态管理存储
- `web/src/router/` - Vue Router 配置（动态路由）
- `web/src/utils/request.js` - 统一的 Axios 封装（添加 JWT token，处理错误）

**核心模式：**
- **仅使用 Composition API** - 不使用 Options API
- **Pinia 状态管理** - 取代 Vuex
- **统一请求处理** - 所有 HTTP 调用通过 `request.js`
- **动态路由** - 根据用户权限加载路由
- **优先使用 UnoCSS** - 原子化 CSS 优于自定义样式

## 关键开发规则

### 后端开发规则

1. **严格的层级分离**
   - **绝不**跨层调用（例如：API 层不能直接调用 Model 层）
   - **绝不**在 Service 层使用 `gin.Context`
   - Service 层返回 `error`，API 层格式化 HTTP 响应

2. **enter.go 模式**
   - 所有模块**必须**使用 `enter.go` 暴露组：
     - `service/enter.go` → `ServiceGroupApp`
     - `api/v1/enter.go` → `ApiGroupApp`
     - `router/enter.go` → `RouterGroupApp`
   - 这可以防止循环依赖

3. **强制 Swagger 文档**
   - 每个 API 端点**必须**有完整的 Swagger 注释：
     ```go
     // CreateXxx 创建XXX
     // @Tags     XxxModule
     // @Summary  创建一个新的XXX
     // @Security ApiKeyAuth
     // @accept   application/json
     // @Produce  application/json
     // @Param    data body request.CreateXxxRequest true "请求参数"
     // @Success  200  {object} response.Response{data=model.Xxx,msg=string} "成功"
     // @Router   /xxx/createXxx [post]
     func (a *XxxApi) CreateXxx(c *gin.Context) { ... }
     ```
   - 添加/修改端点后运行 `swag init`

4. **数据类型一致性**
   - 字段类型在 Model、Request 和 Response 结构体中**必须**保持一致
   - 在模型中使用指针类型（`*string`、`*int`）时，在 Service 层正确处理 nil 检查
   - 常见错误：在 Model 和 Request DTO 中对同一字段使用不同类型

5. **错误处理**
   - Service 层：返回 `error` 对象
   - API 层：使用 `response` 包转换为 HTTP 响应

### 前端开发规则

1. **仅使用 Composition API**
   - 所有组件**必须**使用 Vue 3 Composition API
   - 不允许使用 Options API

2. **API 层隔离**
   - 所有 HTTP 调用**必须**通过 `src/api/` 模块
   - **绝不**在组件中直接调用 `axios`
   - 在 API 模块中使用 `@/utils/request.js` 封装

3. **状态管理**
   - 全局状态**必须**使用 Pinia stores
   - **绝不**直接修改状态 - 使用 actions
   - Store 文件放在 `src/pinia/modules/`

4. **样式规范**
   - **优先**使用 UnoCSS 原子化类
   - **避免**内联样式
   - 使用 CSS 变量进行主题定制

5. **组件指南**
   - 单一职责原则
   - 完整的 props 和 emits 定义
   - 为复杂组件添加 JSDoc 注释

### 插件开发

**后端插件结构** (`server/plugin/{name}/`):
```
plugin/{name}/
├── api/              # API 处理器
│   └── enter.go     # ApiGroup 入口
├── config/          # 插件配置结构
├── initialize/      # 初始化函数
│   ├── api.go      # 注册 API
│   ├── gorm.go     # 数据库迁移（AutoMigrate）
│   ├── menu.go     # 菜单初始化
│   ├── router.go   # 注册路由
│   └── viper.go    # 加载配置
├── model/           # 数据模型
│   └── request/    # 请求 DTO
├── router/          # 路由定义
│   └── enter.go    # RouterGroup 入口
├── service/         # 业务逻辑
│   └── enter.go    # ServiceGroup 入口
└── plugin.go        # 插件入口（实现 system.Plugin 接口）
```

**前端插件结构** (`web/src/plugin/{name}/`):
```
plugin/{name}/
├── api/             # API 服务调用
├── view/           # 插件页面
├── components/     # 插件组件
└── form/           # 表单组件
```

**参考示例**：学习 `server/plugin/announcement/` 和 `web/src/plugin/announcement/` 了解完整的插件示例。

## 认证与授权

- **认证 (Authentication)**：通过 `/middleware/jwt.go` 的 JWT tokens
  - Token 在 `x-token` header 中发送
  - 用户 ID 在 `x-user-id` header 中
  - Token 刷新通过 `new-token` response header

- **授权 (Authorization)**：通过 `/middleware/casbin_rbac.go` 的 Casbin RBAC
  - 基于角色的访问控制
  - 动态 API 权限
  - 前端：使用 `v-auth` 指令进行按钮级权限控制

## 配置管理

### 后端配置 (`server/config.yaml`)
- **system**：服务器端口（8888）、数据库类型、OSS 类型
- **mysql/pgsql/sqlite/oracle/mssql**：数据库连接
- **redis**：缓存配置
- **jwt**：Token 签名和过期时间
- **zap**：日志级别和输出
- **captcha**：验证码设置
- **cors**：CORS 策略
- **autocode**：代码生成器设置
- **mcp**：MCP 服务器配置（端口 8889）
- 云存储配置（阿里云、AWS、腾讯云、七牛云、MinIO 等）

使用 Viper 支持热重载。

### 前端配置
- `.env.development` - 开发环境（API 代理到 http://127.0.0.1:8888）
- `.env.production` - 生产环境
- `web/src/core/config.js` - 应用配置
- `web/vite.config.js` - 构建配置

## 代码生成

项目包含强大的代码生成器：

1. **自动代码生成器**（系统工具 → 代码生成器）
   - 生成完整的 CRUD 功能
   - 创建 Model、Service、API、Router 和前端页面
   - 基于模板（在 `server/resource/template/` 中自定义）

2. **表单生成器**（基于 @form-create/designer）
   - 可视化表单设计器
   - 生成 Vue 表单组件

3. **MCP 集成**（AI 辅助开发）
   - MCP 服务器，端口 8888
   - 端点：`/sse`（SSE）、`/message`
   - 服务名称：GVA-MCP v1.0.0
   - 提供 AI 辅助代码生成支持
   - **重要**：在实现功能之前，使用 GVA Helper MCP 获取开发指导

   **MCP 配置详情**（`server/config.yaml`）:
   ```yaml
   mcp:
     name: GVA-MCP           # MCP服务名称
     version: v1.0.0         # 版本号
     sse_path: /sse          # SSE路径
     message_path: /message  # 消息路径
     url_prefix: ""          # URL前缀
     addr: 8888              # 独立MCP服务端口
     separate: false         # 是否独立运行（false表示集成在主服务中）
   ```

   **使用方式**：
   - GVA Helper MCP 提供项目架构、最佳实践、代码生成建议
   - 在开发新功能前，先咨询 MCP 获取正确的实现方案
   - 确保遵循严格的四层架构和 enter.go 模式
   - MCP 访问地址：`http://47.106.76.182:8888/sse`


## 常见开发工作流程

### 添加新的 API 端点

1. **创建模型** (`server/model/system/xxx.go`):
   ```go
   type Xxx struct {
       global.GVA_MODEL
       Name string `json:"name" gorm:"comment:名称"`
   }
   ```

2. **创建请求 DTO** (`server/model/system/request/xxx.go`):
   ```go
   type XxxSearch struct {
       request.PageInfo
       Name string `json:"name" form:"name"`
   }
   ```

3. **创建 Service** (`server/service/system/xxx.go`):
   ```go
   type XxxService struct{}

   func (s *XxxService) CreateXxx(xxx model.Xxx) error {
       return global.GVA_DB.Create(&xxx).Error
   }
   ```

4. **在 `server/service/enter.go` 中注册 Service**

5. **创建 API** (`server/api/v1/system/xxx.go`):
   ```go
   var xxxService = service.ServiceGroupApp.SystemServiceGroup.XxxService

   // CreateXxx 创建新的 Xxx
   // @Tags Xxx
   // @Summary 创建 Xxx
   // [... 完整的 Swagger 注释 ...]
   func (a *XxxApi) CreateXxx(c *gin.Context) {
       var xxx model.Xxx
       if err := c.ShouldBindJSON(&xxx); err != nil {
           response.FailWithMessage(err.Error(), c)
           return
       }
       if err := xxxService.CreateXxx(xxx); err != nil {
           response.FailWithMessage(err.Error(), c)
       } else {
           response.OkWithMessage("创建成功", c)
       }
   }
   ```

6. **在 `server/api/v1/enter.go` 中注册 API**

7. **创建 Router** (`server/router/system/xxx.go`):
   ```go
   func (r *XxxRouter) InitXxxRouter(Router *gin.RouterGroup) {
       xxxRouter := Router.Group("xxx").Use(middleware.OperationRecord())
       xxxApi := v1.ApiGroupApp.SystemApiGroup.XxxApi
       {
           xxxRouter.POST("createXxx", xxxApi.CreateXxx)
       }
   }
   ```

8. **在 `server/router/enter.go` 和 `server/initialize/router.go` 中注册 Router**

9. **生成 Swagger**：运行 `swag init`

10. **创建前端 API** (`web/src/api/xxx.js`):
    ```javascript
    import service from '@/utils/request'

    export const createXxx = (data) => {
      return service({
        url: '/xxx/createXxx',
        method: 'post',
        data
      })
    }
    ```

11. **创建前端页面** (`web/src/view/xxx/xxx.vue`)，使用 Composition API

### 运行测试

```bash
# 后端
cd server
go test ./service/system/...  # 测试特定包
go test ./...                 # 测试所有

# 前端
cd web
npm run test                  # 如果配置了测试脚本
```

### 数据库迁移

数据库表在服务器启动时通过 GORM 的 `AutoMigrate` 自动创建。要添加迁移：

1. 在 `server/initialize/gorm.go` 中添加模型
2. 在 `server/source/` 中添加初始化数据
3. 运行服务器 - 表将自动创建

## 重要提示

1. **务必阅读项目规则** - 在 `.claude/rules/project_rules.md` 中查看详细的架构指导
2. **使用 GVA Helper MCP** - 在实现新功能之前，获取最佳实践和指导
3. **不要跳过 Swagger 文档** - 所有 API 都必须有文档
4. **保持严格的层级分离** - 不允许跨层调用
5. **使用 enter.go 模式** - 所有模块注册都使用此模式
6. **确保数据类型一致性** - Model、Request 和 Response 结构体中的字段类型必须一致
7. **遵循 Composition API** - 所有 Vue 组件都使用 Composition API
8. **使用统一的 request.js** - 前端所有 HTTP 调用都通过它

## 故障排除

- **Swagger 未更新**：在 server 目录下运行 `swag init`
- **前端无法连接后端**：检查 `vite.config.js` 和 `.env.development` 中的代理设置
- **数据库连接错误**：验证 `config.yaml` 中的数据库设置
- **JWT 认证错误**：检查 `x-token` header 中的 token 和中间件配置
- **类型转换错误**：确保 Model 和 Request/Response 结构体中的数据类型一致
