package main

import (
	"fmt"
	"os"
	"bufio"
	// "math/rand"
	// "time"
	// "path/filepath"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func main() {
	images := []string{
		"http://res.cloudinary.com/millerd/image/upload/v1515493291/Beatinglog/home/CASSINI_THE_GRAND_FINALE_apg5yb.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515493291/Beatinglog/home/dawn_ci59ps.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515493344/Beatinglog/home/curiosity_a4g6je.jpg",
		"http://res.cloudinary.com/millerd/image/upload/v1515494236/Beatinglog/home/space_bhrgqw.jpg",
		"http://res.cloudinary.com/millerd/image/upload/c_scale,q_auto,w_1200/v1515493300/Beatinglog/home/sunrise-spacewalk-png8_qhnw0k.png",
	}
	// connect with redis
	rd, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("failed to establish connection with redis, quiting")
		return
	}
	fmt.Println("redis connected")
	defer rd.Close()

	// detect file
	dirs, err := os.Open("./newfiles/")
	if err != nil {
		fmt.Println("failed to open folder ./newfiles")
		return
	}
	defer dirs.Close()

	files, err := dirs.Readdir(0)
	if err != nil {
		fmt.Println("failed to read file in folder")
		return
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		file, err := os.Open("./newfiles/" + f.Name())
		if err != nil {
			fmt.Println("failed to open file ./newfiles/ " + f.Name())
		}

		scanner := bufio.NewScanner(file)
		var textData string
		var summaryData string
		var count = 0
		for scanner.Scan() {
			text := scanner.Text()
			textData = textData + text + "\n"
			if count < 5 {
				summaryData += text
				count++
			}
		}
		var id_to_update int
		fmt.Println("read file 《"+ f.Name() + "》 success, enter article id to update: ")
		fmt.Scanln(&id_to_update)
		imageURL := images[rand.Intn(len(images))]
		_, err = rd.Do(
			"HMSET", fmt.Sprintf("article:%d", id_to_update),
			"Title", strings.TrimSuffix(f.Name(), ".md"),
			"Summary", summaryData,
			"Content", textData,
			"Poster", imageURL,
		)
		if err != nil {
			fmt.Println("failed to update data in redis")
			return
		}

		err = os.Rename("./newfiles/" + f.Name(), "./files/" + f.Name())
		if err != nil {
			fmt.Println("failed to move file: " + f.Name())
			return
		}

		fmt.Println("script quiting")
	}
}