package gosparkclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// assembleAuthURL generates the authenticated URL for API requests
func (c *SparkClient) assembleAuthURL(httpMethod string, hostURL string) string {
	ul, err := url.Parse(hostURL)
	if err != nil {
		panic(fmt.Sprintf("invalid URL: %v", err))
	}

	date := time.Now().UTC().Format(time.RFC1123)
	signString := []string{
		"host: " + ul.Host,
		"date: " + date,
		httpMethod + " " + ul.Path + " HTTP/1.1",
	}

	signature := hmacSha256ToBase64(strings.Join(signString, "\n"), c.config.ApiSecret)
	authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
		"hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
		c.config.ApiKey,
		"hmac-sha256",
		"host date request-line",
		signature,
	)))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)

	return hostURL + "?" + v.Encode()
}

// hmacSha256ToBase64 generates an HMAC-SHA256 signature and returns it as base64
func hmacSha256ToBase64(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// defaultTransport creates a default HTTP transport with the given timeout
func defaultTransport(timeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
	}
}

// readBody reads and returns the response body as a string
func readBody(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error reading response body: %v", err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
