# ---------- 基础阶段 ----------
FROM node:20-slim AS base

# 配置 pnpm 环境变量
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

WORKDIR /app
COPY . .

# ---------- 仅安装生产依赖 ----------
FROM base AS prod-deps
# 缓存 pnpm store，加速重复构建
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm install --registry=https://registry.npmmirror.com/ --prod

# ---------- 构建阶段 ----------
FROM base AS build

# 从 prod-deps 阶段复制生产依赖
COPY --from=prod-deps /app/node_modules /app/node_modules

# 限制 Node 内存使用（3GB）
ENV NODE_OPTIONS="--max-old-space-size=3072"

# 使用缓存加速依赖安装，并执行构建
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm install --registry=https://registry.npmmirror.com/ && pnpm run build

# ---------- 运行阶段 ----------
FROM nginx:alpine AS final

# 拷贝 nginx 配置文件
COPY --from=base /app/.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf

# 拷贝构建后的静态资源
COPY --from=build /app/dist /usr/share/nginx/html

# 验证构建结果
RUN ls -al /usr/share/nginx/html

# 默认启动 Nginx
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

