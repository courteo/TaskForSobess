package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"task/pkg/admin"
	"task/pkg/data"
	"task/pkg/handlers"
	"task/pkg/middleware"
	"task/pkg/session"
	"task/pkg/sites"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	// основные настройки к базе
	dsn := "root:123456789@tcp(localhost:3306)/sites?"
	// указываем кодировку
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Panicf(err.Error())
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		logger.Panicf(err.Error())
	}

	// основные настройки к базе
	dsn1 := "root:123456789@tcp(localhost:3306)/sites?"
	// указываем кодировку
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db1, err := sql.Open("mysql", dsn1)
	if err != nil {
		logger.Panicf(err.Error())
	}

	db1.SetMaxOpenConns(10)

	err = db1.Ping()
	if err != nil {
		logger.Panicf(err.Error())
	}
	sessionManager := session.NewSessionsManager(db1)
	AdminRepo := admin.NewMemoryRepo(db)
	AdminHandler := &handlers.AdminHandler{
		AdminRepo:      AdminRepo,
		Logger:         logger,
		SessionManager: sessionManager,
	}

	admin := &admin.Admin{
		Login:    "ayta",
		Password: "12345678",
	}
	err1 := AdminHandler.AdminRepo.Add(admin)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	SiteRepo := sites.NewMemoryRepo(db)
	siteHandler := &handlers.SitesHandler{
		SiteRepo: SiteRepo,
		Logger:   logger,
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/site/{SITE_NAME}", siteHandler.GetSite).Methods("GET")
	r.HandleFunc("/api/get_max_access_time", siteHandler.GetMaxAccesTimeSite).Methods("GET")
	r.HandleFunc("/api/get_min_access_time", siteHandler.GetMinAccesTimeSite).Methods("GET")

	r.HandleFunc("/api/admin/get_value", middleware.Auth(siteHandler.GetValue)).Methods("GET")
	r.HandleFunc("/api/admin/login", AdminHandler.Login).Methods("GET")
	r.HandleFunc("/api/admin/register", middleware.Auth(AdminHandler.Login)).Methods("POST") //могут добавлять только другие админы через query параметры
	r.NotFoundHandler = http.HandlerFunc(NotHandler)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(*sync.WaitGroup) {
		defer wg.Done()
		i := 0
		for {
			for _, name := range data.NameOfSites {

				url := "https://" + name + "/"
				accessTime, _ := handlers.GetDurationOfRequestToSite(url, logger)

				if i == 0 {
					err := siteHandler.SiteRepo.Add(&sites.Site{
						Name:       name,
						AccessTime: accessTime,
					})

					if err != nil {
						logger.Panicf(err.Error())
					}
				} else {
					siteHandler.SiteRepo.Update(name, accessTime)
				}

			}
			i++
			time.Sleep(60 * time.Second)
		}
	}(wg)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)
	wg.Wait()
}

func NotHandler(w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadFile("static/html/index.html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
