package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type article struct {
	ID          int
	Title       string
	Subtitle    string
	Content     string
	Views       int
	Source      string
	Publishtime int
	Poster      string
	Summary     string
}

var conn redis.Conn

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		if r.Method == "OPTIONS" {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fn(w, r)
	}
}

func artAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list of all articles")
	artcIds, err := redis.Strings(conn.Do("LRANGE", "articles", "0", "10"))
	if err != nil {
		return
	}

	var list []article
	for i := 0; i < len(artcIds); i++ {
		var tmp article
		artc, _ := redis.StringMap(conn.Do("HGETALL", fmt.Sprintf("article:%s", artcIds[i])))

		// TODO: THIS IS UGLY
		intid, _ := strconv.Atoi(artc["ID"])
		intviews, _ := strconv.Atoi(artc["Views"])
		intpublishtime, _ := strconv.Atoi(artc["Publishtime"])

		tmp.ID = intid
		tmp.Title = artc["Title"]
		tmp.Subtitle = artc["Subtitle"]
		tmp.Views = intviews
		tmp.Source = artc["Source"]
		tmp.Publishtime = intpublishtime
		tmp.Poster = artc["Poster"]
		tmp.Summary = artc["Summary"]

		list = append(list, tmp)
	}
	jsonRes, err := json.Marshal(list)
	if err == nil {
		fmt.Fprint(w, string(jsonRes))
	}
}

func artCategories(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list of all categories")
}

func artDetail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("article detail")
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var body map[string]interface{}
	err := json.Unmarshal(data, &body)

	if err == nil {
		rs, _ := redis.StringMap(conn.Do("HGETALL", "article:"+strconv.FormatFloat(body["ID"].(float64), 'f', -1, 64)))
		jsonres, _ := json.Marshal(rs)
		fmt.Fprint(w, string(jsonres))
	}
}

func main() {
	rtMain := mux.NewRouter()

	// redis
	var err error
	conn, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("can not establish connection to redis, quiting")
		return
	} else {
		fmt.Println("redis connected")
	}
	defer conn.Close()

	// article subrouter
	rtAtcl := rtMain.PathPrefix("/blog").Subrouter()
	rtAtcl.HandleFunc("/all", makeHandler(artAll))
	rtAtcl.HandleFunc("/categories", makeHandler(artCategories))
	rtAtcl.HandleFunc("/detail", makeHandler(artDetail))

	http.ListenAndServe(":8080", rtMain)
}
