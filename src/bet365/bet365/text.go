package bet365

const (
	TEXT_INVALID = `[尴尬] [%s] 
%s
%s
[%02d:%02d] 进球无效，比分:%d-%d`

	TEXT_RED = `[红包] [红] [%s] 
%s
%s
[%02d:%02d] ⚽ 比分:%d-%d`

	TEXT_BLACK_HALF = `[炸弹] [黑] [%s] 
%s
%s
上半场结束，比分:%d-%d`

	TEXT_BLACK = `[炸弹] [黑] [%s] 
%s
%s  
比赛结束，比分:%d-%d`

	TEXT_RULE_MSG = `⚽[%s] 
%s
%s
当前比分:%d-%d
平局概率:%d%%
推荐:%s大%.1f
当前水位:%.1f
id:%s`

	TEXT_ABOVE = "\n[忍者]注意：当前盘口(%.2f)高于推荐盘口,可等水"

	TEXT_NOTICE_ODD_MSG = `[忍者][%s]
%s
%s
降盘啦，快上车
id:%s`

	TEXT_PREVIEW = `%s 
%s  
	胜平负:%.2f,%.2f,%.2f 
	让分:%.2f,%.2f,%.2f 
	大小盘:%.2f,%.2f,%.2f
  上半场:
	胜平负:%.2f,%.2f,%.2f 
	让分:%.2f,%.2f,%.2f
	大小盘:%.2f,%.2f,%.2f
id:%s`

	TEXT_FULL = `%s
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
  大小盘:%.2f,%.2f,%.2f`

	ATTENTION = `[忍者] 关注的球队比赛开始了
%s
%s
%s`

	STATE_INFO = `[%s] %s %s %d-%d 
平局概率:%d%%
id:%s
	`

	GLOAL = `[进球] %s %s %02d:%02d %d-%d 
平局概率:%d%%
id:%s`

	INVALID = `[无效] %s %s %02d:%02d %d-%d 
id:%s`

	ATTENTION_EQ = `[忍者]注意 
%s 
%s 
经评估，上半场破蛋概率较大，请关注。
id:%s`
)
