package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type suggestion struct {
	Text string
	Arc  string
}

type chapter struct {
	Title   string
	Story   []string
	Options []*suggestion
}

var book map[string]chapter

func main() {
	b, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Fatal("while reading gopher.json", err)
	}
	err = json.Unmarshal(b, &book)
	if err != nil {
		log.Fatal("while unmarshal", err)
	}
	//fmt.Println("result\n", book)
	// for key, value := range book {
	// 	fmt.Printf("%q\n %s \n %v \n %v \n\n", key, value.Title, value.Story,
	// 		func() string {
	// 			if len(value.Options) != 0 {
	// 				return value.Options[0].Text
	// 			}
	// 			return fmt.Sprintf("%v", value.Options)
	// 		}())
	// }

	http.HandleFunc("/", homepage)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("while ListenAndServe 8080", err)
	}
}

func homepage(w http.ResponseWriter, r *http.Request) {
	//b, err := json.MarshalIndent(book["intro"], "", "	")
	t := template.Must(template.New("page").Parse(
		func() string {
			b, err := ioutil.ReadFile("template.html")
			if err != nil {
				fmt.Print("ReadFile(\"template.html\")", err)
			}
			return string(b)
		}(),
	))
	// if err != nil {
	// 	fmt.Println("while Marshal(book[\"intro\"])", err)
	// }
	//w.Header().Add("Content-type", "application/html")
	if _, ok := book[strings.TrimLeft(r.URL.Path, "/")]; !ok {
		fmt.Fprint(w, "404 Status Not Found")
		return
	}
	t.Execute(w, book[strings.TrimLeft(r.URL.Path, "/")])
}
