# 🪞 英语学习镜子 — English Learning Mirror

一套**完全免费**的英语学习诊断工具，学生端做诊断拿画像，教师端看全班数据做决策。

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
- 数据实时同步（Firebase 版本）

---

## 📂 项目结构

```
├── 英语诊断工具-学生端.html        ← 学生诊断入口（初中版 · Firebase）
├── 英语诊断工具-教师端.html        ← 教师数据看板（初中版 · Firebase）
├── 英语诊断工具-初升高-学生端.html  ← 学生诊断入口（初升高版 · Firebase）
├── 英语诊断工具-初升高-教师端.html  ← 教师数据看板（初升高版 · Firebase）
├── 英语诊断工具-设计稿.html        ← UI 设计参考
├── 初中英语三年学习规划.html       ← 学习规划工具
│
├── english-diagnostic/             ← Google Sheets 备用版本
│   ├── student.html / teacher.html
│   └── apps-script.js
│
└── gofit/                          ← Go + SQLite 后端版本
    ├── main.go
    └── static/
```

---

## 🚀 快速开始（Firebase 版 · 推荐）

本版本使用 **Firebase Firestore** 实时数据库，无需部署后端。

### 直接使用

1. 用浏览器打开 `英语诊断工具-学生端.html`
2. 填写姓名、年级，完成诊断题目
3. 打开 `英语诊断工具-教师端.html` 查看班级数据

数据通过 Firebase 实时同步，学生提交后教师端秒级更新。

### 用自己的 Firebase

如需更换为自己的 Firebase 项目：

1. 去 [Firebase Console](https://console.firebase.google.com) 创建项目
2. 启用 Firestore Database（测试模式）
3. 修改 HTML 文件中的 `firebaseConfig`

---

## 🌐 GitHub Pages 部署（Google Sheets 版）

如果你希望学生通过网址访问（无需下载 HTML 文件），可以部署为网页。

### 第一步：Fork 仓库

点击右上角 **Fork**，复制到你自己的 GitHub 账号下。

### 第二步：创建 Google Sheet

1. 打开 [Google Sheets](https://sheets.google.com)，新建表格
2. 重命名工作表为 `diagnoses`
3. 第一行填入表头：

```
id | name | grade | date | vocab | grammar | speed | intuition | express | adapt | inputScore | outputScore | gap | typeId | typeName | typeIcon | weakKey | strongKey | cogData | stylePrefs | timestamp
```

### 第三步：部署 Google Apps Script

1. 在 Sheet 中点击 **扩展程序 → Apps Script**
2. 粘贴 `english-diagnostic/apps-script.js` 全部内容
3. 运行一次 `setup` 函数（需要授权），**记下日志中的教师密钥**
4. **部署 → 新建部署 → Web 应用**，执行身份选"我"，访问权限选"任何人"
5. 复制生成的 URL

### 第四步：配置并部署

1. 把 `english-diagnostic/student.html` 和 `teacher.html` 中的 `APPS_SCRIPT_URL` 替换为上一步的 URL
2. 同样替换 `TEACHER_KEY`
3. 把这两个文件复制到仓库根目录（覆盖 Firebase 版本或重命名）
4. GitHub → Settings → Pages → Source: main, / (root) → Save

### 第五步：分享给学生

```
https://你的用户名.github.io/仓库名/
```

---

## 🔧 技术栈

| 版本 | 后端 | 数据库 | 适用场景 |
|------|------|--------|----------|
| 主版本 | Firebase | Firestore | 本地使用，实时同步 |
| 备用版 | Google Apps Script | Google Sheets | GitHub Pages 部署 |
| Go版 | Go 后端 | SQLite | 本地服务器 |

---

## 🛠 常见问题

**Q: 教师端看不到学生数据？**
A: 
- Firebase 版：确保学生端和教师端用的是**同一套 HTML 文件**（都是 Firebase 版）。检查浏览器控制台是否有报错。
- Google Sheets 版：等 15 秒自动刷新，或手动点「🔄 刷新」。检查教师密钥和 APPS_SCRIPT_URL。

**Q: 能自定义题目吗？**
A: 编辑 HTML 文件中 `questions` 数组，修改题目和选项即可。

**Q: 数据安全吗？**
A: Firebase 版数据在你自己 Firebase 项目下。Google Sheets 版数据在你自己的 Google 账号下。都只有你能访问。

---

## 📝 版本记录

- **v2.1** — 修复教师端与学生端后端不一致的问题，统一为 Firebase
- **v2.0** — 新增 Google Sheets 版，支持 GitHub Pages 零成本部署
- **v1.0** — 原始版本（Firebase + Go 后端）
