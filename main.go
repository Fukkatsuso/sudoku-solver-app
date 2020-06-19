package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Fukkatsuso/sudoku"
)

// 数字文字列を9x9のテーブルに変換
func stringToTable(s string) [][]int {
	t := make([][]int, 9)
	for i := 0; i < 9; i++ {
		t[i] = make([]int, 9)
		for j := 0; j < 9; j++ {
			t[i][j] = int(s[i*9+j] - '0')
		}
	}
	return t
}

// 9x9のテーブルを文字列に変換
func tableToString(t [][]int) string {
	bs := make([]byte, 9*9)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bs[i*9+j] = byte('0' + t[i][j])
		}
	}
	return string(bs)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[Error]", err)
	}
}

func imageToSudoku(w http.ResponseWriter, r *http.Request) {
	// 数独画像を解析する
	// 数独データを返す
}

type SolveResponse struct {
	Problem string `json:"problem"`
	Answer  string `json:"answer"`
}

func solveSudoku(w http.ResponseWriter, r *http.Request) {
	// 数独データが送られてくる
	tablestr := r.URL.Query().Get("table")
	table := stringToTable(tablestr)
	s := sudoku.Sudoku9x9{table}
	// 解いて返す
	s.Solve()
	solveRes := SolveResponse{
		Problem: tablestr,
		Answer:  tableToString(s.Table),
	}
	res, err := json.Marshal(solveRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", index)
	http.HandleFunc("/api/image/analyze", imageToSudoku)
	http.HandleFunc("/api/sudoku/solve", solveSudoku)
	http.ListenAndServe(":"+port, nil)
}
