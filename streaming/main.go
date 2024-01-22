package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"io/fs"

	"github.com/gorilla/mux"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Open the target file for writing
	file, err := os.Create("output.txt")
	if err != nil {
		http.Error(w, "Unable to open file for writing", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Stream data from the request body to the file
	// This will copy the request body to the file in small chunks instead of loading it all into memory
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func uploadFileHandlerInMem(w http.ResponseWriter, r *http.Request) {
	// Read the entire request body into memory
	data, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Write the data to a file
	err = os.WriteFile("output1.txt", data, fs.FileMode(0644))
	if err != nil {
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

var RegisterBookStore = func(router *mux.Router) {
	router.HandleFunc("/file", uploadFileHandler).Methods("POST")
	router.HandleFunc("/fileInMem", uploadFileHandlerInMem).Methods("POST")
	// router.HandleFunc("/book/{id}", controllers.GetBookById).Methods("GET")
	// // router.HandleFunc("/book/{id}", controllers.UpdateBook).Methods("PUT")
	// // router.HandleFunc("/book/{id}", controllers.DeleteBook).Methods("DELETE")
	// router.HandleFunc("/book", controllers.GetBook).Methods("GET")
}

func main() {
	r := mux.NewRouter()
	RegisterBookStore(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r))
	// reader, err := os.Open("./machao.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer reader.Close()
	// p := make([]byte, 4)

	// for {
	// 	n, err := reader.Read(p)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Println(string(p[:n])) //should handle any remainding bytes.
	// 			fmt.Println("--------------------------------")
	// 			break
	// 		}
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(string(p[:n]))
	// 	fmt.Println("--------------------------------")
	// }
}
