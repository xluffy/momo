package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/yl2chen/cidranger"

	"github.com/gin-gonic/gin"
)

func main() {
	whitelist := os.Getenv("WHITELIST")
	log.Println("whitelist:", whitelist)
	log.Println("whitelist type:", reflect.TypeOf(whitelist))
	list := strings.Split(strings.TrimSpace(whitelist), ",")

	for _, v := range list {
		log.Println("whitelist item:", v)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"client":  c.ClientIP(),
		})
	})

	r.GET("/check", func(c *gin.Context) {
		ip := c.Query("ip")
		isContain, _ := IsIPInCIDRs(ip, strings.Split(strings.TrimSpace(whitelist), ","))

		c.JSON(http.StatusOK, gin.H{
			"message": isContain,
			"client":  c.ClientIP(),
		})
	})

	r.Run()
}

func IsIPInCIDRs(ip string, cidrs []string) (bool, error) {
	ranger := cidranger.NewPCTrieRanger()
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Println("ParseCIDR error:", err)
		}

		if err := ranger.Insert(cidranger.NewBasicRangerEntry(*ipnet)); err != nil {
			log.Println("Insert error:", err)
		}
	}

	isContain, err := ranger.Contains(net.ParseIP(ip))
	if err != nil {
		log.Println("Contains error:", err)
	}

	return isContain, nil
}
