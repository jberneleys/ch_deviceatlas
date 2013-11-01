package proxy

import "errors"
import "net"
import "net/http"
import "fmt"
import "io/ioutil"

type Config struct {
  Listen string
  AppendHeaders func(*http.Request)
}

func StartProxy(targetHost string, config *Config) {
  l, err := net.Listen("tcp", config.Listen)
  if err != nil {
    panic(err)
  }

  server := http.NewServeMux()
  server.HandleFunc("/", proxyHandler(targetHost, config.AppendHeaders))

  defer l.Close()
  http.Serve(l, server)
}

func noRedirect(req *http.Request, via []*http.Request) error {
  return errors.New("No redirect")
}

func proxyHandler(
  targetHost string,
  appendHeaders func(*http.Request)) func(http.ResponseWriter, *http.Request) {

  handler := func (w http.ResponseWriter, req *http.Request) {
    // Create new request for the target server
    request, err := http.NewRequest(req.Method, "http://" + targetHost + req.URL.Path, nil)
    if err != nil {
      fmt.Println(err)
    }

    // Add headers from original request
    for header, value := range req.Header {
      request.Header.Set(header, value[0])
    }

    // Append headers
    if appendHeaders != nil {
      appendHeaders(request)
    }

    client := http.Client{CheckRedirect: noRedirect}
    resp, err := client.Do(request)
    if err != nil {
      fmt.Println(err)
    }


    // Reply the response
    headers := w.Header()
    for header, value := range resp.Header {
      headers[header] = value
    }
    w.WriteHeader(resp.StatusCode)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err == nil {
      w.Write(body)
    } else {
      fmt.Println(err)
    }
    fmt.Println("Proxied %v", req)
  }
  return handler
}


