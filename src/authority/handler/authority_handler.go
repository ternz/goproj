package handler

import (
	"database/sql"
	"strconv"
	"strings"
	"net/http"
	
	"common/logging"
	"authority/dbaccess"
	"authority/util"
)

func HandleQueryAuthority(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	if strings.Contains(r.URL.Path, "/id") {
		handleQueryAuthorityById(w, r)
	} else if strings.Contains(r.URL.Path, "/code") {
		handleQueryAuthorityByCode(w, r)
	} else if strings.Contains(r.URL.Path, "/level") {
		handleQueryAuthorityByGroup(w, r)
	} else {
		panic("wrong QueryAuthority url pattern")
	}
}

func handleQueryAuthorityById(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get(":id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		logging.Warning("handleQueryAuthorityById got invalid input id:%s", id_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	auth, err := dbaccess.QueryAuthorityById(id)
	if err != nil && err != sql.ErrNoRows {
		logging.Error("handleQueryAuthorityById database error:%", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	util.ResponseJson(w, auth)
}

func handleQueryAuthorityByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(":code")
	auth, err := dbaccess.QueryAuthorityByCode(code)
	if err != nil && err != sql.ErrNoRows {
		logging.Error("handleQueryAuthorityByCode database error:%", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	util.ResponseJson(w, auth)
}

func handleQueryAuthorityByGroup(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get(":group")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		logging.Warning("handleQueryAuthorityByGroup got invalid input id:%s", id_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	auths, err := dbaccess.QueryAuthoritysByGroupId(id)
	if err != nil {
		logging.Error("handleQueryAuthorityByGroup database error:%", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	util.ResponseJson(w, auths)
}

func HandleQueryAuthorityGroupAll(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(":code")
	auths, err := dbaccess.QueryAuthorityGroup(code)
	if err != nil {
		logging.Error("HandleQueryAuthorityGroupAll database error:%", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	
	tree := dbaccess.ConstructAuthorityTree(auths)
	util.ResponseJson(w, tree)
}

func HandleCreateAuthority(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	auth := &dbaccess.Authority{}
	err := util.ParseJsonBody(r.Body, auth)
	if err != nil {
		logging.Error("HandleCreateAuthority ParseJsonBody error:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	//TODO 检查code的前缀是否一致
	err = dbaccess.CreateAuthority(auth)
	if err != nil {
		logging.Error("CreateAuthority database error:%s", err.Error())
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, ERR_STR_DUPLICATE_ENTRY, http.StatusBadRequest)
		} else {
			http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUpdateAuthorityName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	auth := &dbaccess.Authority{}
	err := util.ParseJsonBody(r.Body, auth)
	if err != nil {
		logging.Error("HandleUpdateAuthorityName ParseJsonBody error:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	err = dbaccess.UpdateAuthorityName(auth.Id, auth.Name)
	if err != nil {
		logging.Error("HandleUpdateAuthorityName database error:%s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(REQUEST_SUCCESS))
}