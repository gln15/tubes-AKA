package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// ==========================================
// 1. HTML & CSS (TAMPILAN BARU UNTUK ANALISIS)
// ==========================================
const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Analisis Kompleksitas: Perkalian Digit</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #f0f2f5; padding: 20px; display: flex; justify-content: center; }
        .container { max-width: 800px; width: 100%; }
        
        .card-input { background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); margin-bottom: 20px; text-align: center; }
        input { padding: 12px; width: 60%; border: 1px solid #ddd; border-radius: 6px; font-size: 16px; }
        button { padding: 12px 24px; background: #0070f3; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 16px; margin-left: 10px; transition: 0.2s;}
        button:hover { background: #0051a2; }

        .comparison-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        .result-card { background: white; padding: 20px; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); border-top: 5px solid #ccc; }
        
        .iterative { border-top-color: #10b981; } /* Hijau */
        .recursive { border-top-color: #f59e0b; } /* Oranye */

        h2 { color: #333; margin-bottom: 5px; }
        .badge { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; color: white; margin-bottom: 15px;}
        .bg-green { background: #10b981; }
        .bg-orange { background: #f59e0b; }

        .metric { display: flex; justify-content: space-between; border-bottom: 1px solid #eee; padding: 10px 0; }
        .metric strong { color: #555; }
        .metric span { font-family: monospace; font-size: 1.1em; color: #000; }
        
        .big-o { margin-top: 15px; padding: 10px; background: #f8f9fa; border-radius: 6px; font-size: 0.9em; color: #666; font-style: italic;}
        
        @media (max-width: 600px) { .comparison-grid { grid-template-columns: 1fr; } input { width: 100%; margin-bottom: 10px; } }
    </style>
</head>
<body>

<div class="container">
    <div class="card-input">
        <h1>Analisis Algoritma Perkalian Digit</h1>
        <p>Masukkan angka n untuk membandingkan performa Iteratif vs Rekursif</p>
        <form method="POST">
            <input type="number" name="angka" placeholder="Contoh: 12345" required value="{{.InputAngka}}">
            <button type="submit">Analisis</button>
        </form>
        {{if .ErrorMsg}}<p style="color:red">{{.ErrorMsg}}</p>{{end}}
    </div>

    {{if .ShowResult}}
    <div class="comparison-grid">
        <div class="result-card iterative">
            <span class="badge bg-green">ITERATIF</span>
            <h2>Hasil: {{.ResIter.Value}}</h2>
            
            <div class="metric">
                <strong>Waktu Eksekusi:</strong>
                <span>{{.ResIter.TimeTaken}}</span>
            </div>
            <div class="metric">
                <strong>Jumlah Langkah (Loop):</strong>
                <span>{{.ResIter.Steps}}</span>
            </div>
            
            <div class="big-o">
                <strong>Analisis:</strong> Time Complexity O(d) di mana d adalah jumlah digit. Space Complexity O(1) karena hanya menggunakan satu variabel.
            </div>
        </div>

        <div class="result-card recursive">
            <span class="badge bg-orange">REKURSIF</span>
            <h2>Hasil: {{.ResRec.Value}}</h2>
            
            <div class="metric">
                <strong>Waktu Eksekusi:</strong>
                <span>{{.ResRec.TimeTaken}}</span>
            </div>
            <div class="metric">
                <strong>Jumlah Panggilan (Call):</strong>
                <span>{{.ResRec.Steps}}</span>
            </div>

            <div class="big-o">
                <strong>Analisis:</strong> Time Complexity O(d). Space Complexity O(d) karena memakan Stack Memory sebanyak jumlah digit.
            </div>
        </div>
    </div>
    {{end}}
</div>

</body>
</html>
`

// ==========================================
// 2. STRUKTUR DATA
// ==========================================

// Menyimpan hasil perhitungan + metrik analisis
type ResultMetric struct {
	Value     int    // Hasil hitungan (misal: 120)
	Steps     int    // Jumlah operasi (loop/rekursi)
	TimeTaken string // Waktu dalam nanosecond
}

type PageData struct {
	ShowResult bool
	InputAngka string
	ErrorMsg   string
	ResIter    ResultMetric
	ResRec     ResultMetric
}

// ==========================================
// 3. LOGIKA ALGORITMA (DENGAN COUNTER)
// ==========================================

// Iteratif: Mengembalikan (hasil, jumlah_loop)
func iterativeWithStats(n int) (int, int) {
	if n == 0 {
		return 0, 1
	}
	product := 1
	steps := 0

	// Analisis: Loop berjalan sebanyak jumlah digit
	for n > 0 {
		steps++ // Hitung 1 langkah setiap loop
		digit := n % 10
		product *= digit
		n /= 10
	}
	return product, steps
}

// Rekursif Wrapper: Supaya bisa hitung step total
// Kita butuh fungsi pembantu agar step tidak kureset tiap rekursi
func recursiveWithStats(n int) (int, int) {
	steps := 0
	var solve func(int) int

	solve = func(val int) int {
		steps++ // Hitung 1 langkah setiap panggilan fungsi
		if val < 10 {
			return val
		}
		return (val % 10) * solve(val/10)
	}

	result := solve(n)
	return result, steps
}

// ==========================================
// 4. HANDLER UTAMA
// ==========================================

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("web").Parse(htmlTemplate)
	if err != nil {
		fmt.Fprintf(w, "Error template: %v", err)
		return
	}

	data := PageData{}

	if r.Method == http.MethodPost {
		inputStr := r.FormValue("angka")
		data.InputAngka = inputStr // Simpan biar input gak hilang
		n, err := strconv.Atoi(inputStr)

		if err != nil {
			data.ErrorMsg = "Input harus angka bulat!"
		} else {
			data.ShowResult = true

			// --- UKUR ITERATIF ---
			startIter := time.Now()
			resIter, stepsIter := iterativeWithStats(n)
			durationIter := time.Since(startIter)

			data.ResIter = ResultMetric{
				Value:     resIter,
				Steps:     stepsIter,
				TimeTaken: fmt.Sprintf("%d ns", durationIter.Nanoseconds()),
			}

			// --- UKUR REKURSIF ---
			startRec := time.Now()
			resRec, stepsRec := recursiveWithStats(n)
			durationRec := time.Since(startRec)

			data.ResRec = ResultMetric{
				Value:     resRec,
				Steps:     stepsRec,
				TimeTaken: fmt.Sprintf("%d ns", durationRec.Nanoseconds()),
			}
		}
	}
	// Update terakhir untuk cek Vercel
	tmpl.Execute(w, data)
}
