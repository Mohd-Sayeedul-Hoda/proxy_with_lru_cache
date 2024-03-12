package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"go_web_proxy_with_cache/lru_cache/cache"
)

var defaultConfig = http.DefaultTransport
var lruCache cache.LRUCache

func main(){
  capacity := flag.Uint64("cap", 20, "insalize the size of cache")
  flag.Parse()

  lruCache.Nodes = make(map[string]*cache.Node)
  lruCache.Capacity = *capacity

  log.Println("server running on port 8000")
  http.ListenAndServe(":8000", http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request){
  targetURL := r.URL

	cacheNode, present := lruCache.Get(targetURL.String())

	if present{
		fmt.Println("serving from the cache")
		toTheClient(w, cacheNode.Header, cacheNode.Body)
		return
	}

  proxyRequest, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
  if err != nil{
    http.Error(w, "you request can't be resolve", http.StatusInternalServerError)
    return
  }

  for name, values := range r.Header{
    for _, value := range values{
      r.Header.Set(name, value)
    }
  }

  resp, err := defaultConfig.RoundTrip(proxyRequest)
  if err != nil{
    http.Error(w, "error while sending request", http.StatusInternalServerError)
		return
  }

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil{
		http.Error(w, "Inertal server error", http.StatusInternalServerError)
		return
	}
	lruCache.Put(targetURL.String(), resp.Header, bodyByte)
	fmt.Println("all the way around from server")
	toTheClient(w, resp.Header, bodyByte)
}

func toTheClient(w http.ResponseWriter,header http.Header, bodyBytes []byte){

  for name, values := range header{
    for _, value := range values{
      w.Header().Add(name, value)
    }
  }

  io.Copy(w, bytes.NewReader(bodyBytes))
}
