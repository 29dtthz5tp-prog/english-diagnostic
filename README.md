# 🪞 英语学习镜子 - English Learning Mirror

一个帮助学生诊断英语学习情况、生成个性化学习画像的互动工具。分学生端和教师端，支持实时数据同步。

## 📂 项目结构

```
├── 英语诊断工具-学生端.html        ← 学生诊断入口（初中版）
├── 英语诊断工具-教师端.html        ← 教师数据看板（初中版）
├── 英语诊断工具-初升高-学生端.html  ← 学生诊断入口（初升高版）
├── 英语诊断工具-初升高-教师端.html  ← 教师数据看板（初升高版）
├── 英语诊断工具-设计稿.html        ← UI 设计参考
├── 诊断.html                       ← 早期版本
├── 初中英语三年学习规划.html       ← 学习规划工具
│
├── english-diagnostic/             ← Google Apps Script 备用版本
│   ├── index.html                  ← 学生端
│   ├── student.html                ← 学生端（直接访问）
│   ├── teacher.html                ← 教师端
│   ├── apps-script.js              ← Google Apps Script 后端脚本
│   └── README.md
│
└── gofit/                          ← Go + SQLite 后端版本
    ├── main.go                     ← Go 后端
    ├── static/                     ← 前端静态文件
    │   ├── 学生端.html
    │   ├── 教师端.html
    │   └── 诊断.html
    └── 一键启动.bat
```

## 🚀 快速开始

### 方式一：直接使用（推荐）

1. 用浏览器打开 `英语诊断工具-学生端.html`
2. 填写姓名、年级，完成诊断题目
3. 打开 `英语诊断工具-教师端.html` 查看班级数据

数据通过 **Firebase Firestore** 实时同步，学生提交后教师端秒级更新。

### 方式二：Go 后端版本

进入 `gofit/` 目录，双击 `一键启动.bat`，然后访问 `http://localhost:8080`。

## 🔧 技术栈

| 版本 | 前端 | 后端 | 数据库 |
|------|------|------|--------|
| 主版本 | 原生 HTML/CSS/JS | Firebase Firestore | Firestore |
| 备用版 | 原生 HTML/CSS/JS | Google Apps Script | Google Sheets |
| Go版 | 原生 HTML/CSS/JS | Go + net/http | SQLite |

## 📊 功能特性

- **六维诊断**：词汇、语法、速度、语感、表达、适应能力
- **学员分类**：全能型、输入强势型、输出强势型、词汇/语法/速度等专项薄弱型
- **热力图**：班级各维度一目了然
- **历史对比**：同一学生多次诊断对比
- **教学建议**：根据班级短板自动生成教学决策建议
- **CSV 导出**：一键导出班级数据

## ⚙️ Firebase 配置

本项目使用 Firebase Firestore 作为数据存储。配置信息已内置在 HTML 文件中：

- 项目 ID: `gaotu-diagnostic`
- 集合名: `students`（初中）/ `senior_students`（初升高）

如需更换 Firebase 项目，修改 HTML 文件中的 `firebaseConfig` 即可。
