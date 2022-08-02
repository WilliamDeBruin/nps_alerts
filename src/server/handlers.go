package server

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	helpPrefix  = "help"
	alertPrefix = "alerts "

	helpMessage  = "Welcome to NPS alerts! Here is a list of commands:\n\nHelp: receive this help text\n\nAlerts {state}: Text \"alerts\" followed by the 2-letter state code of the state you would like to see alerts for"
	alertMessage = "Here is the most recent NPS %s alert from %s, published %s:\n\n%s\n\n%s\n\nFor a full list of NPS %s alerts, visit %s"
)

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "all is good!")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) IncomingSmsHandler(w http.ResponseWriter, r *http.Request) {

	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	from := r.FormValue("from")
	if from == "" {
		s.logger.Error("missing field in request body: from")
		w.WriteHeader(http.StatusBadRequest)
	}

	body := r.FormValue("body")
	if body == "" {
		s.logger.Error("missing field in request body: body")
		w.WriteHeader(http.StatusBadRequest)
	}

	if strings.HasPrefix(body, helpPrefix) {
		s.helpHandler(w, r)
		return
	} else if strings.HasPrefix(body, alertPrefix) {
		s.alertHandler(w, r)
		return
	} else {
		s.logger.Error("unhandled text body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) helpHandler(w http.ResponseWriter, r *http.Request) {
	err := s.twilioClient.SendMessage(r.FormValue("from"), helpMessage)
	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.logger.Info("sent help message")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) alertHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	body = strings.TrimSpace(body)
	from := r.FormValue("from")
	words := strings.Split(body, " ")

	if len(words) != 2 {
		err := s.twilioClient.SendMessage(from, `I'm sorry, I couldn't understand your message. Please text "alerts {state}" for recent alerts`)
		if err != nil {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stateCode := words[1]

	alert, err := s.npsClient.GetAlert(stateCode)

	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("alert response", zap.Any("alertResponse", alert))

	message := fmt.Sprintf(alertMessage,
		alert.FullStateName,
		alert.FullParkName,
		alert.RecentAlertDate,
		alert.AlertHeader,
		alert.AlertMessage,
		alert.FullStateName,
		alert.URL)

	err = s.twilioClient.SendMessage(from, message)

	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(alert)
}
