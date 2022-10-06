package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key", k)
		fmt.Println("val:", strings.Join(v, ""))

	}
}

// input verification
func verifyEmpty(field string, r *http.Request) {
	r.ParseForm()
	if len(r.Form[field][0]) == 0 {
		//Handle Empty Field
	}
}

func verifyInt(v_context string, field string, r *http.Request) bool {
	r.ParseForm()
	//getInt, err := strconv.Atoi(r.Form.Get(field))
	// if err != nil {
	//not convertible to int not an int
	// }

	//or Regexpr
	if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age")); !m {

	}
	//Check the context we are validating in
	switch v_context {
	case "range":
		// handle range checking
	case "another context to check":
		// handle another context to check

	}
	return true
}

// ========================================================================================
// <select name="fruit">
// <option value="apple">apple</option>
// <option value="pear">pear</option>
// <option value="banana">banana</option>
// </select>
func verifyDropDown(field string, r *http.Request, acceptables ...string) bool {
	r.ParseForm()
	for _, l := range acceptables {
		if l == r.Form.Get(field) {
			return true
		}

	}
	return false

}

// ========================================================================================

func verifyRadio(field string, r *http.Request, acceptables ...string) bool {
	r.ParseForm()
	for _, l := range acceptables {
		if l == r.Form.Get(field) {
			return true
		}

	}
	return false

}

func verifyCheckBox(field string, r *http.Request, acceptables map[string]bool) bool {
	r.ParseForm()
	for _, l := range r.Form[field] {
		if !acceptables[l] {
			fmt.Printf("error: %s does not match any of the acceptable values", field)
			return false
		}

	}
	return true
	//Single checkbox
	// if r.Form.Get("field") != "checked" {
	// 	return false
	// }
}

func verifyEmail(field string, r *http.Request) bool {
	r.ParseForm()
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("field")); !m {
		return true
	} else {
		return false
	}
}

func verifyString(field string, r *http.Request) bool {
	r.ParseForm()
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("field")); !m {

	}
	return true
}

const MAX_UPLOAD_SIZE = 1024 * 1024

// file upload
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crtime, 10))
		csrfToken := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.tmpl")
		t.Execute(w, csrfToken)
	} else {
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
			return
		}
		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Check cntent type of the bytes in the buffer
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" { {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%v", handler.Header)
		// Create file with permission
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		// Copy file from reader to dst
		io.Copy(f, file)


		//Start of multiple file upload handling
			// 	files := r.MultipartForm.File["file"]

			// for _, fileHeader := range files {
			// 	// Restrict the size of each uploaded file to 1MB.
			// 	// To prevent the aggregate size from exceeding
			// 	// a specified value, use the http.MaxBytesReader() method
			// 	// before calling ParseMultipartForm()
			// 	if fileHeader.Size > MAX_UPLOAD_SIZE {
			// 		http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			// 		return
			// 	}

			// 	// Open the file
			// 	file, err := fileHeader.Open()
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusInternalServerError)
			// 		return
			// 	}

			// 	defer file.Close()

			// 	buff := make([]byte, 512)
			// 	_, err = file.Read(buff)
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusInternalServerError)
			// 		return
			// 	}

			// 	filetype := http.DetectContentType(buff)
			// 	if filetype != "image/jpeg" && filetype != "image/png" {
			// 		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			// 		return
			// 	}

			// 	_, err = file.Seek(0, io.SeekStart)
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusInternalServerError)
			// 		return
			// 	}

			// 	err = os.MkdirAll("./uploads", os.ModePerm)
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusInternalServerError)
			// 		return
			// 	}

			// 	f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusBadRequest)
			// 		return
			// 	}

			// 	defer f.Close()

			// 	_, err = io.Copy(f, file)
			// 	if err != nil {
			// 		http.Error(w, err.Error(), http.StatusBadRequest)
			// 		return
			// 	}
			// }

			// fmt.Fprintf(w, "Upload successful")


	// End of multiple file upload handling





	}
}

// func HTMLEscape(w io.Writer, b []byte) escapes b to w.
// func HTMLEscapeString(s string) string returns a string after escaping from s.
// func HTMLEscaper(args ...interface{}) string returns a string after escaping from multiple arguments.

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Method:", r.Method)
	if r.Method == "GET" {
		crtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crtime, 10))
		csrfToken := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.tmpl")
		t.Execute(w, csrfToken)
	} else {
		r.ParseForm()
		csrfToken := r.Form.Get("csrfToken")
		if csrfToken != "" {
			//
		}

		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username")))
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		//allows to access the form without r.Parseform()
		// r.FormValue("username")
		//prints back the escaped username to clients browser
		template.HTMLEscape(w, []byte(r.Form.Get("username")))
	}
}

func main() {
	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("Listenandserve", err)
	}
}
