package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"

	"github.com/GoSeoTaxi/check_freedomnet/internal/model"
	"github.com/GoSeoTaxi/check_freedomnet/internal/repository"
)

const (
	proxyCheckTimeout = time.Second * 1
)

type FreedomNetService struct {
	Repo       *repository.FreedomNetRepo
	MaxRetries int
	Logger     *zap.Logger
}

func NewFreedomNetService(repo *repository.FreedomNetRepo, maxRetries int, logger *zap.Logger) *FreedomNetService {
	return &FreedomNetService{
		Repo:       repo,
		MaxRetries: maxRetries,
		Logger:     logger,
	}
}

func (s *FreedomNetService) GetFreedomNet() (string, error) {
	var result string
	var err error

	for i := 0; i < s.MaxRetries; i++ {
		result, err = s.Repo.FetchFromServers()
		if err == nil {
			err = s.checkProxy(result)
			if err == nil {
				return result, nil
			}
		}
		s.Logger.Warn("Attempt failed", zap.Int("attempt", i+1), zap.Error(err))
	}

	return "", errors.New("all attempts failed")
}

func (s *FreedomNetService) checkProxy(proxyStr string) error {
	time.Sleep(proxyCheckTimeout)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("Проверяем прокси")
	fmt.Println(proxyStr)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	proxyT, err := parseProxyString(proxyStr)
	if err != nil {
		return fmt.Errorf("не удалось распарсить строку прокси: %w", err)
	}

	var proxyURL *url.URL
	if proxyT.User == "" {
		proxyURL, err = url.Parse(fmt.Sprintf("socks5://%s:%s", proxyT.Host, proxyT.Port))
	} else {
		proxyURL, err = url.Parse(fmt.Sprintf("socks5://%s:%s@%s:%s", proxyT.User, proxyT.Pass, proxyT.Host, proxyT.Port))
	}
	if err != nil {
		return fmt.Errorf("не удалось создать URL прокси: %w", err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return fmt.Errorf("не удалось создать диалер прокси: %w", err)
	}

	restyTransport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := resty.New()
	client.SetTransport(restyTransport)

	// Делаем тестовый запрос через прокси
	resp, err := client.R().Get("https://api.iplocation.net/?cmd=get-ip")
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос через прокси: %w", err)
	}

	if resp.IsSuccess() && len(resp.Body()) > 0 {
		fmt.Printf("Ответ через прокси: %s\n", resp.String())
		return nil
	}

	return fmt.Errorf("прокси не вернул ответа или запрос не успешен")
}

func parseProxyString(proxyStr string) (*model.Proxy, error) {
	hashSplit := strings.Split(proxyStr, "#")
	proxyStr = hashSplit[0]

	protocolSplit := strings.SplitN(proxyStr, "://", 2)
	if len(protocolSplit) != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid proxy string: %s", proxyStr))
	}

	typeProxy := protocolSplit[0]
	proxyStr = protocolSplit[1]

	userInfoSplit := strings.Split(proxyStr, "@")
	var user, pass, hostPortInfo string

	if len(userInfoSplit) == 2 {
		hostPortInfo = userInfoSplit[1]
		userPassSplit := strings.SplitN(userInfoSplit[0], ":", 2)
		if len(userPassSplit) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid user info: %s", userInfoSplit[0]))
		}
		user = userPassSplit[0]
		pass = userPassSplit[1]
	} else {
		hostPortInfo = userInfoSplit[0]
	}

	hostPortSplit := strings.SplitN(hostPortInfo, ":", 2)
	if len(hostPortSplit) != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid host:port info: %s", hostPortInfo))
	}

	host := hostPortSplit[0]
	port := hostPortSplit[1]

	return &model.Proxy{
		Type: typeProxy,
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	}, nil
}
