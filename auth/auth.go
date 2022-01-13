package auth

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bugfan/de"
)

var (
	ServerIP string
	AuthURL  string
	defaults map[string]string
)

func init() {
	defaults = map[string]string{
		"auth_url":     "https://wg.lt53.cn",
		"api_secret":   "jwdlhtrh",
		"ifconfig_url": "http://ifconfig.co",
	}
	de.SetKey(Get("api_secret"))
	de.SetExp(100)
	AuthURL = Get("auth_url")
	ServerIP = getIp()
}
func getIp() string {
	resp, err := http.Get(Get("ifconfig_url"))
	if err != nil {
		fmt.Println("get server ip error:", err)
		return ""
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return strings.Replace(string(data), "\n", "", 1)
}
func getDefault(key string) string {
	return defaults[key]
}

func Get(key string) string {
	env := strings.TrimSpace(os.Getenv(strings.ToUpper(key)))
	if env != "" {
		return env
	}
	return getDefault(key)
}

type Peer struct {
	Address   string
	PublicKey string
}

func Verify(clientPublicKey string) (*Peer, error) {
	bearer, _ := de.EncodeWithBase64()
	header := make(map[string]string)

	header["Wgkey"] = clientPublicKey
	header["Wgtoken"] = string(bearer)
	header["Wgserverip"] = ServerIP

	code, data, err := Request("GET", AuthURL+"/wireguard", header, nil)
	if code >= 300 || err != nil {
		return nil, errors.New(fmt.Sprintf("auth:request to auth server wireguard error:code is %v,error is %v\n", code, err))
	}
	peer := &Peer{}
	json.Unmarshal(data, peer)

	if clientPublicKey != peer.PublicKey {
		fmt.Println("err equals:", clientPublicKey, peer.PublicKey)
		return nil, errors.New("key not equals")
	}

	return peer, nil
}

type Config struct {
	ListenPort string `json:"wg_listen_port"`
	PrivateKey string `json:"wg_private_key"`
}

func GetWireguardConfig() (*Config, error) {
	bearer, _ := de.EncodeWithBase64()
	header := make(map[string]string)

	header["Wgtoken"] = string(bearer)
	header["Wgserverip"] = ServerIP

	code, data, err := Request("GET", AuthURL+"/config", header, nil)
	if code >= 300 || err != nil {
		errString := fmt.Sprintf("auth:request to auth server config error:code is %v,error is %v\n", code, err)
		fmt.Println(errString)
		return nil, errors.New(errString)
	}
	conf := &Config{}
	json.Unmarshal(data, conf)
	conf.PrivateKey = KeyToHex(conf.PrivateKey)
	return conf, nil
}
func atob(data []byte) []byte {
	// Base64 Standard Decoding
	sDec, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return []byte{}
	}
	return sDec
}
func KeyToHex(key string) string {
	data := atob([]byte(key))
	return hex.EncodeToString(data)
}

var (
	client *http.Client
)

func init() {
	client = new(http.Client)
	netTransport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(20))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		DisableKeepAlives:     true,
		MaxIdleConnsPerHost:   20,                              //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(60), //数据收发5秒超时
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client.Timeout = time.Second * 30
	client.Transport = netTransport
}
func NewHttpClient() *http.Client {
	return client
}
func Request(method, target string, headers map[string]string, body io.ReadCloser) (int, []byte, error) {
	req, _ := http.NewRequest(method, target, body)
	req.Header.Add("cache-control", "no-cache")
	req.Close = true
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	cli := NewHttpClient()
	res, err := cli.Transport.RoundTrip(req)
	if err != nil {
		return -1, nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, data, err
}
