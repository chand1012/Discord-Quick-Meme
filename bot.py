#!/usr/bin/env python3
import discord
import logging
import time
from lib import get_post_thing, json_extract, get_post_by_id
token = json_extract('token')
client = discord.Client()
filetypes = ['gif', 'gifv', 'gfycat', 'v.redd.it', 'youtube', 'youtu.be', '.jpg', '.png', '.jpeg']
channels = {}
#logging things
logger = logging.getLogger('discord')
logger.setLevel(logging.INFO)
handler = logging.FileHandler(filename='discord.log', encoding='utf-8', mode='w')
handler.setFormatter(logging.Formatter('%(asctime)s:%(levelname)s:%(name)s: %(message)s'))
logger.addHandler(handler)

@client.event
async def on_message(message):
	# The recv stuff
	recv = message.content
	now = time.time()

	# spam filter code
	if not message.channel in channels: # if a channel is not in the dictionary, check for it.
		channels[message.channel] = [time.time(), 0] # later this should be made a json

	if channels[message.channel][1]>=5 and channels[message.channel][0]+60>now:
		await client.send_message(message.channel, content="You guys have been sending a lot of messages, why don't you slow down a bit?") 
		return
	elif channels[message.channel][1]<5 and channels[message.channel][0]+60>now:
		channels[message.channel][1] += 1
	elif channels[message.channel][1]<5 and channels[message.channel][0]+60<now:
		channels[message.channel][1] = 0
		channels[message.channel][0] = time.time()
	
	
	#Check if nsfw
	nsfw = False
	if 'nsfw' in str(message.channel):  # future reference: so there is a function that adds a nsfw checker,
		nsfw = True						# BUT IT DOESNT FUCKING WORK
	if message.author == client.user: # have the bot ignore its own messages
		return
	if message.content.startswith("!meme"): # for memes 
		if recv[6:] is '': # if the command is just '!meme'
			raw_msg = ""
			while True: # loop so if it fails it can find another post
				raw_msg = get_post_thing(["dankmemes","funny","memes","dank_meme"], nsfw=nsfw) # get a random post from a random choice of this random list of random subreddits
				if not raw_msg[1]==None: # break if it finds a vaild post, marked with a None value
					break
			if "Error" in str(raw_msg[2]): # post the error message if it fails
				print(raw_msg[2])
				print(raw_msg[0])
				print(raw_msg[4])
				print("------")
				await client.send_message(message.channel, content=raw_msg[2])
				await client.send_message(message.channel, content=raw_msg[0])
				await client.send_message(message.channel, content=raw_msg[4])
				return
			elif any(n in raw_msg[0] for n in filetypes):
				print("Posting on {}:".format(message.channel))
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				await client.send_message(message.channel, content=str(raw_msg[0]), tts=False)
				await client.send_message(message.channel, content=raw_msg[2], tts=False)
				await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
			else: # post the meme if it works
				print("Posting on {}:".format(message.channel))
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				embed = discord.Embed(title=raw_msg[2], url=raw_msg[0])
				embed.set_image(url=raw_msg[0])
				await client.send_message(message.channel, embed=embed, tts=False)
				await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
		else: # if there is a subreddit after the command
			raw_msg = ""
			count=0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]], nsfw=nsfw) # get a random post from the subreddit
				if not raw_msg[1]==None:
					break
				if count>=10:
					print("Count exceeded!")
					break
			if "Error" in str(raw_msg[2]): # if the post search fails
				print(raw_msg[2])
				print(raw_msg[0])
				print("------")
				await client.send_message(message.channel, content=raw_msg[2])
				await client.send_message(message.channel, content=raw_msg[0])
				await client.send_message(message.channel, content=raw_msg[4])
				return
			elif count>=10: # also a failsafe
				await client.send_message(message.channel, "Something went wrong, please try again!")
				return
				# await client.send_message(message.channel, "Problem subreddit: https://reddit.com/{}".format(recv[6:]))
			elif any(n in raw_msg[0] for n in filetypes):
				print("Posting on {}:".format(message.channel))
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				await client.send_message(message.channel, content=raw_msg[2], tts=False)
				await client.send_message(message.channel, content=str(raw_msg[0]), tts=False)
				await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
				#this block of code below thats commented out was causing issues when pulling from basically any subreddit other than the default
				#listed ones, so I made that the previous block's duty to pull from unofficial subreddits
			'''else: # if it finds an ok post
				print("Posting on {}:".format(message.channel))
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				embed = discord.Embed(title=raw_msg[2], url="https://reddit.com{}".format(raw_msg[3]))
				embed.set_image(url=raw_msg[0])
				await client.send_message(message.channel, embed=embed, tts=False)
				await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
'''
	if message.content.startswith("!joke") or message.content.startswith("!text"): # for jokes
		if recv[6:] is '': # gets it from the default subreddits
			raw_msg = ""
			while True:
				raw_msg = get_post_thing(["jokes", "darkjokes"], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
			print("Posting on {}:".format(message.channel))#post it after the loop
			print(raw_msg[2])
			print(raw_msg[0]) 
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=True)
			await client.send_message(message.channel, content=raw_msg[0], tts=True)
			await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
		else: # if a sub is specified
			raw_msg = ""
			premsg = "Original Post:"
			count = 0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
				if count>=10: # if it fails
					print("Count exceeded!") 
					raw_msg[2] = "Something went wrong."
					raw_msg[0] = "Please try again"
					raw_msg[3] = "/{}".format(recv[6:])
					premsg = "Subreddit that caused the issue:"
					break

			print("Posting on {}:".format(message.channel)) # executes both if it fails or if it works
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=False)
			await client.send_message(message.channel, content=raw_msg[0], tts=False)
			await client.send_message(message.channel, content="{} https://reddit.com{}".format(premsg, raw_msg[3]), tts=False)
			return

	if message.content.startswith('!news') or message.content.startswith("!link"): # for news
		raw_msg = ""
		count = 0
		while True:
			count+=1
			if not recv[6:]=='':
				raw_msg = get_post_thing(subs=[recv[6:]], nsfw=nsfw)
			else:
				raw_msg = get_post_thing(["UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"])
			if not raw_msg[1]==None:
				break
			if count>=10: # if it fails
				print("Count exceeded!")
				raw_msg[2] = "Something went wrong."
				raw_msg[0] = "Please try again"
				raw_msg[3] = "/{}".format(recv[6:])
				premsg = "Subreddit that caused the issue:"
				break
		print("Posting on {}:".format(message.channel)) # for both if it does and does not fail
		print(raw_msg[2])
		print(raw_msg[0])
		print("------")
		await client.send_message(message.channel, content=raw_msg[2], tts=True)
		await client.send_message(message.channel, content="Link: {}".format(raw_msg[0]), tts=False)
		return

	if message.content.startswith('!5050') or message.content.startswith('!fiftyfifty'):
		raw_msg = ""
		if 'nsfw' in str(message.channel):
			while True:
				raw_msg = get_post_thing(["fiftyfifty"], nsfw=True)
				if not raw_msg[1]==None:
					break
			print("Posting on {}:".format(message.channel))
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=False)
			await client.send_message(message.channel, content=str(raw_msg[0]), tts=False)
			await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
		else:
			while True:
				raw_msg = get_post_thing(["fiftyfifty"], nsfw=False)
				if not raw_msg[1]==None:
					break
			print("Posting on {}:".format(message.channel))
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=False)
			await client.send_message(message.channel, content=str(raw_msg[0]), tts=False)
			await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
	
	if message.content.startswith("!post"):
		raw_msg = ""
		if recv[6:]=='':
			await client.send_message(message.channel, content="Error, no post specified! Please try again.", tts=False)
		else:
			postid = recv[6:]
			raw_msg = get_post_by_id(subid=postid, nsfw=nsfw)
			print("Posting on {}:".format(message.channel))
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=False)
			await client.send_message(message.channel, content=str(raw_msg[0]), tts=False)
			await client.send_message(message.channel, content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return



@client.event # the on_ready event
async def on_ready():
	print('Logged in as')
	print(client.user.name)
	print(client.user.id)
	print('------')

while True: # run the bot forever
	try:
		client.run(token) # the bot run command
		'''
		current = time.time()
		if current>=t:
			log_channels(channel_list)
			t = current+300
		'''
	except Exception as e: # kill it if a keyboard interrupt is invoked
		if "Event loop" in str(e):
			print("\nStopping bot....")
			break
		else:
			print(e)
			continue
