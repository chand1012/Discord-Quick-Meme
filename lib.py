import praw
from json import loads
from random import choice, randint

def json_extract(thing='', filename='data.json'):
    json_file = open(filename)
    json_data = loads(json_file.read())
    json_file.close()
    if not thing is '':
        return json_data[thing]
    else:
        return json_data

reddit_client = json_extract('client_id')
reddit_secret = json_extract('client_secret')
reddit_agent = json_extract('user_agent')

def extract_info(subreddit='all', limit=1, nsfw=False):
    reddit = praw.Reddit(client_id=reddit_client, client_secret=reddit_secret, user_agent=reddit_agent)
    url_things = ['.jpg', '.png', '.jpeg']
    urls = []
    titles = []
    permalinks = []
    nsfw_tags = []
    submissions = reddit.subreddit(subreddit).hot(limit=limit)
    for submission in submissions:
        if not len(submission.selftext)>=2000:
            titles += [submission.title]
            permalinks += [submission.permalink]
            if submission.over_18:
                nsfw_tags += [True]
            else:
                nsfw_tags += [False]
            if not any(n in submission.url for n in url_things) and submission.selftext!='':
                urls += [submission.selftext]
            else:
                urls += [submission.url]

    return [urls, titles, permalinks, nsfw_tags]

def get_post_thing(subs=["funny"], nsfw=False):
	#subreddit = choice(['funny', 'dankmemes','dank_memes', 'jokes', 'darkjokes'])
    subreddit = choice(subs)
    posts = extract_info(subreddit, 30, nsfw)
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
            post_type = True
            if not nsfw:
                if not post_nsfw:
                    break
            if count>=10:
                post_link = "Too many tries to not find NSFW post, maybe that Subreddit is filled with them?"
                post_title = "Error!!!"
                post_permalink = None
                break
    except Exception as e:
		#post_link = 'https://mediaconnectpartners.staticscdn.com/wp-content/uploads/oops-header.png'
        post_title = 'Please try again.'
        post_type = None
        print("Error!")
        print(e)
        print("Retrying....")
        print("------")

    return [post_link, post_type, post_title, post_permalink]

def log_channels(channels):
    currentlist = open("channels.log").readlines().close()
    for channel in channels:
        if not channel in currentlist:
            chlist = open("channel.log", "a")
            chlist.write("{}\n".format(channel))
            chlist.close()
