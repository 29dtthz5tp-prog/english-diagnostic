# 🪞 英语学习镜子 — 免费英语诊断工具

一套**完全免费**的英语学习诊断工具，学生端做诊断拿画像，教师端看全班数据做决策。

基于 [GitHub Pages](https://pages.github.com) 托管 + [Google Sheets](https://sheets.google.com) 存储，零成本运行。

---

## ✨ 功能

### 学生端
- 19 道诊断题，涵盖情景判断、知识点、自我评估、学习偏好
- 六维能力画像：词汇、语法、速度、语感、表达、适应
- 输入/输出能力对比分析
- 9 种学习类型自动分类
- 专属行动指南（含每日练习建议）
- 支持多次诊断 + 历史对比

### 教师端
- 全班数据总览（诊断次数、学员数、类型分布）
- 🔥 六维热力图 —— 一眼看全班的强项和短板
- 学员列表（可按学段筛选）
- 个人详细画像 + 历史对比
- 班级教学决策建议
- CSV 导出

---

## 🚀 部署指南（约 15 分钟）

### 第一步：复制代码到你的 GitHub

1. 点击本仓库右上角的 **Fork**，复制到你自己的 GitHub 账号下
2. 或者直接 **Download ZIP**，然后创建你自己的仓库上传

### 第二步：创建 Google Sheet

1. 打开 [Google Sheets](https://sheets.google.com)，新建一个空白表格
2. 重命名第一个工作表为 `diagnoses`
3. 在第一行填入以下表头：

```
id | name | grade | date | vocab | grammar | speed | intuition | express | adapt | inputScore | outputScore | gap | typeId | typeName | typeIcon | weakKey | strongKey | cogData | stylePrefs | timestamp
```

> 💡 复制上面的表头，在 Sheet 的第一行粘贴（使用「选择性粘贴 → 以纯文本格式粘贴」）

### 第三步：部署 Google Apps Script

1. 在 Google Sheet 中点击 **扩展程序 → Apps Script**
2. 删除编辑器中的默认代码
3. 打开本仓库的 `apps-script.js` 文件，**复制全部内容**
4. 粘贴到 Apps Script 编辑器中
5. 点击顶部工具栏，**运行一次 `setup` 函数**：
   - 选择函数下拉框 → 选 `setup` → 点击 ▶️ 运行
   - 首次运行会要求授权，点击「审核权限」→ 选择你的 Google 账号 → 允许
   - 运行成功后，**记下日志中显示的教师密钥**（类似 `class-a3b5f2x1`）
6. 点击右上角 **部署 → 新建部署**
   - 类型：**Web 应用**
   - 执行身份：**我**
   - 访问权限：**任何人**
   - 点击「部署」
7. **复制生成的 Web 应用 URL**（类似 `https://script.google.com/macros/s/xxx/exec`）

### 第四步：配置 HTML 文件

1. 打开仓库中的 `student.html`
2. 找到第 1 行附近的 `APPS_SCRIPT_URL`，替换为你的 Apps Script URL：
   ```javascript
   const APPS_SCRIPT_URL = 'https://script.google.com/macros/s/你的ID/exec';
   ```
3. 打开 `teacher.html`，同样替换 `APPS_SCRIPT_URL` 和 `TEACHER_KEY`
4. 提交（commit）这些修改到 GitHub

### 第五步：启用 GitHub Pages

1. 进入你的 GitHub 仓库 → **Settings → Pages**
2. Source: **main** 分支，文件夹选 **/ (root)**
3. 点击 **Save**
4. 等待 1-2 分钟，你的网站就上线了！
5. 访问地址：`https://你的用户名.github.io/仓库名/`

### 第六步：分享给学生

把网址发给学生和家长：
```
https://你的用户名.github.io/仓库名/
```

学生在自己手机上打开 → 点「学生入口」→ 填写诊断 → 数据自动汇总到你的 Google Sheet 和教师面板。

---

## 🔑 教师访问

教师入口需要密钥保护：

1. 打开教师面板页面
2. 首次访问会弹出输入框
3. 输入你在第三步设置的教师密钥
4. 之后可以在页面上看到全部学生的数据

> 💡 忘记密钥？在 Google Sheet 中点击 **🪞 诊断工具 → 🔑 查看/修改教师密钥**

---

## 📊 数据说明

| 项目 | 详情 |
|------|------|
| 存储位置 | 你的 Google Sheets → `diagnoses` 工作表 |
| 数据归属 | 完全属于你，存储在你自己 Google 账号下 |
| 隐私 | 学生只需提供名字（可用昵称），不收集任何个人信息 |
| 备份 | 建议定期从 Google Sheets 导出备份 |
| 费用 | **完全免费**，无任何收费环节 |

---

## 🛠 常见问题

**Q: 学生提交后教师面板看不到数据？**
A: 等 15 秒自动刷新，或手动点「🔄 刷新」按钮。检查教师密钥是否正确。

**Q: Apps Script 报错？**
A: 确认 `diagnoses` 工作表存在，且表头与第二步一致。重新运行一次 `setup` 函数。

**Q: 能自定义题目吗？**
A: 可以！编辑 `student.html` 中 `questions` 数组（约第 306 行），修改题目和选项即可。

**Q: 可以把页面嵌入微信公众号吗？**
A: GitHub Pages 的域名在微信中可能会被拦截。可以考虑绑定自定义域名（Settings → Pages → Custom domain）。

**Q: 数据安全吗？**
A: 数据存在你自己的 Google 账号下，只有你能访问。教师面板需要密钥才能查看。

---

## 📝 修改记录

- **v2.0** — 从 Firebase 迁移到 Google Sheets，完全免费
- **v1.0** — 原始版本（Firebase + Go 后端）
