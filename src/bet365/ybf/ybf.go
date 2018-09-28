package ybf

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	STATUS_NONE       = 0   // 未开始
	STATUS_FirstHalf  = 1   // 上半场
	STATUS_MIDDLE     = 2   // 中场
	STATUS_SecondHALF = 3   // 下半场
	STATUS_COMPLETE   = -1  // 完赛
	STATUS_DELAY      = -14 // 推迟

)

type Odds struct {
	Let        float64 `json:"-"`
	LetHm      float64 `json:"-"`
	LetAw      float64 `json:"-"`
	AvgHm      float64 `json:"-"`
	AvgAw      float64 `json:"-"`
	AvgEq      float64 `json:"-"`
	Size       float64 `json:"-"`
	SizeBig    float64 `json:"-"`
	SizeSma    float64 `json:"-"`
	StringJson string  `json:"stringJson"`
}

func (o *Odds) Update(s string) {
	Let, LetHm, LetAw, AvgHm, AvgAw, AvgEq, Size, SizeBig, SizeSma := splitJson(s)
	if !math.IsNaN(Let) {
		o.Let = Let
		o.LetHm = LetHm
		o.LetAw = LetAw
	}
	if !math.IsNaN(AvgHm) {
		o.AvgHm = AvgHm
		o.AvgAw = AvgAw
		o.AvgEq = AvgEq
	}
	if !math.IsNaN(Size) {
		o.Size = Size
		o.SizeBig = SizeBig
		o.SizeSma = SizeSma
	}

	o.StringJson = fmt.Sprintf("%f,%f,%f|%f,%f,%f|%f,%f,%f", o.Let, o.LetHm, o.LetAw, o.Size, o.SizeBig, o.SizeSma, o.AvgHm, o.AvgAw, o.AvgEq)
}

func splitOdd(s string) (x, y, z float64) {
	x, y, z = math.NaN(), math.NaN(), math.NaN()
	s1 := strings.Split(s, ",")
	if len(s1) != 3 {
		return
	}

	if s1[0] != "" {
		f, err := strconv.ParseFloat(s1[0], 32)
		if err == nil {
			x = f
		}
	}
	if s1[1] != "" {
		f, err := strconv.ParseFloat(s1[1], 32)
		if err == nil {
			y = f
		}
	}
	if s1[2] != "" {
		f, err := strconv.ParseFloat(s1[2], 32)
		if err == nil {
			z = f
		}
	}
	return
}
func splitJson(j string) (Let, LetHm, LetAw, AvgHm, AvgAw, AvgEq, Size, SizeBig, SizeSma float64) {
	//1,0.85,0.91|2.5,0.96,0.80|1.46,5.50,3.90
	s := strings.Split(j, "|")
	if len(s) != 3 {
		return
	}

	Let, LetHm, LetAw = splitOdd(s[0])
	Size, SizeBig, SizeSma = splitOdd(s[1])
	AvgHm, AvgAw, AvgEq = splitOdd(s[2])
	return
}

type Match struct {
	MatchId        int    `json:"id"`
	Status         int    `json:"status"`
	Time           string `json:"time"`
	Min            int    `json:"min"`
	LeagueName     string `json:"leagueName"`
	LeagueSimpName string `json:"leagueSimpName"`
	HoTeamName     string `json:"hoTeamName"`
	GuTeamName     string `json:"guTeamName"`
	HoScore        int    `json:"hoScore"`     // 主队进球
	HoHalfScore    int    `json:"hoHalfScore"` // 主队半场得分
	GuScore        int    `json:"guScore"`     // 客队进球
	GuHalfScore    int    `json:"guHalfScore"` // 客队半场得分
	HoRed          int    `json:"hoRed"`       // 主队红牌
	HoYellow       int    `json:"hoYellow"`    // 主队黄牌
	HoCo           int    `json:"hoCo"`        // 主队角球
	HoHfCo         int    `json:"hoHfCo"`      // 主队半场角球
	GuRed          int    `json:"guRed"`       // 客队红牌
	GuYellow       int    `json:"guYellow"`    // 客队黄牌
	GuCo           int    `json:"guCo"`        // 客队角球
	GuHfCo         int    `json:"guHfCo"`      // 客队半场角球
	Odds           Odds   `json:"odds"`        // 赔率
	Firstodds      Odds   `json:"firstodds"`   // 初盘赔率
}

type Filter struct {
	Id             int64
	Rule           string
	MatchId        int
	Score          int
	Status         int
	Time           string
	Min            int
	LeagueName     string
	LeagueSimpName string
	HoTeamName     string
	GuTeamName     string
	HoScore        int
	HoHalfScore    int
	GuScore        int
	GuHalfScore    int
	HoRed          int
	HoYellow       int
	HoCo           int
	HoHfCo         int
	GuRed          int
	GuYellow       int
	GuCo           int
	GuHfCo         int
	Odds           string
	Firstodds      string
	Notice         int
}

func (f *Filter) Exist() bool {
	total, err := engine.Where("rule=? and match_id=?", f.Rule, f.MatchId).Count(f)
	if err != nil {
		return false
	}
	return total > 0
}

func (f *Filter) ExistById() bool {
	total, err := engine.Where("match_id=?", f.MatchId).Count(f)
	if err != nil {
		return false
	}
	return total > 0
}

func (f *Filter) CopyFrom(m *Match) {
	f.Status = m.Status
	f.Time = m.Time
	f.Min = m.Min
	f.LeagueName = m.LeagueName
	f.LeagueSimpName = m.LeagueSimpName
	f.HoTeamName = m.HoTeamName
	f.GuTeamName = m.GuTeamName
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
	f.Odds = m.Odds.StringJson
	f.Firstodds = m.Firstodds.StringJson
}

type SnapShoot struct {
	MatchId        int `xorm:"index"`
	Status         int
	Time           string
	Min            int
	LeagueName     string
	LeagueSimpName string
	HoTeamName     string
	GuTeamName     string
	HoScore        int // 主队进球
	HoHalfScore    int // 主队半场得分
	GuScore        int // 客队进球
	GuHalfScore    int // 客队半场得分
	HoRed          int // 主队红牌
	HoYellow       int // 主队黄牌
	HoCo           int // 主队角球
	HoHfCo         int // 主队半场角球
	GuRed          int // 客队红牌
	GuYellow       int // 客队黄牌
	GuCo           int // 客队角球
	GuHfCo         int // 客队半场角球
	Let            float64
	LetHm          float64
	LetAw          float64
	AvgHm          float64
	AvgAw          float64
	AvgEq          float64
	Size           float64
	SizeBig        float64
	SizeSma        float64
	FirstLet       float64
	FirstLetHm     float64
	FirstLetAw     float64
	FirstAvgHm     float64
	FirstAvgAw     float64
	FirstAvgEq     float64
	FirstSize      float64
	FirstSizeBig   float64
	FirstSizeSma   float64
}

func Status(s int) string {
	switch s {
	case STATUS_NONE:
		return "未赛"
	case STATUS_FirstHalf:
		return "上半"
	case STATUS_MIDDLE:
		return "中场"
	case STATUS_SecondHALF:
		return "下半"
	case STATUS_COMPLETE:
		return "结束"
	case STATUS_DELAY:
		return "推迟"
	}
	return fmt.Sprintf("%d", s)
}

type YBF struct {
	Matches []*Match `json:"matches"`
}

func (y *YBF) Get(id int) *Match {
	for _, m := range y.Matches {
		if m.MatchId == id {
			return m
		}
	}
	return nil
}

func (y YBF) String() string {
	s := ""
	for _, m := range y.Matches {

		s += fmt.Sprintf("[%s] %s %d %svs%s %d-%d 滚球:%s 初盘:%s 角球:%d-%d 红牌:%d-%d 黄牌:%d-%d\n",
			Status(m.Status), m.LeagueName, m.Min, m.HoTeamName, m.GuTeamName, m.HoScore, m.GuScore,
			m.Odds.StringJson, m.Firstodds.StringJson,
			m.HoCo, m.GuCo,
			m.HoRed, m.GuRed,
			m.HoYellow, m.GuYellow)
	}
	return s
}

func ParseList(data []byte) (*YBF, error) {
	y := new(YBF)
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	for _, v := range y.Matches {
		v.Odds.Update(v.Odds.StringJson)
		v.Firstodds.Update(v.Odds.StringJson)
	}
	return y, nil
}

type OverTime struct {
	Ot    string `json:"ot"`
	State string `json:"state"`
}

// type 0 更新时间
// type 1 更新
type BasicUpdate struct {
	Id          int    `json:"id"`
	Type        int    `json:"type"`
	Time        string `json:"time"`
	Status      int    `json:"status"`
	Min         *int   `json:"min"`
	HoScore     *int   `json:"hoScore"`
	HoHalfScore *int   `json:"hoHalfScore"`
	GuScore     *int   `json:"guScore"`
	GuHalfScore *int   `json:"guHalfScore"`
	HoRed       *int   `json:"hoRed"`
	HoYellow    *int   `json:"hoYellow"`
	HoCo        *int   `json:"hoCo"`
	HoHfCo      *int   `json:"hoHfCo"`
	GuRed       *int   `json:"guRed"`
	GuYellow    *int   `json:"guYellow"`
	GuCo        *int   `json:"guCo"`
	GuHfCo      *int   `json:"guHfCo"`
}

type Update struct {
	Base         map[string]BasicUpdate `json:"basic"`
	Odds         map[string]string      `json:"odds"`
	OverTimeInfo map[string]OverTime    `json:"overTimeInfo"`
}

func (u Update) String() string {
	s := ""
	for k, v := range u.Base {
		s += fmt.Sprintf("%s: %d, %v\n", k, v.Type, v)
	}
	/*for k, v := range u.Odds {
		s += fmt.Sprintf("%s:%s\n", k, v)
	}*/
	return s
}

func ParseUpdate(data []byte) (*Update, error) {
	u := new(Update)
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
