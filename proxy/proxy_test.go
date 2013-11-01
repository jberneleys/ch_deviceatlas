package proxy

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
  {url: "/assert/added_forwarded_for"},
  {url: "/assert/appends_to_existing_forwarded_for", prepareRequest: func(req *http.Request) {
    req.Header["X-Forwarded-For"] = []string{"192.268.10.10"}
  }},
}

func handleRequests(t *testing.T) {
  http.HandleFunc("/assert/added_forwarded_for", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Header)
    forwardedFor := r.Header.Get("X-Forwarded-For")
    if forwardedFor == "" {
      t.Error("Missing CH-forwardedFor header")
    } else if forwardedFor != "127.0.0.1" {
      t.Error("CH-forwardedFor expected: 127.0.0.1 Got:", forwardedFor)
    }
  })

  http.HandleFunc("/assert/appends_to_existing_forwarded_for", func(w http.ResponseWriter, r *http.Request) {
    forwardedFor := r.Header.Get("X-Forwarded-For")
    if forwardedFor != "192.268.10.10, 127.0.0.1" {
      t.Error("Wrong X-Other-Header header expected: 192.268.10.10, 127.0.0.1 Got:", forwardedFor)
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
  l, err := net.Listen("tcp", ":2345")
  if err != nil {
    t.Fatal(err)
  }
  handleRequests(t)
  go func() {
    http.Serve(l, nil)
  }()

  go func() {
    config := Config{Listen: ":5432"}
    StartProxy("0.0.0.0:2345", &config)
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
  req, _ := http.NewRequest("GET", "http://0.0.0.0:5432" + test.url, nil)
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

