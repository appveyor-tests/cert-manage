package server

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	srv *http.Server
	wg  sync.WaitGroup

	seed sync.Once

	doneTpl = `<!doctype html>
<html>
<head>
  <title>cert-manage</title>
</head>
<body>
  <h3>All Done!</h3>
  <p>You can close this browser window now.</p>
</body>
</html>
`
)

// Address returns the http://$server:$port/ representation for clients to load
func Address() string {
	if srv == nil {
		return ""
	}
	return fmt.Sprintf("http://%s/", srv.Addr)
}

func Register() {
	if srv != nil {
		return // already initialized
	}

	srv = &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", getPort()),
	}

	http.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		t, err := template.New("contents").Parse(doneTpl)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("ERROR: %v", err))
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("ERROR: %v", err))
		}
	})
}

func getPort() int {
	// prep seed, otherwise math/rand.Intn is deterministic, which means there could be
	// something relying on that port
	seed.Do(func() {
		rand.Seed(time.Now().Unix()) // something
	})

	lowerBound := 1024 // ports below require elevated privileges
	n := rand.Intn(2<<15 - lowerBound)
	return n + lowerBound
}

// Start initializes the http server, binds and accepts connections
func Start() {
	if srv == nil {
		fmt.Fprint(os.Stderr, "ui/server: no http server has been registered")
		return
	}

	wg.Add(1) // force callers to wait

	// spawn off http server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ERROR: failed creating localhost server err=%s\n", err)
		}
	}()
}

// Stop calls for a shutdown of the http server, if it exists
func Stop() error {
	// hold on until the form has been filled out
	wg.Wait()

	if srv != nil {
		return srv.Shutdown(nil)
	}
	return nil
}