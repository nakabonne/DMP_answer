package main

import (
	"net/http"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

type AnswerJson struct {
	QuestionId int `json:"question_id"`
	teamId string `json:"team_id"`
	Answer string `json: answer"`
}

const (
	project   = "ca-intern-201710-team02"
	instance  = "teamb-bigtable1"
	tableName = "latlon-table"
	family    = "Log"
	teamId    = "b"
	//pathToKeyFile = "ca-intern-201710-team02-4d5815ebcb43.json"
)

var (
	ctx = context.Background()
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, TeamB!\n")
	})
	e.GET("/place", readLogs)

	e.Logger.Fatal(e.Start(":8080"))
}

func readLogs(c echo.Context) error {
	// TODO: jsonを読み込んでrowKeyを作る&is_latlonとquestionIDを渡してもらう

	isLatLon := false
	questionId := 1
	rowKey := "2017-07-03 10:57:23 +0000 UTC#5e778b10-40a4-438f-9d62-74e85e7e5c32"
	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		log.Println(err)
	}
	tbl := client.Open(tableName)
	row, err := tbl.ReadRow(ctx, rowKey)
	if err != nil {
		log.Println(err)
	}
	lat := string(row[family][0].Value)
	lon := string(row[family][1].Value)
	log.Print("lat: " +lat)
	log.Print("lon: " +lon)

	if(isLatLon) {
		answer := lat + ":" + lon
		answerJson := &AnswerJson{questionId, teamId, answer}
		return c.JSON(http.StatusOK, answerJson)
	}

	// TODO: 位置計算関数にlat lonを渡す
	return c.String(http.StatusOK, "hogehoge")
}

