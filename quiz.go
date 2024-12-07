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
		Options:  []string{"Soekarno", "Soeharto", "Jokowi", "Habibie"},
		Answer:   "Soekarno",
	},
	{
		Question: "Apa warna bendera Indonesia?",
		Options:  []string{"Merah dan Putih", "Merah dan Kuning", "Biru dan Putih", "Hijau dan Putih"},
		Answer:   "Merah dan Putih",
	},
}

var score int

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/result", resultHandler)

	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Quiz App</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f4;
                color: #333;
                text-align: center;
                padding: 50px;
            }
            a {
                text-decoration: none;
                color: #fff;
                background-color: #007bff;
                padding: 10px 20px;
                border-radius: 5px;
            }
            a:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <h1>Selamat Datang di Quiz</h1>
        <a href="/quiz">Mulai Kuis</a>
    </body>
    </html>
    `
	w.Write([]byte(tmpl))
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		score = 0
		for i, question := range questions {
			answer := r.FormValue("question" + strconv.Itoa(i))
			if answer == question.Answer {
				score++
			}
		}
		http.Redirect(w, r, "/result", http.StatusSeeOther)
		return
	}

	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Quiz</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f4;
                color: #333;
                padding: 20px;
            }
            h1 {
                color: #007bff;
            }
            fieldset {
                border: 2px solid #007bff;
                border-radius: 5px;
                padding: 10px;
                margin-bottom: 20px;
            }
            legend {
                font-weight: bold;
                color: #007bff;
            }
            label {
                display: block;
                margin: 5px 0;
            }
            button {
                background-color: #007bff;
                color: white;
                border: none;
                padding: 10px 20px;
                border-radius: 5px;
                cursor: pointer;
            }
            button:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <h1>Kuis</h1>
        <form method="POST">
            {{range $index, $question := .}}
                <fieldset>
                    <legend>{{$question.Question}}</legend>
                    {{range $option := $question.Options}}
                        <label>
                            <input type="radio" name="question{{ $index }}" value="{{ $option }}" required>
                            {{ $option }}
                        </label>
                    {{end}}
                </fieldset>
            {{end}}
            <button type="submit">Kirim Jawaban</button>
        </form>
    </body>
    </html>
    `
	t := template.New("quiz")
	t, _ = t.Parse(tmpl)
	t.Execute(w, questions)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Hasil Kuis</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f4;
                color: #333;
                text-align: center;
                padding: 20px;
            }
            h1 {
                color: #007bff;
            }
            a {
                text-decoration: none;
                color: #fff;
                background-color: #007bff;
                padding: 10px 20px;
                border-radius: 5px;
            }
            a:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <h1>Hasil Kuis</h1>
        <p>Skor Anda: {{.}}</p>
        <a href="/">Kembali ke Beranda</a>
    </body>
    </html>
    `
	t := template.New("result")
	t, _ = t.Parse(tmpl)
	t.Execute(w, score)
}
