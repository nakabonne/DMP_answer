package main

import (
    "net/http"
    "log"
    "fmt"
    "io/ioutil"
    "os"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "cloud.google.com/go/bigtable"
    "golang.org/x/net/context"
    "google.golang.org/api/option"
    "golang.org/x/oauth2/google"
    "math"
    "strconv"
)
type AnswerJson struct {
    QuestionId int `json:"question_id"`
    TeamId string `json:"team_id"`
    Answer string `json: answer"`
}
type Question struct{
  QuestionId int `json:"question_id"`
  Idfa string `json:"idfa"`
  Timestamp string `json:"timestamp"`
  IsLatlon bool `json:"is_latlon"`
}
const (
    project   = "ca-intern-201710-team02"
    instance  = "teamb-bigtable1"
    tableName = "latlon-table"
    family    = "Log"
    teamId    = "b"
    pathToKeyFile = "ca-intern-201710-team02-4d5815ebcb43.json"
)
var (
    ctx = context.Background()
    client *bigtable.Client
    table *bigtable.Table
)
type Staiton struct {
    Name string
    Lat  float64
    Lon  float64
}
var stations = []*Staiton{
    &Staiton{
        "Osaki",
        35.6197,
        139.728553,
    },
    &Staiton{
        "Gotanda",
        35.626446,
        139.723444,
    },
    &Staiton{
        "Meguro",
        35.633998,
        139.715828,
    },
    &Staiton{
        "Ebisu",
        35.64669,
        139.710106,
    },
    &Staiton{
        "Shibuya",
        35.658517,
        139.701334,
    },
    &Staiton{
        "Harajuku",
        35.670168,
        139.702687,
    },
    &Staiton{
        "Yoyogi",
        35.683061,
        139.702042,
    },
    &Staiton{
        "Shinjuku",
        35.690921,
        139.700258,
    },
    &Staiton{
        "Shinokubo",
        35.701306,
        139.700044,
    },
    &Staiton{
        "Takadanobaba",
        35.712285,
        139.703782,
    },
    &Staiton{
        "Mejiro",
        35.721204,
        139.706587,
    },
    &Staiton{
        "Ikebukuro",
        35.728926,
        139.71038,
    },
    &Staiton{
        "Otsuka",
        35.73159,
        139.729329,
    },
    &Staiton{
        "Sugamo",
        35.733492,
        139.739345,
    },
    &Staiton{
        "Komagome",
        35.736489,
        139.746875,
    },
    &Staiton{
        "Tabata",
        35.738062,
        139.76086,
    },
    &Staiton{
        "Nishinippori",
        35.732135,
        139.766787,
    },
    &Staiton{
        "Nippori",
        35.727772,
        139.770987,
    },
    &Staiton{
        "Uguisudani",
        35.720495,
        139.778837,
    },
    &Staiton{
        "Ueno",
        35.713768,
        139.777254,
    },
    &Staiton{
        "Okachimachi",
        35.707893,
        139.774332,
    },
    &Staiton{
        "Akihabara",
        35.698683,
        139.774219,
    },
    &Staiton{
        "Kanda",
        35.69169,
        139.770883,
    },
    &Staiton{
        "Tokyo",
        35.681382,
        139.766084,
    },
    &Staiton{
        "Yurakucho",
        35.675069,
        139.763328,
    },
    &Staiton{
        "Shimbashi",
        35.665498,
        139.75964,
    },
    &Staiton{
        "Hamamatsucho",
        35.655646,
        139.756749,
    },
    &Staiton{
        "Tamachi",
        35.645736,
        139.747575,
    },
    &Staiton{
        "Shinagawa",
        35.630152,
        139.74044,
    },
}
func authenticate() (*bigtable.Client, error) {
    jsonKey, err := ioutil.ReadFile(pathToKeyFile)
    if err != nil {
        return nil, err
    }
    config, err := google.JWTConfigFromJSON(jsonKey, bigtable.Scope)
    if err != nil {
        return nil, err
    }
    client, err := bigtable.NewClient(ctx, project, instance, option.WithTokenSource(config.TokenSource(ctx)))
    if err != nil {
        return nil, err
    }
    return client, nil
}
func isDevelop() bool {
    return os.Getenv("DEV") == "1"
}
func openBigtable(tableName string) (table *bigtable.Table, err error) {
    if isDevelop() {
        client, err = authenticate()
    } else {
        client, err = bigtable.NewClient(ctx, project, instance)
    }
    if err != nil {
        log.Fatal(err)
    }
    table = client.Open(tableName)
    return
}
func getNearStation(latitude string, longitude string) string {
    var (
        minDistance float64
        nearStation string
    )
    radius := 6378.137 //赤道半径
    inputLat, _ := strconv.ParseFloat(latitude, 64)
    inputLon, _ := strconv.ParseFloat(longitude, 64)
    inputLat = inputLat * math.Pi / 180
    inputLon = inputLon * math.Pi / 180
    for i, station := range stations {
        lat := station.Lat * math.Pi / 180
        lon := station.Lon * math.Pi / 180
        distance := radius * math.Acos(math.Sin(inputLat)*math.Sin(lat)+math.Cos(inputLat)*math.Cos(lat)*math.Cos(lon-inputLon))
        if minDistance > distance || i == 0 {
            minDistance = distance
            nearStation = station.Name
        }
    }
    return nearStation
}

func readLogs(c echo.Context) error {
    // TODO: jsonを読み込んでrowKeyを作る&is_latlonとquestionIDを渡してもらう
    q := new(Question)
    if err := c.Bind(q); err != nil {
        return err
    }
    log.Println(q)
    isLatLon := q.IsLatlon
    questionId := q.QuestionId
    rowKey := q.Timestamp + "#" + q.Idfa
    row, err := table.ReadRow(ctx, rowKey)
    if err != nil {
        log.Println(err)
    }
    if len(row) == 0{
         answerJson := &AnswerJson{questionId, teamId, "Shibuya"}
         return c.JSON(http.StatusOK, answerJson)
    } else {
        lat := string(row[family][0].Value)
        lon := string(row[family][1].Value)
        log.Print("lat: " +lat)
        log.Print("lon: " +lon)
        if(isLatLon) {
            answer := lat + ":" + lon
            answerJson := &AnswerJson{questionId, teamId, answer}
            log.Println(answerJson)
            return c.JSON(http.StatusOK, answerJson)
        } else {
            stationname := getNearStation(lat, lon)
            fmt.Println(stationname)
            answer := stationname
            answerJson := &AnswerJson{questionId, teamId, answer}
            log.Println(answerJson)
            return c.JSON(http.StatusOK, answerJson)
        }
    }
}
func init() {
    var err error
    table, err = openBigtable("latlon-table")
    if err != nil {
        log.Fatal(err)
    }
}
func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, TeamB!\n")
    })
    e.POST("/place", readLogs)
    e.Logger.Fatal(e.Start(":8080"))
}
