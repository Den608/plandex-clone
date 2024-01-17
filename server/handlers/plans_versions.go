package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"plandex-server/db"

	"github.com/gorilla/mux"
	"github.com/plandex/plandex/shared"
)

func ListLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for ListLogsHandler")

	auth := authenticate(w, r, true)
	if auth == nil {
		return
	}
	currentUserId := auth.User.Id

	vars := mux.Vars(r)
	planId := vars["planId"]

	log.Println("planId: ", planId)

	if authorizePlan(w, planId, currentUserId, auth.OrgId) == nil {
		return
	}

	body, shas, err := db.GetGitCommitHistory(auth.OrgId, planId)

	if err != nil {
		log.Println("Error getting logs: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := shared.LogResponse{
		Body: body,
		Shas: shas,
	}

	bytes, err := json.Marshal(res)

	if err != nil {
		log.Println("Error marshalling logs: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)

	log.Println("Successfully processed request for ListLogsHandler")
}

func RewindPlanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for RewindPlanHandler")

	auth := authenticate(w, r, true)
	if auth == nil {
		return
	}
	currentUserId := auth.User.Id

	vars := mux.Vars(r)
	planId := vars["planId"]

	log.Println("planId: ", planId)

	if authorizePlan(w, planId, currentUserId, auth.OrgId) == nil {
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody shared.RewindPlanRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	err = db.GitRewindToSha(auth.OrgId, planId, requestBody.Sha)

	if err != nil {
		log.Println("Error rewinding plan: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.SyncPlanTokens(auth.OrgId, planId)

	if err != nil {
		log.Println("Error syncing plan tokens: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sha, latest, err := db.GetLatestCommit(auth.OrgId, planId)

	if err != nil {
		log.Println("Error getting latest commit: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := shared.RewindPlanResponse{
		LatestSha:    sha,
		LatestCommit: latest,
	}

	bytes, err := json.Marshal(res)

	if err != nil {
		log.Println("Error marshalling response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)

	log.Println("Successfully processed request for RewindPlanHandler")
}
