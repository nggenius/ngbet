# -*- coding: utf-8 -*-

from qqbot.utf8logger import DEBUG
import urllib
import urllib2

# 判断是否是自己的发言 getattr(member, 'uin') == bot.conf.qq

def onQQMessage(bot, contact, member, content):
    if '@ME' in content:
        send_data = {'group':getattr(contact, 'name'),  'content':content, 'from':member.name}
        send_data_urlencode = urllib.urlencode(send_data)

        requrl = "http://localhost:8888/qqbot"

        req = urllib2.Request(url = requrl,data =send_data_urlencode)
        urllib2.urlopen(req, timeout=10).read()
