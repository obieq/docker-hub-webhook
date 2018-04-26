package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/obieq/docker-hub-webhook/config"
	"github.com/obieq/docker-hub-webhook/lib"
	"github.com/obieq/docker-hub-webhook/scripts"
)

// DockerHubHandler => processes all Docker Hub HTTP requests
type DockerHubHandler struct {
	BaseHandler
}

type PushData struct {
	Tag    string `json:"tag"`
	Pusher string `json:"pusher"`
}

type Repository struct {
	Name     string `json:"name"`
	RepoName string `json:"repo_name"`
}

type Webhook struct {
	PushData    PushData   `json:"push_data"`
	CallbackURL string     `json:"callback_url"`
	Repository  Repository `json:"repository"`
}

type Callback struct {
	State       string `json:"state"`
	Context     string `json:"context"`
	Description string `json:"description"`
}

type TemplateData struct {
	Vhost    string
	RepoName string
	Name     string
	Tag      string
	Params   string
}

// Handle => verfies and runs all supported DockerHub HTTP methods
func (h *DockerHubHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.create(w, r)
		break
	default:
		h.writeError(w, "invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
}

// create => processes a docker hub web hook request
func (h *DockerHubHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload Webhook

	config.Log("Calling DockerHub create()")
	log.Printf("%s - %s %s", r.RemoteAddr, r.Method, r.URL.EscapedPath())

	token := r.FormValue("token")

	// verify query string token
	if token != config.Cfg().SecurityToken {
		h.writeError(w, "invalid token", http.StatusUnauthorized)

		// send asynchronous slack error message
		go lib.SendSlackSwarmRedeployFailureMessage("unknown repo", "unknown tag", "invalid token")

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, err.Error(), http.StatusUnprocessableEntity)

		// send asynchronous slack error message
		go lib.SendSlackSwarmRedeployFailureMessage("unknown repo", "unknown tag", err.Error())

		return
	}
	defer r.Body.Close()

	config.Log(payload.CallbackURL)
	config.Log(payload.PushData)
	config.Log(payload.Repository)

	// update swarm service(s)
	if err := runSwarmUpdateScripts(payload.Repository.RepoName, payload.PushData.Tag); err != nil {
		h.writeError(w, "error occurred while trying to update swarm service: "+err.Error(), http.StatusUnprocessableEntity)

		// send asynchronous slack error message
		go lib.SendSlackSwarmRedeployFailureMessage(payload.Repository.RepoName, payload.PushData.Tag, err.Error())

		return
	}

	// set headers
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusCreated)

	// write response
	b, _ := json.Marshal(payload)
	config.Log(string(b))
	io.Writer.Write(w, b)
}

func runSwarmUpdateScripts(repoName string, tag string) error {
	// get list of swarm services
	rs := scripts.ScriptListSwarmServices()
	response, err := rs.Execute()
	if err != nil {
		return err
	}
	if response.Results[0].Stderr != "" {
		return errors.New(response.Results[0].Stderr)
	}

	config.Log("'docker service ls' Results")
	config.Log(response)

	// NOTE: example of docker service ls result can be found within this project:
	//       test_scripts.go => ScriptTestDockerServiceLS

	// find matching services and restart each one
	arrServices := strings.Split(response.Results[0].Stdout, "\n")

	// get columns
	arrCols := strings.Fields(arrServices[0]) // column names are the first "row" when docker service ls is invoked

	// build map of column names and respective index positions
	m := map[string]int{}
	for index, col := range arrCols {
		m[col] = index
	}

	// loop through each service and see if its image matches that of the docker hub webhook
	// if it matches, then update service
	for _, service := range arrServices {
		if len(service) == 0 { // last service is blank due to split on new line
			break
		}

		arrData := strings.Fields(service)
		image := arrData[m["IMAGE"]]
		config.Log(arrData)

		if image == repoName+":"+tag {
			name := arrData[m["NAME"]]
			log.Printf("Updating swarm service => Name: %s, Image: %s", name, image)

			rs := scripts.ScriptUpdateSwarmService(name, image)
			response, err := rs.Execute()
			if err != nil {
				return err
			}
			if response.Results[0].Stderr != "" {
				return errors.New(response.Results[0].Stderr)
			}

			config.Log("Update swarm service response")
			config.Log(response)

			// send asynchronous slack success message
			go lib.SendSlackSwarmRedeploySuccessMessage(repoName, tag)
		}
	}

	return nil
}
