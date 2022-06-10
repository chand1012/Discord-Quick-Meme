import { memes } from './subs.json';

/**
 * Reach out to the reddit API, and get the first page of results from
 * r/aww. Filter out posts without readily available images or videos,
 * and return a random result.
 * @returns The url of an image or video which is cute.
 */

const formatPost = (post, subreddit) => {
  return {
    title: post.data.title,
    permalink: `https://reddit.com${post.data.permalink}`,
    content:
      post.data?.media?.reddit_video?.fallback_url ||
      post.data?.secure_media?.reddit_video?.fallback_url ||
      post.data?.url,
    nsfw: post.dzata.over_18,
    sub: subreddit,
    score: post.data.score,
  };
};

export const getRandomTextPost = async (subreddit) => {
  const response = await fetch(
    `https://www.reddit.com/r/${subreddit}/hot.json`
  );
  const data = await response.json();
  const posts = data.data.children.map((post) => {
    if (post.selftext) {
      return formatPost(post, subreddit);
    }
  });
  const post = posts[Math.floor(Math.random() * posts.length)];
  return post;
};

export const getRandomMediaPost = async (subreddit) => {
  const response = await fetch(
    `https://www.reddit.com/r/${subreddit}/hot.json`,
    {
      headers: {
        'User-Agent': 'quickmeme-bot',
      },
    }
  );
  const data = await response.json();
  const posts = data.data.children.map((post) => {
    if (post.is_gallery) {
      return '';
    }
    if (
      post.data?.media?.reddit_video?.fallback_url ||
      post.data?.secure_media?.reddit_video?.fallback_url ||
      post.data?.url
    ) {
      return formatPost(post, subreddit);
    }
  });
  const randomIndex = Math.floor(Math.random() * posts.length);
  const randomPost = posts[randomIndex];
  return randomPost;
};

export const getRandomLinkPost = async (subreddit) => {
  const response = await fetch(
    `https://www.reddit.com/r/${subreddit}/hot.json`,
    {
      headers: {
        'User-Agent': 'quickmeme-bot',
      },
    }
  );
  const data = await response.json();
  const posts = data.data.children
    .map((post) => {
      const domain = post.data.domain;
      if (
        domain.includes('redd.it') ||
        domain.includes('reddit.com') ||
        domain.includes('imgur.com')
      ) {
        return null;
      }
      return formatPost(post, subreddit);
    })
    .filter((post) => post !== null);
  const randomIndex = Math.floor(Math.random() * posts.length);
  const randomPost = posts[randomIndex];
  return randomPost;
};

export const getCuteUrl = async () => {
  return getRandomMediaPost('aww');
};

export const getMemeUrl = async () => {
  // choose a random subreddit from the memes list
  const randomSubreddit = memes[Math.floor(Math.random() * memes.length)];
  return getRandomMediaPost(randomSubreddit);
};
