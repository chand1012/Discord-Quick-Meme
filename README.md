# Discord-Quick-Meme Rewrite

This is the serverless rewrite of Discord-Quick-Meme. This version aims to make Discord-Quick-Meme up to modern Discord bot standards, as well as being faster and more reliable.

## Serverless?

Currently, Discord-Quick-Meme is a "Gateway Bot", meaning it has to run all the time on a computer somewhere. Over the years, its migrated between a few different machines, from my crappy old laptop I used as a server, to my Raspberry Pi, to an actual server, back to the Pi, repeat. Currently its being hosted on a Linode server, however I would like to change this.

Discord-Quick-Meme will be moving to a serverless architecture, so it doesn't just run in one single place, it will run on Cloudflare's entire Workers network. Using Workers' high-speed KV storage as a cache, it'll make Discord-Quick-Meme even faster!

## When will it be done?

When it's ready.

# To Do
* Fix video embeds.
* Get all the commands working as they are on the mainline branch.
* Cache all subreddit data in [Workers KV](https://developers.cloudflare.com/workers/runtime-apis/kv/) for one to two hours.
* Get queue working using [Workers Scheduling and cron](https://developers.cloudflare.com/workers/examples/multiple-cron-triggers/).
  + This would run once every 5 to 15 minutes and see if the time is right for a post to be made.
* Convert to Typescript.
