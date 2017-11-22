package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sys-cat/linelogin"
	"github.com/sys-cat/linelogin/token"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Println("access index")
	fmt.Fprintf(w, "<h1>This is sys-cat test site</h1>")
	url := linelogin.New()
	err := url.Parameters(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"), os.Getenv("REDIRECT_URL"))
	if err != nil {
		fmt.Fprintf(w, "<p style=\"color:red;\">link url build error !</p>")
	}
	log.Printf("%+v\n", url.OutputURL())
	link := fmt.Sprintf("<a href=\"%s\">Line Login</a>", url.OutputURL())
	io.WriteString(w, link)
}

func line_login(w http.ResponseWriter, r *http.Request) {
	log.Println("access line")
	urlParam := linelogin.New()
	err := urlParam.Parameters(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"), os.Getenv("REDIRECT_URL"))
	if err != nil {
		data := fmt.Sprintf("<p>missing render: %s</p>", err.Error)
		fmt.Fprintf(w, data)
	}
	http.Redirect(w, r, urlParam.OutputURL(), 301)
}

func line_login_test(w http.ResponseWriter, r *http.Request) {
	log.Println("access build")
	urlParam := linelogin.New()
	err := urlParam.Parameters(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"), os.Getenv("REDIRECT_URL"))
	if err != nil {
		data := fmt.Sprintf("<p>missing render: %s</p>", err.Error)
		fmt.Fprintf(w, data)
	}
	data := fmt.Sprint(urlParam.OutputURL())
	log.Println(data)
	fmt.Fprintf(w, data)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%+v\n", r.URL.Query())
	code := r.URL.Query().Get("code")
	fmt.Fprintf(w, "%+v\n", code)
	state := r.URL.Query().Get("state")
	fmt.Fprintf(w, "%+v\n", state)
	if state != "" {
		fmt.Fprintf(w, "Invalid access\n")
	}
	newToken := token.New()
	err := newToken.Parameters(code, os.Getenv("REDIRECT_URL"), os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"))
	if err != nil {
		fmt.Fprintf(w, "Invalid parameters\n")
	}
	res, err := token.GetToken(newToken)
	if err != nil {
		fmt.Fprintf(w, "Get Token miss %s\n", err.Error)
	}
	io.WriteString(w, fmt.Sprintf("%+v\n", res))
}

func main() {
	log.Println("------start server-----")
	http.HandleFunc("/index", index)
	http.HandleFunc("/line", line_login)
	http.HandleFunc("/build", line_login_test)
	http.HandleFunc("/redirect", redirect)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
