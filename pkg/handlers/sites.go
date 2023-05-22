package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	// "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"task/pkg/forms"
	"task/pkg/sites"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)


var (
	ErrCantMarshal = errors.New("cant Marshal")
	ErrCantDelete  = errors.New("cant Delete")
)

type SitesHandler struct {
	Logger         *zap.SugaredLogger
	SiteRepo 		sites.SiteRepo
	Stats			[3]int
}

// All getters

func (h *SitesHandler) GetSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars["SITE_NAME"])
	res, err := h.SiteRepo.FindSite(vars["SITE_NAME"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetSite: "+err.Error(), h.Logger)
		return
	}
	SendTimeRequest(w, "GetSite: ", res, http.StatusOK, h.Logger)
	h.Stats[0]++
	h.Logger.Infof("GetSite: %v", http.StatusOK)
	return
}

func (h *SitesHandler) GetMinAccesTimeSite(w http.ResponseWriter, r *http.Request) {
	res, err := h.SiteRepo.FindMinAccessTimeSite()
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetMinAccessTimeSite: "+err.Error(), h.Logger)
		return
	}
	SendNameRequest(w, "GetMinAccessTimeSite: ", res, http.StatusOK, h.Logger)

	h.Logger.Infof("GetMinAccessTimeSite: %v", http.StatusOK)
	h.Stats[1]++
	return
}

func (h *SitesHandler) GetMaxAccesTimeSite(w http.ResponseWriter, r *http.Request) {
	res, err := h.SiteRepo.FindMaxAccessTimeSite()
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetMaxAccessTimeSite: "+err.Error(), h.Logger)
		return
	}
	SendNameRequest(w, "GetMaxAccessTimeSite: ", res, http.StatusOK, h.Logger)
	h.Stats[2]++
	h.Logger.Infof("GetMaxAccessTimeSite: %v", http.StatusOK)
	return
}

func (h *SitesHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(map[string]int{
		"GetMaxAccessTimeSite": h.Stats[2],
		"GetinAccessTimeSite": h.Stats[1],
		"GetSite": h.Stats[0],
	})
	if err != nil {
		JsonError(w, http.StatusBadRequest, ErrCantMarshal.Error(), h.Logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}


// Send Request

func SendNameRequest(w http.ResponseWriter, errStr string, site *sites.Site, status int, Logger *zap.SugaredLogger) {
	form := forms.SendNameForm{Name: site.Name}
	resp, err := json.Marshal(form)
	if err != nil {
		JsonError(w, http.StatusBadRequest, errStr+ErrCantMarshal.Error(), Logger)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func SendTimeRequest(w http.ResponseWriter, errStr string, site *sites.Site, status int, Logger *zap.SugaredLogger) {
	form := forms.SendTimeForm{Time: site.AccessTime}
	fmt.Println(site.AccessTime, site.Name)
	resp, err := json.Marshal(form)
	if err != nil {
		JsonError(w, http.StatusBadRequest, errStr+ErrCantMarshal.Error(), Logger)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

// JsonError

func JsonError(w io.Writer, status int, msg string, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})

	if err != nil {
		w.Write([]byte("bad request"))
		return
	}

	w.Write(resp)
}

// Get Duration of Request to site (every minute)

func GetDurationOfRequestToSite(u string,  Logger *zap.SugaredLogger) (time.Duration, bool) {
	t := time.Now()
    if u == "" {
        return time.Since(t), false
    }
    
	client := http.Client{}
    req, _ := http.NewRequest("GET", u, nil)
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
    
	resp, err := client.Do(req)
    if err != nil {
        Logger.Infof("error: %v ", err)
        return time.Since(t), false
    } 
    
	_, err = ioutil.ReadAll(resp.Body)
    if err == nil {
        if resp.StatusCode == 200 {
            Logger.Infof( "%v : ok", u)
        } else {
            Logger.Infof("Site %v returned error code: %v", u, resp.StatusCode)   
        }
    }
    
	defer resp.Body.Close()
    return time.Since(t), true
}