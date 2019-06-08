package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/golang-remote-debug/wizFind/wizards"
	"go.uber.org/zap"
)

var l *zap.SugaredLogger

const usage = `{"usage":{"search":{"endpoint":"/wizards/:field/:filter","description":"basic wizard search","fields":{"name":"wizard name","origin":"whence the wizard came","school":"school where the wizard studied"},"example":"/wizards/school/"},"return all":{"endpoint":"/wizards","description":"return all wizards"}}}`

func main() {
	zp, _ := zap.NewProduction()
	l = zp.Sugar()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	r := httprouter.New()
	r.GET("/wizards", getWizardsHandler)
	r.GET("/wizards/*filter", filterWizardHandler)
	s := &s{}
	r.NotFound = s

	hs := http.Server{Addr: ":" + port, Handler: r}

	errChan := make(chan error, 1)
	go func() {
		l.Infow("starting http server", "port", port)
		if err := hs.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		l.Infow("http err", "error", err.Error())
	case <-sigs:
	}

	l.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hs.Shutdown(ctx)
	l.Info("finished")
}

func getWizardsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	wz := wizards.GetAll()
	b, err := json.Marshal(wz)
	if err != nil {
		l.Info("error marshalling all wizards response", "error", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"message": "error marshalling response", "error": "%s"}`, err.Error())))
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func filterWizardHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	psSplit := strings.Split(ps.ByName("filter"), "/")
	if len(psSplit) != 3 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(usage))
		return
	}

	filter := strings.ToLower(psSplit[1])
	value := strings.ToLower(psSplit[2])

	if filter != "name" && filter != "origin" && filter != "school" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(usage))
		return
	}

	wizs := wizards.Find(wizards.GetAll(), func(w wizards.Wizard) bool {
		var fval string
		switch filter {
		case "name":
			fval = strings.ToLower(w.Name)
		case "origin":
			fval = strings.ToLower(w.Origin)
		case "school":
			fval = strings.ToLower(w.School)
		}
		return fval == value
	})
	b, err := json.Marshal(wizs)
	if err != nil {
		l.Info("error marshalling all wizards response", "error", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"message": "error marshalling response", "error": "%s"}`, err.Error())))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(b))
}

type s struct{}

func (s *s) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(404)
	w.Write([]byte(usage))
}
