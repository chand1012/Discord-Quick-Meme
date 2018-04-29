#!/usr/bin/env python3
import discord
from lib import *
from random import choice, randint

url_things = ['.jpg', '.png', '.jpeg']
def get_post_thing(subs=["funny"]):
	#subreddit = choice(['funny', 'dankmemes','dank_memes', 'jokes', 'darkjokes'])
	subreddit = choice(subs)
	posts = extract_info(subreddit, 30)
	num_of_posts = len(posts[0]) - 1
	post_number = randint(0, num_of_posts)
	pic_e = ['.jpg', '.png', '.jpeg']
	vid_e = ['.gif', '.gifv']
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
		print("Retrying....\n")

	return [post_link, post_type, post_title, post_permalink]


token = json_extract('token')

client = discord.Client()

@client.event
async def on_message(message):
	if message.author == client.user:
		return
	if message.content.startswith("!meme"):
		recv = message.content
		if recv[6:] is '':
			raw_msg = ""
			while True:
				raw_msg = get_post_thing(["dankmemes","funny","memes","dank_memes"])
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
		else:
			raw_msg = ""
			count=0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]])
				if not raw_msg[1]==None:
					break
				if count>=10:
					print("Count exceeded!")
					break
			print("Posting:")
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("\n")
			embed = discord.Embed(title=raw_msg[2], url="https://reddit.com{}".format(raw_msg[3]))
			embed.set_image(url=raw_msg[0])
			if not count>=10:
				await client.send_message(message.channel, embed=embed, tts=False)
				await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)
			else:
				await client.send_message(message.channel, "Something went wrong, please try again!")
				await client.send_message(message.channel, "Problem subreddit: https://reddit.com/{}".format(recv[6:]))

	if message.content.startswith("!joke"):
		recv = message.content
		if recv[6:] is '':
			raw_msg = ""
			while True:
				raw_msg = get_post_thing(["jokes", "darkjokes"])
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
		else:
			raw_msg = ""
			premsg = "Original Post:"
			count = 0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]])
				if not raw_msg[1]==None:
					break
				if count>=10:
					print("Count exceeded!")
					raw_msg[2] = "Something went wrong."
					raw_msg[0] = "Please try again"
					raw_msg[3] = "/{}".format(recv[6:])
					premsg = "Subreddit that caused the issue:"
					break

			print("Posting:")
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("\n")
			await client.send_message(message.channel, content=raw_msg[2], tts=True)
			await client.send_message(message.channel, content=raw_msg[0], tts=True)
			await client.send_message(message.channel, content="{} https://reddit.com{}".format(premsg, raw_msg[3]), tts=False)

	if message.content.startswith('!news'):
		recv = message.content
		raw_msg = ""
		count = 0
		while True:
			count+=1
			if not recv[6:]=='':
				raw_msg = get_post_thing([recv[6:]], 1, True)
			else:
				raw_msg = get_post_thing(["UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"])
			if not raw_msg[1]==None:
				break
			if count>=10:
				print("Count exceeded!")
				raw_msg[2] = "Something went wrong."
				raw_msg[0] = "Please try again"
				raw_msg[3] = "/{}".format(recv[6:])
				premsg = "Subreddit that caused the issue:"
				break
		print("Posting:")
		print(raw_msg[2])
		print(raw_msg[0])
		await client.send_message(message.channel, content=raw_msg[2], tts=True)
		await client.send_message(message.channel, content="Link: {}".format(raw_msg[0]), tts=False)
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
		if "Event loop" in str(e):
			print("\nStopping bot....")
			break
		else:
			print(e)
			continue
