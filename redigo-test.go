package main

import (
	"fmt"
	"os"
	"bufio"
	"math/rand"

	"github.com/garyburd/redigo/redis"
)

// write new article data to redis
// Title: filename
// Content: ~
// Filepath: file should be in the same folder with this script
// redis key: article:1511858529959
func newArticle(rd redis.Conn) {
	fileTitle := "什么是CORS"
	file, err := os.Open("./files/" + fileTitle + ".md")
	if err != nil {
		fmt.Println("os open file failed")
		return
	}

	var textData string
	var summary string
	sumIndex := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		textData = textData + text + "\n"
		if sumIndex < 5 {
			summary = summary + text
			sumIndex++
		}
	}

	_, er := rd.Do("HMSET", "article:1511858529959", "Content", textData, "Title", fileTitle, "Summary", summary)
	if er != nil {
		fmt.Println("write new article data to redis failed, quiting")
		return
	}
}

// fake_redis_data
// Id, Title, Subtitle, Content, Views, Source, Poster, Publishtime
// num: 10
// unix time stamp: 1511858529950 to 1511858529959
// images: random in
// 		CASSINI_THE_GRAND_FINALE.jpg'
// 		dawn.jpg'
// 		curiosity.jpg'
// 		space.jpg'
// 		sunrise-spacewalk.jpg'
func fake_redis_data(rd redis.Conn) {
	images := []string{
		"CASSINI_THE_GRAND_FINALE.11b33571.jpg",
		"dawn.56e3457b.jpg",
		"curiosity.f596b6bb.jpg",
		"space.c5f4a997.jpg",
		"sunrise-spacewalk.028439bb.jpg",
	}
	for i := 0; i < 10; i++ {
		imageURL := "http://www.beating.io/static/media/" + images[rand.Intn(len(images))]
		_, err := rd.Do("HMSET",
			fmt.Sprintf("article:%d", 1511858529950+i),
			"Id",
			1511858529950+i,
			"Title",
			fmt.Sprintf("Article:%d", i),
			"Subtitle",
			fmt.Sprintf("ArticleSubtitle:%d", i),
			"Content",
			fmt.Sprintf("ArticleContent:%d", i),
			"Views",
			i,
			"Source",
			fmt.Sprintf("ArticleSource:%d", i),
			"Poster",
			imageURL,
			"Publishtime",
			1511858529950+i,
		)
		_, err = rd.Do("LPUSH", "articles", 1511858529950+i)
		if err != nil {
			fmt.Println("writer fake data to redis failed, quiting")
			return
		}
	}
}

func main() {
	rd, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis connection failed, quiting")
		return
	}
	defer rd.Close()

	newArticle(rd)
	// fake_redis_data(rd)
	fmt.Println("operation finish")
}
