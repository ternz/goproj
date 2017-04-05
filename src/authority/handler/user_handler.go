package handler

import (
	"strconv"
	"common/logging"
	"net/http"
	"strings"
	
	"authority/util"
	"authority/dbaccess"
)

type paramIdCode struct {
	Id 	uint64 `json:"id"`
	Code string `json:"code"`
}

type paramUserAuthorityGrant struct {
	UserId uint64 `json:"user_id"`
	AuthAdd []uint64 `json:"auth_add"`
	AuthRemove []uint64 `json:"auth_remove"`
}

type paramUserAuthorityCheck struct {
	UserId uint64 `json:"user_id"`
	CheckList []string 	`json:"check_list"`
}

func HandleUserQuery(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	logging.Debug("request: %+v", r)
	
	queryParams := r.URL.Query()
	idx_str := queryParams.Get("idx")
	size_str := queryParams.Get("size")
	idx, err := strconv.Atoi(idx_str)
	if err != nil || idx < 0 {
		logging.Warning("HandleUserQuery got invalid paramter idx:%s", idx_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(size_str)
	if err != nil || size < 0 {
		logging.Warning("HandleUserQuery got invalid paramter size:%s", size_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	
	count, err := dbaccess.QueryUserCount()
	if err != nil {
		logging.Error("dbaccess.QueryUserCount error:%s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	
	users, err := dbaccess.QueryUsers(idx, size)
	if err != nil {
		logging.Error("dbaccess.QueryUsers error:%s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	
	rspMap := make(map[string]interface{})
	rspMap["users"] = users
	rspMap["next_idx"] = idx+len(users)
	rspMap["total_size"] = count
	util.ResponseJson(w, rspMap)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	user := &dbaccess.User{}
	err := util.ParseJsonBody(r.Body, user)
	if err != nil {
		logging.Error("HandlerUserCreate invalid input body:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	err = dbaccess.CreateUser(user)
	if err != nil {
		logging.Error("dbaccess.CreateUser error: %s", err.Error())
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, ERR_STR_DUPLICATE_ENTRY, http.StatusBadRequest)
		} else {
			http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	invalidInput := func() {
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
	}
	
	param := paramIdCode{}
	err := util.ParseJsonBody(r.Body, &param)
	if err != nil {
		logging.Error("HandleUserDelete invalid input body:%s", err.Error())
		invalidInput()
		return
	}
	if param.Id > 0 {
		err = dbaccess.DeleteUserById(param.Id)
	} else if param.Code != "" {
		err = dbaccess.DeleteUserByCode(param.Code)
	} else {
		invalidInput()
		return
	}
	if err != nil {
		logging.Error("db delete user failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserDeleteById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	id_str := r.URL.Query().Get(":id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil || id == 0 {
		logging.Error("invalid input id:%s", id_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	err = dbaccess.DeleteUserById(id)
	if err != nil {
		logging.Error("db delete user failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserDeleteByCode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	code := r.URL.Query().Get(":code")
	err := dbaccess.DeleteUserByCode(code)
	if err != nil {
		logging.Error("db delete user failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	user := &dbaccess.User{}
	err := util.ParseJsonBody(r.Body, user)
	if err != nil {
		logging.Error("HandleUserUpdate invalid input body:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	if user.Id == 0 || user.Code == "" {
		logging.Info("HandleUserUpdate invalid id or code (id=%d, code=%s)", user.Id, user.Code)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	err = dbaccess.UpdateUser(user)
	if err != nil {
		logging.Error("dbaccess.CreateUser error: %s", err.Error())
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, ERR_STR_DUPLICATE_ENTRY, http.StatusBadRequest)
		} else {
			http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserAuthorityGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	id_str := r.URL.Query().Get(":user_id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil || id == 0 {
		logging.Warning("HandleUserAuthorityGet invalid input id:%s", id_str)
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	auths, err := dbaccess.QueryUserAuthority(id)
	if err != nil {
		logging.Error("HandleUserAuthorityGet QueryUserAuthority failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	outMap := make(map[string]interface{})
	outMap["user_id"]=id
	outMap["auths"]=auths
	util.ResponseJson(w, outMap)
}

/*
input:
{
	"user_id":xxx,
	"auth_add":[xxx,xxx,...],
	"auth_remove":[xxx,xxx,...]
}
*/
func HandleUserAuthorityGrant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	para := &paramUserAuthorityGrant{}
	err := util.ParseJsonBody(r.Body, para)
	logging.Debug("paramter: %+v", para)
	if err != nil {
		logging.Error("HandleUserAuthorityGrant invalid input body:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	tx, err := dbaccess.DB().Begin()
	if err != nil {
		logging.Error("HandleUserAuthorityGrant begin transaction failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
	}
	if len(para.AuthRemove) != 0 {
		for _, id := range para.AuthRemove {
			err = dbaccess.TxDeleteUserAuthority(para.UserId, id, tx)
			if err != nil {
				tx.Rollback()
				logging.Error("HandleUserAuthorityGrant TxDeleteUserAuthority failed: %s", err.Error())
				http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
				return
			}
		}
	}
	if len(para.AuthAdd) != 0 {
		for _, id := range para.AuthAdd {
			err = dbaccess.TxGrantUserAuthority(para.UserId, id, tx)
			if err != nil {
				tx.Rollback()
				logging.Error("HandleUserAuthorityGrant TxGrantUserAuthority failed: %s", err.Error())
				if strings.Contains(err.Error(), "Duplicate entry") {
					http.Error(w, ERR_STR_DUPLICATE_ENTRY, http.StatusBadRequest)
				} else {
					http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
				}
				return
			}
		}
	}
	tx.Commit()
	w.Write([]byte(REQUEST_SUCCESS))
}

func HandleUserAuthorityCheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Debug("request: %+v", r)
	
	para := &paramUserAuthorityCheck{}
	err := util.ParseJsonBody(r.Body, para)
	logging.Debug("paramter: %+v", para)
	if err != nil {
		logging.Error("HandleUserAuthorityCheck invalid input body:%s", err.Error())
		http.Error(w, ERR_STR_INVALID_INPUT, http.StatusBadRequest)
		return
	}
	user_auths, err := dbaccess.QueryUserAuthority(para.UserId)
	if err != nil {
		logging.Error("HandleUserAuthorityCheck QueryUserAuthority failed: %s", err.Error())
		http.Error(w, ERR_STR_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}
	outMap := make(map[string]bool)
	for _, check := range para.CheckList {
		flag := false
		for _, auth := range user_auths {
			if strings.Contains(check, auth.Code) {
				flag = true
				break
			}
		}
		outMap[check] = flag
	}
	util.ResponseJson(w, outMap)
}