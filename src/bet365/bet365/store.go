package bet365

import (
	"chat"
	"fmt"
	"log"
	"math"
	"sort"
)

type SnapShot struct {
	Id               int64
	It               string `xorm:"index"`
	Mid              string
	LeagueName       string
	TeamName         string
	Min              int
	Sec              int
	State            int
	HoScore          int     // 主队进球
	HoHalfScore      int     // 主队半场得分
	GuScore          int     // 客队进球
	GuHalfScore      int     // 客队半场得分
	HoRed            int     // 主队红牌
	HoYellow         int     // 主队黄牌
	HoCo             int     // 主队角球
	HoHfCo           int     // 主队半场角球
	GuRed            int     // 客队红牌
	GuYellow         int     // 客队黄牌
	GuCo             int     // 客队角球
	GuHfCo           int     // 客队半场角球
	HalfLet          float64 // 上半场数据
	HalfLetHm        float64
	HalfLetAw        float64
	HalfAvgHm        float64
	HalfAvgAw        float64
	HalfAvgEq        float64
	HalfSize         float64
	HalfSizeBig      float64
	HalfSizeSma      float64
	HalfFirstLet     float64
	HalfFirstLetHm   float64
	HalfFirstLetAw   float64
	HalfFirstAvgHm   float64
	HalfFirstAvgAw   float64
	HalfFirstAvgEq   float64
	HalfFirstSize    float64
	HalfFirstSizeBig float64
	HalfFirstSizeSma float64
	Let              float64 // 下半场数据
	LetHm            float64
	LetAw            float64
	AvgHm            float64
	AvgAw            float64
	AvgEq            float64
	Size             float64
	SizeBig          float64
	SizeSma          float64
	FirstLet         float64
	FirstLetHm       float64
	FirstLetAw       float64
	FirstAvgHm       float64
	FirstAvgAw       float64
	FirstAvgEq       float64
	FirstSize        float64
	FirstSizeBig     float64
	FirstSizeSma     float64
}

func (s *SnapShot) Insert() {
	_, err := engine.Insert(s)
	if err != nil {
		panic(err)
	}
}

func (s *SnapShot) CopyFromMatch(m *Match) {
	s.It = m.It
	s.Mid = m.Mid
	s.LeagueName = m.LeagueName
	s.TeamName = m.TeamName
	s.Min = m.Min
	s.Sec = m.Sec
	s.State = m.State
	s.HoScore = m.HoScore
	s.HoHalfScore = m.HoHalfScore
	s.GuScore = m.GuScore
	s.GuHalfScore = m.GuHalfScore
	s.HoRed = m.HoRed
	s.HoYellow = m.HoYellow
	s.HoCo = m.HoCo
	s.HoHfCo = m.HoHfCo
	s.GuRed = m.GuRed
	s.GuYellow = m.GuYellow
	s.GuCo = m.GuCo
	s.GuHfCo = m.GuHfCo
	s.HalfLet = m.HalfLet
	s.HalfLetHm = m.HalfLetHm
	s.HalfLetAw = m.HalfLetAw
	s.HalfAvgHm = m.HalfAvgHm
	s.HalfAvgAw = m.HalfAvgAw
	s.HalfAvgEq = m.HalfAvgEq
	s.HalfSize = m.HalfSize
	s.HalfSizeBig = m.HalfSizeBig
	s.HalfSizeSma = m.HalfSizeSma
	s.HalfFirstLet = m.HalfFirstLet
	s.HalfFirstLetHm = m.HalfFirstLetHm
	s.HalfFirstLetAw = m.HalfFirstLetAw
	s.HalfFirstAvgHm = m.HalfFirstAvgHm
	s.HalfFirstAvgAw = m.HalfFirstAvgAw
	s.HalfFirstAvgEq = m.HalfFirstAvgEq
	s.HalfFirstSize = m.HalfFirstSize
	s.HalfFirstSizeBig = m.HalfFirstSizeBig
	s.HalfFirstSizeSma = m.HalfFirstSizeSma

	s.Let = m.Let
	s.LetHm = m.LetHm
	s.LetAw = m.LetAw
	s.AvgHm = m.AvgHm
	s.AvgAw = m.AvgAw
	s.AvgEq = m.AvgEq
	s.Size = m.Size
	s.SizeBig = m.SizeBig
	s.SizeSma = m.SizeSma
	s.FirstLet = m.FirstLet
	s.FirstLetHm = m.FirstLetHm
	s.FirstLetAw = m.FirstLetAw
	s.FirstAvgHm = m.FirstAvgHm
	s.FirstAvgAw = m.FirstAvgAw
	s.FirstAvgEq = m.FirstAvgEq
	s.FirstSize = m.FirstSize
	s.FirstSizeBig = m.FirstSizeBig
	s.FirstSizeSma = m.FirstSizeSma
}

type Filter struct {
	Id               int64
	Rule             string
	It               string `xorm:"index"`
	Mid              string
	LeagueName       string
	TeamName         string
	Min              int
	Sec              int
	State            int
	HoScore          int     // 主队进球
	HoHalfScore      int     // 主队半场得分
	GuScore          int     // 客队进球
	GuHalfScore      int     // 客队半场得分
	HoRed            int     // 主队红牌
	HoYellow         int     // 主队黄牌
	HoCo             int     // 主队角球
	HoHfCo           int     // 主队半场角球
	GuRed            int     // 客队红牌
	GuYellow         int     // 客队黄牌
	GuCo             int     // 客队角球
	GuHfCo           int     // 客队半场角球
	HalfLet          float64 // 上半场数据
	HalfLetHm        float64
	HalfLetAw        float64
	HalfAvgHm        float64
	HalfAvgAw        float64
	HalfAvgEq        float64
	HalfSize         float64
	HalfSizeBig      float64
	HalfSizeSma      float64
	HalfFirstLet     float64
	HalfFirstLetHm   float64
	HalfFirstLetAw   float64
	HalfFirstAvgHm   float64
	HalfFirstAvgAw   float64
	HalfFirstAvgEq   float64
	HalfFirstSize    float64
	HalfFirstSizeBig float64
	HalfFirstSizeSma float64
	Let              float64 // 下半场数据
	LetHm            float64
	LetAw            float64
	AvgHm            float64
	AvgAw            float64
	AvgEq            float64
	Size             float64
	SizeBig          float64
	SizeSma          float64
	FirstLet         float64
	FirstLetHm       float64
	FirstLetAw       float64
	FirstAvgHm       float64
	FirstAvgAw       float64
	FirstAvgEq       float64
	FirstSize        float64
	FirstSizeBig     float64
	FirstSizeSma     float64
	FilterState      int
	HalfState        int
	WaitOdd          bool
	Inactive         bool  // 未激活
	Created          int64 `xorm:"created"`
	extra            int
}

func (f *Filter) CheckActive(m *Match) {
	if f.Inactive {
		switch f.Rule {
		case RULE_HALF_05:
			if m.Min == 30 &&
				m.Score() == 0 {
				f.Inactive = false
				f.Update("inactive")
				msg := f.MakeRuleMessage()
				log.Println(msg)
				chat.SendToRecommend(msg)
			}
		case RULE_HALF_EQ:
			if m.Min >= 20 &&
				m.State == STATUS_FIRSTHALF &&
				m.Score() == 0 && m.HalfSize < 0.51 && m.HalfSizeBig > 1.9 &&
				math.Abs(m.HalfLet) > 0.24 {
				f.Inactive = false
				f.Update("inactive")
				msg := f.MakeRuleMessage()
				log.Println(msg)
				chat.SendToRecommend(msg)
			}
		case RULE_LZ_001:
			if m.Min >= 70 && m.Min <= 85 {
				sub := math.Abs(m.AvgHm - m.AvgAw)
				if sub < 0.001 || (sub > 0.99 && sub < 1.001) {
					if f.extra == 1 {
						f.Inactive = false
						f.HoScore = m.HoScore
						f.GuScore = m.GuScore
						f.Update("inactive", "ho_score", "gu_score")
						msg := f.MakeRuleMessage()
						log.Println(msg)
						chat.SendToRecommend(msg)
						return
					}
					f.extra = 1
				}

			}
		}

	}
}

func (f *Filter) MakeResultMessage(reset bool, m *Match) string {
	if reset {
		return fmt.Sprintf(`/流泪 [%s] 
%s
%s
[%02d:%02d] 进球无效，比分:%d-%d`, f.Rule, f.LeagueName, f.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore)
	}

	if f.FilterState == FILTER_STATE_RED {
		return fmt.Sprintf(`/红包 [红] [%s] 
%s
%s
[%02d:%02d] /足球 比分:%d-%d`, f.Rule, f.LeagueName, f.TeamName, m.Min, m.Sec, m.HoScore, m.GuScore)
	}
	if f.FilterState == FILTER_STATE_BLACK {
		if f.HalfState == STATUS_FIRSTHALF {
			return fmt.Sprintf(`/炸弹 [黑] [%s] 
%s
%s
上半场结束，比分:%d-%d`, f.Rule, f.LeagueName, f.TeamName, m.HoScore, m.GuScore)
		}
		return fmt.Sprintf(`/炸弹 [黑] [%s] 
%s
%s  
比赛结束，比分:%d-%d`, f.Rule, f.LeagueName, f.TeamName, m.HoScore, m.GuScore)
	}
	return ""
}

func (f *Filter) Dogfall() int {
	ratio := math.Abs(f.FirstAvgEq - 3.3)
	if ratio > 2.3 {
		ratio = (ratio - 2.3) / 20
		if ratio > 1 {
			ratio = 0.99
		}
		ratio = (1 - ratio) * 50
	} else {
		ratio = 50 + (1-(ratio/2.3))*40
	}

	return int(ratio)
}

func (f *Filter) AboveSize() float64 {
	return float64(f.HoScore+f.GuScore) + 0.5
}

func (f *Filter) MakeRuleMessage() string {
	rescore := f.AboveSize()
	var half string
	if f.HalfState == STATUS_FIRSTHALF {
		half = "半场"
	} else {
		half = "全场"
	}
	s := fmt.Sprintf(`/足球[%s] 
%s
%s
当前比分:%d-%d
平局概率:%d%%
推荐:%s大%.1f
id:%s`, f.Rule, f.LeagueName, f.TeamName, f.HoScore, f.GuScore, f.Dogfall(), half, rescore, f.It)
	if f.HalfState == STATUS_SECONDHALF && f.Size-rescore > 0.1 { // 当前盘比目标盘大
		f.WaitOdd = true
		s += fmt.Sprintf("\n/闪电注意：当前盘口(%.2f)高于推荐盘口,可等水", f.Size)
	}
	return s
}

func (f *Filter) MakeNoticeOdd() string {
	return fmt.Sprintf(`/开车[%s]
%s
%s
降盘啦，快上车
id:%s`, f.Rule, f.LeagueName, f.TeamName, f.It)
}

func (f *Filter) Insert() {
	_, err := engine.Insert(f)
	if err != nil {
		panic(err)
	}
}

func (f *Filter) Update(col ...string) {
	engine.Id(f.Id).Cols(col...).Update(f)
}

func (f *Filter) Score() int {
	return f.HoScore + f.GuScore
}

func (f *Filter) LoadFromDB(it string, rule string) bool {
	b, err := engine.Where("it=? and rule=?", it, rule).Get(f)
	if err != nil {
		return false
	}
	return b
}

func (f *Filter) CopyFromMatch(m *Match) {
	f.It = m.It
	f.Mid = m.Mid
	f.LeagueName = m.LeagueName
	f.TeamName = m.TeamName
	f.Min = m.Min
	f.Sec = m.Sec
	f.State = m.State
	f.HoScore = m.HoScore
	f.HoHalfScore = m.HoHalfScore
	f.GuScore = m.GuScore
	f.GuHalfScore = m.GuHalfScore
	f.HoRed = m.HoRed
	f.HoYellow = m.HoYellow
	f.HoCo = m.HoCo
	f.HoHfCo = m.HoHfCo
	f.GuRed = m.GuRed
	f.GuYellow = m.GuYellow
	f.GuCo = m.GuCo
	f.GuHfCo = m.GuHfCo
	f.HalfLet = m.HalfLet
	f.HalfLetHm = m.HalfLetHm
	f.HalfLetAw = m.HalfLetAw
	f.HalfAvgHm = m.HalfAvgHm
	f.HalfAvgAw = m.HalfAvgAw
	f.HalfAvgEq = m.HalfAvgEq
	f.HalfSize = m.HalfSize
	f.HalfSizeBig = m.HalfSizeBig
	f.HalfSizeSma = m.HalfSizeSma
	f.HalfFirstLet = m.HalfFirstLet
	f.HalfFirstLetHm = m.HalfFirstLetHm
	f.HalfFirstLetAw = m.HalfFirstLetAw
	f.HalfFirstAvgHm = m.HalfFirstAvgHm
	f.HalfFirstAvgAw = m.HalfFirstAvgAw
	f.HalfFirstAvgEq = m.HalfFirstAvgEq
	f.HalfFirstSize = m.HalfFirstSize
	f.HalfFirstSizeBig = m.HalfFirstSizeBig
	f.HalfFirstSizeSma = m.HalfFirstSizeSma
	f.Let = m.Let
	f.LetHm = m.LetHm
	f.LetAw = m.LetAw
	f.AvgHm = m.AvgHm
	f.AvgAw = m.AvgAw
	f.AvgEq = m.AvgEq
	f.Size = m.Size
	f.SizeBig = m.SizeBig
	f.SizeSma = m.SizeSma
	f.FirstLet = m.FirstLet
	f.FirstLetHm = m.FirstLetHm
	f.FirstLetAw = m.FirstLetAw
	f.FirstAvgHm = m.FirstAvgHm
	f.FirstAvgAw = m.FirstAvgAw
	f.FirstAvgEq = m.FirstAvgEq
	f.FirstSize = m.FirstSize
	f.FirstSizeBig = m.FirstSizeBig
	f.FirstSizeSma = m.FirstSizeSma
}

type Match struct {
	Id               int64
	It               string `xorm:"index"`
	Mid              string
	LeagueName       string
	TeamName         string
	Min              int
	Sec              int
	State            int
	HoScore          int     // 主队进球
	HoHalfScore      int     // 主队半场得分
	GuScore          int     // 客队进球
	GuHalfScore      int     // 客队半场得分
	HoRed            int     // 主队红牌
	HoYellow         int     // 主队黄牌
	HoCo             int     // 主队角球
	HoHfCo           int     // 主队半场角球
	GuRed            int     // 客队红牌
	GuYellow         int     // 客队黄牌
	GuCo             int     // 客队角球
	GuHfCo           int     // 客队半场角球
	HalfLet          float64 // 上半场数据
	HalfLetHm        float64
	HalfLetAw        float64
	HalfAvgHm        float64
	HalfAvgAw        float64
	HalfAvgEq        float64
	HalfSize         float64
	HalfSizeBig      float64
	HalfSizeSma      float64
	HalfFirstLet     float64
	HalfFirstLetHm   float64
	HalfFirstLetAw   float64
	HalfFirstAvgHm   float64
	HalfFirstAvgAw   float64
	HalfFirstAvgEq   float64
	HalfFirstSize    float64
	HalfFirstSizeBig float64
	HalfFirstSizeSma float64
	Let              float64 // 全场数据
	LetHm            float64
	LetAw            float64
	AvgHm            float64
	AvgAw            float64
	AvgEq            float64
	Size             float64
	SizeBig          float64
	SizeSma          float64
	FirstLet         float64
	FirstLetHm       float64
	FirstLetAw       float64
	FirstAvgHm       float64
	FirstAvgAw       float64
	FirstAvgEq       float64
	FirstSize        float64
	FirstSizeBig     float64
	FirstSizeSma     float64
}

func (m *Match) Score() int {
	return m.HoScore + m.GuScore
}

func (m *Match) Dogfall() int {
	ratio := math.Abs(m.FirstAvgEq - 3.3)
	if ratio > 2.3 {
		ratio = (ratio - 2.3) / 20
		if ratio > 1 {
			ratio = 0.99
		}
		ratio = (1 - ratio) * 50
	} else {
		ratio = 50 + (1-(ratio/2.3))*40
	}

	return int(ratio)
}

func (m *Match) Load(it string) bool {
	b, err := engine.Where("it=?", it).Get(m)
	if err != nil {
		return false
	}
	return b
}

func (m *Match) Insert() {
	_, err := engine.Insert(m)
	if err != nil {
		panic(err)
	}
}

func (m *Match) Preview() string {
	return fmt.Sprintf(`%s 
%s  
	胜平负:%.2f,%.2f,%.2f 
	让分:%.2f,%.2f,%.2f 
	大小盘:%.2f,%.2f,%.2f
  上半场:
	胜平负:%.2f,%.2f,%.2f 
	让分:%.2f,%.2f,%.2f
	大小盘:%.2f,%.2f,%.2f
id:%s`,
		m.LeagueName, m.TeamName,
		m.FirstAvgHm, m.FirstAvgEq, m.FirstAvgAw,
		m.FirstLet, m.FirstLetHm, m.FirstLetAw,
		m.FirstSize, m.FirstSizeBig, m.FirstSizeSma,
		m.HalfFirstAvgHm, m.HalfFirstAvgEq, m.HalfFirstAvgAw,
		m.HalfFirstLet, m.HalfFirstLetHm, m.HalfFirstLetAw,
		m.HalfFirstSize, m.HalfFirstSizeBig, m.HalfFirstSizeSma,
		m.It,
	)
}

func (m *Match) String() string {
	return fmt.Sprintf(`%s
%s 
时间:[%02d:%02d] 比分:%d-%d
胜平负:%.2f,%.2f,%.2f
让分:%.2f,%.2f,%.2f
大小盘:%.2f,%.2f,%.2f
初盘：
 全场:
  胜平负:%.2f,%.2f,%.2f
  让分:%.2f,%.2f,%.2f
  大小盘:%.2f,%.2f,%.2f
 上半场:
  胜平负:%.2f,%.2f,%.2f
  让分:%.2f,%.2f,%.2f
  大小盘:%.2f,%.2f,%.2f`,
		m.LeagueName, m.TeamName,
		m.Min, m.Sec,
		m.HoScore, m.GuScore,
		m.AvgHm, m.AvgEq, m.AvgAw,
		m.Let, m.LetHm, m.LetAw,
		m.Size, m.SizeBig, m.SizeSma,
		m.FirstAvgHm, m.FirstAvgEq, m.FirstAvgAw,
		m.FirstLet, m.FirstLetHm, m.FirstLetAw,
		m.FirstSize, m.FirstSizeBig, m.FirstSizeSma,
		m.HalfFirstAvgHm, m.HalfFirstAvgEq, m.HalfFirstAvgAw,
		m.HalfFirstLet, m.HalfFirstLetHm, m.HalfFirstLetAw,
		m.HalfFirstSize, m.HalfFirstSizeBig, m.HalfFirstSizeSma,
	)
}

func (m *Match) Init(d *Bet365Data, node *Node) bool {
	pnode := node.parent
	if pnode != nil && pnode.tag == "CT" {
		m.LeagueName = pnode.Attr("NA")
	} else {
		m.LeagueName = node.Attr("CT")
	}

	m.TeamName = node.Attr("NA")
	m.It = node.It()
	id := node.Attr("ID")
	id = id[0 : len(id)-len("A_10_0")]
	m.Mid = id

	m.Min, m.Sec = d.MatchTime(node)
	m.State = node.State()
	m.HoScore, m.GuScore = node.SS()

	if m.State < STATUS_MIDDLE {
		for _, ot := range HT_ODDS {
			mait := fmt.Sprintf("OVM2-%s-%d", m.Mid, ot)
			ma := d.FindNode(mait)
			if ma == nil {
				continue
			}
			oddnode := d.ChildByType(ma, "PA")
			if len(oddnode) == 0 {
				break
			}
			ns := NewSortNode("ID", oddnode)
			sort.Sort(ns)
			switch ot {
			case HT_RESULT: // 胜平负
				m.HalfFirstAvgHm = oddnode[0].Odd()
				m.HalfFirstAvgEq = oddnode[1].Odd()
				m.HalfFirstAvgAw = oddnode[2].Odd()
			case HT_HANDICAP: // 让球
				m.HalfFirstLet = oddnode[0].Float("HA")
				m.HalfFirstLetHm = oddnode[0].Odd()
				m.HalfFirstLetAw = oddnode[1].Odd()
			case HT_GOALS:
				m.HalfFirstSize = oddnode[0].Float("HA")
				m.HalfFirstSizeBig = oddnode[0].Odd()
				m.HalfFirstSizeSma = oddnode[1].Odd()
			}
		}
	}

	noodds := true
	for _, ot := range FT_ODDS {
		mait := fmt.Sprintf("OVM3-%s-%d", m.Mid, ot)
		ma := d.FindNode(mait)
		if ma == nil {
			continue
		}
		noodds = false
		oddnode := d.ChildByType(ma, "PA")
		if len(oddnode) == 0 {
			continue
		}
		ns := NewSortNode("ID", oddnode)
		sort.Sort(ns)
		switch ot {
		case FT_RESULT: // 胜平负
			m.FirstAvgHm = oddnode[0].Odd()
			m.FirstAvgEq = oddnode[1].Odd()
			m.FirstAvgAw = oddnode[2].Odd()
		case FT_HANDICAP: // 让球
			m.FirstLet = oddnode[0].Float("HA")
			m.FirstLetHm = oddnode[0].Odd()
			m.FirstLetAw = oddnode[1].Odd()
		case FT_GOALS:
			m.FirstSize = oddnode[0].Float("HA")
			m.FirstSizeBig = oddnode[0].Odd()
			m.FirstSizeSma = oddnode[1].Odd()
		}
	}

	if noodds {
		return false
	}

	m.HalfAvgHm = m.HalfFirstAvgHm
	m.HalfAvgEq = m.HalfFirstAvgEq
	m.HalfAvgAw = m.HalfFirstAvgAw
	m.HalfLet = m.HalfFirstLet
	m.HalfLetHm = m.HalfFirstLetHm
	m.HalfLetAw = m.HalfFirstLetAw
	m.HalfSize = m.HalfFirstSize
	m.HalfSizeBig = m.HalfFirstSizeBig
	m.HalfSizeSma = m.HalfFirstSizeSma

	m.AvgHm = m.FirstAvgHm
	m.AvgEq = m.FirstAvgEq
	m.AvgAw = m.FirstAvgAw
	m.Let = m.FirstLet
	m.LetHm = m.FirstLetHm
	m.LetAw = m.FirstLetAw
	m.Size = m.FirstSize
	m.SizeBig = m.FirstSizeBig
	m.SizeSma = m.FirstSizeSma

	return true
}

func (m *Match) Update(d *Bet365Data, node *Node) []int {

	var event []int
	min, sec := d.MatchTime(node)
	if m.Min != min {
		event = append(event, EVENT_UPDATE_TIME)
	}
	m.Min = min
	m.Sec = sec
	s := node.State()
	if m.State != s {
		m.State = s
		event = append(event, EVENT_CHANGE_STATE)
	}

	oldhs, oldgs := node.SS()
	if m.HoScore != oldhs || m.GuScore != oldgs {
		if m.HoScore+m.GuScore > oldhs+oldgs {
			event = append(event, EVENT_CANCEL_GOAL) // 进球无效
		} else {
			event = append(event, EVENT_GOAL) // 进球
		}
		m.HoScore, m.GuScore = oldhs, oldgs
	}
	if m.State < STATUS_MIDDLE {
		for _, id := range HT_ODDS {
			mait := fmt.Sprintf("OVM2-%s-%d", m.Mid, id)
			ma := d.FindNode(mait)
			if ma == nil {
				break
			}
			oddnode := d.ChildByType(ma, "PA")
			if len(oddnode) == 0 {
				break
			}
			ns := NewSortNode("ID", oddnode)
			sort.Sort(ns)
			switch id {
			case HT_RESULT: // 胜平负
				m.HalfAvgHm = oddnode[0].Odd()
				m.HalfAvgEq = oddnode[1].Odd()
				m.HalfAvgAw = oddnode[2].Odd()
			case HT_HANDICAP: // 让球
				m.HalfLet = oddnode[0].Float("HA")
				m.HalfLetHm = oddnode[0].Odd()
				m.HalfLetAw = oddnode[1].Odd()
			case HT_GOALS:
				m.HalfSize = oddnode[0].Float("HA")
				m.HalfSizeBig = oddnode[0].Odd()
				m.HalfSizeSma = oddnode[1].Odd()
			}
		}
	}

	for _, id := range FT_ODDS {
		mait := fmt.Sprintf("OVM3-%s-%d", m.Mid, id)
		ma := d.FindNode(mait)
		if ma == nil {
			return event
		}
		oddnode := d.ChildByType(ma, "PA")
		if len(oddnode) == 0 {
			continue
		}
		ns := NewSortNode("ID", oddnode)
		sort.Sort(ns)
		switch id {
		case FT_RESULT: // 胜平负
			m.AvgHm = oddnode[0].Odd()
			m.AvgEq = oddnode[1].Odd()
			m.AvgAw = oddnode[2].Odd()
		case FT_HANDICAP: // 让球
			m.Let = oddnode[0].Float("HA")
			m.LetHm = oddnode[0].Odd()
			m.LetAw = oddnode[1].Odd()
		case FT_GOALS:
			m.Size = oddnode[0].Float("HA")
			m.SizeBig = oddnode[0].Odd()
			m.SizeSma = oddnode[1].Odd()
		}
	}

	return event
}

type Attention struct {
	Id   int64
	Team string `xorm:"index"`
}

func (a *Attention) Insert() {
	_, err := engine.Insert(a)
	if err != nil {
		panic(err)
	}
}

func (a *Attention) Remove() {
	engine.Where("team=?", a.Team).Delete(a)
}

func (a *Attention) Find(team string) bool {
	b, err := engine.Where("team=?", team).Get(a)
	if err != nil {
		return false
	}
	return b
}
