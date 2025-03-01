package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"github.com/tsawler/toolbox"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var tools toolbox.Tools

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.writeLog(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
    case "ping":
	default:
		tools.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	// create some json we'll send to the mail microservice
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://svc-mailer/send", bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err)
        return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
        tools.ErrorJSON(w, errors.New("error calling mail service"))
        return
    }

    var payload toolbox.JSONResponse
    payload.Error = false
    payload.Message = "Mail sent"
    tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) writeLog(w http.ResponseWriter, l LogPayload) {

	// create some json we'll send to the logging microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://svc-logger/log", bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("error calling logging service"))
		return
	}

	var payload toolbox.JSONResponse 
	payload.Error = false
	payload.Message = "Logged"
	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://svc-auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService toolbox.JSONResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
