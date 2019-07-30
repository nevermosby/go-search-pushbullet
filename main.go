package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nevermosby/go-search-pushbullet/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

// const vars
// const (
// 	PB_TOKEN    = ""
// 	PB_PUSH_URL = "https://api.pushbullet.com/v2/pushes"
// )

var (
	cfg = pflag.StringP("config", "c", "", "config file path")
)

func init() {
	// set log config
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	r := setupRouter()
	r.Run(":8080")
	log.Debugln("Started...")
}

// setup the router engine
func setupRouter() *gin.Engine {
	r := gin.Default()

	// ping for test
	r.GET("/ping", func(c *gin.Context) {
		log.Debugln("ping...")
		c.JSON(200, gin.H{
			"code":    0,
			"message": "pongpongpong",
		})
	})

	// search for keyword
	r.GET("/search/:keyword", func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		// entry validation
		if token == viper.GetString("inner_token") {

		} else {
			log.Println("token", token)
			c.JSON(401, gin.H{
				"code":    1,
				"message": "wrong inner token",
			})
		}
		keyword := c.Param("keyword")
		if keyword != "" {
			log.Debugln("keyword", keyword)
			var ret string
			// call the pushbullet api
			// build http client for adding token
			transPort := &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			// create new http client with the defined transport above
			client := &http.Client{Transport: transPort}
			req, err := http.NewRequest("GET", viper.GetString("pb_push_url"), nil)
			if err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
			req.Header.Set("Access-Token", viper.GetString("pb_token"))
			resp, err := client.Do(req)

			if err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
			//log.Println("return body", body)
			// TODO: make the return object as  a struct and unmarshall it to json
			ret = string(body)
			c.JSON(200, gin.H{
				"code":    0,
				"message": ret,
			})
		} else {
			c.JSON(400, gin.H{
				"code":    1,
				"message": "keyword cannot be empty",
			})
		}
	})
	return r
}
