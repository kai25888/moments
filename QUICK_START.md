# Moments 项目快速开始指南

## 🚀 团队技术提升快速启动

本文档帮助团队快速开始使用已建立的技术规范。

---

## 📚 必读文档

| 文档 | 说明 | 优先级 |
|-----|------|--------|
| [TECH_IMPROVEMENT_PLAN.md](./TECH_IMPROVEMENT_PLAN.md) | 技术提升整体方案 | ⭐⭐⭐ |
| [CODING_STANDARDS.md](./CODING_STANDARDS.md) | 代码规范详细说明 | ⭐⭐⭐ |
| [CODE_REVIEW_CHECKLIST.md](./CODE_REVIEW_CHECKLIST.md) | 代码审查清单 | ⭐⭐ |

---

## ⚡ 立即行动项

### 1. 修复关键安全问题（今天完成）

#### SQL 注入修复
**文件**: `backend/handler/memo.go:140`

```go
// ❌ 修改前
tx = tx.Where("content like ?", "%"+req.ContentContains+"%")

// ✅ 修改后
if len(req.ContentContains) > 100 {
    return FailRespWithMsg(c, ParamError, "搜索内容过长")
}
tx = tx.Where("content LIKE ?", "%"+req.ContentContains+"%")
```

#### 错误处理完善
**文件**: `backend/handler/memo.go:333`

```go
// ❌ 修改前
if err = m.base.db.First(&memo, req.ID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
    return FailResp(c, ParamError)
}

// ✅ 修改后
if err = m.base.db.First(&memo, req.ID).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return FailRespWithMsg(c, NotFound, "memo不存在")
    }
    m.base.log.Error().Err(err).Int("memo_id", req.ID).Msg("查询memo失败")
    return FailResp(c, InternalError)
}
```

### 2. 添加 JSON 解析错误处理（今天完成）

**文件**: `front/components/Memo.vue:328`

```typescript
// ❌ 修改前
const extJSON = computed(() => {
  return JSON.parse(props.memo.ext || "{}") as ExtDTO;
});

// ✅ 修改后
const extJSON = computed<ExtDTO>(() => {
  try {
    return JSON.parse(props.memo.ext || "{}") as ExtDTO;
  } catch (error) {
    console.error("Failed to parse ext JSON:", error);
    return {} as ExtDTO;
  }
});
```

---

## 📁 重构示例参考

已提供完整的重构示例，位于 `examples/refactored/` 目录：

```
examples/refactored/
├── memo_service.go          # 分层架构示例
├── MemoCard.vue             # 组件重构示例
├── composables/
│   └── useMemo.ts          # 组合式函数示例
└── utils/
    └── validation.ts       # 验证工具示例
```

### 如何使用示例

1. **学习架构模式**：查看 `memo_service.go` 了解分层架构
2. **组件重构参考**：参考 `MemoCard.vue` 学习组件设计
3. **状态管理**：参考 `useMemo.ts` 学习组合式函数

---

## 🛠️ 开发工具配置

### VSCode 推荐配置

创建 `.vscode/settings.json`：

```json
{
  // Go 配置
  "go.formatTool": "gofmt",
  "go.lintTool": "golangci-lint",
  "go.toolsManagement.autoUpdate": true,
  
  // Vue/TypeScript 配置
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "typescript.preferences.importModuleSpecifier": "relative",
  
  // 文件关联
  "files.associations": {
    "*.vue": "vue"
  }
}
```

### 推荐插件

**Go 开发**
- Go (官方插件)
- Go Test Explorer
- Go Outliner

**Vue/前端**
- Volar (Vue 3 官方)
- TypeScript Vue Plugin
- ESLint
- Prettier

**通用**
- GitLens
- Error Lens
- Todo Tree

---

## 🔄 Git 工作流

### 分支命名规范

```bash
# 功能分支
feature/memo-image-upload
feature/user-profile-edit

# 修复分支
fix/sql-injection-memo-list
fix/memory-leak-comment

# 重构分支
refactor/extract-memo-service
refactor/simplify-auth-middleware

# 文档分支
docs/api-documentation
docs/code-style-guide
```

### 提交信息规范

```bash
# 格式: type(scope): subject

feat(memo): 添加图片上传功能
fix(auth): 修复 token 过期问题
docs(api): 更新用户接口文档
refactor(service): 提取 memo 业务逻辑
test(memo): 添加单元测试
style(frontend): 统一代码格式
chore(deps): 更新依赖版本
```

### 完整工作流示例

```bash
# 1. 创建功能分支
git checkout -b feature/memo-search

# 2. 开发并提交
git add .
git commit -m "feat(memo): 添加内容搜索功能

- 支持按关键词搜索 memo
- 添加搜索结果高亮
- 限制搜索长度防止滥用

Closes #123"

# 3. 推送分支
git push origin feature/memo-search

# 4. 创建 PR，等待审查

# 5. 审查通过后合并
```

---

## 🧪 本地开发环境

### 后端启动

```bash
cd moments-fork/backend

# 安装依赖
go mod download

# 启动开发服务器
go run main.go

# 或使用 air 热重载
air
```

### 前端启动

```bash
cd moments-fork/front

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 类型检查
pnpm type-check

# 代码检查
pnpm lint
```

### Docker 启动（完整环境）

```bash
cd moments-fork

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

## 📝 代码审查流程

### 提交者检查清单

提交 PR 前自检：

```markdown
- [ ] 代码可以编译/运行
- [ ] 本地测试通过
- [ ] 添加了必要的测试
- [ ] 代码符合规范
- [ ] 无敏感信息泄露
- [ ] Commit 信息规范
```

### 审查者检查清单

审查时关注：

```markdown
## 安全
- [ ] 输入验证完整
- [ ] 无 SQL 注入风险
- [ ] 权限检查正确

## 质量
- [ ] 错误处理完整
- [ ] 代码结构清晰
- [ ] 命名规范

## 测试
- [ ] 测试覆盖充分
- [ ] 边界条件有测试
```

---

## 🎓 学习资源

### Go 进阶

**必读**
- [Effective Go](https://go.dev/doc/effective_go) - Go 官方最佳实践
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) - Uber 代码规范
- [Go Clean Architecture](https://github.com/bxcodec/go-clean-arch) - 干净架构示例

**推荐书籍**
- 《Go 语言高级编程》
- 《Go 语言设计与实现》

### Vue/Nuxt 进阶

**必读**
- [Vue 3 组合式 API](https://vuejs.org/guide/extras/composition-api-faq.html)
- [Nuxt 3 文档](https://nuxt.com/docs)
- [Vue Style Guide](https://vuejs.org/style-guide/)

**推荐**
- [VueUse](https://vueuse.org/) - 常用组合式函数库
- [Pinia](https://pinia.vuejs.org/) - 状态管理

### 工程化

**必读**
- 《重构：改善既有代码的设计》
- 《代码整洁之道》
- 《领域驱动设计》

---

## 🎯 本周任务

### 第一天：规范学习
- [ ] 阅读 CODING_STANDARDS.md
- [ ] 配置开发环境
- [ ] 修复关键安全问题

### 第二天：代码审查
- [ ] 学习 CODE_REVIEW_CHECKLIST.md
- [ ] 参与一次代码审查
- [ ] 提交第一个符合规范的 PR

### 第三天：架构学习
- [ ] 学习 examples/refactored/ 示例
- [ ] 理解分层架构
- [ ] 尝试小范围重构

### 第四天：测试实践
- [ ] 学习测试规范
- [ ] 为现有代码添加测试
- [ ] 运行测试覆盖率检查

### 第五天：总结分享
- [ ] 整理本周学习心得
- [ ] 团队技术分享
- [ ] 制定下周改进计划

---

## 🆘 常见问题

### Q: 如何开始重构现有代码？
**A**: 建议按以下顺序：
1. 先修复安全问题和 Bug
2. 添加单元测试（确保重构安全）
3. 小步重构，每次只改一个函数/组件
4. 每次重构后运行测试

### Q: 如何处理遗留代码？
**A**: 
1. 使用 `// TODO:` 标记需要改进的地方
2. 创建技术债务 Issue 跟踪
3. 新功能按新规范编写
4. 逐步重构旧代码

### Q: 代码审查意见冲突怎么办？
**A**:
1. 保持开放心态，理解对方观点
2. 引用规范文档作为依据
3. 必要时团队讨论决定
4. 记录决策，更新规范

---

## 📞 获取帮助

遇到问题？

1. **查看文档**：先查阅相关规范文档
2. **参考示例**：查看 examples/ 目录
3. **团队讨论**：在团队群讨论
4. **技术分享**：申请专题技术分享

---

*快速开始指南 v1.0*
*祝编码愉快！*
