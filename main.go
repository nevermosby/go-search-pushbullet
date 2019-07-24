package main

import (
	"log"
	"github.com/nevermosby/"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		log.Println("ping...")
		// TODO: make the return object as  a struct and unmarshall it to json
		rsp := CreateResponse{
			Username: r.Username,
		}
	
		c.JSON(200, gin.H{
			"code":  0,
			"message": "pongpongpong",
		})
	})
	r.GET("/search/:keyword",func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		log.Println("token",token)
		if token != nil {

		}

		c.JSON(401, gin.H{
			"code": 1,
			"message": "",
		})
		name

	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
	log.Println("Started...")
}

//package main
//
//import "fmt"
//import "os"
//import "os/signal"
//import "syscall"
//
//func main() {
//
//	// Go signal notification works by sending `os.Signal`
//	// values on a channel. We'll create a channel to
//	// receive these notifications (we'll also make one to
//	// notify us when the program can exit).
//	sigs := make(chan os.Signal, 1)
//	done := make(chan bool, 1)
//
//	// `signal.Notify` registers the given channel to
//	// receive notifications of the specified signals.
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//	// This goroutine executes a blocking receive for
//	// signals. When it gets one it'll print it out
//	// and then notify the program that it can finish.
//	go func() {
//		sig := <-sigs
//		fmt.Println()
//		fmt.Println(sig)
//		done <- true
//	}()
//
//	// The program will wait here until it gets the
//	// expected signal (as indicated by the goroutine
//	// above sending a value on `done`) and then exit.
//	fmt.Println("awaiting signal")
//	<-done
//	fmt.Println("exiting")
//}
