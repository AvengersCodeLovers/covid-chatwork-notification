package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AvengersCodeLovers/covid-chatwork-notification/util"
)

var (
	chatworkBaseUrl = "https://api.chatwork.com/v2/"
	covidBaseUrl    = "https://api.coronatracker.com/v3/stats/worldometer/"
)

type CovidModel struct {
	TotalConfirmed int `json:"totalConfirmed"`
	TotalDeaths    int `json:"totalDeaths"`
	TotalRecovered int `json:"totalRecovered"`
	TotalCritical  int `json:"totalCritical"`
	ActiveCases    int `json:"activeCases"`
	DailyConfirmed int `json:"dailyConfirmed"`
}

func init() {
	util.LoadEnvVars()
}

func main() {
	apiKey := util.GetEnv("CHATWORK_TOKEN", "")
	roomID := util.GetEnv("CHATWORK_ROOM", "")
	currentTime := time.Now()

	chatworkClient := util.New(apiKey, chatworkBaseUrl)
	covidClient := util.New("", covidBaseUrl)

	covdRes := covidClient.Get("/country", map[string]string{
		"countryCode": "VN",
	})

	var covidModel []CovidModel
	json.Unmarshal(covdRes, &covidModel)

	var (
		totalConfirmed = covidModel[0].TotalConfirmed
		totalDeaths    = covidModel[0].TotalDeaths
		totalRecovered = covidModel[0].TotalRecovered
		totalCritical  = covidModel[0].TotalCritical
		dailyConfirmed = covidModel[0].DailyConfirmed
		activeCases    = covidModel[0].ActiveCases
	)

	chatworkClient.Post("/rooms/"+roomID+"/messages", map[string]string{
		"body": fmt.Sprintf(`>> All [info][title]Tình hình dịch bệnh covid tại việt nam ngày %v [/title]- Số ca nhiễm thời điểm hiện tại : %d người
- Số ca tử vong : %d người
- Số ca bình phục : %d người
- Số người có chuyển biến nặng : %d người
- Số ca tăng trong ngày : %d người
- Số ca đang điều trị: %d người[/info]
		[info][title]Các biện pháp phòng chống COVID – 19[/title]- Hạn chế tiếp xúc trực tiếp với người bị bệnh đường hô hấp cấp tính (sốt, ho, khó thở); khi cần thiết phải đeo khẩu trang y tế đúng cách và giữ khoảng cách trên 02 mét khi tiếp xúc.
- Người có các triệu chứng sốt, ho, khó thở không nên đi du lịch hoặc đến nơi tập trung đông người. Thông báo ngay cho cơ quan y tế khi có các triệu chứng kể trên.
- Giữ ấm cơ thể, tăng cường sức khỏe bằng ăn uống, nghỉ ngơi, sinh hoạt hợp lý, luyện tập thể thao.
- Hiện tại, người dân có thể liên hệ 2 số điện thoại đường dây nóng của Bộ Y tế để cung cấp thông tin về bệnh Viêm đường hô hấp cấp Covid – 19 là: 1900 3228 và 1900 9095.[/info]
`, currentTime.Format("02-01-2006"), totalConfirmed, totalDeaths, totalRecovered, totalCritical, dailyConfirmed, activeCases),
	})
}
