package controllers

import (
	"fmt"
	"goproxy/internal/infrastructure/config"
	"goproxy/internal/infrastructure/repositories"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var (
	tagRegex = regexp.MustCompile(
		`<(?:a|abbr|acronym|address|applet|area|audioscope|b|base|basefront|bdo|bgsound|bi
g|blackface|blink|blockquote|body|bq|br|button|caption|center|cite|code|col|colgroup
|comment|dd|del|dfn|dir|div|dl|dt|em|embed|fieldset|fn|font|form|frame|frameset|h
1|head|hr|html|i|iframe|ilayer|img|input|ins|isindex|kdb|keygen|label|layer|legend|li|l
imittext|link|listing|map|marquee|menu|meta|multicol|nobr|noembed|noframes|noscri
pt|nosmartquotes|object|ol|optgroup|option|p|param|plaintext|pre|q|rt|ruby|s|samp|
script|select|server|shadow|sidebar|small|spacer|span|strike|strong|style|sub|sup|tabl
e|tbody|td|textarea|tfoot|th|thead|title|tr|tt|u|ul|var|wbr|xml|xmp)\\W`,
	)
	replaceIp = regexp.MustCompile("[^0-9]")
)

type ProxyController struct {
	repo  *repositories.BlockedIpsRepository
	log   *zap.Logger
	proxy *httputil.ReverseProxy
	cfg   config.AppConfig
	host  string
}

func NewProxyController(
	repo *repositories.BlockedIpsRepository,
	log *zap.Logger,
	cfg config.AppConfig,
) (*ProxyController, error) {
	url, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return &ProxyController{repo: repo, log: log, proxy: proxy, cfg: cfg, host: url.Host}, nil
}

func (p ProxyController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ips, err := p.repo.ListBlockedIPs()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if ips[replaceIp.ReplaceAllString(ip, "")] {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if strings.HasPrefix(r.Host, "www.") {
		r.URL.Path = "/site/www" + r.URL.Path
	}
	fmt.Println(r.URL.RawQuery)
	if tagRegex.MatchString(r.URL.RawQuery) || checkBody(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = p.host
	p.proxy.ServeHTTP(w, r)
}

func checkBody(r *http.Request) bool {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		return tagRegex.Match(body)
	}
	return false
}
