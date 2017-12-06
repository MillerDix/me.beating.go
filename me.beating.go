package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
)

type Article struct {
  Id int
  Title string
  Subtitle string
  Content string
  Views int
  Source string
  Publishtime int
}

var fakeData []Article

func fake_Data() {
  var list []Article
  for i := 0; i < 10; i++ {
    list = append(list, Article{
      Id: 10000 + i,
      Title: fmt.Sprintf("文章%d", i),
      Subtitle: fmt.Sprintf("文章%d副标题", i),
      Content: fmt.Sprintf("文章%d正文内容", i),
      Views: i,
      Source: fmt.Sprintf("文章%d来源", i),
      Publishtime: 1511858529950,
    })
  }
  fakeData = list
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUST, DELETE")
    fn(w, r)
  }
}

func art_all(w http.ResponseWriter, r *http.Request) {
  fmt.Println("list of all articles")
  jsonRes, err := json.Marshal(fakeData)
  if err == nil {
    fmt.Fprint(w, string(jsonRes))
  }
}

func art_categories(w http.ResponseWriter, r *http.Request) {
  fmt.Println("list of all categories")
}

func art_detail(w http.ResponseWriter, r *http.Request) {
  fmt.Println("article detail")
}

func main() {
  rt_main := mux.NewRouter()
  fake_Data()

  // article subrouter
  rt_atcl := rt_main.PathPrefix("/articles").Subrouter()


  rt_atcl.HandleFunc("/all", makeHandler(art_all)).Methods("POST")
  rt_atcl.HandleFunc("/categories", makeHandler(art_categories))
  rt_atcl.HandleFunc("/{category}/{id}", makeHandler(art_detail))

  http.ListenAndServe("beating.io:8080", rt_main)
}
