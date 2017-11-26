package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
			b, err := ioutil.ReadAll(file)
			if err != nil {
				c.JSON(500, "Failed to read file")
			}
			filename := strconv.FormatInt(time.Now().Unix(), 16) + ".wav"
			ioutil.WriteFile("./tmp/"+filename, b, 0755)
			cmd := exec.Command("/home/fronbasal/.env/tensorflow/bin/python3", "/home/fronbasal/Desktop/dialektor-tf/predict.py", "/home/fronbasal/go/src/github.com/fronbasal/dialektor-web/tmp/"+filename)
			var out bytes.Buffer
			cmd.Stderr = &out
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				c.JSON(500, err.Error()+"\n\n"+out.String())
				return
			}
			x := strings.Split(out.String(), "\n")
			s := x[len(x)-2]
			c.JSON(200, strings.Split(s, " ")[1])
		})
	}
	r.Run(":5000")
}
