package main

import (
    "io/ioutil"
    "fmt"
)

type Article struct {
  Title string `json:"title"`
  Content []byte `json:"content"`
  Author string `json:"author"`
  Version string `json:"version"`
  Saved bool `json:"saved"`
  Locked bool `json:"locked"`
}

var articleList map[string]*Article;

func initFileSystem(){
    fmt.Println("Init FileSystem")
    articleList = make(map[string]*Article)
}

func GetContent(title string)(string, error){
  article, inList := articleList[title];
	if !inList{
    article = (&Article{Title:title, Saved:true})
		_, err := article.Load();
		if err != nil {
			return "", err
		}
		articleList[title] = article
	}
	return string(article.Content), nil;
}

func SaveArticle(article *Article)error{
  //article.Saved=false;
  articleList[article.Title] = article
  return nil
}

func SaveContent(title string, content string)error{
  article, ok := articleList[title];
  if !ok {
    article = &Article{Title:title, Saved:false}
  }
  article.Content = []byte(content);
  articleList[title] = article
  return nil
}

func (a *Article) Save() error {
    filename := config.path["data"] + a.Title
    return ioutil.WriteFile(filename, a.Content, 0600)
}

func (a *Article) Load() (*Article, error) {
    filename := config.path["data"] + a.Title
    content, err := ioutil.ReadFile(filename)
    a.Content = content
    if err != nil {
        return nil, err
    }
    return a, nil
}
