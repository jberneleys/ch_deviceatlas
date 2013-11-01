package deviceatlas

import (
  "testing"
  "net"
  "net/http"
  "fmt"
)

type testpair struct {
  url string
  prepareRequest func(req *http.Request)
  responseAssertions func(t *testing.T, resp *http.Response)
}

var tests = []testpair{
  {url: "/assert/keeps_existing_header", prepareRequest: func(req *http.Request) {
    req.Header["CH-DPR"] = []string{"Existing DPR header"}
    req.Header["CH-RW"] = []string{"Existing RW header"}
  }},
  {url: "/assert/added_header", prepareRequest: func(req *http.Request) {
    // iPhone 5 UA
    req.Header["User-Agent"] = []string{uaiPhone5}
  }},
  {url: "/assert/respects_other_headers", prepareRequest: func(req *http.Request) {
    req.Header["X-Other-Header"] = []string{"passed"}
  }},
}

func handleRequests(t *testing.T) {
  http.HandleFunc("/assert/added_header", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Header)
    dpr := r.Header.Get("CH-DPR")
    if dpr == "" {
      t.Error("Missing CH-DPR header")
    } else if dpr != "2.0" {
      t.Error("CH-DPR expected: 123 Got:", dpr)
    }
    rw := r.Header.Get("CH-RW")
    if rw == "" {
      t.Error("Missing CH-RW header")
    } else if rw != "320" {
      t.Error("CH-RW expected: Added header Got:", rw)
    }
  })

  http.HandleFunc("/assert/keeps_existing_header", func(w http.ResponseWriter, r *http.Request) {
    dpr := r.Header.Get("CH-DPR")
    if dpr != "Existing DPR header" {
      t.Error("Wrong CH-DPR header expected: Existing DPR header Got:", dpr)
    }
    rw := r.Header.Get("CH-RW")
    if rw != "Existing RW header" {
      t.Error("Wrong CH-RW header expected: Existing RW header Got:", rw)
    }
  })

  http.HandleFunc("/assert/respects_other_headers", func(w http.ResponseWriter, r *http.Request) {
    other_header := r.Header.Get("X-Other-Header")
    if other_header != "passed" {
      t.Error("Wrong X-Other-Header header expected: passed Got:", other_header)
    }
  })
}


func TestProxy(t *testing.T) {
  l, err := net.Listen("tcp", ":1234")
  if err != nil {
    t.Fatal(err)
  }
  handleRequests(t)
  go func() {
    http.Serve(l, nil)
  }()

  go func() {
    Proxy("0.0.0.0:1234", "19920eabfccb458bfdd96ef395e659eb", ":4321")
  }()

  for _, test := range tests {
    request(test, t)
  }

  // Close listener
  if err = l.Close(); err != nil {
    panic(err)
  }
}

func request(test testpair, t *testing.T) {
  req, _ := http.NewRequest("GET", "http://0.0.0.0:4321" + test.url, nil)
  if test.prepareRequest != nil {
    test.prepareRequest(req)
  }

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    t.Error(test, err)
  }

  if resp.StatusCode != 200 {
    t.Error("No URL", test.url, "status code:", resp.StatusCode)
  }

  if test.responseAssertions != nil {
    test.responseAssertions(t, resp)
  }
}

const uaiPhone5 = "Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3"

