package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Fukkatsuso/sudoku"
	"github.com/Fukkatsuso/sudoku-solver-app/lib/ocr"
)

func table9x9() [][]int {
	t := make([][]int, 9)
	for i := range t {
		t[i] = make([]int, 9)
	}
	return t
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[Error]", err)
	}
}

type AnalyzeResponse struct {
	Table [][]int `json:"table"`
}

func imageToSudoku(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("[Error]", "Request must be POST method")
		http.Error(w, "Request must be POST method", http.StatusMethodNotAllowed)
		return
	}
	file, _, err := r.FormFile("sudokuImage")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		file.Close()
		r.MultipartForm.RemoveAll()
	}()

	// 数独画像を解析する
	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch format {
	case "png", "jpg", "jpeg":
	default:
		err := fmt.Errorf("File format error: %s is cannot use", format)
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	table := table9x9()
	if err := ocr.ImageToSudoku(img, table, "log/images"); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 数独データを返す
	analyzeRes := AnalyzeResponse{
		Table: table,
	}
	res, err := json.Marshal(analyzeRes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("[Response]", string(res))
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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

func solveSudoku(w http.ResponseWriter, r *http.Request) {
	// 数独のJSONデータがPOSTで送られてくる
	if r.Method != "POST" {
		fmt.Println("[Error]", "Request must be POST method")
		http.Error(w, "Request must be POST method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req SolveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("[Request]", req)
	s := sudoku.Sudoku9x9{table9x9()}
	for i := range s.Table {
		for j := range s.Table[i] {
			s.Table[i][j] = req.Table[i][j]
		}
	}

	// 解けるか判定
	if !s.Solvable() {
		fmt.Println("[Error]", "This puzzle is not solvable")
		http.Error(w, "This puzzle is not solvable", http.StatusBadRequest)
		return
	}

	// 解いて返す
	s.Solve()
	solveRes := SolveResponse{
		Problem: req.Table,
		Answer:  s.Table,
	}
	res, err := json.Marshal(solveRes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("[Response]", string(res))
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("[Info]", "listen and serve on port:", port)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/api/analyze/image", imageToSudoku)
	http.HandleFunc("/api/solve/sudoku", solveSudoku)
	http.ListenAndServe(":"+port, nil)
}
