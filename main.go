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
			log.Printf("error decoding job\n")
			http.Error(writer, "error decoding job", http.StatusBadRequest)

			return
		}

		jobId := scheduler.AddJob(&job)

		response, _ := json.Marshal(struct {
			Id string
		}{
			jobId,
		})

		log.Printf("adding job: %s\n", jobId)

		writer.Header().Set("Content-type", "application/json")
		writer.Write(response)
	}
}

func getJobs(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		_ httprouter.Params) {
		ids := make([]string, 0, len(scheduler.Jobs))

		for k := range scheduler.Jobs {
			ids = append(ids, k)
		}

		response, _ := json.Marshal(struct {
			Ids []string
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

		job, ok := scheduler.Jobs[params.ByName("id")]

		if !ok {
			return
		}

		response, _ := json.Marshal(job)

		writer.Header().Set("Content-type", "application/json")
		writer.Write(response)
	}
}

func removeJob(scheduler *Scheduler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request,
		params httprouter.Params) {

		log.Printf("removing job: %s\n", params.ByName("id"))
		log.Printf("CODE ME: %s\n", params.ByName("id"))
	}
}

func main() {
	log.Printf("starting apimaker: version: %s\n", VERSION)

	bindAddr := flag.String("addr", "0.0.0.0:8080", "address to bind apimaker server to")

	flag.Parse()

	scheduler := NewScheduler()
	router := httprouter.New()

	router.GET("/jobs", getJobs(scheduler))
	router.POST("/jobs", addJob(scheduler))
	router.GET("/jobs/:id", getJob(scheduler))
	router.DELETE("/jobs/:id", removeJob(scheduler))

	log.Printf("starting server on: %s\n", *bindAddr)

	log.Fatal(http.ListenAndServe(*bindAddr, router))
}
