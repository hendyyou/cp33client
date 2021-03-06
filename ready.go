package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/henrylee2cn/surfer"
)

func getCron(gameId int) {
	//	dt := dataTime{}
	//	err := Db.Model(&dt).Where("type=? and action_time<?", gameId, time.Now().Format("15:04:05")).Order("action_time DESC").Limit(1).Select()
	//	if err == nil && dt.Type >= 0 {
	//		//fmt.Println("getCron() 20:", dt.ActionTime)
	//	} else if err != nil {
	//		fmt.Println("getCron() 28:" + err.Error())
	//		return
	//	}
	//	//fmt.Println(dt.ActionTime)
	//	var ttActionTime time.Time
	//	ttActionTime, err = time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" "+dt.ActionTime, time.Local)
	//	if err != nil {
	//		fmt.Println("getCron() 29:", err.Error())
	//		return
	//	}
	//	diffTime := time.Now().Local().Unix() - ttActionTime.Unix()
	//	if diffTime > 120 || diffTime < -15 { //2分钟采集不到就别浪费时间了
	//		//fmt.Println("getCron() 34:"+"不在采集时间内跳过！", dt.Type, "	", diffTime)
	//		return
	//	}
	var err error
	var resp *http.Response
	resp, err = surfer.Download(&surfer.Request{
		Url: "http://127.0.0.1:8080/apiMyself/" + "1",
		//DownloaderID: 1,
	})
	if err != nil {
		fmt.Println("getCron()43:", err.Error())
		return
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("getCron()50:", err.Error())
		return
	}

	var openInfo OpenInfo
	err = json.Unmarshal(b, &openInfo)
	if err != nil {
		fmt.Println("getCron()61:", err.Error())
		return
	}
	if openInfo.Last_open != "" {
		return
	}
	//fmt.Println("getCron()64:", openInfo.Last_period, "	", diffTime)
	openInfo.checkIsGot(strconv.Itoa(openInfo.Last_period))
}

func (self *OpenInfo) checkIsGot(period string) {
	if redisClient.Exists("Client_"+strconv.Itoa(self.Type)+"_"+strconv.Itoa(self.Last_period)).Val() == 0 { //不存在
		switch self.Type {
		case 1:
			fmt.Println("checkIsGot():1		", time.Now(), "		period=", period)
			go cqssc_163_com(period)
			go cqssc_cqcp_net()
			time.Sleep(time.Second * 4)
		default:
			fmt.Println("checkIsGot():default")
		}
	}
}
