package bet365

import (
	"bytes"
	"chat"
	"config"
	"container/list"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
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

const (
	WAIT_TIME      = 1
	WAIT_HALF_SIZE = 2 // 半场大小
	WAIT_FULL_SIZE = 3 // 全场大小
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
	lock    sync.Mutex
	conn    *bet365conn
	data    *Bet365Data
	sendovm bool
	matchs  map[string]*Match
	filter  map[string]map[string]*Filter
	notify  map[string]*list.List
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

type Notify struct {
	Group    string
	Member   string
	WaitType int
	WaitTime int
	WaitSize float64
	WaitBig  float64
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
	b.notify = make(map[string]*list.List)
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

func (b *Bet365) Match(it string) *Match {
	if m, ok := b.matchs[it]; ok {
		return m
	}
	return nil
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

	if _, ok := b.notify[m.It]; !ok {
		b.notify[m.It] = list.New()
	}
}

func (b *Bet365) AddNotify(group, member string, it string, typ int, time int, size float64, big float64) string {
	b.lock.Lock()
	defer b.lock.Unlock()
	n := new(Notify)
	n.Group = group
	n.Member = member
	n.WaitType = typ
	n.WaitTime = time
	n.WaitSize = size
	n.WaitBig = big
	if l, ok := b.notify[it]; ok {
		m := b.Match(it)
		if m == nil {
			return "[error] 通知增加失败，比赛没有找到"
		}
		if n.WaitType != WAIT_TIME {
			if m.State == STATUS_FIRSTHALF {
				n.WaitType = WAIT_HALF_SIZE
			} else {
				n.WaitType = WAIT_FULL_SIZE
			}
		}
		l.PushBack(n)
		return fmt.Sprintf("增加 %s %s 通知成功", m.LeagueName, m.TeamName)
	}
	return "[error] 通知增加失败，比赛没有找到"
}

func (b *Bet365) CheckNotify(m *Match) {
	if l, ok := b.notify[m.It]; ok {
		ele := l.Front()
		for ele != nil {
			e := ele
			ele = ele.Next()
			n := e.Value.(*Notify)
			switch n.WaitType {
			case WAIT_TIME:
				if m.Min >= n.WaitTime {
					msg := fmt.Sprintf("@%s\n%s\n%s\n当前时间[%d:%d]", n.Member, m.LeagueName, m.TeamName, m.Min, m.Sec)
					chat.SendQQMessage(msg, n.Group)
					l.Remove(e)
				}
			case WAIT_HALF_SIZE:
				if m.State == STATUS_FIRSTHALF && m.HalfSize <= n.WaitSize+0.01 &&
					m.HalfSizeBig >= n.WaitBig {
					msg := fmt.Sprintf("@%s\n%s\n%s\n上半场[%d:%d]\n大小盘:%.2f, %.2f", n.Member, m.LeagueName, m.TeamName, m.Min, m.Sec, m.HalfSize, m.HalfSizeBig)
					chat.SendQQMessage(msg, n.Group)
					l.Remove(e)
				}
			case WAIT_FULL_SIZE:
				if m.State == STATUS_SECONDHALF && m.Size <= n.WaitSize+0.01 &&
					m.SizeBig >= n.WaitBig {
					msg := fmt.Sprintf("@%s\n%s\n%s\n全场[%d:%d]\n大小盘:%.2f, %.2f", n.Member, m.LeagueName, m.TeamName, m.Min, m.Sec, m.Size, m.SizeBig)
					chat.SendQQMessage(msg, n.Group)
					l.Remove(e)
				}
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
				b.StateFilter(e, match)
				b.CheckFilter(e, match)
			}
		}

		b.CheckNotify(match)

	}

	dels := b.data.GetDel()
	if len(dels) > 0 {
		for _, it := range dels {
			if match, ok := b.matchs[it]; ok {
				delete(b.matchs, it)
				delete(b.filter, it)
				if match.State == STATUS_COMPLETE {
					delete(b.notify, it)
				}
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

func (b *Bet365) StateFilter(e int, m *Match) {
	switch e {
	case EVENT_CHANGE_STATE:
		switch m.State {
		case STATUS_FIRSTHALF:
			for k, rs := range config.Setting.Rule.State {
				if k == "firsthalf" {
					for _, r := range rs {
						switch r {
						case RULE_HALF_EQ:
							b.rulehalfeq(m)
						}
					}
				}
			}
		}
	}
}

func (b *Bet365) CheckFilter(e int, m *Match) {
	switch e {
	case EVENT_CHANGE_STATE:
		b.CheckBlack(m)
		if m.State != STATUS_UNKNOWN {
			msg := fmt.Sprintf("[%s] %s %s %d-%d 平局概率:%d%%\nid:%s", State(m.State), m.LeagueName, m.TeamName, m.HoScore, m.GuScore, m.Dogfall(), m.It)
			log.Println(msg)
			chat.SendToBroadcast(msg)
		}
	case EVENT_GOAL:
		b.CheckRed(m)
		msg := fmt.Sprintf("[进球] %s %s %d:%d %d-%d 平局概率:%d%%\nid:%s", m.LeagueName, m.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore, m.Dogfall(), m.It)
		log.Println(msg)
		chat.SendToBroadcast(msg)
	case EVENT_CANCEL_GOAL:
		b.ResetRed(m)
		msg := fmt.Sprintf("[无效] %s %s %d:%d %d-%d \nid:%s", m.LeagueName, m.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore, m.It)
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
	for _, r := range config.Setting.Rule.Update {
		switch r {
		case RULE_HALF_05:
			b.rulehalf05(m)
		case RULE_334:
			b.rule334(m)
		case RULE_7091:
			b.rule7091(m)
		case RULE_757:
			b.rule757(m)
		}
	}
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
			msg := fmt.Sprintf("/闪电注意 \n%s \n%s \n 经评估，上半场破蛋概率较大，请关注。\nid:%s", m.LeagueName, m.TeamName, m.It)
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

var (
	bet *Bet365
)

func Stat() string {
	var result []string
	n := time.Now()
	last := n.AddDate(0, 0, -1)
	lzero, _ := time.ParseInLocation("2006-01-02", last.Format("2006-01-02"), time.Local)
	l24, _ := time.ParseInLocation("2006-01-02", n.Format("2006-01-02"), time.Local)
	result = append(result, "今日:")
	for _, v := range RULES {
		var f Filter
		total, err := engine.Where("rule=? and created > ? and Inactive=0", v, l24.Unix()).Count(&f)
		if err != nil {
			continue
		}
		red, err := engine.Where("rule=? and filter_state=? and created > ? and Inactive=0", v, FILTER_STATE_RED, l24.Unix()).Count(&f)
		result = append(result, fmt.Sprintf("[%s] 总:%d 红:%d 命中率：%.1f", v, total, red, float64(red)/float64(total)*100))
	}
	result = append(result, "昨日:")
	for _, v := range RULES {
		var f Filter
		total, err := engine.Where("rule=? and created > ? and created < ? and Inactive=0", v, lzero.Unix(), l24.Unix()).Count(&f)
		if err != nil {
			continue
		}
		red, err := engine.Where("rule=? and filter_state=? and created > ? and created < ? and Inactive=0", v, FILTER_STATE_RED, lzero.Unix(), l24.Unix()).Count(&f)
		result = append(result, fmt.Sprintf("[%s] 总:%d 红:%d 命中率：%.1f", v, total, red, float64(red)/float64(total)*100))
	}

	result = append(result, "总评:")
	for _, v := range RULES {
		var f Filter
		total, err := engine.Where("rule=? and Inactive=0", v).Count(&f)
		if err != nil {
			continue
		}
		red, err := engine.Where("rule=? and filter_state=? and Inactive=0", v, FILTER_STATE_RED).Count(&f)
		result = append(result, fmt.Sprintf("[%s] 总:%d 红:%d 命中率：%.1f", v, total, red, float64(red)/float64(total)*100))
	}

	return strings.Join(result, "\n")
}

func AddTimeNotify(group, member string, it string, time string) string {
	i, err := strconv.Atoi(time)
	if err != nil {
		return err.Error()
	}
	return bet.AddNotify(group, member, it, WAIT_TIME, i, 0, 0)
}

func AddSizeNotify(group, member string, it string, size, big string) string {
	fsize, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return err.Error()
	}
	fbig, err := strconv.ParseFloat(big, 64)
	if err != nil {
		return err.Error()
	}
	return bet.AddNotify(group, member, it, -1, 0, fsize, fbig)
}

func Run(addr string, origin string, getcookieurl string) {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./bet365.db")
	if err != nil {
		panic(err)
	}

	engine.Sync2(new(Match), new(Filter), new(SnapShot))

	engine.DatabaseTZ = time.Local
	engine.TZLocation = time.Local

	chat.SendToRecommend("初始化")
	bet = NewBet365()
	delay := time.Second * 3
	retrys := 0
	for {
		err := bet.conn.Connect(addr, origin, getcookieurl)
		if err != nil {
			log.Printf("connect %s failed, err:%s", addr, err)
			retrys++
			if retrys > 3 {
				delay = delay + time.Second
				if delay > time.Minute {
					delay = time.Minute
				}
			}
			time.Sleep(delay)
			continue
		}

		retrys = 0
		//chat.SendQQMessage("连接成功")
		log.Println("connected")
		err = bet.work()
		if err != nil {
			bet.conn.close()
			log.Printf("catch err: %s, waiting 3 seconds to reconnect", err)
			continue
		}
	}
}
