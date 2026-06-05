package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// =================== 数据模型 ===================

type DiagnosisSubmit struct {
	Name       string `json:"name"`
	Grade      string `json:"grade"`
	Scores     DimensionScores `json:"scores"`
	InputScore  int    `json:"inputScore"`
	OutputScore int    `json:"outputScore"`
	ActualRead  int    `json:"actualRead"`
	SelfRating  int    `json:"selfRating"`
	GapAdvice   string `json:"gapAdvice"`
	PrefAdvice  string `json:"prefAdvice"`
	MatrixData  []MatrixRecord `json:"matrixData"`
}

type DimensionScores struct {
	Grammar      int `json:"grammar"`
	Vocabulary   int `json:"vocabulary"`
	Intuition    int `json:"intuition"`
	Willingness  int `json:"willingness"`
	Adaptability int `json:"adaptability"`
	MetaCognition int `json:"metaCognition"`
}

type MatrixRecord struct {
	Qid           int    `json:"qid"`
	IsCorrect     bool   `json:"isCorrect"`
	CognitionType string `json:"cognitionType"`
	Skill         string `json:"skill"`
}

type StudentSummary struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	LastTest  string    `json:"lastTest"`
	AvgScore  int       `json:"avgScore"`
	Level     string    `json:"level"`
}

type StudentDetail struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Grade        string          `json:"grade"`
	Diagnoses    []DiagnosisResult `json:"diagnoses"`
}

type DiagnosisResult struct {
	ID            int             `json:"id"`
	Scores        DimensionScores `json:"scores"`
	InputScore    int             `json:"inputScore"`
	OutputScore   int             `json:"outputScore"`
	ActualRead    int             `json:"actualRead"`
	SelfRating    int             `json:"selfRating"`
	GapAdvice     string          `json:"gapAdvice"`
	PrefAdvice    string          `json:"prefAdvice"`
	MatrixData    []MatrixRecord  `json:"matrixData"`
	CreatedAt     string          `json:"createdAt"`
}

type DashboardStats struct {
	TotalStudents  int              `json:"totalStudents"`
	TotalTests     int              `json:"totalTests"`
	AvgScores      DimensionScores  `json:"avgScores"`
	LevelDist      map[string]int   `json:"levelDist"`
	RecentTests    []StudentSummary `json:"recentTests"`
}

var db *sql.DB

// =================== 数据库初始化 ===================

func initDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		grade TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS diagnosis_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		grammar_score INTEGER DEFAULT 0,
		vocabulary_score INTEGER DEFAULT 0,
		intuition_score INTEGER DEFAULT 0,
		willingness_score INTEGER DEFAULT 0,
		adaptability_score INTEGER DEFAULT 0,
		meta_cognition_score INTEGER DEFAULT 0,
		input_score INTEGER DEFAULT 0,
		output_score INTEGER DEFAULT 0,
		actual_read_percent INTEGER DEFAULT 0,
		self_rating_read INTEGER DEFAULT 0,
		vocab_method TEXT DEFAULT '',
		fav_activity TEXT DEFAULT '',
		matrix_data TEXT DEFAULT '[]',
		gap_advice TEXT DEFAULT '',
		pref_advice TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (student_id) REFERENCES students(id)
	);

	CREATE INDEX IF NOT EXISTS idx_diag_student ON diagnosis_results(student_id);
	`
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}
	return nil
}

// =================== HTTP 处理器 ===================

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// POST /api/submit
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "仅支持POST", http.StatusMethodNotAllowed)
		return
	}

	var data DiagnosisSubmit
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "JSON解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(data.Name) == "" {
		data.Name = "匿名学生"
	}

	// 查找或创建学生
	studentID, err := findOrCreateStudent(data.Name, data.Grade)
	if err != nil {
		http.Error(w, "保存学生失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 保存诊断结果
	matrixJSON, _ := json.Marshal(data.MatrixData)
	_, err = db.Exec(`
		INSERT INTO diagnosis_results (
			student_id, grammar_score, vocabulary_score, intuition_score,
			willingness_score, adaptability_score, meta_cognition_score,
			input_score, output_score, actual_read_percent, self_rating_read,
			vocab_method, fav_activity, matrix_data, gap_advice, pref_advice
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		studentID,
		data.Scores.Grammar, data.Scores.Vocabulary, data.Scores.Intuition,
		data.Scores.Willingness, data.Scores.Adaptability, data.Scores.MetaCognition,
		data.InputScore, data.OutputScore, data.ActualRead, data.SelfRating,
		data.PrefAdvice, "", // vocab_method and fav_activity extracted from prefAdvice
		string(matrixJSON), data.GapAdvice, data.PrefAdvice,
	)
	if err != nil {
		http.Error(w, "保存诊断失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"studentId": studentID,
		"message": "诊断结果已保存",
	})
}

func findOrCreateStudent(name, grade string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM students WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		result, err := db.Exec("INSERT INTO students (name, grade) VALUES (?, ?)", name, grade)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	}
	if err != nil {
		return 0, err
	}
	// 更新年级
	db.Exec("UPDATE students SET grade = ? WHERE id = ?", grade, id)
	return id, nil
}

// GET /api/students
func handleStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT s.id, s.name, s.grade,
			COALESCE(MAX(d.created_at), '') as last_test,
			CAST(COALESCE(AVG(
				d.grammar_score + d.vocabulary_score + d.intuition_score +
				d.willingness_score + d.adaptability_score + d.meta_cognition_score
			) / 6.0, 0) AS INTEGER) as avg_score
		FROM students s
		LEFT JOIN diagnosis_results d ON s.id = d.student_id
		GROUP BY s.id
		ORDER BY last_test DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	students := make([]StudentSummary, 0)
	for rows.Next() {
		var s StudentSummary
		var lastTest string
		if err := rows.Scan(&s.ID, &s.Name, &s.Grade, &lastTest, &s.AvgScore); err != nil {
			continue
		}
		if lastTest != "" {
			t, _ := time.Parse("2006-01-02 15:04:05", lastTest)
			s.LastTest = t.Format("01-02 15:04")
		} else {
			s.LastTest = "-"
		}
		s.Level = calcLevel(s.AvgScore)
		students = append(students, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GET /api/student?id=xxx  or  DELETE /api/student?id=xxx
func handleStudentRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleStudentDetail(w, r)
	case "DELETE":
		handleDeleteStudent(w, r)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "不支持的方法", http.StatusMethodNotAllowed)
	}
}

func handleStudentDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "缺少学生ID", http.StatusBadRequest)
		return
	}

	var detail StudentDetail
	err := db.QueryRow("SELECT id, name, grade FROM students WHERE id = ?", idStr).
		Scan(&detail.ID, &detail.Name, &detail.Grade)
	if err == sql.ErrNoRows {
		http.Error(w, "学生不存在", http.StatusNotFound)
		return
	}

	rows, err := db.Query(`
		SELECT id, grammar_score, vocabulary_score, intuition_score,
			willingness_score, adaptability_score, meta_cognition_score,
			input_score, output_score, actual_read_percent, self_rating_read,
			matrix_data, gap_advice, pref_advice, created_at
		FROM diagnosis_results
		WHERE student_id = ?
		ORDER BY created_at DESC
	`, idStr)
	if err != nil {
		http.Error(w, "查询失败", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d DiagnosisResult
		var matrixStr, createdAt string
		if err := rows.Scan(&d.ID,
			&d.Scores.Grammar, &d.Scores.Vocabulary, &d.Scores.Intuition,
			&d.Scores.Willingness, &d.Scores.Adaptability, &d.Scores.MetaCognition,
			&d.InputScore, &d.OutputScore, &d.ActualRead, &d.SelfRating,
			&matrixStr, &d.GapAdvice, &d.PrefAdvice, &createdAt,
		); err != nil {
			continue
		}
		json.Unmarshal([]byte(matrixStr), &d.MatrixData)
		d.CreatedAt = createdAt
		detail.Diagnoses = append(detail.Diagnoses, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

// GET /api/stats
func handleStats(w http.ResponseWriter, r *http.Request) {
	var stats DashboardStats
	stats.LevelDist = make(map[string]int)

	// 总学生数
	db.QueryRow("SELECT COUNT(*) FROM students").Scan(&stats.TotalStudents)

	// 总诊断数
	db.QueryRow("SELECT COUNT(*) FROM diagnosis_results").Scan(&stats.TotalTests)

	// 平均分
	row := db.QueryRow(`
		SELECT
			COALESCE(AVG(grammar_score), 0),
			COALESCE(AVG(vocabulary_score), 0),
			COALESCE(AVG(intuition_score), 0),
			COALESCE(AVG(willingness_score), 0),
			COALESCE(AVG(adaptability_score), 0),
			COALESCE(AVG(meta_cognition_score), 0)
		FROM diagnosis_results
	`)
	row.Scan(&stats.AvgScores.Grammar, &stats.AvgScores.Vocabulary,
		&stats.AvgScores.Intuition, &stats.AvgScores.Willingness,
		&stats.AvgScores.Adaptability, &stats.AvgScores.MetaCognition)

	// 等级分布
	rows, _ := db.Query(`
		SELECT CAST(
			(grammar_score + vocabulary_score + intuition_score +
			willingness_score + adaptability_score + meta_cognition_score) / 6.0 AS INTEGER
		) as avg_score
		FROM diagnosis_results
	`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var avg int
			rows.Scan(&avg)
			stats.LevelDist[calcLevel(avg)]++
		}
	}

	// 最近测试
	recentRows, _ := db.Query(`
		SELECT s.id, s.name, s.grade,
			COALESCE(MAX(d.created_at), '') as last_test,
			(SELECT CAST((grammar_score+vocabulary_score+intuition_score+willingness_score+adaptability_score+meta_cognition_score)/6.0 AS INTEGER)
			 FROM diagnosis_results WHERE student_id = s.id ORDER BY created_at DESC LIMIT 1) as latest_avg
		FROM students s
		LEFT JOIN diagnosis_results d ON s.id = d.student_id
		GROUP BY s.id
		ORDER BY last_test DESC
		LIMIT 10
	`)
	if recentRows != nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var s StudentSummary
			var lt string
			recentRows.Scan(&s.ID, &s.Name, &s.Grade, &lt, &s.AvgScore)
			if lt != "" {
				t, _ := time.Parse("2006-01-02 15:04:05", lt)
				s.LastTest = t.Format("01-02 15:04")
			} else {
				s.LastTest = "-"
			}
			s.Level = calcLevel(s.AvgScore)
			stats.RecentTests = append(stats.RecentTests, s)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func calcLevel(avg int) string {
	switch {
	case avg >= 85:
		return "A-优秀"
	case avg >= 70:
		return "B-良好"
	case avg >= 55:
		return "C-中等"
	default:
		return "D-待提升"
	}
}

// DELETE /api/diagnosis?id=xxx
func handleDeleteDiagnosis(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "仅支持DELETE", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "缺少诊断ID", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("DELETE FROM diagnosis_results WHERE id = ?", idStr)
	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DELETE /api/student?id=xxx
func handleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "仅支持DELETE", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "缺少学生ID", http.StatusBadRequest)
		return
	}
	// 先删除诊断记录，再删除学生
	db.Exec("DELETE FROM diagnosis_results WHERE student_id = ?", idStr)
	db.Exec("DELETE FROM students WHERE id = ?", idStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// =================== 静态文件服务 ===================

func getStaticDir() string {
	execPath, _ := os.Executable()
	staticDir := filepath.Join(filepath.Dir(execPath), "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		staticDir = filepath.Join(cwd, "static")
	}
	return staticDir
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(getStaticDir()))
	fs.ServeHTTP(w, r)
}

// =================== 主函数 ===================

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	// 确定数据库路径
	execPath, _ := os.Executable()
	dbDir := filepath.Dir(execPath)
	dbPath := filepath.Join(dbDir, "gofit.db")

	// 开发模式：如果exe还没编译好，用当前目录
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		dbPath = filepath.Join(cwd, "gofit.db")
	}

	if err := initDB(dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	localIP := getLocalIP()
	shareURL := "http://localhost:8877"
	if localIP != "" {
		shareURL = "http://" + localIP + ":8877"
	}

	fmt.Println("==================================")
	fmt.Println("  GoFit · 英语诊断工具 后端服务")
	fmt.Println("==================================")
	fmt.Printf("  数据库: %s\n", dbPath)
	fmt.Println()
	fmt.Println("  🌐 公网链接（发给学生/家长）：")
	fmt.Printf("     🧑‍🎓 学生入口: %s/s\n", shareURL)
	fmt.Printf("     👩‍🏫 教师入口: %s/t\n", shareURL)
	fmt.Println()
	fmt.Println("  本机: http://localhost:8877")
	fmt.Println("==================================")

	// API 路由
	http.HandleFunc("/api/submit", corsMiddleware(handleSubmit))
	http.HandleFunc("/api/students", corsMiddleware(handleStudents))
	http.HandleFunc("/api/student", corsMiddleware(handleStudentRoute))
	http.HandleFunc("/api/stats", corsMiddleware(handleStats))
	http.HandleFunc("/api/diagnosis", corsMiddleware(handleDeleteDiagnosis))

	// 短链接入口
	http.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(getStaticDir(), "学生端.html"))
	})
	http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(getStaticDir(), "教师端.html"))
	})

	// 静态文件
	http.HandleFunc("/", serveStatic)

	if err := http.ListenAndServe(":8877", nil); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
