# PostgreSQL 数据库 Schema 设计

## 设计原则

1. **读写分离**：主库写入，从库读取
2. **索引优化**：覆盖 90% 查询场景
3. **JSON 存储**：灵活扩展配置字段
4. **级联删除**：合理配置外键约束

---

## 表结构

### 1. 用户表 (User)

```sql
CREATE TABLE "User" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    nickname VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    slogan TEXT,
    cover_url VARCHAR(500),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    -- S3 配置
    enable_s3 BOOLEAN DEFAULT false,
    s3_domain VARCHAR(500),
    bucket VARCHAR(255),
    region VARCHAR(100),
    access_key VARCHAR(255),
    secret_key VARCHAR(255),
    endpoint VARCHAR(500),
    thumbnail_suffix VARCHAR(100),
    favicon VARCHAR(500),
    title VARCHAR(255) DEFAULT '极简朋友圈',
    beian_no VARCHAR(100),
    css TEXT,
    js TEXT
);

-- 索引
CREATE INDEX idx_user_username ON "User"(username);
CREATE INDEX idx_user_created_at ON "User"(created_at DESC);
```

### 2. 动态表 (Memo)

```sql
CREATE TABLE "Memo" (
    id SERIAL PRIMARY KEY,
    content TEXT,
    imgs TEXT,  -- 逗号分隔，查询时用 string_to_array
    fav_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    user_id INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    music_163_url VARCHAR(500),
    bilibili_url VARCHAR(500),
    location VARCHAR(255),
    external_url VARCHAR(500),
    external_title VARCHAR(500),
    external_favicon VARCHAR(500) DEFAULT '/favicon.png',
    pinned BOOLEAN DEFAULT false,
    ext JSONB DEFAULT '{}',
    show_type INTEGER DEFAULT 1,  -- 1=公开, 2=好友, 3=私密
    tags TEXT  -- 逗号分隔，查询时用 string_to_array
);

-- 🔥 关键索引：覆盖 90% 查询场景
CREATE INDEX idx_memo_user_time ON "Memo"(user_id, created_at DESC);           -- 用户时间线
CREATE INDEX idx_memo_time_public ON "Memo"(created_at DESC) WHERE show_type = 1;  -- 公开时间线
CREATE INDEX idx_memo_pinned ON "Memo"(pinned DESC, created_at DESC);        -- 置顶优先
CREATE INDEX idx_memo_tag ON "Memo" USING gin(string_to_array(tags, ','));    -- 标签搜索
CREATE INDEX idx_memo_content_gin ON "Memo" USING gin(to_tsvector('simple', content)); -- 全文搜索
CREATE INDEX idx_memo_updated ON "Memo"(updated_at DESC);  -- 更新排序
```

### 3. 评论表 (Comment)

```sql
CREATE TABLE "Comment" (
    id SERIAL PRIMARY KEY,
    memo_id INTEGER NOT NULL REFERENCES "Memo"(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_comment_memo ON "Comment"(memo_id, created_at DESC);
CREATE INDEX idx_comment_user ON "Comment"(user_id);
CREATE INDEX idx_comment_created ON "Comment"(created_at DESC);
```

### 4. 系统配置表 (SysConfig)

```sql
CREATE TABLE "SysConfig" (
    id SERIAL PRIMARY KEY,
    content JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 单行配置使用 UPSERT
INSERT INTO "SysConfig" (id, content) VALUES (1, '{}')
ON CONFLICT (id) DO NOTHING;
```

### 5. 好友链接表 (Friend)

```sql
CREATE TABLE "Friend" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL,
    logo VARCHAR(500),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_friend_user ON "Friend"(user_id);
```

### 6. 点赞记录表 (Like) - 可选扩展

```sql
CREATE TABLE "Like" (
    id SERIAL PRIMARY KEY,
    memo_id INTEGER NOT NULL REFERENCES "Memo"(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(memo_id, user_id)  -- 防止重复点赞
);

CREATE INDEX idx_like_memo ON "Like"(memo_id);
CREATE INDEX idx_like_user ON "Like"(user_id);
```

---

## 常见查询优化

### 用户动态列表 (主查询)

```sql
-- 优化前: 全表扫描
SELECT * FROM "Memo" WHERE user_id = 1 ORDER BY created_at DESC LIMIT 20;

-- 优化后: 使用复合索引
SELECT m.*, u.username, u.nickname
FROM "Memo" m
JOIN "User" u ON u.id = m.user_id
WHERE m.user_id = 1
ORDER BY m.pinned DESC, m.created_at DESC
LIMIT 20;
-- 使用索引: idx_memo_user_time
```

### 公开时间线 (首页)

```sql
-- 优化前: 全表过滤
SELECT * FROM "Memo" WHERE show_type = 1 AND created_at <= NOW() ORDER BY created_at DESC;

-- 优化后: 部分索引
SELECT m.*, u.username, u.nickname
FROM "Memo" m
JOIN "User" u ON u.id = m.user_id
WHERE m.show_type = 1
  AND m.created_at <= NOW()
ORDER BY m.pinned DESC, m.created_at DESC
LIMIT 20 OFFSET 0;
-- 使用索引: idx_memo_time_public
```

### 标签搜索

```sql
-- 使用 PostgreSQL 数组
SELECT * FROM "Memo"
WHERE string_to_array(tags, ',') && ARRAY['旅游', '摄影']
ORDER BY created_at DESC
LIMIT 20;
-- 使用索引: idx_memo_tag (GIN 索引)
```

### 全文搜索

```sql
-- 搜索内容包含关键词
SELECT * FROM "Memo"
WHERE to_tsvector('simple', content) @@ to_tsquery('simple', '旅游 & 摄影')
ORDER BY created_at DESC
LIMIT 20;
-- 使用索引: idx_memo_content_gin
```

---

## 慢查询分析

```sql
-- 启用查询统计
ALTER SYSTEM SET track_activities = on;
ALTER SYSTEM SET track_counts = on;
ALTER SYSTEM SET track_io_timing = on;

-- 查看慢查询
SELECT
    query,
    calls,
    mean_exec_time,
    total_exec_time,
    rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;
```

---

## 连接池配置

```sql
-- PostgreSQL 连接池 (PgBouncer)
[databases]
moments = host=postgres-primary port=5432 dbname=moments

[pgbouncer]
pool_mode = transaction
max_client_conn = 200
default_pool_size = 20
min_pool_size = 5
reserve_pool_size = 5
reserve_pool_timeout = 5
```
