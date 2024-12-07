package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Question struct {
	Question string
	Options  []string
	Answer   string
}

var questions = []Question{
	{
		Question: "Apa ibukota Indonesia?",
		Options:  []string{"Jakarta", "Bandung", "Surabaya", "Medan"},
		Answer:   "Jakarta",
	},
	{
		Question: "Siapa presiden pertama Indonesia?",
		Options:  []string{"Habibie", "Soeharto", "Jokowi", "Soekarno"},
		Answer:   "Soekarno",
	},
	{
		Question: "Apa warna bendera Indonesia?",
		Options:  []string{"Merah dan Kuning", "Merah dan Putih", "Biru dan Putih", "Hijau dan Putih"},
		Answer:   "Merah dan Putih",
	},
}

var score int

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/style.css", stylesHandler)

	http.ListenAndServe(":8080", nil)
}

func stylesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Quiz Indonesia</title>
        <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap" rel="stylesheet">
        <link rel="stylesheet" href="/style.css">
    </head>
    <body class="home-bg">
        <div class="container">
            <div class="card">
                <h1>🇮🇩 Quiz Indonesia</h1>
                <p>Uji pengetahuanmu tentang Indonesia!</p>
                <a href="/quiz" class="btn btn-primary">Mulai Kuis</a>
            </div>
        </div>
    </body>
    </html>
    `
	w.Write([]byte(tmpl))
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		currentIndex, _ := strconv.Atoi(r.FormValue("currentIndex"))
		answer := r.FormValue("answer")

		if answer == questions[currentIndex].Answer {
			score++
		}

		currentIndex++

		if currentIndex >= len(questions) {
			http.Redirect(w, r, "/result", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/quiz?index="+strconv.Itoa(currentIndex), http.StatusSeeOther)
		return
	}

	currentIndex, _ := strconv.Atoi(r.URL.Query().Get("index"))
	if currentIndex < 0 || currentIndex >= len(questions) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	timerDuration := 30 // seconds

	tmpl := `
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Quiz Indonesia</title>
        <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap" rel="stylesheet">
        <link rel="stylesheet" href="/style.css">
    </head>
    <body class="quiz-bg">
        <div class="container">
            <div class="quiz-card">
                <div class="timer" id="timer">Waktu tersisa: ` + strconv.Itoa(timerDuration) + ` detik</div>
                <form id="quizForm" method="POST">
                    <input type="hidden" name="currentIndex" value="` + strconv.Itoa(currentIndex) + `">
                    <h2>` + questions[currentIndex].Question + `</h2>
                    <div class="options">
                        {{range .Options}}
                            <label class="option">
                                <input type="radio" name="answer" value="{{.}}" required>
                                <span>{{.}}</span>
                            </label>
                        {{end}}
                    </div>
                    <button type="submit" class="btn btn-submit">Kirim Jawaban</button>
                </form>
            </div>
        </div>
        <script>
            let timeLeft = ` + strconv.Itoa(timerDuration) + `;
            const timerElement = document.getElementById('timer');
            const countdown = setInterval(() => {
                if (timeLeft <= 0) {
                    clearInterval(countdown);
                    document.getElementById('quizForm').submit();
                } else {
                    timerElement.innerHTML = 'Waktu tersisa: ' + timeLeft + ' detik';
                }
                timeLeft -= 1;
            }, 1000);
        </script>
    </body>
    </html>
    `
	t := template.New("quiz")
	t, _ = t.Parse(tmpl)
	t.Execute(w, questions[currentIndex])
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Hasil Quiz</title>
        <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap" rel="stylesheet">
        <link rel="stylesheet" href="/style.css">
    </head>
    <body class="result-bg">
        <div class="container">
            <div class="result-card">
                <h1>Hasil Kuis</h1>
                <div class="score">
                    <p>Skor Anda:</p>
                    <span class="score-value">{{.}}/3</span>
                </div>
                <div class="result-actions">
                    <a href="/" class="btn btn-primary">Kembali ke Beranda</a>
                    <a href="/quiz" class="btn btn-secondary">Coba Lagi</a>
                </div>
            </div>
        </div>
    </body>
    </html>
    `
	t := template.New("result")
	t, _ = t.Parse(tmpl)
	t.Execute(w, score)
}
