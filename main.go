package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

type SolveRequest struct {
	Table [][]int `json:"table"`
}

type SolveResponse struct {
	// Problem string `json:"problem"`
	// Answer  string `json:"answer"`
	Problem [][]int `json:"problem"`
	Answer  [][]int `json:"answer"`
}

func table9x9() [][]int {
	t := make([][]int, 9)
	for i := range t {
		t[i] = make([]int, 9)
	}
	return t
}

func solveSudoku(w http.ResponseWriter, r *http.Request) {
	// 数独のJSONデータがPOSTで送られてくる
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var req SolveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[Request]", req)
	s := sudoku.Sudoku9x9{table9x9()}
	for i := range s.Table {
		for j := range s.Table[i] {
			s.Table[i][j] = req.Table[i][j]
		}
	}

	// 解いて返す
	s.Solve()
	solveRes := SolveResponse{
		Problem: req.Table,
		Answer:  s.Table,
	}
	res, err := json.Marshal(solveRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("[Response]", string(res))
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
