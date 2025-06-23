package netbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/holonet/core/logger"
)

type NetBoxStatusResponse struct {
	NetboxVersion string                 `json:"netbox-version"`
	Plugins       map[string]interface{} `json:"plugins"`
	DjangoVersion string                 `json:"django-version"`
	InstalledApps map[string]string      `json:"installed-apps"`
}

type Client struct {
	Host               string
	Token              string
	Client             *http.Client
	hasLoggedAvailable bool
	hasLoggedMutex     sync.Mutex
}

func New(httpClient *http.Client) (*Client, error) {
	host := os.Getenv("NETBOX_HOST")
	if host == "" {
		return nil, fmt.Errorf("NETBOX_HOST environment variable is not set")
	}

	token := os.Getenv("NETBOX_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("NETBOX_API_TOKEN environment variable is not set")
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	client := &Client{
		Host:               host,
		Token:              token,
		Client:             httpClient,
		hasLoggedAvailable: false,
	}

	logger.Info("NetBox client created for host: %s", host)
	return client, nil
}

func (c *Client) CheckConnection() bool {
	logger.Info("Checking NetBox connection for host: %s", c.Host)
	return c.IsAvailable()
}

func (c *Client) IsAvailable() bool {
	host := c.Host
	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	url := fmt.Sprintf("%s/api/status/", host)
	logger.Debug("Checking NetBox availability and version at URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Failed to create request for %s: %v", url, err)
		return false
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Token "+c.Token)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Error("NetBox availability check failed for %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	logger.Debug("NetBox availability check response status: %d for %s", resp.StatusCode, url)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body: %v", err)
			return false
		}

		var statusResponse NetBoxStatusResponse
		if err := json.Unmarshal(body, &statusResponse); err != nil {
			logger.Error("Failed to parse NetBox status response: %v", err)
			return false
		}

		requiredVersion := "4.3.2"
		if statusResponse.NetboxVersion != requiredVersion {
			logger.Error("NetBox version mismatch: got %s, required %s", statusResponse.NetboxVersion, requiredVersion)
			return false
		}

		requiredApps := map[string]string{
			"django_filters":          "25.1",
			"django_prometheus":       "2.3.1",
			"django_rq":               "3.0.1",
			"django_tables2":          "2.7.5",
			"drf_spectacular":         "0.28.0",
			"drf_spectacular_sidecar": "2025.6.1",
			"mptt":                    "0.17.0",
			"rest_framework":          "3.16.0",
			"social_django":           "5.4.3",
			"taggit":                  "6.1.0",
			"timezone_field":          "7.1",
		}

		for app, requiredVersion := range requiredApps {
			installedVersion, ok := statusResponse.InstalledApps[app]
			if !ok {
				logger.Error("Required app %s is not installed", app)
				return false
			}
			if installedVersion != requiredVersion {
				logger.Error("App %s version mismatch: got %s, required %s", app, installedVersion, requiredVersion)
				return false
			}
		}

		c.hasLoggedMutex.Lock()
		if !c.hasLoggedAvailable {
			logger.Info("NetBox is available and token is valid at %s", url)
			c.hasLoggedAvailable = true
		} else {
			logger.Debug("NetBox is available and token is valid at %s", url)
		}
		c.hasLoggedMutex.Unlock()

		return true
	} else if resp.StatusCode == 401 || resp.StatusCode == 403 {
		logger.Error("NetBox token authentication failed for %s: status code %d", url, resp.StatusCode)
		return false
	}

	url = fmt.Sprintf("%s/api/", host)
	logger.Debug("Trying fallback NetBox availability check at URL: %s", url)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Failed to create request for fallback %s: %v", url, err)
		return false
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Token "+c.Token)
	}

	resp, err = c.Client.Do(req)
	if err != nil {
		logger.Error("NetBox fallback availability check failed for %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	logger.Debug("NetBox fallback availability check response status: %d for %s", resp.StatusCode, url)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Warn("NetBox is available but couldn't check version and installed apps")
		return true
	}

	logger.Error("NetBox is unavailable: all endpoints failed")
	return false
}
