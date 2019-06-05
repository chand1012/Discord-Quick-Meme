from json import loads, dumps
from datetime import datetime

def log_channels(channels):
    currentlist = open("channels.log").readlines().close()
    for channel in channels:
        if not channel in currentlist:
            chlist = open("channel.log", "a")
            chlist.write("{}\n".format(channel))
            chlist.close()

def check_blacklist(channel, postlink):
    with open("posts.json") as postfile:
        rawdata = loads(postfile.read())
        channeldata = rawdata[channel]
        if postlink in channeldata:
            now = time.time()
            post = channeldata[postlink]
            if post<now:
                return False
            else:
                return True
        else:
            return False

def add_blacklist(channel, postlink):
    rawdata = None
    with open('posts.json') as postfile:
        rawdata = loads(postfile.read())
    if not channel in rawdata:
        rawdata[channel] = {}
    channeldata = rawdata[channel]
    post = time.time() + 86400 # seconds in a day
    channeldata[postlink] = post
    rawdata[channel] = channeldata
    with open('posts.json', "w") as postfile:
        dumps(rawdata, postfile)
