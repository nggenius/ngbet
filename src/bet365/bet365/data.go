package bet365

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

// MD=1 上半场结束
// TS=1 下半场
// TT=0 时间暂停
// VS=-1 隐藏比赛
// FS=0
func splitMsg(data []byte) (string, map[string]string) {
	s := bytes.Split(data, []byte{';'})
	if len(s) == 0 {
		return "", nil
	}
	kw := string(s[0])

	kv := make(map[string]string)
	for _, s1 := range s[1:] {
		s2 := bytes.Split(s1, []byte{'='})
		if len(s2) == 1 {
			kv[string(s2[0])] = ""
		}
		if len(s2) == 2 {
			kv[string(s2[0])] = string(s2[1])
		}
	}
	return kw, kv
}

func split(data []byte) map[string]string {
	s := bytes.Split(data, []byte{';'})
	kv := make(map[string]string)
	for _, s1 := range s {
		s2 := bytes.Split(s1, []byte{'='})
		if len(s2) == 1 {
			kv[string(s2[0])] = ""
		}
		if len(s2) == 2 {
			kv[string(s2[0])] = string(s2[1])
		}
	}
	return kv
}

type Node struct {
	tag    string
	child  map[string]*Node
	attrs  map[string]string
	parent *Node
}

func (n *Node) Update(kv map[string]string) {
	for k, v := range kv {
		n.attrs[k] = v
	}
}

func (n *Node) It() string {
	if it, ok := n.attrs["IT"]; ok {
		return it
	}

	return ""
}

func (n *Node) State() int {
	if n.tag != "EV" {
		panic("node type error")
	}

	if n.Int("FS") == 0 {
		return STATUS_NONE
	}

	md := n.Int("MD")
	tt := n.Int("TT")
	tm := n.Int("TM")
	if md == 0 && tt == 1 {
		return STATUS_FIRSTHALF
	}
	if md == 1 && tt == 0 && tm == 45 {
		return STATUS_MIDDLE
	}
	if md == 1 && tt == 1 {
		return STATUS_SECONDHALF
	}
	if md == 1 && tt == 0 && tm == 90 {
		return STATUS_COMPLETE
	}

	return STATUS_UNKNOWN
}

func (n *Node) Attr(key string) string {
	if v, ok := n.attrs[key]; ok {
		return v
	}

	return ""
}

func (n *Node) Float(key string) float64 {
	v := n.Attr(key)
	if v == "" {
		return 0
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}

	return f
}

func (n *Node) Int(key string) int {
	v := n.Attr(key)
	if v == "" {
		return 0
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return i
}

func (n *Node) SS() (a int, b int) {
	fmt.Sscanf(n.Attr("SS"), "%d-%d", &a, &b)
	return
}

func (n *Node) Odd() float64 {
	od := n.Attr("OD")
	if od == "" {
		return 0
	}

	var a, b float64
	_, err := fmt.Sscanf(od, "%f/%f", &a, &b)
	if err != nil {
		return 0
	}

	return a/b + 1
}

func (n *Node) Remove(it string) {
	delete(n.child, it)
}

func (n *Node) AddChild(node *Node) {
	n.child[node.It()] = node
}

func (n *Node) Path() string {
	c := n
	p := n.It()
	for c != nil {
		p = c.It() + "/" + p
		c = c.parent
	}

	return p
}

func NewNode(t string) *Node {
	n := new(Node)
	n.tag = t
	n.attrs = make(map[string]string)
	n.child = make(map[string]*Node)
	return n
}

func NewSimpleNode(t string) *Node {
	n := new(Node)
	n.tag = t
	n.attrs = make(map[string]string)
	n.child = make(map[string]*Node)
	n.attrs["IT"] = t
	return n
}

type Bet365Data struct {
	sync.Mutex
	RootName  string
	Root      *Node
	ItHash    map[string]*Node
	TI        string
	time      time.Time
	localtime time.Time
	del       []string
}

func NewBet365Data(name string) *Bet365Data {
	d := new(Bet365Data)
	d.RootName = name
	d.ItHash = make(map[string]*Node)
	d.Root = NewSimpleNode(name)
	d.del = make([]string, 0, 32)
	return d
}

func formatTime(s string) time.Time {
	if len(s) == 17 {
		t1 := s[:14] + "." + s[14:]
		t, err := time.Parse("20060102150405.000", t1)
		if err != nil {
			log.Fatalln("parse time error ", err)
		}

		return t
	}

	if len(s) == 14 {
		t, err := time.Parse("20060102150405", s)
		if err != nil {
			log.Fatalln("parse time error ", err)
		}
		return t
	}

	log.Fatalln("parse time error")
	return time.Now()
}

func (d *Bet365Data) MatchTime(node *Node) (m, s int) {
	tu := node.Attr("TU")
	if tu == "" || tu == "19000101000000" {
		return 0, 0
	}

	if node.Attr("TT") == "0" { // 时间暂停
		m, _ = strconv.Atoi(node.Attr("TM"))
		s, _ = strconv.Atoi(node.Attr("TS"))
		return
	}

	t := formatTime(tu)
	d1 := t.Sub(d.time)
	d2 := time.Now().Sub(d.localtime)

	om, _ := strconv.Atoi(node.Attr("TM"))
	os, _ := strconv.Atoi(node.Attr("TS"))
	dur := int((d2 - d1).Seconds()) + om*60 + os // 偏移量秒
	m = dur / 60
	s = dur % 60
	return
}

func (d *Bet365Data) parseTime(data []byte) {
	infos := bytes.Split(data, []byte{124})
	switch infos[0][0] {
	case 'F':
		kw, kv := splitMsg(infos[1])
		if kw == "IN" {
			TI := kv["TI"]
			d.TI = TI
			d.time = formatTime(TI)
			d.localtime = time.Now()
			log.Println("time:", d.time)
		}
	}
}

func (d *Bet365Data) FindNode(it string) *Node {
	if n, ok := d.ItHash[it]; ok {
		return n
	}

	return nil
}

func (d *Bet365Data) ChildByType(node *Node, t string) []*Node {
	var ret []*Node
	for _, c := range node.child {
		if c != nil && c.tag == t {
			ret = append(ret, c)
		}

		if len(c.child) > 0 {
			r1 := d.ChildByType(c, t)
			if len(r1) > 0 {
				ret = append(ret, r1...)
			}
		}
	}
	return ret
}

func (d *Bet365Data) AddNode(parent *Node, node *Node) *Node {
	node.parent = parent
	it := node.It()
	if it == "" {
		panic("IT is empty")
	}

	d.ItHash[it] = node

	parent.AddChild(node)
	return node
}

func (d *Bet365Data) RemoveAllChild(node *Node) {
	for k, v := range node.child { // 循环遍历所有子结点
		if v != nil {
			d.Remove(k)
		}
	}
}

func (d *Bet365Data) Remove(it string) {
	node := d.FindNode(it)
	if node == nil {
		return
	}

	d.RemoveAllChild(node)

	p := node.parent
	if p != nil {
		p.Remove(it) //移除子结点
	}

	node.parent = nil
	delete(d.ItHash, it)
	d.Lock()
	d.del = append(d.del, it)
	d.Unlock()
}

func (d *Bet365Data) GetDel() []string {
	var ds []string
	d.Lock()
	if len(d.del) > 0 {
		for _, v := range d.del {
			ds = append(ds, v)
		}
		d.del = d.del[:0]
	}
	d.Unlock()
	return ds
}

func parseInPlayNode(d *Bet365Data, parent *Node, data [][]byte) {
	stack := NewStack()
	stack.Push(parent)

	for _, line := range data {
		kw, kv := splitMsg(line)
		switch kw {
		case "CL": // 分类标签
			top := stack.Top().(*Node)
			for top.tag != parent.tag {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("CL")
			node.attrs = kv
			d.AddNode(top, node)
			stack.Push(node)
		case "CT": // 分组标签，联赛名
			top := stack.Top().(*Node)
			for top.tag != "CL" {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("CT")
			node.attrs = kv
			d.AddNode(top, node)
			stack.Push(node)
		case "EV": // 比赛信息
			top := stack.Top().(*Node)
			for top.tag != "CT" && top.tag != "CL" { // 可以没有CT
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("EV")
			node.attrs = kv
			d.AddNode(top, node)
			stack.Push(node)
		case "MA": // 盘口大类
			top := stack.Top().(*Node)
			for top.tag != "EV" {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("MA")
			node.attrs = kv
			d.AddNode(top, node)
			stack.Push(node)
		case "PA": // 盘口赔率
			top := stack.Top().(*Node)
			for top.tag != "MA" {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("PA")
			node.attrs = kv
			d.AddNode(top, node)
		}
	}

}

func ParseInPlay(d *Bet365Data, path [][]byte, data []byte) {
	node := d.FindNode("OVInPlay_10_0")
	if node == nil {
		node = d.AddNode(d.Root, NewNode("OVInPlay_10_0"))
	}
	infos := bytes.Split(data, []byte{124})
	switch infos[0][0] {
	case 'F':
		d.RemoveAllChild(node)
		parseInPlayNode(d, node, infos[1:])
		log.Println("recv inplay full data")
	case 'U':
		it := string(path[len(path)-1])
		log.Fatalln("recv update inplay", it, string(data))
	case 'I':
		if len(path) < 2 {
			log.Fatalln("path is too short")
		}
		it := string(path[len(path)-2])
		pnode := d.FindNode(it)
		if pnode == nil {
			log.Fatalln("parent not found", it)
			return
		}
		parseInPlayNode(d, pnode, infos[1:])
	case 'D':
		it := string(path[len(path)-1])
		node := d.FindNode(it)
		if node != nil {
			d.Remove(it)
		}
	}
}

func parseOVMNode(d *Bet365Data, parent *Node, data [][]byte) {
	stack := NewStack()
	stack.Push(parent)

	for _, line := range data {
		kw, kv := splitMsg(line)
		switch kw {
		case "EV":
			continue
		case "MA": // 盘口分类
			top := stack.Top().(*Node)
			for top.tag != parent.tag {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("MA")
			node.attrs = kv
			d.AddNode(top, node)
			stack.Push(node)
		case "PA": // 盘口赔率
			top := stack.Top().(*Node)
			for top.tag != "MA" {
				stack.Pop()
				top = stack.Top().(*Node)
			}
			node := NewNode("PA")
			node.attrs = kv
			d.AddNode(top, node)
		}
	}
}

func ParseOVM(it string, d *Bet365Data, path [][]byte, data []byte) {
	node := d.FindNode(it)
	if node == nil {
		node = d.AddNode(d.Root, NewNode(it))
	}

	infos := bytes.Split(data, []byte{124})
	switch infos[0][0] {
	case 'F':
		d.RemoveAllChild(node)
		parseOVMNode(d, node, infos[1:])
		log.Println("recv ", it)
	case 'U':
		it := string(path[len(path)-1])
		node := d.FindNode(it)
		if node != nil {
			kv := split(data)
			node.Update(kv)
		}
	case 'I':
		if len(path) < 2 {
			log.Fatalln("path is too short")
		}
		it := string(path[len(path)-2])
		pnode := d.FindNode(it)
		if pnode == nil {
			log.Fatalln("parent not found", it)
			return
		}
		parseOVMNode(d, pnode, infos[1:])
	case 'D':
		it := string(path[len(path)-1])
		node := d.FindNode(it)
		if node != nil {
			d.Remove(it)
		}
	}

}

func updateWithIt(d *Bet365Data, path [][]byte, data []byte) {
	infos := bytes.Split(data, []byte{124})
	op := infos[0][0]
	switch op {
	case 'F':
	case 'U':
		it := string(path[len(path)-1])
		node := d.FindNode(it)
		if node != nil {
			kv := split(infos[1])
			node.Update(kv)
		}

	case 'I':
		if len(path) < 2 {
			log.Fatalln("path is too short")
		}
		it := string(path[len(path)-2])
		pnode := d.FindNode(it)
		if pnode == nil {
			return
		}

		kw, kv := splitMsg(infos[1])
		node := NewNode(kw)
		node.attrs = kv
		d.AddNode(pnode, node)
	case 'D':
		it := string(path[len(path)-1])
		node := d.FindNode(it)
		if node != nil {
			d.Remove(it)
			log.Println("delete ", it)
		}
	default:
		log.Fatalf("unsupport op %s", string(op))
	}
}

func ParseData(d *Bet365Data, path []byte, data []byte) error {
	p := bytes.Split(path, []byte{'/'})
	if len(p) == 0 || len(p[0]) == 0 {
		return nil
	}

	it := string(p[0])

	switch it {
	case "OVInPlay_10_0":
		ParseInPlay(d, p, data)
	case "OVM1", "OVM2", "OVM3":
		ParseOVM(it, d, p, data)
	case "__time":
		d.parseTime(data)
	default:
		updateWithIt(d, p, data)
	}
	return nil
}
