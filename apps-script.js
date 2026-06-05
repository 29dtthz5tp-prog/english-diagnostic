/**
 * 英语学习镜子 - Google Apps Script 后端
 *
 * 部署步骤：
 * 1. 在 Google Sheets 中：扩展程序 > Apps Script
 * 2. 粘贴本文件全部内容
 * 3. 运行一次 setup() 设置教师密钥
 * 4. 部署 > 新建部署 > Web 应用 > 执行身份：我 / 访问权限：任何人
 * 5. 复制部署 URL，填入 student.html 和 teacher.html 的 APPS_SCRIPT_URL
 */

var SHEET_NAME = 'diagnoses';
var PROP_KEY = 'TEACHER_KEY';

// ==================== 辅助函数 ====================

function getSheet_() {
  var ss = SpreadsheetApp.getActiveSpreadsheet();
  var sheet = ss.getSheetByName(SHEET_NAME);
  if (!sheet) {
    sheet = ss.insertSheet(SHEET_NAME);
    sheet.appendRow([
      'id', 'name', 'grade', 'date',
      'vocab', 'grammar', 'speed', 'intuition', 'express', 'adapt',
      'inputScore', 'outputScore', 'gap',
      'typeId', 'typeName', 'typeIcon',
      'weakKey', 'strongKey',
      'cogData', 'stylePrefs', 'timestamp'
    ]);
    // 冻结表头
    sheet.setFrozenRows(1);
  }
  return sheet;
}

function getTeacherKey_() {
  return PropertiesService.getScriptProperties().getProperty(PROP_KEY) || 'changeme';
}

function rowToObject_(row, rowNum) {
  return {
    id: rowNum,
    name: row[1] || '',
    grade: row[2] || '',
    date: row[3] || '',
    vocab: Number(row[4]) || 0,
    grammar: Number(row[5]) || 0,
    speed: Number(row[6]) || 0,
    intuition: Number(row[7]) || 0,
    express: Number(row[8]) || 0,
    adapt: Number(row[9]) || 0,
    inputScore: Number(row[10]) || 0,
    outputScore: Number(row[11]) || 0,
    gap: Number(row[12]) || 0,
    typeId: row[13] || '',
    typeName: row[14] || '',
    typeIcon: row[15] || '',
    weakKey: row[16] || '',
    strongKey: row[17] || '',
    cogData: row[18] || '[]',
    stylePrefs: row[19] || '[]',
    timestamp: row[20] || ''
  };
}

function json_(obj) {
  return ContentService.createTextOutput(JSON.stringify(obj))
    .setMimeType(ContentService.MimeType.JSON);
}

function error_(msg) {
  return json_({ error: msg });
}

// ==================== 学生提交 (POST) ====================

function doPost(e) {
  try {
    var data = JSON.parse(e.postData.contents);
    var sheet = getSheet_();

    var row = [
      '',                             // A: id (auto = row number)
      data.name || '',                // B
      data.grade || '',               // C
      data.date || '',                // D
      data.vocab || 0,                // E
      data.grammar || 0,              // F
      data.speed || 0,                // G
      data.intuition || 0,            // H
      data.express || 0,              // I
      data.adapt || 0,                // J
      data.inputScore || 0,           // K
      data.outputScore || 0,          // L
      data.gap || 0,                  // M
      data.typeId || '',              // N
      data.typeName || '',            // O
      data.typeIcon || '',            // P
      data.weakKey || '',             // Q
      data.strongKey || '',           // R
      data.cogData || '[]',           // S
      data.stylePrefs || '[]',        // T
      data.timestamp || new Date().toISOString()  // U
    ];

    sheet.appendRow(row);
    var newRowNum = sheet.getLastRow();

    return json_({ success: true, row: newRowNum });
  } catch (err) {
    return error_('提交失败: ' + err.message);
  }
}

// ==================== 教师查询 (GET) ====================

function doGet(e) {
  if (!e || !e.parameter) {
    return error_('缺少参数');
  }

  var action = e.parameter.action;
  var key = e.parameter.key;

  // 需要教师密钥的操作
  var needsAuth = ['list', 'delete', 'clear', 'stats'];
  if (needsAuth.indexOf(action) >= 0) {
    if (key !== getTeacherKey_()) {
      return error_('密钥错误，无访问权限');
    }
  }

  switch (action) {
    case 'list':
      return listStudents_();
    case 'delete':
      return deleteRow_(e.parameter.id);
    case 'clear':
      return clearAll_();
    case 'history':
      return getHistory_(e.parameter.name, e.parameter.grade);
    case 'stats':
      return getStats_();
    default:
      return error_('未知操作: ' + action);
  }
}

// ==================== 具体操作 ====================

function listStudents_() {
  var sheet = getSheet_();
  var data = sheet.getDataRange().getValues();
  var result = [];
  // 跳过表头 (row 0)
  for (var i = 1; i < data.length; i++) {
    result.push(rowToObject_(data[i], i + 1));
  }
  // 最新的排最前
  result.reverse();
  return json_(result);
}

function deleteRow_(id) {
  var sheet = getSheet_();
  var rowNum = parseInt(id, 10);
  if (isNaN(rowNum) || rowNum < 2) {
    return error_('无效的ID');
  }
  sheet.deleteRow(rowNum);
  return json_({ success: true });
}

function clearAll_() {
  var sheet = getSheet_();
  var lastRow = sheet.getLastRow();
  if (lastRow > 1) {
    sheet.deleteRows(2, lastRow - 1);
  }
  return json_({ success: true });
}

function getHistory_(name, grade) {
  var sheet = getSheet_();
  var data = sheet.getDataRange().getValues();
  // 从最新开始查找
  for (var i = data.length - 1; i >= 1; i--) {
    var row = data[i];
    if (row[1] === name && row[2] === grade) {
      return json_(rowToObject_(row, i + 1));
    }
  }
  return json_(null);
}

function getStats_() {
  var sheet = getSheet_();
  var data = sheet.getDataRange().getValues();
  var students = {};
  var totalTests = 0;
  var typeCount = {};

  for (var i = 1; i < data.length; i++) {
    totalTests++;
    var name = data[i][1];
    var grade = data[i][2];
    var typeName = data[i][14];
    var key = name + '|||' + grade;
    students[key] = true;
    typeCount[typeName] = (typeCount[typeName] || 0) + 1;
  }

  return json_({
    totalTests: totalTests,
    totalStudents: Object.keys(students).length,
    typeCount: typeCount
  });
}

// ==================== 初始化 ====================

function setup() {
  // 创建 Sheet 并写入表头
  getSheet_();

  // 生成一个随机教师密钥
  var key = 'class-' + Math.random().toString(36).substring(2, 10);
  PropertiesService.getScriptProperties().setProperty(PROP_KEY, key);

  Logger.log('✅ 设置完成！');
  Logger.log('📋 教师密钥: ' + key);
  Logger.log('⚠ 请复制保存这个密钥，它只会显示这一次');
  Logger.log('📊 数据表: ' + SHEET_NAME);

  // 在 Google Sheets 中创建菜单
  var ui = SpreadsheetApp.getUi();
  ui.createMenu('🪞 诊断工具')
    .addItem('🔑 查看/修改教师密钥', 'showKey')
    .addItem('🗑 清空所有数据', 'clearDataMenu')
    .addSeparator()
    .addItem('📖 使用帮助', 'showHelp')
    .addToUi();
}

function showKey() {
  var key = getTeacherKey_();
  var ui = SpreadsheetApp.getUi();
  var result = ui.prompt('教师密钥', '当前密钥（可修改）：', ui.ButtonSet.OK_CANCEL);
  if (result.getSelectedButton() === ui.Button.OK) {
    var newKey = result.getResponseText().trim();
    if (newKey) {
      PropertiesService.getScriptProperties().setProperty(PROP_KEY, newKey);
      ui.alert('✅ 密钥已更新为: ' + newKey);
    }
  }
}

function clearDataMenu() {
  var ui = SpreadsheetApp.getUi();
  var result = ui.alert('确认清空', '确定要清空全部诊断数据吗？此操作不可恢复。', ui.ButtonSet.YES_NO);
  if (result === ui.Button.YES) {
    clearAll_();
    ui.alert('已清空所有数据');
  }
}

function showHelp() {
  var ui = SpreadsheetApp.getUi();
  ui.alert('使用帮助',
    '1. 教师密钥在 teacher.html 中使用，用于访问教师面板\n' +
    '2. 学生提交不需要密钥\n' +
    '3. 数据存储在「' + SHEET_NAME + '」工作表中，请勿手动修改表头\n' +
    '4. 部署为 Web 应用后，将 URL 填入 HTML 文件的 APPS_SCRIPT_URL\n' +
    '5. 有任何问题请查看 README.md',
    ui.ButtonSet.OK);
}
