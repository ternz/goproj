// It includes skill, equipment, card and so on
package main

import (
	"fmt"
	"time"

	"net/http"
	"runtime"

	"common/libutil"
	"common/logging"

	_ "github.com/go-sql-driver/mysql"
	"github.com/drone/routes"
	
	"authority/handler"
	"authority/dbaccess"
)

func init() {
	//配置解析
	InitConfigure("conf/config.json")
	//mysql
	dbaccess.InitMysql(Cfg.Server.Mysql)
}

//注册http回调
func registerHttpHandle() {
	//http.HandleFunc("/test", app.HandleTest)
	mux := routes.New()
	mux.Get("/authority/user/list", handler.HandleUserQuery)
	mux.Post("/authority/user/register", handler.HandleUserCreate)
	mux.Del("/authority/user/delete", handler.HandleUserDelete)
	mux.Del("/authority/user/delete/id/:id", handler.HandleUserDeleteById)
	mux.Del("/authority/user/delete/code/:code", handler.HandleUserDeleteByCode)
	mux.Get("/authority/user/:user_id/authority", handler.HandleUserAuthorityGet)
	mux.Post("/authority/user/authority/grant", handler.HandleUserAuthorityGrant)
	mux.Post("/authority/user/authority/check", handler.HandleUserAuthorityCheck)
	mux.Get("/authority/authority/id/:id",handler.HandleQueryAuthority)
	mux.Get("/authority/authority/code/:code",handler.HandleQueryAuthority)
	mux.Get("/authority/authority/level/:group",handler.HandleQueryAuthority)
	mux.Get("/authority/authority/group/:code",handler.HandleQueryAuthorityGroupAll)
	mux.Get("/authority/authority/group",handler.HandleQueryAuthorityGroupAll)
	mux.Post("/authority/authority/register", handler.HandleCreateAuthority)
	mux.Post("/authority/authority/update/name", handler.HandleUpdateAuthorityName)
	http.Handle("/", mux)
}

func main() {
	
	//日志
	if err := libutil.TRLogger(Cfg.Log.File, Cfg.Log.Level, Cfg.Log.Name, Cfg.Log.Suffix, Cfg.Prog.Daemon); err != nil {
		fmt.Printf("init time rotate logger error: %s\n", err.Error())
		return
	}
	if Cfg.Prog.CPU == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU()) //配0就用所有核
	} else {
		runtime.GOMAXPROCS(Cfg.Prog.CPU)
	}

	logging.Debug("server start")

	libutil.InitSignal()

	logging.Debug("server init finish")
	/*go func() {
		err := http.ListenAndServe(Cfg.Prog.HealthPort, nil)

		fmt.Printf("err:%+v", err)
		if err != nil {
			logging.Error("ListenAndServe: %s", err.Error())
		}
	}()*/

	registerHttpHandle()

	go func() {
		err := http.ListenAndServe(Cfg.Server.PortInfo, nil)
		//err := http.ListenAndServeTLS(cfg.Server.PortInfo, "cert_server/server.crt",
		//"cert_server/server.key", nil)
		if err != nil {
			logging.Error("ListenAndServe port:%s failed", Cfg.Server.PortInfo)
		}
	}()

	file, err := libutil.DumpPanic("gsrv")
	if err != nil {
		logging.Error("init dump panic error: %s", err.Error())
	}

	defer func() {
		logging.Info("server stop...:%d", runtime.NumGoroutine())
		time.Sleep(time.Second)
		logging.Info("server stop...,ok")
		if err := libutil.ReviewDumpPanic(file); err != nil {
			logging.Error("review dump panic error: %s", err.Error())
		}

	}()
	<-libutil.ChanRunning

}
