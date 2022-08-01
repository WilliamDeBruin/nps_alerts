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
)

type msgCtx struct {
	body string
	from string
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "all is good!")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) IncomingSmsHandler(w http.ResponseWriter, r *http.Request) {

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	r.ParseForm()

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
	err := s.twilioClient.SendHelp(r.FormValue("from"))
	if err != nil {
		s.logger.Error(err.Error())
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
		err := s.twilioClient.SendAlertErr(from)
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

	err = s.twilioClient.SendAlert(from, alert)

	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(alert)
}
