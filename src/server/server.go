package main

import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"log"
  "strconv"
	"strings"
	"net"
	"net/http"
	"time"
)

func formatRequest(r *http.Request) string {
 // Create return string
 var request []string
 // Add the request string
 url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
 request = append(request, url)
 //  Add the Remote Address
 request = append(request, fmt.Sprintf("Remote Address: %v", r.RemoteAddr))
 // Add the host
 request = append(request, fmt.Sprintf("Host: %v", r.Host))
 // Add headers
 for name, headers := range r.Header {
   for _, h := range headers {
     request = append(request, fmt.Sprintf("%v: %v", name, h))
   }
 }
  // Return the request as a string
  return strings.Join(request, ", ")
}

func root(w http.ResponseWriter, r *http.Request) {
	log.Printf(formatRequest(r))
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        query := r.URL.Query()
        delayparam, present := query["delay"] // delay=10
        if present {
          delay, _ := strconv.Atoi(delayparam[0])
          time.Sleep(time.Duration(delay) * time.Second)
        }
	fmt.Fprintf(w, "Hello, world!")
}

func version(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("VERSION")
    failOnError(err, "Could not get version")
    ver := string(b)
	fmt.Fprintf(w, ver)
}

func ping (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func headers (w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "%v: %v\n", k, v)
	}
}

func ipAddress (w http.ResponseWriter, r *http.Request) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		addresses := strings.Split(xForwardedFor, ",")
		fmt.Fprintf(w, "Forwarded Source IP address: %s\n", strings.TrimSpace(addresses[0]))
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Error parsing IP address: %v\n", err)
		return
	}
	fmt.Fprintf(w, "Source IP address: %s\n", ip)

}

func getUrl(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
  url := query.Get("url")
  if url == "" {
    fmt.Fprintf(w, "Must provide url query param")
    return
  }
	fmt.Fprintf(os.Stderr, "Fetching URL %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "Request failed with status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading response body: %v\n", err)
		return
	}

	fmt.Fprintf(w, string(body))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/version", version)
	http.HandleFunc("/ipAddress", ipAddress)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/getUrl", getUrl)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	fmt.Fprintf(os.Stderr, fmt.Sprintf("Listening on port %s...\n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
