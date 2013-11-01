package deviceatlas

import "github.com/ernesto-jimenez/ch_deviceatlas/proxy"
import "net/http"
import "net/url"
import "encoding/json"
import "io/ioutil"
import "strconv"
import "fmt"

const deviceatlas_url = "http://region0.deviceatlascloud.com/v1/detect/properties?licencekey=%s&useragent=%s"

type DeviceInfo struct {
  DPR string
  RW string
}

type daProperties struct {
  DisplayWidth int
  DisplayPpi int
  DevicePixelRatio string
}

type daReply struct {
  Properties daProperties
}

func GetInfo(useragent string, key string) DeviceInfo {
  if useragent == "" {
    return DeviceInfo{}
  }
  useragent = url.QueryEscape(useragent)
  url := fmt.Sprintf(deviceatlas_url, key, useragent)

  res, _ := http.Get(url)
  //defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)

  var data daReply
  _ = json.Unmarshal(body, &data)

  return DeviceInfo{
    DPR: data.Properties.DevicePixelRatio,
    RW: strconv.Itoa(data.Properties.DisplayWidth),
  }
}

func Proxy(targetHost string, apiKey string, listen string) {
  config := new(proxy.Config)
  config.Listen = listen
  config.AppendHeaders = func(request *http.Request) {
    info := GetInfo(request.Header.Get("User-Agent"), apiKey)
    if request.Header.Get("CH-DPR") == "" {
      request.Header.Set("CH-DPR", info.DPR)
    }
    if request.Header.Get("CH-RW") == "" {
      request.Header.Set("CH-RW", info.RW)
    }
  }
  proxy.StartProxy(targetHost, config)
}

