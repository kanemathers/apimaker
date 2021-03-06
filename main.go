package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	VERSION = "0.1"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func addJob(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		_ httprouter.Params) {
		defer request.Body.Close()

		var job Job

		decoder := json.NewDecoder(request.Body)

		if err := decoder.Decode(&job); err != nil {
			log.Printf("error decoding job: %s\n", err)
			http.Error(writer, "error decoding job", http.StatusBadRequest)

			return
		}

		task := scheduler.NewTask(&job)

		response, _ := json.Marshal(struct {
			Id string `json:"id"`
		}{
			task.Id,
		})

		log.Printf("adding job: %s\n", task.Id)

		writer.Header().Set("Content-type", "application/json")
		writer.Write(response)
	}
}

func getJobs(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		_ httprouter.Params) {
		ids := make([]string, 0, len(scheduler.Tasks))

		for k := range scheduler.Tasks {
			ids = append(ids, k)
		}

		response, _ := json.Marshal(struct {
			Ids []string `json:"ids"`
		}{
			ids,
		})

		writer.Header().Set("Content-type", "application/json")
		writer.Write(response)
	}
}

func getJob(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		params httprouter.Params) {

		task, ok := scheduler.Tasks[params.ByName("id")]

		if !ok {
			return
		}

		response, _ := json.Marshal(task.Job)

		writer.Header().Set("Content-type", "application/json")
		writer.Write(response)
	}
}

func removeJob(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		params httprouter.Params) {

		log.Printf("removing job: %s\n", params.ByName("id"))

		if err := scheduler.RemoveTask(params.ByName("id")); err != nil {
			http.Error(writer, "unknown job id", http.StatusNotFound)

			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func main() {
	log.Printf("starting apimaker: version: %s\n", VERSION)

	bindAddr := flag.String("addr", "0.0.0.0:8080", "address to bind apimaker server to")

	flag.Parse()

	scheduler := NewScheduler()
	router := httprouter.New()

	scheduler.Start()

	router.GET("/jobs", getJobs(scheduler))
	router.POST("/jobs", addJob(scheduler))
	router.GET("/jobs/:id", getJob(scheduler))
	router.DELETE("/jobs/:id", removeJob(scheduler))

	log.Printf("starting server on: %s\n", *bindAddr)

	log.Fatal(http.ListenAndServe(*bindAddr, router))
}
