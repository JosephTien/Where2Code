package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "html/template"
    "regexp"
    "strconv"
)

type RawHandlerFunc http.HandlerFunc

type Article struct {
    Title string
    Body  []byte
}

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
func (a *Article) Save() error {
    filename := config.path["data"] + a.Title + ".txt"
    return ioutil.WriteFile(filename, a.Body, 0600)
}

func LoadArticle(title string) (*Article, error) {
    filename := config.path["data"] + title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Article{Title: title, Body: body}, nil
}

func RenderPage(w http.ResponseWriter, cont string, a *Article) {
	err := pages[cont].ExecuteTemplate(w, "base", a)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

//*******************************************************************************************
//**handler function*************************************************************************
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

func GenHandler_View() RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		title := validPath.FindStringSubmatch(r.URL.Path)[2]
		a, err := LoadArticle(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		RenderPage(w, "view", a)
	}
}

func GenHandler_Edit() RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		title := validPath.FindStringSubmatch(r.URL.Path)[2]
		a, err := LoadArticle(title)
		if err != nil {
			a = &Article{Title: title}
		}
		RenderPage(w, "edit", a)
	}
}

func GenHandler_Save() RawHandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		title := validPath.FindStringSubmatch(r.URL.Path)[2]
		body := r.FormValue("body")
		a := &Article{Title: title, Body: []byte(body)}
		err := a.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
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
	//http.HandleFunc("/", GenHandler("index", &Article{}).Done())
	http.HandleFunc("/", GenHandler_Index().Done())
    http.HandleFunc("/view/", GenHandler_View().Validate().Done())
    http.HandleFunc("/edit/", GenHandler_Edit().Validate().Done())
    http.HandleFunc("/save/", GenHandler_Save().Validate().Done())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.path["static"]))))
    fmt.Println("Server Started on Port ", config.port)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), nil))
}
