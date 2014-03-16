package handlers

import (
	"log"
	"net/http"
    "strings"
    "strconv"
    "fmt"

	"file-storage-system/adapters"
	. "file-storage-system/core"
	"github.com/gorilla/mux"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{}

	log.Println("Request received")

	err := r.ParseMultipartForm(104857600)
	if err != nil {
		res.Set(http.StatusBadRequest, "Bad Format: "+err.Error())
		res.Write(w)
		return
	}
	if r.MultipartForm == nil {
		res.Set(http.StatusBadRequest, "Bad Multipart format")
		res.Write(w)
		return
	}

	files := r.MultipartForm.File

	for _, v := range files["file"] {
		f, _ := v.Open()

		id, _ := adapters.InsertFileGlobal("files", f)
		fileObj := File{Id: id}
		res.Data.Files = append(res.Data.Files, fileObj)
	}

	res.Set(http.StatusOK, "OK")

	res.Write(w)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := adapters.FindFileGlobal("files", mux.Vars(r)["id"])
    if rstr := r.Header.Get("Range"); rstr!=""{
        var start, end int;
        size := len(data)
        ra := strings.Split(strings.Split(rstr, "=")[1], "-")
        if ra[0]!="" {
            start, _ = strconv.Atoi(ra[0])
        } else {
            start = 0
        }
        if ra[1]!="" {
            end, _ = strconv.Atoi(ra[1])
        } else {
            end = size-1
        }
        data = data[start:end+1]
        w.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, size))
    }

	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Accept-Ranges", "bytes")

	if name := r.URL.Query().Get("name"); name!=""{
		w.Header().Set("Content-disposition", "attachment;filename=" + name)
	}
	if r.URL.Query().Get("type") == "video"{
		w.Header().Set("Content-Type", "video/mp4")
	}

	if data == nil {
		w.WriteHeader(404)
	} else if r.Header.Get("Range") != "" {
		w.WriteHeader(206)
		w.Write(data)
	} else {
		w.WriteHeader(200)
		w.Write(data)
    }
}

func KnockHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{}

	res.Set(http.StatusOK, "Welcome")
	res.Write(w)
}

func Start(port string) {
	log.Printf("File storage system running [\033[0;32;1mOK%+v\033[0m] \n", port)

	go func() {
		for {
			r := mux.NewRouter()

			r.HandleFunc("/", PostHandler).Methods("POST")
			r.HandleFunc("/", KnockHandler).Methods("GET", "OPTIONS")
			r.HandleFunc("/{id:[0-9a-z]+}", GetHandler).Methods("GET")

			http.Handle("/", r)

			err := http.ListenAndServe(port, nil)

			panic(err)
		}
	}()
}
