# -*- coding: utf-8 -*-

import urllib
import urllib2

# 判断是否是自己的发言 getattr(member, 'uin') == bot.conf.qq

def onQQMessage(bot, contact, member, content):
    if content.startswith('#'):
        send_data = {'group':getattr(contact, 'name'),  'content':content[1:], 'from':member.name}
        send_data_urlencode = urllib.urlencode(send_data)
        result = urllib2.urlopen("http://localhost:8888/qqbot", data=send_data_urlencode,timeout=1000).read()    
        bot.SendTo(contact, '@'+member.name+'\n'+result)
