package ybf

import (
	"chat"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	UpdateTime = time.Second * 3
)

var (
	list   *YBF
	engine *xorm.Engine
)

func GetList(cookie string) []byte {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://live.13322.com/allScore/list", strings.NewReader("lang=zh"))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Host", "live.13322.com")
	req.Header.Add("Origin", "https://live.13322.com")
	req.Header.Add("Referer", "https://live.13322.com/jsbf.html")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(string(body))
	//ioutil.WriteFile("list.json", body, 0x666)
	return body
}

func UpdateInfo(cookie string) []byte {
	client := &http.Client{}
	var r http.Request
	r.ParseForm()
	date := time.Now().Format("2006-01-02")
	r.Form.Add("date", date)
	bodystr := strings.TrimSpace(r.Form.Encode())

	req, err := http.NewRequest("POST", "https://live.13322.com/common/ajaxRefreshByDate", strings.NewReader(bodystr))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Host", "live.13322.com")
	req.Header.Add("Origin", "https://live.13322.com")
	req.Header.Add("Referer", "https://live.13322.com/jsbf.html")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}

func GetCookie() string {
	r, err := http.Get("https://live.13322.com/jsbf.html")
	if err != nil {
		panic(err)
	}

	c := r.Cookies()
	cookie := ""
	sep := ""
	for k, v := range c {
		if k > 0 {
			sep = ";"
		}
		cookie += fmt.Sprintf("%s%s=%s", sep, v.Name, v.Value)
	}
	return cookie
}

func updateList(cookie string) {
	for {
		data := GetList(cookie)
		if data == nil {
			fmt.Println("get list failed")
			time.Sleep(time.Second * 1)
			continue
		}

		y, err := ParseList(data)
		if err != nil {
			fmt.Println("parse list failed", err)
			continue
		}

		list = y
		break
	}
}

func process(m *Match, u BasicUpdate) {

	halfscore := false
	score := false
	oldScore := m.HoScore + m.GuScore
	if u.HoScore != nil {
		if m.HoScore != *u.HoScore { // 主队进球
			m.HoScore = *u.HoScore
			score = true
		}
	}
	if u.HoHalfScore != nil {
		if m.HoHalfScore != *u.HoHalfScore { // 主队半场进球
			m.HoHalfScore = *u.HoHalfScore
			halfscore = true
		}
	}
	if u.GuScore != nil {
		if m.GuScore != *u.GuScore { // 客队进球
			m.GuScore = *u.GuScore
			score = true
		}
	}
	if u.GuHalfScore != nil {
		if m.GuHalfScore != *u.GuHalfScore { // 客队半场进球
			m.GuHalfScore = *u.GuHalfScore
			halfscore = true
		}
	}

	if oldScore < m.HoScore+m.GuScore && (score || halfscore) {
		notice(m, m.Status == STATUS_FirstHalf)
	}

	if u.HoRed != nil {
		if m.HoRed != *u.HoRed {
			m.HoRed = *u.HoRed
		}
	}
	if u.HoYellow != nil {
		if m.HoYellow != *u.HoYellow {
			m.HoYellow = *u.HoYellow
		}
	}
	if u.HoCo != nil {
		if m.HoCo != *u.HoCo {
			m.HoCo = *u.HoCo
		}
	}
	if u.HoHfCo != nil {
		if m.HoHfCo != *u.HoHfCo {
			m.HoHfCo = *u.HoHfCo
		}
	}
	if u.GuRed != nil {
		if m.GuRed != *u.GuRed {
			m.GuRed = *u.GuRed
		}
	}
	if u.GuYellow != nil {
		if m.GuYellow != *u.GuYellow {
			m.GuYellow = *u.GuYellow
		}
	}
	if u.GuCo != nil {
		if m.GuCo != *u.GuCo {
			m.GuCo = *u.GuCo
		}
	}
	if u.GuHfCo != nil {
		if m.GuHfCo != *u.GuHfCo {
			m.GuHfCo = *u.GuHfCo
		}
	}
	if u.Min != nil {
		if m.Min != *u.Min {
			m.Min = *u.Min
			filter(m)
			switch m.Min {
			case 10, 20, 30, 40, 60, 70, 75, 80, 85:
				snapshoot(m)
			}
		}
	}
	if m.Status != u.Status {
		m.Status = u.Status
		switch m.Status {
		case STATUS_FirstHalf, STATUS_MIDDLE, STATUS_SecondHALF, STATUS_COMPLETE:
			noticeLose(m, m.Status)
			snapshoot(m)
		}
	}
}

func noticeLose(m *Match, state int) {

	fs := make([]Filter, 0)
	if err := engine.Where("match_id=? and notice=?", m.MatchId, 0).Find(&fs); err != nil {
		return
	}

	switch state {
	case STATUS_MIDDLE:
		for _, f := range fs {
			if f.Rule == "334" {
				chat.SendQQMessage(fmt.Sprintf("[黑] %s %s VS %s 上半场结束，上半场比分：%d-%d", f.LeagueName, f.HoTeamName, f.GuTeamName, m.HoHalfScore, m.GuHalfScore))
			}
		}
	case STATUS_COMPLETE:
		for _, f := range fs {
			if f.Rule != "334" {
				chat.SendQQMessage(fmt.Sprintf("[黑] %s %s VS %s 比赛结束，最终比分：%d-%d", f.LeagueName, f.HoTeamName, f.GuTeamName, m.HoScore, m.GuScore))
			}
		}
	}
}

func notice(m *Match, half bool) {
	fs := make([]Filter, 0)
	if err := engine.Where("match_id=? and notice=?", m.MatchId, 0).Find(&fs); err != nil {
		return
	}

	if half {
		for _, f := range fs {
			if f.Rule == "334" && f.HoScore+f.GuScore < m.HoScore+m.GuScore {
				chat.SendQQMessage(fmt.Sprintf("[红] %s %s VS %s 进球， 比分：%d-%d", f.LeagueName, f.HoTeamName, f.GuTeamName, m.HoScore, m.GuScore))
				f.Notice = 1
				engine.Id(f.Id).Cols("notice").Update(f)
				time.Sleep(time.Millisecond * 100)
			}
		}
		return
	}

	for _, f := range fs {
		if f.Rule != "334" && f.HoScore+f.GuScore < m.HoScore+m.GuScore {
			chat.SendQQMessage(fmt.Sprintf("[红] %s %s VS %s 进球， 比分：%d-%d", f.LeagueName, f.HoTeamName, f.GuTeamName, m.HoScore, m.GuScore))
			f.Notice = 1
			engine.Id(f.Id).Cols("notice").Update(f)
			time.Sleep(time.Second)
		}
	}

}

func filter(m *Match) {
	rule334(m)
	rule7091(m)
	rule757(m)
}

func makeMsg(f *Filter, info string) string {
	msg := fmt.Sprintf("[%s] 评分:%d %s %s VS %s 比分:%d-%d 红牌:%d-%d \n推荐：%s\n", f.Rule, f.Score, f.LeagueName, f.HoTeamName, f.GuTeamName,
		f.HoScore, f.GuScore,
		f.HoRed, f.GuRed, info)
	return msg
}

func rule757(m *Match) {
	if m.Min != 75 && m.Min != 76 {
		return
	}

	f := new(Filter)
	f.Rule = "757"
	f.MatchId = m.MatchId

	if f.Exist() {
		return
	}

	if m.HoScore+m.GuScore != 6 {
		return
	}

	f.Score = 50
	f.CopyFrom(m)

	engine.Insert(f)
	msg := makeMsg(f, "全场大6.5")
	fmt.Println(msg)
	chat.SendQQMessage(msg)

}

func rule334(m *Match) {
	if m.Min != 30 && m.Min != 31 {
		return
	}

	f := new(Filter)
	f.Rule = "334"
	f.MatchId = m.MatchId

	if f.Exist() {
		return
	}

	if m.HoScore+m.GuScore != 3 {
		return
	}

	f.Score = 50
	f.CopyFrom(m)

	engine.Insert(f)
	msg := makeMsg(f, "半场大3.5")
	fmt.Println(msg)
	chat.SendQQMessage(msg)

}

func rule7091(m *Match) {
	if m.Min != 70 && m.Min != 71 {
		return
	}

	f := new(Filter)
	f.Rule = "7091"
	f.MatchId = m.MatchId

	if f.Exist() {
		return
	}

	score := 50
	fLet, _, _, _, _, fAvgEq, _, _, _ := splitJson(m.Firstodds.StringJson)
	Let, LetHm, _, _, _, _, Size, SizeBig, _ := splitJson(m.Odds.StringJson)
	if LetHm < 1.3 && SizeBig > 0.85 && SizeBig < 1 && math.Abs(float64(m.HoScore-m.GuScore)) < 1.0001 {
		if Let < -0.25 {
			score += 5
		}

		if fLet < -0.251 {
			score += 5
		}

		if fAvgEq > 3.31 {
			score += 5
		}

		if Let < -0.51 {
			score += 5
		}

		if Size > 0.51 {
			score += 5
		}

		if m.HoCo+m.GuCo > 7 {
			score += 5
		}
	} else {
		return
	}

	f.Score = score
	f.CopyFrom(m)
	engine.Insert(f)
	msg := makeMsg(f, fmt.Sprintf("全场大%d.5", f.HoScore+f.GuScore))
	fmt.Println(msg)
	chat.SendQQMessage(msg)
}

func snapshoot(m *Match) {
	ss := new(SnapShoot)
	ss.MatchId = m.MatchId
	ss.Status = m.Status
	ss.Min = m.Min
	ss.LeagueName = m.LeagueName
	ss.LeagueSimpName = m.LeagueSimpName
	ss.HoTeamName = m.HoTeamName
	ss.GuTeamName = m.GuTeamName
	ss.HoScore = m.HoScore
	ss.HoHalfScore = m.HoHalfScore
	ss.GuScore = m.GuScore
	ss.GuHalfScore = m.GuHalfScore
	ss.HoRed = m.HoRed
	ss.HoYellow = m.HoYellow
	ss.HoCo = m.HoCo
	ss.HoHfCo = m.HoHfCo
	ss.GuRed = m.GuRed
	ss.GuYellow = m.GuYellow
	ss.GuCo = m.GuCo
	ss.GuHfCo = m.GuHfCo
	ss.Let, ss.LetHm, ss.LetAw, ss.AvgHm, ss.AvgAw, ss.AvgEq, ss.Size, ss.SizeBig, ss.SizeSma = splitJson(m.Odds.StringJson)
	ss.FirstLet, ss.FirstLetHm, ss.FirstLetAw, ss.FirstAvgHm, ss.FirstAvgAw, ss.FirstAvgEq, ss.FirstSize, ss.FirstSizeBig, ss.FirstSizeSma = splitJson(m.Firstodds.StringJson)
	_, err := engine.Insert(ss)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("insert to db ", id)
}

func processUpdate(cookie string) error {
	udata := UpdateInfo(cookie)
	if udata == nil {
		return fmt.Errorf("get update failed")
	}
	u, err := ParseUpdate(udata)
	if err != nil {
		return fmt.Errorf("get update failed")
	}

	needUpdateList := false
	for _, v := range u.Base {
		m := list.Get(v.Id)
		if m == nil {
			needUpdateList = true
			continue
		}

		sid := fmt.Sprintf("%d", v.Id)
		if odd, ok := u.Odds[sid]; ok {
			m.Odds.Update(odd)
			//m.Odds.StringJson = odd // 更新盘口数据
		}
		process(m, v)
	}

	if needUpdateList {
		updateList(cookie)
	}

	return nil
}

func Run() {
	chat.SendQQMessage("初始化")
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./ybf.db")
	if err != nil {
		panic(err)
	}

	engine.Sync2(new(SnapShoot), new(Filter))
	cookie := GetCookie()
	updateList(cookie)
	chat.SendQQMessage(fmt.Sprintf("更新比赛信息完成，共%d场比赛,启用规则:[334,7091,757]", len(list.Matches)))
	time.Sleep(time.Second * 3)
	chat.SendQQMessage("初始化完成，开始扫盘")
	for {
		processUpdate(cookie)
		time.Sleep(UpdateTime)
	}
}
