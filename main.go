package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID              int    `json : "id"`
	Title           string `json : "title"`
	Auther          string `json : "auther"`
	publicationYear int   `json:"publication_year"`
}

// to store data
var (
	books = make(map[int]Book)
	mutex =& sync.Mutex{}
)

func main() {

	http.HandleFunc("/books",getBooks)
	http.HandleFunc("/books/",func(w http.ResponseWriter,r *http.Request){
		switch r.Method {
		case http.MethodGet:
			getBookById(w,r)
		case http.MethodPost:
			createBook(w,r)
		case http.MethodPut:
			updateBook(w,r)
		case http.MethodDelete:
			deleteBook(w,r)
		default:
			http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		}

	})
	fmt.Println("server is running")
	http.ListenAndServe(":8080",nil) 


	
}

func getBooks(w http.ResponseWriter,r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("content-Type","application/json")
	json.NewEncoder(w).Encode(books) 
}

func getBookById(w http.ResponseWriter,r * http.Request){
	mutex.Lock()
	defer mutex.Unlock()

	idStr :=strings.TrimPrefix(r.URL.Path,"/books/")

	id,err :=strconv.Atoi(idStr)
	if err !=nil{
		http.Error(w,"invalid id",http.StatusBadRequest)
		return
	}
	book,exists :=books[id] 
	if !exists{
		http.Error(w,"book not found",http.StatusNotFound)
		return
	}
	w.Header().Set("content-Type","application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter,r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()

	var book Book
	if err :=json.NewDecoder(r.Body).Decode(&book); err !=nil{
		http.Error(w,"invalid input",http.StatusBadRequest)
		return
	}
	if book.ID ==0{ 
		http.Error(w,"book id is erquired",http.StatusBadRequest)
		return
	} 
	books[book.ID]=book
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book) 
}

func updateBook(w http.ResponseWriter,r * http.Request){
	mutex.Lock()
	defer mutex.Unlock()

	idStr:=strings.TrimPrefix(r.URL.Path,"/books/")
	id,err :=strconv.Atoi(idStr) 
	if err !=nil{
		http.Error(w,"invalid id",http.StatusBadRequest)
		return
	}
	var updatdeBook Book
	if err:=json.NewDecoder(r.Body).Decode(&updatdeBook) ; err !=nil{
		http.Error(w,"invalid input",http.StatusBadRequest)
		return
	}
	if updatdeBook.ID !=id{
		http.Error(w,"book id mismatch",http.StatusBadRequest)
		return
	}
	books[id]=updatdeBook
	w.Header().Set("content-Type","application/json")
	json.NewEncoder(w).Encode(updatdeBook)
}
func deleteBook(w http.ResponseWriter,r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()



	idStr:=strings.TrimPrefix(r.URL.Path,"/books/")

	id,err :=strconv.Atoi(idStr)
	if err !=nil{
		http.Error(w,"invalid id",http.StatusBadRequest)
		return
	}

	if _,exists:=books[id]; !exists{
		http.Error(w,"book not found",http.StatusNotFound)
		return
	}

	delete(books,id)

	w.WriteHeader(http.StatusNoContent)


	   

}
