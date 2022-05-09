package main

import (
	"GO-2/cotroller"
	"GO-2/repository"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := Init("./data/"); err != nil {
		os.Exit(100)
	}
	r := gin.Default()
	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := cotroller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})
	err := r.Run()
	if err != nil {
		return
	}
}

func Init(filePath string) error {
	if err := repository.Init(filePath); err != nil {
		return err
	}
	return nil
}
