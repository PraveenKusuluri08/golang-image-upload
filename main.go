package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Praveenkusuluri08/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func init() {
	helpers.EnvInitilizer()
}

func main() {
	helpers.EnvInitilizer()
	router := gin.Default()
	Region := os.Getenv("AWS_REGION")
	conf := aws.Config{Region: aws.String(Region)}
	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)
	router.Static("/web/assets", "./web/assets")
	router.LoadHTMLGlob("web/templates/*")
	router.MaxMultipartMemory = 8 << 20

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/uploadtodisc", func(c *gin.Context) {
		file, err := c.FormFile("image")
		fmt.Println(file.Filename)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload file",
			})
			return
		}

		fmt.Println(file.Filename)

		f, e := file.Open()
		if e != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload file",
			})
			return
		}

		// uploadToDisc,uploadError:= c.SaveUploadedFile(file,"/web/assests/uploads")
		// if uploadError != nil {
		// 	c.JSON(http.StatusOK,gin.H{
		// 		"error":"Failed to upload file"
		// 	})
		// }

		result, uploadError := svc.Upload(&s3manager.UploadInput{
			Bucket: aws.String("golang-image-upload"),
			Key:    aws.String(file.Filename),
			Body:   f,
		})
		defer f.Close()

		if uploadError != nil {
			log.Fatal(uploadError)
		}
		c.JSON(http.StatusOK, gin.H{
			"Status": result.Location,
		})
	})

	PORT := os.Getenv("PORT")
	fmt.Println(PORT)
	if PORT == "" {
		PORT = "8080"
	}

	router.Run(":" + PORT)
}
