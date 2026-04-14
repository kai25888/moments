# Moments 项目代码规范

## 🎯 总则

本文档定义了 Moments 项目的代码编写规范，所有开发者必须遵守。

---

## 📝 Go 代码规范

### 1. 项目结构

```
backend/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go
├── internal/               # 私有代码
│   ├── domain/            # 领域模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   ├── handler/           # HTTP处理层
│   ├── middleware/        # 中间件
│   └── pkg/               # 内部工具包
├── pkg/                    # 公共库（可被外部使用）
├── configs/                # 配置文件
├── scripts/                # 脚本
└── tests/                  # 测试文件
```

### 2. 命名规范

#### 文件命名
- 全部小写，使用下划线分隔
- 测试文件以 `_test.go` 结尾
- 示例: `user_handler.go`, `memo_service_test.go`

#### 变量命名
```go
// ✅ 正确
var userCount int
var isEnabled bool
var HTTPClient *http.Client

// ❌ 错误
var user_count int  // 不使用下划线
var UserCount int   // 非导出变量不使用大写开头
```

#### 常量命名
```go
// ✅ 正确
const (
    MaxRetryCount = 3
    DefaultTimeout = 30 * time.Second
)

// 枚举类型
const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)
```

#### 接口命名
```go
// ✅ 正确
// 单一方法接口：方法名 + er
 type Reader interface {
     Read(p []byte) (n int, err error)
 }

// 多方法接口：描述性名称
type UserRepository interface {
    GetByID(id int64) (*User, error)
    Create(user *User) error
    Update(user *User) error
}
```

### 3. 错误处理

#### 错误创建
```go
// ✅ 使用 errors.New 或 fmt.Errorf
var ErrUserNotFound = errors.New("user not found")

// 包装错误时添加上下文
if err != nil {
    return fmt.Errorf("failed to get user by id %d: %w", id, err)
}
```

#### 错误处理规范
```go
// ✅ 正确示例
func (h *UserHandler) GetUser(c echo.Context) error {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        h.log.Warn().Err(err).Str("param", c.Param("id")).Msg("invalid user id")
        return FailResp(c, ParamError)
    }

    user, err := h.userService.GetByID(c.Request().Context(), id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return FailRespWithMsg(c, NotFound, "用户不存在")
        }
        h.log.Error().Err(err).Int64("user_id", id).Msg("failed to get user")
        return FailResp(c, InternalError)
    }

    return SuccessResp(c, user)
}
```

#### 错误处理禁止项
```go
// ❌ 不要忽略错误
_ = db.Create(&user)

// ❌ 不要直接暴露内部错误
return err

// ❌ 不要使用 panic
if err != nil {
    panic(err)
}
```

### 4. 日志规范

#### 日志级别使用
```go
// Debug - 开发调试信息
log.Debug().Str("query", sql).Msg("executing query")

// Info - 重要业务事件
log.Info().Int64("user_id", user.ID).Msg("user logged in")

// Warn - 警告但不影响功能
log.Warn().Str("ip", c.RealIP()).Msg("suspicious request detected")

// Error - 错误需要处理
log.Error().Err(err).Msg("failed to process payment")

// Fatal - 致命错误，程序退出
log.Fatal().Err(err).Msg("failed to connect to database")
```

#### 日志字段规范
```go
// ✅ 使用结构化字段
log.Info().
    Str("request_id", requestID).
    Int64("user_id", userID).
    Str("action", "create_memo").
    Dur("duration", time.Since(start)).
    Msg("memo created")

// ❌ 不要拼接字符串
log.Info().Msgf("user %d created memo in %v", userID, duration)
```

### 5. 数据库操作

#### GORM 最佳实践
```go
// ✅ 使用 Context
tx := db.WithContext(ctx).Model(&User{})

// ✅ 分页查询
func (r *UserRepo) List(ctx context.Context, page, size int) ([]*User, error) {
    var users []*User
    offset := (page - 1) * size
    
    err := r.db.WithContext(ctx).
        Order("created_at DESC").
        Limit(size).
        Offset(offset).
        Find(&users).Error
    
    return users, err
}

// ✅ 事务处理
func (r *UserRepo) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&Account{}).Where("user_id = ?", fromID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }
        if err := tx.Model(&Account{}).Where("user_id = ?", toID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }
        return nil
    })
}
```

#### SQL 安全
```go
// ✅ 使用参数化查询
tx.Where("username = ?", username)

// ❌ 禁止字符串拼接
tx.Where(fmt.Sprintf("username = '%s'", username))
```

### 6. HTTP Handler

#### Handler 结构
```go
type UserHandler struct {
    userService UserService
    log         zerolog.Logger
    validator   *validator.Validate
}

func NewUserHandler(svc UserService, log zerolog.Logger) *UserHandler {
    return &UserHandler{
        userService: svc,
        log:         log,
        validator:   validator.New(),
    }
}

// Handler 方法
func (h *UserHandler) Create(c echo.Context) error {
    // 1. 参数绑定和验证
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return FailResp(c, ParamError)
    }
    
    if err := h.validator.Struct(req); err != nil {
        return FailRespWithMsg(c, ParamError, err.Error())
    }
    
    // 2. 调用服务层
    user, err := h.userService.Create(c.Request().Context(), &req)
    if err != nil {
        h.log.Error().Err(err).Msg("failed to create user")
        return FailResp(c, InternalError)
    }
    
    // 3. 返回响应
    return SuccessResp(c, user)
}
```

---

## 🎨 Vue/TypeScript 代码规范

### 1. 项目结构

```
front/
├── components/            # 组件
│   ├── common/           # 通用组件
│   ├── business/         # 业务组件
│   └── layout/           # 布局组件
├── composables/          # 组合式函数
├── pages/                # 页面
├── stores/               # 状态管理
├── types/                # 类型定义
├── utils/                # 工具函数
├── api/                  # API 接口
├── assets/               # 静态资源
└── middleware/           # 中间件
```

### 2. 组件规范

#### 文件命名
```
// ✅ 正确
UserProfile.vue
MemoCard.vue
CommentList.vue

// ❌ 错误
user-profile.vue  // 使用大驼峰
userProfile.vue   // 不使用小驼峰
```

#### 组件结构
```vue
<script setup lang="ts">
// 1. 类型导入
import type { User } from '~/types'

// 2. Vue API 导入
import { computed, ref, watch } from 'vue'

// 3. 第三方库
import { useRoute } from 'vue-router'

// 4. 本地模块
import { useUserStore } from '~/stores/user'
import UserAvatar from '~/components/UserAvatar.vue'

// 5. Props 定义
interface Props {
  user: User
  showActions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true
})

// 6. Emits 定义
const emit = defineEmits<{
  (e: 'update', user: User): void
  (e: 'delete', id: number): void
}>()

// 7. 状态定义
const isLoading = ref(false)
const userStore = useUserStore()

// 8. 计算属性
const displayName = computed(() => {
  return props.user.nickname || props.user.username
})

// 9. 方法
async function handleUpdate() {
  isLoading.value = true
  try {
    await userStore.update(props.user)
    emit('update', props.user)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="user-card">
    <UserAvatar :src="user.avatarUrl" />
    <span class="name">{{ displayName }}</span>
    <button v-if="showActions" @click="handleUpdate">更新</button>
  </div>
</template>

<style scoped>
.user-card {
  @apply flex items-center gap-2 p-4;
}
</style>
```

### 3. TypeScript 规范

#### 类型定义
```typescript
// ✅ 使用 interface 定义对象类型
interface User {
  id: number
  username: string
  email: string
  createdAt: Date
}

// ✅ 使用 type 定义联合类型
type Status = 'pending' | 'active' | 'inactive'
type ID = string | number

// ✅ 泛型使用
interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

// ✅ 函数类型
interface UserService {
  getById(id: number): Promise<User | null>
  create(user: CreateUserDTO): Promise<User>
}
```

#### 禁止使用 any
```typescript
// ❌ 错误
function process(data: any): any {
  return data.value
}

// ✅ 正确
function process<T extends { value: unknown }>(data: T): T['value'] {
  return data.value
}

// 或者使用 unknown
function process(data: unknown): void {
  if (typeof data === 'string') {
    // 类型收窄后使用
    console.log(data.toUpperCase())
  }
}
```

### 4. Composables 规范

```typescript
// composables/useUser.ts
import { ref, computed } from 'vue'
import type { User } from '~/types'

export function useUser(userId: number) {
  // 状态
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)
  
  // 计算属性
  const isLoggedIn = computed(() => user.value !== null)
  
  // 方法
  async function fetchUser() {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await $fetch<User>(`/api/users/${userId}`)
      user.value = response
    } catch (e) {
      error.value = e as Error
    } finally {
      isLoading.value = false
    }
  }
  
  // 返回
  return {
    user: readonly(user),
    isLoading: readonly(isLoading),
    error: readonly(error),
    isLoggedIn,
    fetchUser
  }
}
```

### 5. 状态管理

```typescript
// stores/user.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '~/types'

export const useUserStore = defineStore('user', () => {
  // State
  const currentUser = ref<User | null>(null)
  const users = ref<User[]>([])
  
  // Getters
  const isLoggedIn = computed(() => currentUser.value !== null)
  const userCount = computed(() => users.value.length)
  
  // Actions
  async function login(credentials: LoginDTO) {
    const user = await $fetch<User>('/api/login', {
      method: 'POST',
      body: credentials
    })
    currentUser.value = user
    return user
  }
  
  function logout() {
    currentUser.value = null
  }
  
  return {
    currentUser,
    users,
    isLoggedIn,
    userCount,
    login,
    logout
  }
})
```

### 6. API 层

```typescript
// api/user.ts
import type { User, CreateUserDTO, UpdateUserDTO } from '~/types'

export const userApi = {
  async getById(id: number): Promise<User> {
    return useMyFetch<User>(`/users/${id}`)
  },
  
  async create(data: CreateUserDTO): Promise<User> {
    return useMyFetch<User>('/users', {
      method: 'POST',
      body: data
    })
  },
  
  async update(id: number, data: UpdateUserDTO): Promise<User> {
    return useMyFetch<User>(`/users/${id}`, {
      method: 'PUT',
      body: data
    })
  },
  
  async delete(id: number): Promise<void> {
    return useMyFetch<void>(`/users/${id}`, {
      method: 'DELETE'
    })
  }
}
```

---

## 🧪 测试规范

### Go 测试

```go
// user_service_test.go
package service

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    svc := NewUserService(mockRepo)
    
    ctx := context.Background()
    req := &CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*User")).
        Return(nil)
    
    // Act
    user, err := svc.Create(ctx, req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Username, user.Username)
    mockRepo.AssertExpectations(t)
}
```

### Vue 组件测试

```typescript
// UserCard.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import UserCard from './UserCard.vue'

describe('UserCard', () => {
  it('renders user name correctly', () => {
    const user = {
      id: 1,
      username: 'testuser',
      nickname: 'Test User'
    }
    
    const wrapper = mount(UserCard, {
      props: { user }
    })
    
    expect(wrapper.text()).toContain('Test User')
  })
  
  it('emits update event when button clicked', async () => {
    const wrapper = mount(UserCard, {
      props: { user: { id: 1, username: 'test' } }
    })
    
    await wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted()).toHaveProperty('update')
  })
})
```

---

## 📝 Git 提交规范

### 提交信息格式
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型
- `feat`: 新功能
- `fix`: 修复
- `docs`: 文档
- `style`: 格式（不影响代码运行）
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建过程或辅助工具的变动

### 示例
```
feat(memo): 添加图片上传功能

- 支持多图片同时上传
- 添加图片压缩
- 生成缩略图

Closes #123
```

---

## 🔍 代码审查检查清单

### 提交前自检
- [ ] 代码可以编译/运行
- [ ] 所有测试通过
- [ ] 无 lint 错误
- [ ] 代码符合规范
- [ ] 已添加必要的注释
- [ ] 敏感信息已移除

### 审查者检查清单
- [ ] 代码逻辑正确
- [ ] 错误处理完整
- [ ] 性能考虑周全
- [ ] 安全风险已排除
- [ ] 测试覆盖充分
- [ ] 命名清晰合理

---

*文档版本: v1.0*
*最后更新: 2026-04-14*
