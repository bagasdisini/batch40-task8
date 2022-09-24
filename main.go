package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/detail-project/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/add-project", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{index}", deleteBlog).Methods("GET")

	fmt.Println("Server berjalan di port 8080")

	http.ListenAndServe("localhost:8080", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/index.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	response := map[string]interface{}{
		"Project": dataProject,
	}

	template.Execute(w, response)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/contact.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	template.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/detail-project.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	var BlogDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if i == index {
			BlogDetail = Project{
				ProjectName: data.ProjectName,
				Description: data.Description,
			}
		}
	}

	data := map[string]interface{}{
		"dataProject": BlogDetail,
	}

	template.Execute(w, data)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/add-project.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	template.Execute(w, nil)
}

type Project struct {
	ProjectName string
	Description string
}

var dataProject = []Project{}

func addBlog(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	var projectName = r.PostForm.Get("projectName")
	var deskripsi = r.PostForm.Get("deskripsi")

	var dataBlog = Project{
		ProjectName: projectName,
		Description: deskripsi,
	}

	dataProject = append(dataProject, dataBlog)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
