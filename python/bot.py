#!/usr/bin/env python3
import discord
import logging
import time
import os
from lib import get_rnd_post, get_post_by_id, trump_counter
from json_extract import json_extract
token = json_extract('token')
client = discord.Client()
filetypes = ['gif', 'gifv', 'gfycat', 'v.redd.it', 'youtube', 'youtu.be', '.jpg', '.png', '.jpeg']
postfile = 'posts.json'
channels = {}
memesubs = ["dankmemes", "funny", "memes", "dank_meme", "comedyheaven", "CyanideandHappiness", "therewasanattempt", "wholesomememes", "instant_regret"]
commands = ["!meme", "!joke", "!5050", "!news", "!fiftyfifty", "!text", "!post"]
#logging things
logger = None
def initialize_logger(output_dir):
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
     
    # create console handler and set level to info
    handler = logging.StreamHandler()
    handler.setLevel(logging.INFO)
    formatter = logging.Formatter("%(levelname)s - %(message)s")
    handler.setFormatter(formatter)
    logger.addHandler(handler)
 
    # create error file handler and set level to error
    handler = logging.FileHandler(os.path.join(output_dir, "error.log"),"w", encoding=None, delay="true")
    handler.setLevel(logging.ERROR)
    formatter = logging.Formatter("%(levelname)s - %(message)s")
    handler.setFormatter(formatter)

initialize_logger('.')

@client.event
async def on_message(message):
	# The recv stuff
	recv = message.content
	now = time.time()
	if not message.channel in channels: # if a channel is not in the dictionary, put it in the dictionary.
			channels[message.channel] = [time.time(), 0, True] # later this should be made a json
	# spam filter code
	if any(message.content.startswith(thing) for thing in commands):
		if channels[message.channel][1]>=5 and channels[message.channel][0]+60>now and channels[message.channel][2]:
			channels[message.channel][2] = False
			await message.channel.send(content="You guys have been sending a lot of messages, why don't you slow down a bit?") 
			return
		elif channels[message.channel][1]<5 and channels[message.channel][0]+60>now:
			channels[message.channel][1] += 1
		elif channels[message.channel][1]<5 and channels[message.channel][0]+60<now:
			channels[message.channel][1] = 0
			channels[message.channel][0] = time.time()
			channels[message.channel][2] = True
	
	
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
				raw_msg = get_rnd_post(subs=memesubs, nsfw=nsfw) # get a random post from a random choice of this random list of random subreddits
				if not raw_msg[1]==None: # break if it finds a vaild post, marked with a None value
					break
			if "Error" in str(raw_msg[2]): # post the error message if it fails
				logging.info(raw_msg[2])
				logging.info(raw_msg[0])
				logging.info(raw_msg[4])
				logging.info("------")
				await message.channel.send(content=raw_msg[2])
				await message.channel.send(content=raw_msg[0])
				await message.channel.send(content=raw_msg[4])
				return
			elif any(n in raw_msg[0] for n in filetypes):
				logging.info("Posting on {}:".format(message.channel))
				logging.info(raw_msg[2])
				logging.info(raw_msg[0])
				logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
				logging.info("------")
				await message.channel.send(content=str(raw_msg[0]), tts=False)
				await message.channel.send(content=raw_msg[2], tts=False)
				await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
			else: # post the meme if it works
				logging.info("Posting on {}:".format(message.channel))
				logging.info(raw_msg[2])
				logging.info(raw_msg[0])
				logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
				logging.info("------")
				embed = discord.Embed(title=raw_msg[2], url=raw_msg[0])
				embed.set_image(url=raw_msg[0])
				await message.channel.send(embed=embed, tts=False)
				await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
		else: # if there is a subreddit after the command
			raw_msg = ""
			count=0
			while True:
				count+=1
				raw_msg = get_rnd_post([recv[6:]], nsfw=nsfw) # get a random post from the subreddit
				if not raw_msg[1]==None:
					break
				if count>=10:
					logging.info("Count exceeded!")
					break
			if "Error" in str(raw_msg[2]): # if the post search fails
				logging.info(raw_msg[2])
				logging.info(raw_msg[0])
				logging.info("------")
				await message.channel.send(content=raw_msg[2])
				await message.channel.send(content=raw_msg[0])
				await message.channel.send(content=raw_msg[4])
				return
			elif count>=10: # also a failsafe
				await message.channel.send( "Something went wrong, please try again!")
				return
			elif any(n in raw_msg[0] for n in filetypes):
				logging.info("Posting on {}:".format(message.channel))
				logging.info(raw_msg[2])
				logging.info(raw_msg[0])
				logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
				logging.info("------")
				await message.channel.send(content="{} from r/{}".format(raw_msg[2], raw_msg[5]), tts=False)
				await message.channel.send(content=str(raw_msg[0]), tts=False)
				await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
				return
				
	if message.content.startswith("!joke") or message.content.startswith("!text"): # for jokes
		if recv[6:] is '': # gets it from the default subreddits
			raw_msg = ""
			while True:
				raw_msg = get_rnd_post(["jokes", "darkjokes"], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
			logging.info("Posting on {}:".format(message.channel))#post it after the loop
			logging.info(raw_msg[2])
			logging.info(raw_msg[0]) 
			logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
			logging.info("------")
			await message.channel.send(content=raw_msg[2], tts=True)
			await message.channel.send(content=raw_msg[0], tts=True)
			await message.channel.send(content=raw_msg[5], tts=True)
			await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
		else: # if a sub is specified
			raw_msg = ""
			premsg = "Original Post:"
			count = 0
			while True:
				count+=1
				raw_msg = get_rnd_post([recv[6:]], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
				if count>=10: # if it fails
					logging.info("Count exceeded!") 
					raw_msg[2] = "Something went wrong."
					raw_msg[0] = "Please try again"
					raw_msg[3] = "/{}".format(recv[6:])
					premsg = "Subreddit that caused the issue:"
					break

			logging.info("Posting on {}:".format(message.channel)) # executes both if it fails or if it works
			logging.info(raw_msg[2])
			logging.info(raw_msg[0])
			logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
			logging.info("------")
			await message.channel.send(content=raw_msg[2], tts=False)
			await message.channel.send(content=raw_msg[0], tts=False)
			await message.channel.send(content=raw_msg[5], tts=True)
			await message.channel.send(content="{} https://reddit.com{}".format(premsg, raw_msg[3]), tts=False)
			return

	if message.content.startswith('!news') or message.content.startswith("!link"): # for news
		raw_msg = ""
		count = 0
		while True:
			count+=1
			if not recv[6:]=='':
				raw_msg = get_rnd_post(subs=[recv[6:]], nsfw=nsfw)
			else:
				raw_msg = get_rnd_post(["UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"])
			if not raw_msg[1]==None:
				break
			if count>=10: # if it fails
				logging.info("Count exceeded!")
				raw_msg[2] = "Something went wrong."
				raw_msg[0] = "Please try again"
				raw_msg[3] = "/{}".format(recv[6:])
				premsg = "Subreddit that caused the issue:"
				break
		logging.info("Posting on {}:".format(message.channel)) # for both if it does and does not fail
		logging.info(raw_msg[2])
		logging.info(raw_msg[0])
		logging.info("------")
		await message.channel.send(content=raw_msg[2], tts=True)
		await message.channel.send(content=raw_msg[5], tts=True)
		await message.channel.send(content="Link: {}".format(raw_msg[0]), tts=False)
		return

	if message.content.startswith('!5050') or message.content.startswith('!fiftyfifty'):
		raw_msg = ""
		if 'nsfw' in str(message.channel):
			while True:
				raw_msg = get_rnd_post(["fiftyfifty"], nsfw=True)
				if not raw_msg[1]==None:
					break
			logging.info("Posting on {}:".format(message.channel))
			logging.info(raw_msg[2])
			logging.info(raw_msg[0])
			logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
			logging.info("------")
			await message.channel.send(content=raw_msg[2], tts=False)
			await message.channel.send(content=str(raw_msg[0]), tts=False)
			await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
		else:
			while True:
				raw_msg = get_rnd_post(["fiftyfifty"], nsfw=False)
				if not raw_msg[1]==None:
					break
			logging.info("Posting on {}:".format(message.channel))
			logging.info(raw_msg[2])
			logging.info(raw_msg[0])
			logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
			logging.info("------")
			await message.channel.send(content=raw_msg[2], tts=False)
			await message.channel.send(content=str(raw_msg[0]), tts=False)
			await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
	
	if message.content.startswith("!post"):
		raw_msg = ""
		if recv[6:]=='':
			await message.channel.send(content="Error, no post specified! Please try again.", tts=False)
		else:
			postid = recv[6:]
			raw_msg = get_post_by_id(subid=postid, nsfw=nsfw)
			logging.info("Posting on {}:".format(message.channel))
			logging.info(raw_msg[2])
			logging.info(raw_msg[0])
			logging.info("Original post: https://reddit.com{}".format(raw_msg[3]))
			logging.info("------")
			await message.channel.send(content="{} from r/{}".format(raw_msg[2], raw_msg[5]), tts=False)
			await message.channel.send(content=str(raw_msg[0]), tts=False)
			await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
			return
	
	if message.content.startswith("!hentai"):
		hentai = ['ahegao', 'Artistic_Hentai', 'Hentai', 'MonsterGirl', 'slimegirls', 'wholesomehentai', 'quick_hentai', 'HentaiParadise'] # my friends are sick fucks
		raw_msg = get_rnd_post(subs=hentai, nsfw=True)
		await message.channel.send(content="{} from r/{}".format(raw_msg[2], raw_msg[5]), tts=False)
		await message.channel.send(content=str(raw_msg[0]), tts=False)
		await message.channel.send(content="Score: {}\nOriginal post: https://reddit.com{}".format(raw_msg[4], raw_msg[3]), tts=False)
		return
	
	if message.content.startswith("!trumpcount"):
		tcount_raw = trump_counter()
		tcount = tcount_raw[0]
		sample = tcount_raw[1]
		samplesize = tcount_raw[2]
		await message.channel.send(content="As of a few seconds ago, the number of times Trump or Russia was mentioned on r/politics, r/POLITIC, r/news, and r/neutralnews is {}.".format(tcount))
		await message.channel.send(content="This was done with a total sample of {}, with {} per subreddit.".format(sample, samplesize))

@client.event # the on_ready event
async def on_ready():
	logging.info('Logged in as')
	logging.info(client.user.name)
	logging.info(client.user.id)
	logging.info('------')

while True: # run the bot forever
	try:
		client.run(token) # the bot run command
	except Exception as e: # kill it if a keyboard interrupt is invoked
		if "Event loop" in str(e):
			logging.info("\nStopping bot....")
			break
		else:
			logging.info(e)
			continue
