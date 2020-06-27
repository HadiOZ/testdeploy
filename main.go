package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Payload struct {
	Status string
	Value  interface{}
}

type Recipe struct {
	Name        string
	Picture     string
	ingredients []string
	Steps       []string
}

type Info struct {
	Title    string
	Desc     string
	Img      string
	Category string
	Code     string
	Point    int
	Time     string
	Visibel  string
}

type Question struct {
	Img    string
	Desc   string
	OpsiA  string
	OpsiB  string
	OpsiC  string
	OpsiD  string
	Answer string
}

type Quiz struct {
	QuizInfo Info
	ListQuiz []Question
}

func main() {
	http.HandleFunc("/recipes", corsmiddleware(recipesHandel))
	http.HandleFunc("/recipe", corsmiddleware(recipeHandel))
	http.HandleFunc("/quizs", corsmiddleware(allQuizHandel))
	http.HandleFunc("/create", corsmiddleware(createQuizHandel))
	http.HandleFunc("/quiz", corsmiddleware(quizHandel))
	fmt.Println("server run on port 8086")
	if err := http.ListenAndServe(":8086", nil); err != nil {
		log.Fatal(err)
	}
}

func corsmiddleware(hendel http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		hendel(w, r)
	}
}

func success(w http.ResponseWriter, data interface{}) {
	p := Payload{
		Status: "success 200",
		Value:  data,
	}
	bytes, _ := json.Marshal(p)
	w.Write(bytes)
}

func failed(w http.ResponseWriter, m string, code int) {
	p := Payload{
		Status: "failed " + string(code),
		Value:  m,
	}
	bytes, _ := json.Marshal(p)
	w.Write(bytes)
}

func createQuizHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		failed(w, "method not allowed ! use POST metode for this endpoin", http.StatusBadRequest)
		return
	}

	var quiz Quiz
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&quiz); err != nil {
		failed(w, "someting wrong in data stuctur", http.StatusBadRequest)
		return
	}
	if err := writeJSON(quiz); err != nil {
		failed(w, "failed to store data", http.StatusInternalServerError)
		return
	}

	success(w, "data was stored")
}

func allQuizHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		failed(w, "method not allowed ! use GET metode for this endpoin", http.StatusBadRequest)
		return
	}

	data, err := readJSON()
	if err != nil {
		failed(w, "can't read json data", http.StatusInternalServerError)
		return
	}
	success(w, data)
}

func recipesHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		failed(w, "method not allowed ! use GET metode for this endpoin", http.StatusBadRequest)
		return
	}

	data, err := readRecipeJSON()
	if err != nil {
		failed(w, "can't read json data", http.StatusInternalServerError)
		return
	}
	success(w, data)
}

func quizHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		failed(w, "method not allowed ! use GET metode for this endpoin", http.StatusBadRequest)
		return
	}

	var quiz Quiz
	data, err := readJSON()
	if err != nil {
		failed(w, "can't read json data", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	for _, i := range data {
		if i.QuizInfo.Code == code {
			quiz = i
			break
		}
	}

	if quiz.QuizInfo.Code == "" {
		failed(w, "data quiz not found", http.StatusBadRequest)
		return
	} else {
		success(w, quiz)
	}

}

func recipeHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		failed(w, "method not allowed ! use GET metode for this endpoin", http.StatusBadRequest)
		return
	}

	var recipe Recipe
	data, err := readRecipeJSON()
	if err != nil {
		failed(w, "can't read json data", http.StatusInternalServerError)
		return
	}

	name := r.URL.Query().Get("name")
	for _, i := range data {
		if i.Name == name {
			recipe = i
			break
		}
	}

	if recipe.Name == "" {
		failed(w, "data quiz not found", http.StatusBadRequest)
		return
	} else {
		success(w, recipe)
	}

}

func readJSON() ([]Quiz, error) {
	var data []Quiz

	file, err := os.Open("data.json")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	bites, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bites, &data)

	return data, nil
}

func readRecipeJSON() ([]Recipe, error) {
	var data []Recipe

	file, err := os.Open("recipes.json")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	bites, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bites, &data)

	return data, nil
}

func writeJSON(item Quiz) error {
	var data []Quiz

	//read data
	file, err := os.Open("data.json")
	defer file.Close()
	if err != nil {
		return err
	}
	bites, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	json.Unmarshal(bites, &data)

	//write data
	data = append(data, item)
	bit, err := json.Marshal(data)
	err = ioutil.WriteFile("data.json", bit, 0777)
	if err != nil {
		return err
	}

	return nil
}
