package main

import "github.com/ernesto-jimenez/ch_deviceatlas/deviceatlas"
import "os"

func main() {
  deviceatlas.Proxy(
    os.Getenv("PROXY_TO"),
    os.Getenv("DEVICEATLAS_KEY"),
    ":"+os.Getenv("PORT"),
  )
}

