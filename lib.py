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

reddit_client = json_extract('client_id')
reddit_secret = json_extract('client_secret')
reddit_agent = json_extract('user_agent')

def extract_info(subreddit='all', limit=1):
    reddit = praw.Reddit(client_id=reddit_client, client_secret=reddit_secret, user_agent=reddit_agent)
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

def get_post_thing(subs=["funny"]):
	#subreddit = choice(['funny', 'dankmemes','dank_memes', 'jokes', 'darkjokes'])
	subreddit = choice(subs)
	posts = extract_info(subreddit, 30)
	num_of_posts = len(posts[0]) - 1
	post_number = randint(0, num_of_posts)
	post_type = None
	try:
		post_link = posts[0][post_number]
		post_title = posts[1][post_number]
		post_permalink = posts[2][post_number]
		post_type = True
	except Exception as e:
		#post_link = 'https://mediaconnectpartners.staticscdn.com/wp-content/uploads/oops-header.png'
		post_title = 'Please try again.'
		post_type = None
		print("Error!")
		print(e)
		print("Retrying....")
		print("------")

	return [post_link, post_type, post_title, post_permalink]
