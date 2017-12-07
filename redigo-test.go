package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func main() {
	rd, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis connection failed, quiting")
		return
	}
	defer rd.Close()
	//
	for i := 1; i < 10; i++ {
		rd.Do("HMSET",
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
			"poster",
			"http://www.beating.io/static/media/CASSINI_THE_GRAND_FINALE.11b33571.jpg",
			"Publishtime",
			1511858529950+i,
		)
	}

}
