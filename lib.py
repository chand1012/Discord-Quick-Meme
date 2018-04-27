import praw
from json import loads

def json_extract(thing='', filename='data.json'):
    json_file = open(filename)
    json_data = loads(json_file.read())
    json_file.close()
    if not thing is '':
        return json_data[thing]
    else:
        return json_data

def extract_info(subreddit='all', limit=1):
    reddit = praw.Reddit(client_id=json_extract('client_id'), client_secret=json_extract('client_secret'), user_agent=json_extract('user_agent'))
    url_things = ['.jpg', '.png', '.jpeg']
    urls = []
    titles = []
    permalinks = []
    submissions = reddit.subreddit(subreddit).hot(limit=limit)
    for submission in submissions:
        if not len(submission.selftext)>=2000:
            titles += [submission.title]
            permalinks += [submission.permalink]
            if not any(n in submission.url for n in url_things) and submission.selftext!='':
                urls += [submission.selftext]
            else:
                urls += [submission.url]

    return [urls, titles, permalinks]
