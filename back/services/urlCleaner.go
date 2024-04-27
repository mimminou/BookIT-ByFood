package services

import (
	"encoding/json"
	. "github.com/mimminou/BookIT-ByFood/back/models"
	"github.com/mimminou/BookIT-ByFood/back/utils"
	"net/http"
)

// routes requests that have ID based on HTTP

// @Summary		Process URL
// @Description	Processes URLs depending on the requested operation
// @Tags			url
// @Accept			json
// @Produce		json
// @Param			RequestStruct	body	RequestStruct	true	"Request Body"
// @Success		200 {object}	ResponseStruct
// @Failure		400 {object}	ErrMessage
// @Failure		405
// @Router			/url/ [post]
func ProcessUrl(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Error : Request Body is empty"})
		w.Write(jsonResponse)
		return
	}

	var request RequestStruct
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid request format"})
		w.Write(jsonResponse)
		return
	}

	if utils.IsUrl(request.Url) == false {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Url format invalid"})
		w.Write(jsonResponse)
		return
	}

	// We are sure it's a URL from this point

	//check operation type and run it, if not valid, return BAD REQUEST
	switch request.Operation {
	case "canonical":
		processedUrl, err := utils.GetCanonicalUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return
	case "redirection":
		processedUrl, err := utils.GetRedirectionUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return

	case "all":
		processedUrl, err := utils.GetCanonicalUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}

		processedUrl, err = utils.GetRedirectionUrl(processedUrl) //reprocess url
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return

	default:
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid operation"})
		w.Write(jsonResponse)
		return
	}
}
