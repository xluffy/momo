package main

import (
	"log"
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

func main() {
	whitelist := "119.17.229.173/32,119.17.229.173/32,101.53.53.147/32,115.79.197.218/32"
	ip := "101.53.53.147"

	isContain, _ := IsIPInCIDRs(ip, strings.Split(strings.TrimSpace(whitelist), ","))
	log.Println("isContain:", isContain)

	a := strings.Split(strings.TrimSpace(whitelist), ",")
	for _, v := range a {
		log.Println("v:", v)
	}
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
