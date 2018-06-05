package main

import (
    "fmt"
    "log"
    "net/http"
    "html/template"
    "regexp"
	"strconv"
)

type RawHandlerFunc http.HandlerFunc

type Config struct{
    port int
    path map[string]string
}

var(
    config = Config{
        port : 8080,
        path : map[string]string{
            "server" : "server/",
            "static" : "client/static/",
            "view"   : "client/view/",
            "data"   : "data/",
        },
    }
    htmls = []string{
		config.path["view"]+"base.html",
		config.path["view"]+"index.html",
		config.path["view"]+"editor.html",
    }
    validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	pages = map[string]*template.Template{
		"index" : template.Must(template.ParseFiles(htmls[0],htmls[1],htmls[2])),
	}
)
//*******************************************************************************************
//**handler function*************************************************************************

func RenderPage(w http.ResponseWriter, cont string, a *Article) {
	err := pages[cont].ExecuteTemplate(w, "base", a)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func GenHandler(cont string, a *Article)RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
        RenderPage(w, cont, a)
    }
}

func GenHandler_Index() RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl,err := template.ParseFiles(htmls[0],htmls[1],htmls[2])
			if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "base", &Article{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GenHandler_Save() RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		title := validPath.FindStringSubmatch(r.URL.Path)[2]
		body := r.FormValue("body")
    fmt.Println(body)
		article := &Article{Title: title, Content: []byte(body)}
		err := SaveArticle(article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (Fn RawHandlerFunc)Validate()RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
		Fn(w, r)
    }

}

func (Fn RawHandlerFunc)Done()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		Fn(w,r)
	}
}

//*******************************************************************************************
func main() {
  InitSocket()
  initFileSystem();
  //http.HandleFunc("/", GenHandler("index", &Article{}).Done())
  http.HandleFunc("/", GenHandler_Index().Done())
  http.HandleFunc("/save/", GenHandler_Save().Done())
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.path["static"]))))
  fmt.Println("Server Started on Port ", config.port)
  log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), nil))
}
