package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// ==========================================
// 1. HTML STRING (LANGSUNG DI SINI BIAR AMAN)
// ==========================================
const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Web Perkalian Digit (Vercel)</title>
    <style>
        body { font-family: sans-serif; background: #f4f4f9; display: flex; justify-content: center; padding-top: 50px; }
        .card { background: white; padding: 2rem; border-radius: 10px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); width: 100%; max-width: 400px; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 5px; box-sizing: border-box;}
        button { width: 100%; padding: 10px; background: #000; color: white; border: none; border-radius: 5px; cursor: pointer; }
        button:hover { background: #333; }
        .result { margin-top: 20px; padding: 15px; background: #e8f6f3; border-left: 5px solid #1abc9c; }
        .error { color: red; margin-top: 10px; }
        h3 { margin-top: 0; color: #2c3e50; }
    </style>
</head>
<body>
<div class="card">
    <h2 style="text-align:center;">Perkalian Digit (Go)</h2>
    <form method="POST">
        <label>Masukkan Angka (n):</label>
        <input type="number" name="angka" placeholder="Contoh: 1234" required>
        <button type="submit">Hitung</button>
    </form>
    {{if .ErrorMsg}}
        <div class="error">{{.ErrorMsg}}</div>
    {{end}}
    {{if .ShowResult}}
        <div class="result">
            <h3>Hasil Perhitungan:</h3>
            <p>Input: <strong>{{.InputAngka}}</strong></p>
            <hr>
            <p>Iteratif: <strong>{{.HasilIteratif}}</strong></p>
            <p>Rekursif: <strong>{{.HasilRekursif}}</strong></p>
        </div>
    {{end}}
</div>
</body>
</html>
`

// ==========================================
// 2. LOGIKA MATEMATIKA (KODE ANDA)
// ==========================================

func iterative(n int) int {
	if n == 0 {
		return 0
	}
	product := 1
	for n > 0 {
		digit := n % 10
		product *= digit
		n /= 10
	}
	return product
}

func recursive(n int) int {
	if n < 10 {
		return n
	}
	return (n % 10) * recursive(n/10)
}

// ==========================================
// 3. HANDLER VERCEL (PENTING!)
// ==========================================

// Struktur data
type PageData struct {
	ShowResult    bool
	InputAngka    int
	HasilIteratif int
	HasilRekursif int
	ErrorMsg      string
}

// Perhatikan: Nama fungsinya Handler (Huruf besar H)
// Vercel akan mencari fungsi ini secara otomatis.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Parse HTML dari string constant di atas
	tmpl, err := template.New("web").Parse(htmlTemplate)
	if err != nil {
		fmt.Fprintf(w, "Error template: %v", err)
		return
	}

	var data PageData

	// Logika POST (Saat tombol ditekan)
	if r.Method == http.MethodPost {
		inputStr := r.FormValue("angka")
		n, err := strconv.Atoi(inputStr)

		if err != nil {
			data.ErrorMsg = "Input harus angka bulat!"
		} else {
			data.ShowResult = true
			data.InputAngka = n
			data.HasilIteratif = iterative(n)
			data.HasilRekursif = recursive(n)
		}
	}

	// Tampilkan
	tmpl.Execute(w, data)
}