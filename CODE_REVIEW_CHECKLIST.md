# Moments 项目代码审查清单

> 本清单用于指导代码审查过程，确保代码质量和一致性。

---

## 📋 审查前准备

### 审查者准备
- [ ] 理解本次 PR 的业务目标和改动范围
- [ ] 阅读相关的需求文档或 Issue
- [ ] 确认本地可以运行并测试

### 提交者准备
- [ ] 代码可以编译/运行
- [ ] 自测通过，无明显 Bug
- [ ] 添加了必要的单元测试
- [ ] 更新了相关文档
- [ ] Commit 信息规范

---

## 🔴 安全审查（必须项）

### 输入验证
- [ ] 所有用户输入都经过验证
- [ ] 字符串长度有合理限制
- [ ] 数值范围有边界检查
- [ ] 文件上传有类型和大小限制

### SQL 安全
- [ ] 使用参数化查询，无字符串拼接 SQL
- [ ] ORM 查询条件无注入风险
- [ ] 敏感操作有权限验证

### 认证授权
- [ ] 敏感接口需要登录
- [ ] 权限检查完整（不能仅依赖前端）
- [ ] Token 有过期机制
- [ ] 密码使用安全算法存储（bcrypt/argon2）

### 敏感信息
- [ ] 无硬编码密钥、密码
- [ ] 日志中不包含敏感信息
- [ ] 错误信息不暴露内部细节

---

## 🟠 代码质量审查

### Go 代码

#### 错误处理
- [ ] 所有错误都被处理，不忽略
- [ ] 错误信息清晰，有帮助
- [ ] 错误日志包含足够上下文
- [ ] 使用 `fmt.Errorf` 包装错误时添加 `%w`

#### 代码结构
- [ ] 函数职责单一
- [ ] 函数长度合理（< 50 行）
- [ ] 参数数量合理（< 5 个）
- [ ] 无重复代码（DRY 原则）

#### 命名规范
- [ ] 变量名有意义，避免单字母（循环除外）
- [ ] 接口名符合 Go 惯例（-er 后缀或描述性）
- [ ] 常量使用驼峰命名
- [ ] 包名简短、有意义

#### 并发安全
- [ ] 共享数据有适当的锁保护
- [ ] Channel 使用正确，避免死锁
- [ ] Context 正确传递
- [ ] Goroutine 有退出机制

#### 性能
- [ ] 避免不必要的内存分配
- [ ] 数据库查询有索引
- [ ] 避免 N+1 查询
- [ ] 外部调用有超时设置

### Vue/TypeScript 代码

#### 类型安全
- [ ] 无 `any` 类型（特殊情况需注释说明）
- [ ] Props 类型定义完整
- [ ] 返回值类型明确
- [ ] 使用严格模式

#### 组件设计
- [ ] 组件职责单一
- [ ] Props 向下传递，Events 向上传递
- [ ] 避免过度使用 ref
- [ ] 合理使用 computed 缓存

#### 代码组织
- [ ] 逻辑复用使用 composables
- [ ] 复杂逻辑抽离为独立函数
- [ ] 魔法数字提取为常量
- [ ] 样式使用 scoped

#### 性能优化
- [ ] 大列表使用虚拟滚动
- [ ] 图片懒加载
- [ ] 避免不必要的重渲染
- [ ] 合理使用 watch

---

## 🟡 测试审查

### 测试覆盖
- [ ] 核心逻辑有单元测试
- [ ] 边界条件有测试用例
- [ ] 错误路径有测试
- [ ] 测试覆盖率达标（> 60%）

### 测试质量
- [ ] 测试用例独立，无依赖
- [ ] 测试数据清晰，可维护
- [ ] 测试名称描述清楚
- [ ] 使用表驱动测试（Go）

### 测试示例

**Go 测试**
```go
// ✅ 好的测试
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateUserRequest
        wantErr bool
    }{
        {
            name: "valid user",
            input: &CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
        {
            name: "empty username",
            input: &CreateUserRequest{
                Username: "",
                Email:    "test@example.com",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

**Vue 测试**
```typescript
// ✅ 好的测试
describe('UserCard', () => {
  it('renders user nickname when provided', () => {
    const user = { id: 1, username: 'test', nickname: 'Test User' }
    const wrapper = mount(UserCard, { props: { user } })
    expect(wrapper.text()).toContain('Test User')
  })
  
  it('renders username when nickname is empty', () => {
    const user = { id: 1, username: 'testuser', nickname: '' }
    const wrapper = mount(UserCard, { props: { user } })
    expect(wrapper.text()).toContain('testuser')
  })
})
```

---

## 🟢 文档和注释

### 代码注释
- [ ] 复杂逻辑有注释说明
- [ ] 公共函数有文档注释
- [ ] 接口和类型有说明
- [ ] TODO 有对应的 Issue 链接

### 文档更新
- [ ] API 文档已更新
- [ ] README 需要时更新
- [ ] 变更日志已记录

### 注释规范

**Go 文档注释**
```go
// ✅ 好的文档注释
// CreateUser 创建新用户
// 
// 参数:
//   - ctx: 上下文
//   - req: 创建请求，包含用户名、邮箱、密码
// 
// 返回:
//   - *User: 创建成功的用户信息
//   - error: 可能的错误：ErrInvalidInput, ErrUserExists
// 
// 示例:
//   user, err := svc.CreateUser(ctx, &CreateUserRequest{
//       Username: "john",
//       Email:    "john@example.com",
//       Password: "secure123",
//   })
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 实现
}
```

---

## 🔵 风格和格式

### 代码格式
- [ ] 代码格式化（gofmt/prettier）
- [ ] 无 Lint 错误
- [ ] 导入顺序正确
- [ ] 文件末尾有空行

### Git 规范
- [ ] Commit 信息规范
- [ ] 提交粒度合理
- [ ] 无不必要的文件变更
- [ ] 分支已 rebase 到最新 main

---

## 🟣 业务逻辑审查

### 功能正确性
- [ ] 实现符合需求
- [ ] 边界条件处理正确
- [ ] 异常流程处理完善
- [ ] 状态转换正确

### 用户体验
- [ ] 错误提示友好
- [ ] 加载状态有反馈
- [ ] 操作有确认提示（危险操作）
- [ ] 响应时间可接受

---

## 📊 审查反馈模板

### 批准（Approve）
```
✅ LGTM (Looks Good To Me)

代码质量良好，符合规范，可以合并。
```

### 有条件批准（Approve with comments）
```
⚠️ 基本 OK，有一些小建议：

1. 建议将 XXX 提取为常量
2. 可以优化一下这里的错误处理
3. 考虑添加一个边界条件的测试

这些问题可以在后续 PR 中处理。
```

### 需要修改（Request changes）
```
🔴 需要修改：

**必须修复：**
1. [安全] SQL 注入风险，请使用参数化查询
2. [功能] 权限检查不完整，需要验证用户是否有操作权限

**建议优化：**
1. 这个函数可以拆分为更小的函数
2. 建议添加单元测试
3. 变量命名可以更清晰一些

请修复后重新提交审查。
```

---

## 🎯 审查优先级

### P0 - 阻塞问题（必须修复）
- 安全漏洞
- 功能错误
- 性能严重问题
- 数据丢失风险

### P1 - 重要问题（建议修复）
- 代码重复
- 错误处理不完整
- 测试覆盖不足
- 命名不清晰

### P2 - 建议优化（可选）
- 代码风格
- 注释完善
- 小重构建议
- 性能微优化

---

## 💡 审查技巧

### 审查者
1. **先理解后审查**：先了解改动背景和目的
2. **分批审查**：大 PR 分多次审查，每次专注一个方面
3. **使用工具**：利用 IDE 和 lint 工具辅助
4. **保持尊重**：提出建设性意见，避免人身攻击
5. **及时响应**：尽快完成审查，不阻塞他人

### 提交者
1. **小步提交**：PR 粒度适中，便于审查
2. **自测充分**：提交前充分自测
3. **描述清晰**：PR 描述清楚改动内容和原因
4. **积极沟通**：对审查意见及时响应
5. **保持开放**：接受合理的改进建议

---

## 📈 审查度量

### 团队指标
- 平均审查时间
- 审查发现问题数
- 返工率
- 审查参与度

### 个人成长
- 代码质量提升
- 审查技能提高
- 知识分享次数

---

*文档版本: v1.0*
*最后更新: 2026-04-14*
