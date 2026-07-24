package main

import (
	"net/http"

	litedb "ytdpldownloader/liteDB"
	"ytdpldownloader/utils"
	"ytdpldownloader/web"

	"github.com/charmbracelet/log"
)

var conf utils.AppData

func init() {
	utils.StartLogger()
	conf = utils.LoadConf()
	litedb.DB(conf.DBPath)
	utils.CheckBin(conf)
	utils.CheckUpdates()
	utils.SaveConfig(conf)
}

func main() {
	defer litedb.Dbcon.Close()
	router := web.NewRouter()

	log.Info("Starting server on", conf.ListenAddr)

	if err := http.ListenAndServe(conf.ListenAddr, router); err != nil {
		log.Fatal(err)
	}
}
