package main

import (
	"io/ioutil"
	"os"
	"time"

	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("ui/*.html")

	r.Static("/assets", "assets")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	api := r.Group("/api/v1")
	{
		api.POST("/predict", func(c *gin.Context) {
			data, err := c.FormFile("file")
			if err != nil {
				c.JSON(500, "Failed to get data")
			}
			file, err := data.Open()
			if err != nil {
				c.JSON(500, "Failed to open file")
			}
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				c.JSON(500, "Failed to read file")
			}
			filename := time.Now().String() + ".wav"
			ioutil.WriteFile("./tmp/"+filename, bytes, 0755)

			cmd := exec.Command("python3", "~/dialektor-tensorflow/predict.py", "~/dialektor-web/tmp/"+filename)
			if cmd.Run() != nil {
				os.Remove("./tmp/" + filename)
				c.JSON(500, "Failed to run command! ")
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				os.Remove("./tmp/" + filename)
				c.JSON(500, "Failed to get output of CMD!")
			}
			c.JSON(200, string(output[:]))
		})
	}
	r.Run(":5000")
}
