package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// write new article data to redis
// Title: filename
// Content: ~
// Filepath: file should be in the same folder with this script
// redis key: article:1511858529959
// const cloudinaries = {
//   Cassini: 'http://res.cloudinary.com/millerd/image/upload/v1515493291/CASSINI_THE_GRAND_FINALE_apg5yb.jpg',
//   Dawn: 'http://res.cloudinary.com/millerd/image/upload/v1515493291/500px210695375_ci59ps.jpg',
//   Curiosity: 'http://res.cloudinary.com/millerd/image/upload/v1515493344/curiosity_a4g6je.jpg',
//   SunriseSpacewalk: 'http://res.cloudinary.com/millerd/image/upload/v1515493300/sunrise-spacewalk-png8_qhnw0k.png'
// }

// export default cloudinaries;
func newArticle(rd redis.Conn) {
	images := []string{
		"http://res.cloudinary.com/millerd/image/upload/v1515493291/Beatinglog/home/CASSINI_THE_GRAND_FINALE_apg5yb.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515493291/Beatinglog/home/dawn_ci59ps.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515493344/Beatinglog/home/curiosity_a4g6je.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515494236/Beatinglog/home/space_bhrgqw.jpg",
		"http://res.cloudinary.com/millerd/image/upload/c_scale,q_auto,w_1200/v1515493300/Beatinglog/home/sunrise-spacewalk-png8_qhnw0k.png",
	}

	dirs, err := os.Open("./newfiles/")
	if err != nil {
		fmt.Println("failed to open folder: ./newfiles/, quiting")
	}
	defer dirs.Close()
	files, err := dirs.Readdir(0)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
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

		imageURL := images[rand.Intn(len(images))]
		t := time.Now()
		_, err = rd.Do(
			"HMSET", fmt.Sprintf("article:%d", rs),
			"ID", rs,
			"Title", strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
			"Subtitle", fmt.Sprintf("Article %d Subtitle", rs),
			"Summary", summary,
			"Content", textData,
			"Views", rs,
			"Source", fmt.Sprintf("Article %d Source", rs),
			"Poster", imageURL,
			"Publishtime", t.Unix()*1000,
		)
		if err != nil {
			fmt.Println("failed to write new article data to redis, quiting")
			return
		}

		fmt.Println("wrote file " + f.Name() + " to redis success")

		err = os.Rename("./newfiles/"+f.Name(), "./files/"+f.Name())
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
