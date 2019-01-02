package ssq

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LuckBall struct {
	Red  [6]int `json:"red"`
	Blue int    `json:"blue"`
}

type Result struct {
	Status int
	Last   History
	Lucky  []LuckBall
}

type OpenCode struct {
	Expect   string `json:"expect"`
	OpenCode string `json:"opencode"`
	history  History
}

func (o *OpenCode) Parse() {
	e, err := strconv.Atoi(o.Expect)
	if err != nil {
		panic(err)
	}

	o.history.Expect = e
	ball := strings.Split(o.OpenCode, "+")
	if len(ball) != 2 {
		panic("code format error")
	}

	reds := strings.Split(ball[0], ",")
	if len(reds) != 6 {
		panic("red format error")
	}

	for k := 0; k < 6; k++ {
		v, err := strconv.Atoi(reds[k])
		if err != nil {
			panic("parse number error")
		}
		o.history.Red[k] = v
	}

	b, err := strconv.Atoi(ball[1])
	if err != nil {
		panic(err)
	}

	o.history.Blue = b

}

type OnlineHistory struct {
	Data []OpenCode
}

func (h *OnlineHistory) ParseAll() {
	for k := range h.Data {
		h.Data[k].Parse()
	}
}

func ReadFile(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return ioutil.ReadAll(file)
}

type History struct {
	Expect int    `json:"expect"`
	Red    [6]int `json:"red"`
	Blue   int    `json:"blue"`
}

type BallHot struct {
	Ball int
	Hot  int
}

type BallRatio []BallHot

func (p BallRatio) Len() int           { return len(p) }
func (p BallRatio) Less(i, j int) bool { return p[i].Hot < p[j].Hot }
func (p BallRatio) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func analyzeHot(redRatio, blueRatio BallRatio, loadHistory []History) {
	for k := range redRatio {
		redRatio[k].Ball = k + 1
		redRatio[k].Hot = 0
	}

	for k := range blueRatio {
		blueRatio[k].Ball = k + 1
		blueRatio[k].Hot = 0
	}

	for k := range loadHistory {
		for _, v := range loadHistory[k].Red {
			ball := v - 1
			if redRatio[ball].Hot != 0 {
				continue
			}
			redRatio[ball].Hot = 100 - k
		}
		blueball := loadHistory[k].Blue - 1
		if blueRatio[blueball].Hot == 0 {
			blueRatio[blueball].Hot = 100 - k
		}
	}
	sort.Sort(redRatio)
	sort.Sort(blueRatio)
}

func getLocalExcept() int {
	f, err := os.Open("history.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	l, _, err := r.ReadLine()
	for err == nil {
		if strings.TrimSpace(string(l)) == "" {
			l, _, err = r.ReadLine()
			continue
		}
		s := strings.Split(string(l), "\t")
		if len(s) != 8 {
			panic("line error")
		}

		e, err := strconv.Atoi(s[0])
		if err != nil {
			panic("parse expect error")
		}
		return e
	}

	return 0
}

func Histroy(update bool) []History {
	var loadHistory []History
	f, err := os.OpenFile("history.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	l, _, err := r.ReadLine()
	for err == nil {
		if strings.TrimSpace(string(l)) == "" {
			l, _, err = r.ReadLine()
			continue
		}
		var h History
		s := strings.Split(string(l), "\t")
		if len(s) != 8 {
			panic("line error")
		}

		e, err := strconv.Atoi(s[0])
		if err != nil {
			panic("parse expect error")
		}

		h.Expect = e
		for k := 0; k < 6; k++ {
			v, err := strconv.Atoi(s[k+1])
			if err != nil {
				panic("parse number error")
			}
			h.Red[k] = v
		}

		v, err := strconv.Atoi(s[7])
		if err != nil {
			panic("parse number error")
		}

		h.Blue = v
		loadHistory = append(loadHistory, h)

		l, _, err = r.ReadLine()

	}

	for update {
		resp, err := http.Get("http://f.apiplus.net/ssq-20.json")
		if err != nil {
			break
		}

		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}

		var hd OnlineHistory
		err = json.Unmarshal(data, &hd)
		if err != nil {
			break
		}

		hd.ParseAll()
		if len(loadHistory) == 0 || hd.Data[0].history.Expect > loadHistory[0].Expect {
			for k := len(hd.Data) - 1; k >= 0; k-- {
				if len(loadHistory) == 0 || hd.Data[k].history.Expect > loadHistory[0].Expect {
					loadHistory = append(loadHistory, hd.Data[k].history)
					copy(loadHistory[1:], loadHistory[0:])
					loadHistory[0] = hd.Data[k].history
				}
			}

			f.Truncate(0)
			w := bufio.NewWriter(f)
			f.Seek(0, io.SeekStart)
			for _, l := range loadHistory {
				w.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", l.Expect, l.Red[0], l.Red[1], l.Red[2], l.Red[3], l.Red[4], l.Red[5], l.Blue))
			}
			w.Flush()
		}
		break
	}
	if len(loadHistory) > 50 {
		loadHistory = loadHistory[:50]
	}
	return loadHistory
}

func Millionaire() (res Result) {
	redRatio := make(BallRatio, 33)
	blueRatio := make(BallRatio, 16)
	loadHistory := Histroy(false)
	analyzeHot(redRatio, blueRatio, loadHistory)
	r := make([]int, 6)
	copy(r[:], loadHistory[0].Red[:])
	b := loadHistory[0].Blue

	res.Status = 200
	res.Last = loadHistory[0]

	rand.Seed(time.Now().Unix())
	redball := make([]int, 33)
	for k := range redball {
		redball[k] = k + 1
	}
	blueball := make([]int, 16)
	for k := range blueball {
		blueball[k] = k + 1
	}

	//fmt.Println(redball, blueball)
	//fmt.Println("__________________________________________")

	for k := range r {
		for j := range redball {
			if r[k] == redball[j] {
				copy(redball[j:], redball[j+1:])
				break
			}
		}
	}
	redball = redball[:len(redball)-6]

	for k := range blueball {
		if blueball[k] == b {
			copy(blueball[k:], blueball[k+1:])
			break
		}
	}

	blueball = blueball[:len(blueball)-1]

	for t := 0; t < 4; t++ {
		for k := range r {
			idx := rand.Intn(len(redball))
			r[k] = redball[idx]
			copy(redball[idx:], redball[idx+1:])
			redball = redball[:len(redball)-1]
		}

		idx := rand.Intn(len(blueball))
		b = blueball[idx]
		copy(blueball[idx:], blueball[idx+1:])
		blueball = blueball[:len(blueball)-1]
		sort.Ints(r[:])
		var lucky LuckBall
		copy(lucky.Red[:], r[:])
		lucky.Blue = b
		res.Lucky = append(res.Lucky, lucky)
	}

	r1 := rand.Perm(11)
	r2 := rand.Perm(11)
	r3 := rand.Perm(11)

	b1 := rand.Perm(10)

	var specialRed []int

	for k := 0; k < 3; k++ {
		specialRed = append(specialRed, redRatio[r1[k]].Ball)
	}
	for k := 0; k < 2; k++ {
		specialRed = append(specialRed, redRatio[11+r2[k]].Ball)
	}
	specialRed = append(specialRed, redRatio[22+r3[0]].Ball)
	sort.Ints(specialRed)

	var lucky LuckBall
	copy(lucky.Red[:], specialRed[:])
	lucky.Blue = blueRatio[b1[0]].Ball
	res.Lucky = append(res.Lucky, lucky)
	return
}
