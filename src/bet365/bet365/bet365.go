package bet365

import (
	"bytes"
	"chat"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	EVENT_UPDATE_TIME  = iota + 1 // 时间变动，以分为单位
	EVENT_CHANGE_STATE            // 状态改变
	EVENT_GOAL                    // 进球
	EVENT_CANCEL_GOAL             // 进球无效
)

const (
	STATUS_UNKNOWN    = -1
	STATUS_NONE       = 0 // 未开始
	STATUS_FIRSTHALF  = 1 // 上半场
	STATUS_MIDDLE     = 2 // 中场
	STATUS_SECONDHALF = 3 // 下半场
	STATUS_COMPLETE   = 4 // 完赛
)

const (
	RULE_334     = "334"
	RULE_7091    = "7091"
	RULE_757     = "757"
	RULE_HALF_05 = "half0.5"
	RULE_HALF_EQ = "halfeq"
)

var (
	FT_RESULT   = 1777  // 全场胜平负
	FT_HANDICAP = 10147 // 全场亚洲让分盘
	FT_GOALS    = 10148 // 全场大小球
	FT_ODDS     = []int{FT_RESULT, FT_HANDICAP, FT_GOALS}

	HT_RESULT   = 10161 // 半场胜平负
	HT_HANDICAP = 10170 // 半场亚洲让分盘
	HT_GOALS    = 10171 // 半场大小球
	HT_ODDS     = []int{HT_RESULT, HT_HANDICAP, HT_GOALS}

	RULES  = []string{RULE_334, RULE_7091, RULE_757, RULE_HALF_05, RULE_HALF_EQ}
	engine *xorm.Engine
)

const (
	FILTER_STATE_NONE  = 0
	FILTER_STATE_RED   = 1
	FILTER_STATE_BLACK = 2
)

type Bet365 struct {
	conn    *bet365conn
	data    *Bet365Data
	sendovm bool
	matchs  map[string]*Match
	filter  map[string]map[string]*Filter
}

func State(s int) string {
	switch s {
	case STATUS_UNKNOWN:
		return "未知"
	case STATUS_NONE:
		return "未开始"
	case STATUS_FIRSTHALF:
		return "上半场"
	case STATUS_MIDDLE:
		return "中场"
	case STATUS_SECONDHALF:
		return "下半场"
	case STATUS_COMPLETE:
		return "完赛"
	default:
		return "未知"
	}
}

func NewBet365() *Bet365 {
	b := new(Bet365)
	b.data = NewBet365Data("")
	b.data.AddNode(b.data.Root, NewSimpleNode("OVInPlay_10_0"))
	b.data.AddNode(b.data.Root, NewSimpleNode("OVM1"))
	b.data.AddNode(b.data.Root, NewSimpleNode("OVM2"))
	b.data.AddNode(b.data.Root, NewSimpleNode("OVM3"))
	b.conn = new(bet365conn)
	b.matchs = make(map[string]*Match)
	b.filter = make(map[string]map[string]*Filter)
	return b
}

func (b *Bet365) onMessage(data []byte) error {
	xx := bytes.Split(data, []byte(_DELIMITERS_MESSAGE))
	if xx[0][0] == '1' {
		err := b.conn.subscibe("CONFIG_10_0,OVInPlay_10_0")
		if err != nil {
			return err
		}
		return nil
	}

	for _, s := range xx {
		ss := bytes.Split(s, []byte(_DELIMITERS_RECORD))
		path := ss[0][1:]
		ParseData(b.data, path, ss[1])
		switch string(path) {
		case "OVInPlay_10_0":
			if !b.sendovm {
				for _, s := range []string{"OVM2", "OVM3"} {
					err := b.conn.subscibe(s) // 亚洲让分盘
					if err != nil {
						return err
					}
					time.Sleep(time.Second)
				}
				b.sendovm = true
			}
		}

	}

	return nil
}

func (b *Bet365) work() error {
	b.sendovm = false
	ch := make(chan struct{}, 0)
	go b.updateMatch(ch)
	var err error
	for {
		_, message, e := b.conn.ReadMessage()
		if e != nil {
			err = e
			break
		}
		err = b.onMessage(message)
		if err != nil {
			break
		}
	}

	close(ch)
	return err
}

func (b *Bet365) updateMatch(ch chan struct{}) {
	timer := time.NewTicker(time.Second)
L:
	for {
		select {
		case <-timer.C:
			b.process()
		case <-ch:
			break L
		}
	}
	timer.Stop()
}

func matchTime(lt time.Time, timestamp time.Time, tu, tm, ts string) (m, s int) {
	if tu == "" || tu == "19000101000000" {
		return 0, 0
	}

	t := formatTime(tu)
	d1 := t.Sub(timestamp)
	d2 := time.Now().Sub(lt)

	om, _ := strconv.Atoi(tm)
	os, _ := strconv.Atoi(ts)
	d := int((d2 - d1).Seconds()) + om*60 + os // 偏移量秒
	m = d / 60
	s = d % 60
	return
}

func (b *Bet365) addMatch(m *Match) {
	b.matchs[m.It] = m
	if _, ok := b.filter[m.It]; !ok {
		b.filter[m.It] = make(map[string]*Filter)
		for _, r := range RULES {
			f := new(Filter)
			if f.LoadFromDB(m.It, r) {
				b.filter[m.It][r] = f
			}
		}
	}
}

func (b *Bet365) process() {
	b.data.RLock()
	defer b.data.RUnlock()
	node := b.data.FindNode("OV_1_10_0")
	if node == nil {
		return
	}

	matchs := b.data.ChildByType(node, "EV")
	for _, m := range matchs {
		it := m.It()
		var match *Match
		var ok bool
		if match, ok = b.matchs[it]; !ok {
			match = new(Match)
			if match.Load(it) { // 尝试从数据库加载
				log.Printf("[加载] %s %s", match.LeagueName, match.TeamName)
				b.addMatch(match)
				match.Update(b.data, m)
				continue
			}

			if !match.Init(b.data, m) {
				//log.Println("match init failed")
				continue
			}

			b.addMatch(match)
			match.Insert()
			msg := fmt.Sprintf("[新增] %s %s", match.LeagueName, match.TeamName)
			log.Println(msg)
			chat.SendToBroadcast(msg)
			chat.SendToBroadcast(match.Preview())
			continue
		}

		events := match.Update(b.data, m)
		if m.Attr("VS") == "-1" { //比赛隐藏了
			continue
		}

		for _, e := range events {
			b.Snapshoot(e, match)
			switch e {
			case EVENT_UPDATE_TIME:
				b.Filter(match)
				b.CheckActive(match)
			default:
				b.CheckFilter(e, match)
			}
		}

	}

	dels := b.data.GetDel()
	if len(dels) > 0 {
		for _, it := range dels {
			if match, ok := b.matchs[it]; ok {
				delete(b.matchs, it)
				delete(b.filter, it)
				msg := fmt.Sprintf("[删除] %s %s %d-%d %s", match.LeagueName, match.TeamName, match.HoScore, match.GuScore, State(match.State))
				log.Println(msg)
				//chat.SendQQMessage(msg, "天气预报")
			}
		}
	}

	b.CheckOdd()
}

func (b *Bet365) CheckOdd() {
	for it, fs := range b.filter {
		match := b.matchs[it]
		if match == nil {
			continue
		}
		for _, f := range fs {
			if f.HalfState == STATUS_SECONDHALF && f.WaitOdd && f.FilterState == FILTER_STATE_NONE { // 只判断下半场
				if match.Size-f.AboveSize() > 0.1 {
					continue
				}
				f.WaitOdd = false
				f.Update("wait_odd")
				msg := f.MakeNoticeOdd()
				chat.SendToRecommend(msg)
			}
		}
	}
}

func (b *Bet365) Snapshoot(e int, m *Match) {
	switch e {
	case EVENT_UPDATE_TIME:
		if m.Min%5 == 0 {
			s := new(SnapShot)
			s.CopyFromMatch(m)
			s.Insert()
		}
	case EVENT_CHANGE_STATE:
		switch m.State {
		case STATUS_FIRSTHALF, STATUS_MIDDLE, STATUS_SECONDHALF, STATUS_COMPLETE:
			s := new(SnapShot)
			s.CopyFromMatch(m)
			s.Insert()
		}
	}
}

func (b *Bet365) CheckFilter(e int, m *Match) {
	switch e {
	case EVENT_CHANGE_STATE:
		if m.State == STATUS_FIRSTHALF {
			b.rulehalfeq(m)
		}
		b.CheckBlack(m)
		if m.State != STATUS_UNKNOWN {
			msg := fmt.Sprintf("[%s] %s %s %d-%d 平局概率:%d%%", State(m.State), m.LeagueName, m.TeamName, m.HoScore, m.GuScore, m.Dogfall())
			log.Println(msg)
			chat.SendToBroadcast(msg)
		}
	case EVENT_GOAL:
		b.CheckRed(m)
		msg := fmt.Sprintf("[进球] %s %s %d:%d %d-%d 平局概率:%d%%", m.LeagueName, m.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore, m.Dogfall())
		log.Println(msg)
		chat.SendToBroadcast(msg)
	case EVENT_CANCEL_GOAL:
		b.ResetRed(m)
		msg := fmt.Sprintf("[无效] %s %s %d:%d %d-%d", m.LeagueName, m.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore)
		log.Println(msg)
		chat.SendToBroadcast(msg)
	}
}

func (b *Bet365) CheckRed(m *Match) {
	if fs, ok := b.filter[m.It]; ok {
		for _, v := range fs {
			if v.HalfState == m.State && v.FilterState == FILTER_STATE_NONE && !v.Inactive {
				if m.Score() > v.Score() {
					v.FilterState = FILTER_STATE_RED
					v.Update("filter_state")
					msg := v.MakeResultMessage(false, m)
					log.Println(msg)
					chat.SendToRecommend(msg)
				}
			}
		}
	}
}

func (b *Bet365) ResetRed(m *Match) {
	if fs, ok := b.filter[m.It]; ok {
		for _, v := range fs {
			if v.HalfState == m.State && v.FilterState != FILTER_STATE_NONE {
				if m.Score() == v.Score() {
					v.FilterState = FILTER_STATE_NONE
					v.Update("filter_state")
					msg := v.MakeResultMessage(true, m)
					log.Println(msg)
					chat.SendToRecommend(msg)
				}
			}
		}
	}
}

func (b *Bet365) CheckBlack(m *Match) {
	if fs, ok := b.filter[m.It]; ok {
		for _, v := range fs {
			if m.State == STATUS_MIDDLE && v.HalfState == STATUS_FIRSTHALF { // 中场
				if v.FilterState == FILTER_STATE_NONE && !v.Inactive {
					if m.Score() > v.Score() {
						v.FilterState = FILTER_STATE_RED
					} else {
						v.FilterState = FILTER_STATE_BLACK
					}
					v.Update("filter_state")
					msg := v.MakeResultMessage(false, m)
					log.Println(msg)
					chat.SendToRecommend(msg)
				}
			}

			if m.State == STATUS_COMPLETE && v.HalfState == STATUS_SECONDHALF { // 完赛
				if v.FilterState == FILTER_STATE_NONE && !v.Inactive {
					if m.Score() > v.Score() {
						v.FilterState = FILTER_STATE_RED
					} else {
						v.FilterState = FILTER_STATE_BLACK
					}
					v.Update("filter_state")
					msg := v.MakeResultMessage(false, m)
					log.Println(msg)
					chat.SendToRecommend(msg)
				}
			}
		}
	}
}

func (b *Bet365) Filter(m *Match) {
	b.rulehalf05(m)
	b.rule334(m)
	b.rule7091(m)
	b.rule757(m)
}

func (b *Bet365) CheckActive(m *Match) {
	if fs, ok := b.filter[m.It]; ok {
		for _, v := range fs {
			if !v.Inactive {
				continue
			}
			v.CheckActive(m)
		}
	}
}

func (b *Bet365) rulehalfeq(m *Match) {
	if m.Min != 0 {
		return
	}

	f := new(Filter)
	f.Rule = RULE_HALF_EQ
	if f.LoadFromDB(m.It, f.Rule) {
		return
	}

	f.Inactive = true

	if m.AvgEq > 3.0 && m.AvgEq < 3.7 {
		avgeq := m.AvgEq * 0.618
		if m.HalfAvgEq > 2.05 && m.HalfAvgEq-avgeq > 0.15 && math.Abs(m.HalfLet) > 0.24 {
			f.HalfState = STATUS_FIRSTHALF
			f.FilterState = FILTER_STATE_NONE
			f.CopyFromMatch(m)
			f.Insert()
			b.filter[m.It][f.Rule] = f
			msg := fmt.Sprintf("/闪电注意 \n%s \n%s \n 经评估，上半场破蛋概率较大，请关注。", m.LeagueName, m.TeamName)
			log.Println(msg)
			chat.SendToRecommend(msg)
		}
	}

}

func (b *Bet365) rulehalf05(m *Match) {
	if m.Min != 20 || m.Score() != 0 {
		return
	}

	f := new(Filter)
	f.Rule = RULE_HALF_05
	if f.LoadFromDB(m.It, f.Rule) {
		return
	}

	//f.Inactive = true
	// 大小盘中水以上，大小盘2球以上， 降一个盘以内， 初盘大小球2.25-3.0之间,
	// 初盘让0.5以上
	if m.SizeBig > 1.85 &&
		m.Size > 2.0 &&
		m.FirstSize-m.Size < 0.251 &&
		m.FirstSize > 2.249 &&
		m.FirstSize < 3.01 &&
		math.Abs(m.FirstLet) > 0.49 {
		f.HalfState = STATUS_FIRSTHALF
		f.FilterState = FILTER_STATE_NONE
		f.CopyFromMatch(m)
		f.Insert()
		b.filter[m.It][f.Rule] = f
		msg := f.MakeRuleMessage()
		log.Println(msg)
		chat.SendToRecommend(msg)
	}

}

func (b *Bet365) rule334(m *Match) {
	if m.Min != 30 || m.Score() != 3 {
		return
	}
	f := new(Filter)
	f.Rule = RULE_334
	if f.LoadFromDB(m.It, f.Rule) {
		return
	}

	f.HalfState = STATUS_FIRSTHALF
	f.FilterState = FILTER_STATE_NONE
	f.CopyFromMatch(m)
	f.Insert()
	b.filter[m.It][f.Rule] = f
	msg := f.MakeRuleMessage()
	log.Println(msg)
	chat.SendToRecommend(msg)
}

func (b *Bet365) rule7091(m *Match) {
	if m.Min != 70 {
		return
	}

	if math.Abs(m.Let) > 0.249 &&
		math.Abs(float64(m.HoScore-m.GuScore)) < 1.01 &&
		m.SizeBig > 1.85 &&
		m.SizeBig < 2.0 {

		f := new(Filter)
		f.Rule = RULE_7091
		if f.LoadFromDB(m.It, f.Rule) {
			return
		}

		f.HalfState = STATUS_SECONDHALF
		f.FilterState = FILTER_STATE_NONE
		f.CopyFromMatch(m)
		f.Insert()
		b.filter[m.It][f.Rule] = f
		msg := f.MakeRuleMessage()
		log.Println(msg)
		chat.SendToRecommend(msg)
	}
}

func (b *Bet365) rule757(m *Match) {
	if m.Min != 75 || m.HoScore+m.GuScore != 6 {
		return
	}
	if m.HoScore+m.GuScore == 6 {
		f := new(Filter)
		f.Rule = RULE_757
		if f.LoadFromDB(m.It, f.Rule) {
			return
		}

		f.HalfState = STATUS_SECONDHALF
		f.FilterState = FILTER_STATE_NONE
		f.CopyFromMatch(m)
		f.Insert()
		b.filter[m.It][f.Rule] = f
		msg := f.MakeRuleMessage()
		log.Println(msg)
		chat.SendToRecommend(msg)
	}
}

func Run(addr string, origin string, getcookieurl string) {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./bet365.db")
	if err != nil {
		panic(err)
	}

	engine.Sync2(new(Match), new(Filter), new(SnapShot))

	chat.SendToRecommend("初始化，数据源:365, 规则:334, 7091, 757, half0.5(测试) 测试模式")
	bet := NewBet365()
	for {
		err := bet.conn.Connect(addr, origin, getcookieurl)
		if err != nil {
			log.Printf("connect %s failed, waiting 3 seconds to retry", addr)
			time.Sleep(time.Second * 3)
			continue
		}

		//chat.SendQQMessage("连接成功")
		log.Println("connected")
		err = bet.work()
		if err != nil {
			bet.conn.close()
			log.Printf("catch err: %s, waiting 3 seconds to reconnect", err)
			//chat.SendQQMessage("系统异常，3秒后重新连接")
			time.Sleep(time.Second * 3)
			continue
		}
	}
}
