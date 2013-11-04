package main

import "github.com/ernesto-jimenez/ch_deviceatlas/deviceatlas"
import "flag"
import "regexp"
import "fmt"
import "os"

func main() {
  proxyTo := flag.String("proxy_to", "localhost:8080",
    "Target host:port for the proxy")
  deviceatlasKey := flag.String("deviceatlas_key", "",
    "License Key from deviceatlas.com Cloud API")
  listen := flag.String("listen", ":8000",
    "Interface the proxy will be listening to")
  flag.Parse()

  match, err := regexp.MatchString("^[0-9]+$", *listen)
  if match && err == nil {
    listento := ":" + *listen
    listen = &listento
  }

  if *deviceatlasKey == "" {
    fmt.Fprintf(os.Stderr, "Deviceatlas key is missing. Ussage:\n")
    flag.PrintDefaults()
    os.Exit(-1)
  }

  deviceatlas.Proxy(
    *proxyTo,
    *deviceatlasKey,
    *listen,
  )
}

