package main

import (
	_ "embed"
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

// PreferIPv4 表示在DNS服务器的域名转换为IP地址过程中优先返回IPv4
var PreferIPv4 = true
var LogFormat = 0            // 默认 0:"text",可选 1:"json"
var LogLevel = log.InfoLevel // 默认 4: InfoLevel，可选 0-6

var GeoDB *geoip2.Reader

var serversPath string

func initLog() {
	if LogFormat == 0 {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05", // 简短时间格式
			DisableQuote:    true,                  // 禁用字段引号
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetLevel(log.Level(LogLevel))
}

func initGeoDB() {
	var err error
	GeoDB, err = GetGeoData()
	if err != nil {
		log.Fatal("无法打开GeoIP数据库:", err)
	}
}

func main() {
	initLog()
	initGeoDB()
	defer GeoDB.Close()

	testCheck()
}

func testCheck() {
	servers := []string{
		"1.1.1.1:53",
		"114.114.114.114",
		"192.168.1.1",
		"119.29.29.29",
		"2402:4e00::",
		"https://dns.google/dns-query",
		"tls://dns.cloudflare.com",
		"quic://dns.google:853",
		"https://dns.goooogle/dns-query",
	}
	CheckServers(servers)
}

func CheckServers(servers []string) {
	for _, server := range servers {
		ip, geoCode, err := CheckIPAndGeo(server)
		if err != nil {
			log.WithFields(log.Fields{
				"目标": server,
				"错误": err,
			}).Errorf("\x1b[31m解析 %s 失败\x1b[0m", server)
			continue
		}
		log.WithFields(log.Fields{
			"目标": server,
			"IP": ip,
			"代码": geoCode,
		}).Infof("\x1b[32m解析 %s 成功\x1b[32m", server)

	}
}

func CheckIPAndGeo(_server string) (string, string, error) {
	server := strings.TrimSpace(_server)
	server = strings.TrimSuffix(server, "/")
	server = strings.TrimSuffix(server, ":53")
	server = strings.TrimSuffix(server, ":853")
	if server == "" {
		return "0.0.0.0", "PRIVATE", fmt.Errorf("服务器地址为空")
	}
	if !strings.Contains(server, "://") {
		// 处理普通IP
		ip := net.ParseIP(server)
		if ip == nil {
			return "0.0.0.0", "PRIVATE", fmt.Errorf("无效的IP地址")
		}
		geoCode, err := checkIPGeo(ip)
		return ip.String(), geoCode, err
	}
	// 处理加密DNS地址
	server = strings.TrimSuffix(server, "/dns-query")
	server = strings.TrimPrefix(server, "https://")
	server = strings.TrimPrefix(server, "tls://")
	server = strings.TrimPrefix(server, "quic://")

	ips, err := net.LookupIP(server)
	ipc := len(ips)
	if err != nil || ipc == 0 {
		// 无法解析IP地址
		return "0.0.0.0", "PRIVATE", fmt.Errorf("无法解析IP地址")
	}
	if ipc == 1 {
		// 只有一个IP地址
		geoCode, err := checkIPGeo(ips[0])
		return ips[0].String(), geoCode, err
	}
	// 多个 IP 地址
	var oneIP net.IP
	if PreferIPv4 {
		for _, ip := range ips {
			if ip.To4() != nil {
				oneIP = ip
				break
			}
		}
		if oneIP == nil {
			oneIP = ips[0]
		}
	} else {
		oneIP = ips[0]
	}
	geoCode, err := checkIPGeo(oneIP)
	return oneIP.String(), geoCode, err
}

func checkIPGeo(ip net.IP) (string, error) {
	record, err := GeoDB.Country(ip)
	if err != nil {
		return "CDN", err
	}
	return record.Country.IsoCode, nil
}
