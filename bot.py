import discord
from lib import *
from random import choice, randint

url_things = ['.jpg', '.png', '.jpeg']
def get_post_thing(subs=["funny"]):
	#subreddit = choice(['funny', 'dankmemes','dank_memes', 'jokes', 'darkjokes'])
	subreddit = choice(subs)
	posts = extract_info(subreddit, 30)
	post_number = randint(0, 27)
	pic_e = ['.jpg', '.png', '.jpeg']
	vid_e = ['.gif', '.gifv']
	post_type = None
	try:
		post_link = posts[0][post_number]
		post_title = posts[1][post_number]
		post_permalink = posts[2][post_number]
		post_type = True
	except Exception as e:
		post_link = 'https://mediaconnectpartners.staticscdn.com/wp-content/uploads/oops-header.png'
		post_title = 'Please try again.'
		post_type = None
		print("Error!")
		print(e)
		print("Retrying....\n")

	return [post_link, post_type, post_title, post_permalink]


token = json_extract('token')

client = discord.Client()

@client.event
async def on_message(message):
	if message.author == client.user:
		return
	if message.content.startswith("!meme"):
		raw_msg = ""
		while True:
			raw_msg = get_post_thing(["dankmemes","funny","memes","dank_memes","CringeAnarchy"])
			if not raw_msg[1]==None:
				break
		print("Posting:")
		print(raw_msg[2])
		print(raw_msg[0])
		print("Original post: https://reddit.com{}".format(raw_msg[3]))
		print("\n")
		embed = discord.Embed(title=raw_msg[2], url=raw_msg[0])
		embed.set_image(url=raw_msg[0])
		await client.send_message(message.channel, embed=embed, tts=False)
		await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)

	if message.content.startswith("!joke"):
		raw_msg = ""
		while True:
			raw_msg = get_post_thing(["jokes", "meanjokes", "darkjokes"])
			if not raw_msg[1]==None:
				break
		print("Posting:")
		print(raw_msg[2])
		print(raw_msg[0])
		print("Original post: https://reddit.com{}".format(raw_msg[3]))
		print("\n")
		await client.send_message(message.channel, content=raw_msg[2], tts=True)
		await client.send_message(message.channel, content=raw_msg[0], tts=True)
		await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)
@client.event
async def on_ready():
	print('Logged in as')
	print(client.user.name)
	print(client.user.id)
	print('------')

while True:
	try:
		client.run(token)
	except Exception as e:
		print(e)
		continue
