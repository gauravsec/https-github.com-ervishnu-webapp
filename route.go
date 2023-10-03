package http

import (
	"app/config"
	"app/database"
	"app/utils"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	user_id := utils.UserIdFromCookie(r, store)
	if user_id == 0 {
		log.Print("issue with finding user_id from cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	route_id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	route := database.GetRouteByRouteId(db, route_id)
	if route.Id == 0 {
		log.Printf("issue with finding route  with id = %d", route_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user_info := database.GetUser(db, route.UserId)
	if user_info.FName == "" {
		log.Printf("issue with finding user with user_id = %d", user_id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parsedTemplate := template.Must(template.ParseFiles(
		"templates/route.tmpl",
		"templates/base_layout.tmpl",
	))
	s := struct {
		Style        string
		DashboardUrl string
		SettingsUrl  string
		Title        string
		RouteId      int64
		FName        string
		LName        string
		Location     string
		Distance     float64
		Calories     float64
		Time         string
	}{
		Style:        config.BaseCSS,
		DashboardUrl: config.DashboardAddress,
		SettingsUrl:  config.SettingsAddress,
		Title:        route.Title,
		RouteId:      route_id,
		FName:        user_info.FName,
		LName:        user_info.LName,
		Location:     user_info.Location,
		Distance:     route.Distance,
		Calories:     route.Calories,
		Time:         route.Time,
	}
	err = parsedTemplate.ExecuteTemplate(w, "base", s)
	if err != nil {
		log.Printf("error loading templates for landing - %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
