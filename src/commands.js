/**
 * Share command metadata from a common spot to be used for both runtime
 * and registration.
 */

// export const AWW_COMMAND = {
//   name: 'awwww',
//   description: 'Drop some cuteness on this channel.',
// };

// export const INVITE_COMMAND = {
//   name: 'invite',
//   description: 'Get an invite link to add the bot to your server',
// };

export const MEME_COMMAND = {
  name: 'meme',
  description: 'Get a random meme from reddit.',
  options: [
    {
      name: 'subreddit',
      description: 'The subreddit to get a meme from.',
      required: false,
      type: 3,
    },
  ],
};

export const NEWS_COMMAND = {
  name: 'news',
  description: 'Get the latest news from reddit.',
  options: [
    {
      name: 'subreddit',
      description: 'The subreddit to get news from.',
      required: false,
      type: 3,
    },
  ],
};

export const JOKE_COMMAND = {
  name: 'joke',
  description: 'Get a random joke from reddit.',
  options: [
    {
      name: 'subreddit',
      description: 'The subreddit to get a joke from.',
      required: false,
      type: 3,
    },
  ],
};

export const FIFTY_50_COMMAND = {
  name: '5050',
  description: 'Get a random 50/50 from reddit.',
};

export const FIFTY_FIFTY_COMMAND = {
  name: 'fiftyfifty',
  description: 'Get a random 50/50 from reddit.',
};

export const ALL_COMMAND = {
  name: 'all',
  description: 'Get a random post from reddit.',
};

export const TEXT_COMMAND = {
  name: 'text',
  description: 'Get a random text post from reddit.',
  options: [
    {
      name: 'subreddit',
      description: 'The subreddit to get a text post from.',
      required: false,
      type: 3,
    },
  ],
};

export const BUZZWORD_COMMAND = {
  name: 'buzzword',
  description: 'Get a random buzzword from reddit.',
};

export const LINK_COMMAND = {
  name: 'link',
  description: 'Get a random link from reddit.',
  options: [
    {
      name: 'subreddit',
      description: 'The subreddit to get a link from.',
      required: false,
      type: 3,
    },
  ],
};

export const HENTAI_COMMAND = {
  name: 'hentai',
  description: 'Get a random hentai from reddit. Requires an NSFW channel.',
};

export default [
  MEME_COMMAND,
  NEWS_COMMAND,
  JOKE_COMMAND,
  FIFTY_50_COMMAND,
  FIFTY_FIFTY_COMMAND,
  ALL_COMMAND,
  TEXT_COMMAND,
  BUZZWORD_COMMAND,
  LINK_COMMAND,
  HENTAI_COMMAND,
];
