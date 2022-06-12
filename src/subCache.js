// caching system using workers k

// these just don't work for some reason

/* eslint-disable no-undef */

export const setSubCache = async (subreddit, data) => {
  await REDDIT_CACHE.put(subreddit, data, { expirationTtl: 3600 });
};

export const getSubCache = async (subreddit) => {
  return REDDIT_CACHE.get(subreddit);
};
