import praw
from json import loads, dumps
from random import choice, randint
from datetime import datetime
import time
import threading

def json_extract(thing='', filename='data.json'):
    json_file = open(filename)
    json_data = loads(json_file.read())
    json_file.close()
    if not thing is '':
        return json_data[thing]
    else:
        return json_data

# this needs to be implimented

reddit_client = json_extract('client_id')
reddit_secret = json_extract('client_secret')
reddit_agent = json_extract('user_agent')
url_things = ['.jpg', '.png', '.jpeg', '.gif', '.gifv', 'gfycat', 'youtube', 'youtu.be'] # will only get the link if these are in it

def extract_info(subreddit='all', limit=1, channel=None): # grabs the info from the sub
    urls = []
    titles = []
    permalinks = []
    nsfw_tags = []
    scores = []
    try:
        reddit = praw.Reddit(client_id=reddit_client, client_secret=reddit_secret, user_agent=reddit_agent)
    except:
        urls = [None]
        titles = [None]
        permalinks = [None]
    else:
        submissions = reddit.subreddit(subreddit).hot(limit=limit) # get the hot ones
        mods = reddit.subreddit(subreddit).moderator()
        for submission in submissions: # loop through the submissions
            if any(str(submission.author)==str(n) for n in mods): # ignores a mod's post
                continue
            elif len(submission.selftext)<=2000:
                titles += [submission.title]
                permalinks += [submission.permalink]
                scores += [submission.score]
                if submission.over_18: # check for nsfw
                    nsfw_tags += [True]
                else:
                    nsfw_tags += [False]
                if not any(n in submission.url for n in url_things) and submission.selftext!='':
                    urls += [submission.selftext]
                else:
                    urls += [submission.url]

    return [urls, titles, permalinks, nsfw_tags, scores]

def get_post_thing(subs=["funny"], nsfw=False, limit=100): #grabs a random post from the extract_info def
    subreddit = choice(subs)
    posts = extract_info(subreddit, limit)
    num_of_posts = len(posts[0]) - 1
    post_type = None
    count = 0
    try:
        while True:
            count += 1
            post_number = randint(0, num_of_posts)
            post_link = posts[0][post_number]
            post_title = posts[1][post_number]
            post_permalink = posts[2][post_number]
            post_nsfw = posts[3][post_number]
            post_score = posts[4][post_number]
            post_type = True
            if not nsfw and post_nsfw:
                pass
            else:
                break
            if count>=len(posts):
                post_link = "Too many tries to not find NSFW post, maybe that Subreddit is filled with them?"
                post_title = "Error!!!!"
                post_permalink = "/r/{}".format(choice(subs))
                post_score = "{} Tries".format(count)
                break
    except Exception as e:
        if "403" in str(e):
            post_type = "403"
            post_permalink = False
            post_link = "The server denied our request!"
            post_title = 'Error HTTP 403'
        else:
            post_title = 'Please try again.'
            post_link = str(e)
            post_type = None
        logging.error("Error!")
        logging.error(str(e))
        logging.error("------")

    return [post_link, post_type, post_title, post_permalink, post_score]

def log_channels(channels):
    currentlist = open("channels.log").readlines().close()
    for channel in channels:
        if not channel in currentlist:
            chlist = open("channel.log", "a")
            chlist.write("{}\n".format(channel))
            chlist.close()

def get_post_by_id(subid=None, nsfw=False):
    reddit = praw.Reddit(client_id=reddit_client, client_secret=reddit_secret, user_agent=reddit_agent)
    url = ""
    try:
        submission = reddit.submission(id=subid)
    except Exception as e:
        return [str(e), None, "Error!", None, "nil"]
    else:
        if submission.over_18 and not nsfw:
            return ["Error: Post is NSFW being posted in an SFW chat.", None, "Error!", "/error404", "Tries: 1"]
        else:
            if submission.is_self:
                url = submission.selftext
            else:
                url = submission.url
            return [url, True, submission.title, submission.permalink, submission.score]

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
        dump(rawdata, postfile)





