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
	conf = utils.LoadConf()
	utils.StartLogger(conf)
	litedb.DB(conf.DBPath)
	utils.CheckOutDir(conf)
	utils.CheckBin(conf)
	utils.CheckUpdates()
	utils.SaveConfig(conf)
}

func main() {
	defer litedb.Dbcon.Close()
	router := web.NewRouter()

	log.Info("Starting server on", "\nAddress:\t", conf.ListenAddr)

	if err := http.ListenAndServe(conf.ListenAddr, router); err != nil {
		log.Fatal(err)
	}
}
