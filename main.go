package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nevermosby/go-search-pushbullet/config"
	pb "github.com/nevermosby/go-search-pushbullet/pushbullet"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
	"time"
)

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

func fetchPushes() *pb.Pushes {
	var cursor string
	var retPushes pb.Pushes
	var pushesResponse pb.Pushes

	// build the http client
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

	// loop
	loop := viper.GetInt("pb_push_loop")
	for i := 0; i <= loop; i++ {
		log.Debugln("enter loop:", i)
		log.Debugln("current cursor:", cursor)
		if i > 0 {
			time.Sleep(2 * time.Second)
		}
		q := req.URL.Query()
		q.Del("cursor")
		q.Add("cursor", cursor)
		req.URL.RawQuery = q.Encode()
		resp, err := client.Do(req)

		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&pushesResponse)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Debugln("new cursor:", pushesResponse.Cursor)
		cursor = pushesResponse.Cursor
		retPushes.PushItems = append(retPushes.PushItems, pushesResponse.PushItems...)
		log.Debugln("return pushes size:", len(retPushes.PushItems))

	}

	return &retPushes
}

// checkInnerToken for check inner token as entry
func checkInnerToken(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	// entry validation
	if token == viper.GetString("inner_token") {
		return true
	}
	log.Println("token", token)
	c.JSON(401, gin.H{
		"code":    1,
		"message": "wrong inner token",
	})
	return false
}

// findMatchedPush for return the matched push item with keyword
func findMatchedPush(pushes *pb.Pushes, keyword string) *pb.Pushes {
	var retPushes pb.Pushes
	for _, ele := range pushes.PushItems {
		if strings.Contains(strings.ToLower(ele.Body), keyword) || strings.Contains(strings.ToLower(ele.Title), keyword) || strings.Contains(strings.ToLower(ele.URL), keyword) {
			retPushes.PushItems = append(retPushes.PushItems, ele)
		}
	}
	return &retPushes
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

	r.GET("/pb/me", func(c *gin.Context) {
		ret := checkInnerToken(c)
		if !ret {
			return
		}
		// call the pushbullet api
		// build http client for adding token
		transPort := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// create new http client with the defined transport above
		client := &http.Client{Transport: transPort}
		req, err := http.NewRequest("GET", viper.GetString("pb_me_url"), nil)
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
		var userResponse pb.User
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&userResponse)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Debugln("user responose: %v", &userResponse)
		c.JSON(200, gin.H{
			"code":    0,
			"message": &userResponse,
		})
	})

	// search for keyword
	r.GET("/pb/search/:keyword", func(c *gin.Context) {
		ret := checkInnerToken(c)
		if !ret {
			return
		}
		keyword := c.Param("keyword")
		if keyword != "" {
			log.Debugln("keyword", keyword)
			pushesResponse := fetchPushes()
			log.Debugln("push responose:", len(pushesResponse.PushItems))
			retPushesItems := findMatchedPush(pushesResponse, keyword)
			c.JSON(200, gin.H{
				"code":    0,
				"message": retPushesItems,
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
