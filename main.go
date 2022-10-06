package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
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
