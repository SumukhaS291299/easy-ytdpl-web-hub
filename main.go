package main

import (
	"log"
	"net/http"

	litedb "ytdpldownloader/liteDB"
	"ytdpldownloader/utils"
	"ytdpldownloader/web"
)

func init() {
	utils.StartLogger()
	conf := utils.LoadConf()
	litedb.DB(conf.DBPath)
	utils.CheckBin(conf)
	utils.CheckUpdates()
	utils.SaveConfig(conf)
}

func main() {
	defer litedb.Dbcon.Close()
	router := web.NewRouter()

	log.Println("Starting server on http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
