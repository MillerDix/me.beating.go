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
	images := []string{
		"CASSINI_THE_GRAND_FINALE.11b33571.jpg",
		"dawn.56e3457b.jpg",
		"curiosity.f596b6bb.jpg",
		"space.c5f4a997.jpg",
		"sunrise-spacewalk.028439bb.jpg",
	}

	dirs, err := os.Open("./newfiles/")
	if err != nil {
		fmt.Println("failed to open folder: ./newfiles/, quiting")
	}
	defer dirs.Close()
	files, err := dirs.Readdir(0)

	for _, f := range files {
		file, err := os.Open("./newfiles/" + f.Name())
		if err != nil {
			fmt.Println("failed to read file: " + f.Name() + " quiting")
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
		fmt.Println("read file: " + f.Name() + " success")
		
		// incr id for this article
		rs, err := rd.Do("INCR", "next_article_id")
		if err == nil {
			fmt.Println(rs)
		}

		_, err = rd.Do("LPUSH", "articles", rs)
		if err != nil {
			fmt.Println("failed to write data to redis, quiting")
			return
		}

		imageURL := "http://www.beating.io/static/media/" + images[rand.Intn(len(images))]
		_, err = rd.Do(
			"HMSET", fmt.Sprintf("article:%d", rs),
			"Id", rs,
			"Title", f.Name(),
			"Subtitle", fmt.Sprintf("Article %d Subtitle", rs),
			"Summary", summary,
			"Content", textData,
			"Views", rs,
			"Source", fmt.Sprintf("Article %d Source", rs),
			"Poster", imageURL,
		)
		if err != nil {
			fmt.Println("failed to write new article data to redis, quiting")
			return
		}

		fmt.Println("wrote file " + f.Name() + " to redis sccess")

		err = os.Rename("./newfiles/" + f.Name(), "./files/" + f.Name())
		if err != nil {
			fmt.Println("failed to move file: " + f.Name())
			return
		}
		fmt.Println("move file " + f.Name() + " from ./newfiles. to ./files/ success")
	}
}

func main() {
	rd, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("failed to establish connection with redis, quiting")
		return
	}
	fmt.Println("redis connected")
	defer rd.Close()

	newArticle(rd)
	defer fmt.Println("script quit")
	fmt.Println("operation finish")
}
