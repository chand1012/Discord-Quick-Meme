#!/usr/bin/env python3
import discord
from lib import *
#import time
# get automatic channel logging working
#t = time.time()
url_things = ['.jpg', '.png', '.jpeg']
token = json_extract('token')
#channel_list = []
client = discord.Client()

@client.event
async def on_message(message):
	# getting nsfw filter working is first priority

	#channel = str(message.channel)
	#if not channel in channel_list:
	#	channel_list += [channel]
	recv = message.content
	channel = message.channel
	nsfw = False
	if str(message.channel) == "nsfw":
		nsfw = True
	if message.author == client.user:
		return
	if message.content.startswith("!meme"):
		if recv[6:] is '':
			raw_msg = ""
			while True:
				raw_msg = get_post_thing(["dankmemes","funny","memes","dank_meme"], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
			if nsfw==False and raw_msg[3]==False:
				print(raw_msg[2])
				print(raw_msg[0])
				print("------")
				await client.send_message(message.channel, content=raw_msg[2])
				await client.send_message(message.channel, content=raw_msg[0])
			else:
				print("Posting:")
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				embed = discord.Embed(title=raw_msg[2], url=raw_msg[0])
				embed.set_image(url=raw_msg[0])
				await client.send_message(message.channel, embed=embed, tts=False)
				await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)
		else:
			raw_msg = ""
			count=0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
				if count>=10:
					print("Count exceeded!")
					break
			if nsfw==False and raw_msg[3]==False:
				print(raw_msg[2])
				print(raw_msg[0])
				print("------")
				await client.send_message(message.channel, content=raw_msg[2])
				await client.send_message(message.channel, content=raw_msg[0])
			elif count>=10:
				await client.send_message(message.channel, "Something went wrong, please try again!")
				await client.send_message(message.channel, "Problem subreddit: https://reddit.com/{}".format(recv[6:]))
			else:
				print("Posting:")
				print(raw_msg[2])
				print(raw_msg[0])
				print("Original post: https://reddit.com{}".format(raw_msg[3]))
				print("------")
				embed = discord.Embed(title=raw_msg[2], url="https://reddit.com{}".format(raw_msg[3]))
				embed.set_image(url=raw_msg[0])
				await client.send_message(message.channel, embed=embed, tts=False)
				await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)



	if message.content.startswith("!joke"):
		if recv[6:] is '':
			raw_msg = ""
			while True:
				raw_msg = get_post_thing(["jokes", "darkjokes"], nsfw=nsfw)
				if not raw_msg[1]==None:
					break
			print("Posting:")
			print(raw_msg[2])
			print(raw_msg[0])
			print("Original post: https://reddit.com{}".format(raw_msg[3]))
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=True)
			await client.send_message(message.channel, content=raw_msg[0], tts=True)
			await client.send_message(message.channel, content="Original post: https://reddit.com{}".format(raw_msg[3]), tts=False)
		else:
			raw_msg = ""
			premsg = "Original Post:"
			count = 0
			while True:
				count+=1
				raw_msg = get_post_thing([recv[6:]], nsfw=nsfw)
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
			print("------")
			await client.send_message(message.channel, content=raw_msg[2], tts=True)
			await client.send_message(message.channel, content=raw_msg[0], tts=True)
			await client.send_message(message.channel, content="{} https://reddit.com{}".format(premsg, raw_msg[3]), tts=False)

	if message.content.startswith('!news'):
		raw_msg = ""
		count = 0
		while True:
			count+=1
			if not recv[6:]=='':
				raw_msg = get_post_thing([recv[6:]], 1, True, nsfw)
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
		print("------")
		await client.send_message(message.channel, content=raw_msg[2], tts=True)
		await client.send_message(message.channel, content="Link: {}".format(raw_msg[0]), tts=False)

@client.event
async def on_ready():
	print('Logged in as')
	print(client.user.name)
	print(client.user.id)
	print('------')

#shows server restart
#log_channels(['--------------'])

while True:
	try:
		client.run(token)
		'''
		current = time.time()
		if current>=t:
			log_channels(channel_list)
			t = current+300
		'''
	except Exception as e:
		if "Event loop" in str(e):
			print("\nStopping bot....")
			break
		else:
			print(e)
			continue
