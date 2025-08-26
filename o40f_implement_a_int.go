package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type Pipeline struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Stages []Stage `json:"stages"`
}

type Stage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
}

type Controller struct {
	pipelines map[string]Pipeline
	upgrader websocket.Upgrader
}

func NewController() *Controller {
	return &Controller{
		pipelines: make(map[string]Pipeline),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (c *Controller) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var payload struct {
			Action string `json:"action"`
			PipelineID string `json:"pipeline_id"`
			StageID string `json:"stage_id"`
		}

		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		switch payload.Action {
		case "start_pipeline":
			c.startPipeline(payload.PipelineID)
		case "stop_pipeline":
			c.stopPipeline(payload.PipelineID)
		case "deploy_stage":
			c.deployStage(payload.PipelineID, payload.StageID)
		default:
			log.Println("unknown action")
		}
	}
}

func (c *Controller) startPipeline(pipelineID string) {
	pipeline, ok := c.pipelines[pipelineID]
	if !ok {
		log.Println("pipeline not found")
		return
	}

	// start pipeline logic here
	log.Println("starting pipeline", pipeline.Name)
}

func (c *Controller) stopPipeline(pipelineID string) {
	pipeline, ok := c.pipelines[pipelineID]
	if !ok {
		log.Println("pipeline not found")
		return
	}

	// stop pipeline logic here
	log.Println("stopping pipeline", pipeline.Name)
}

func (c *Controller) deployStage(pipelineID string, stageID string) {
	pipeline, ok := c.pipelines[pipelineID]
	if !ok {
		log.Println("pipeline not found")
		return
	}

	for _, stage := range pipeline.Stages {
		if stage.ID == stageID {
			// deploy stage logic here
			log.Println("deploying stage", stage.Name)
			return
		}
	}

	log.Println("stage not found")
}

func main() {
	controller := NewController()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleWebsocket(w, r)
	})

	http.HandleFunc("/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var pipelines []Pipeline
		for _, pipeline := range controller.pipelines {
			pipelines = append(pipelines, pipeline)
		}

		json.NewEncoder(w).Encode(pipelines)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}