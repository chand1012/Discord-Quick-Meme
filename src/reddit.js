import { SUBS } from './constants';
import { IMAGE_FILE_ENDINGS } from './constants';
import { getSubCache, setSubCache } from './subCache';

/**
 * Reach out to the reddit API, and get the first page of results from
 * r/aww. Filter out posts without readily available images or videos,
 * and return a random result.
 * @returns The url of an image or video which is cute.
 */

const guessPostHint = (post) => {
  // check if the url has any of the image file endings
  const hasImageFileEnding = IMAGE_FILE_ENDINGS.some((ending) =>
    post.url.endsWith(ending)
  );

  if (hasImageFileEnding) {
    return 'image';
  }

  if (post.url.includes('youtube.com') || post?.is_video) {
    return 'video';
  }
  if (post.selftext !== '') {
    return 'text';
  }
  return 'link';
};

const formatPost = (post, subreddit) => {
  return {
    title: post.data.title,
    permalink: `https://reddit.com${post.data.permalink}`,
    content:
      post.data?.selftext ||
      post.data?.media?.reddit_video?.fallback_url ||
      post.data?.secure_media?.reddit_video?.fallback_url ||
      post.data?.url,
    media_url:
      post.data?.media?.reddit_video?.fallback_url ||
      post.data?.secure_media?.reddit_video?.fallback_url ||
      post.data?.url,
    nsfw: post.data.over_18,
    sub: subreddit,
    score: post.data.score,
    is_video: post.data?.is_video,
    hint: post.data?.post_hint || guessPostHint(post.data),
  };
};

export const getRedditData = async (subreddit, env) => {
  const cachedData = await getSubCache(subreddit, env);
  if (cachedData) {
    return JSON.parse(cachedData);
  }

  const response = await fetch(
    `https://www.reddit.com/r/${subreddit}/hot.json`,
    {
      headers: {
        'User-Agent': 'quickmeme-bot',
      },
    }
  );
  const data = await response.json();
  await setSubCache(subreddit, JSON.stringify(data), env);
  return data;
};

export const getRandomTextPost = async (subreddit, env) => {
  const data = await getRedditData(subreddit, env);
  const posts = data.data.children.map((post) => {
    if (post.selftext) {
      return formatPost(post, subreddit);
    }
  });
  const post = posts[Math.floor(Math.random() * posts.length)];
  return post;
};

export const getRandomMediaPost = async (subreddit, env) => {
  const data = await getRedditData(subreddit, env);
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

export const getRandomLinkPost = async (subreddit, env) => {
  const data = await getRedditData(subreddit, env);
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

export const getCuteUrl = async (env) => {
  return getRandomMediaPost('aww', env);
};

export const getMeme = async (env) => {
  // choose a random subreddit from the memes list
  const { memes } = SUBS;
  const randomSubreddit = memes[Math.floor(Math.random() * memes.length)];
  return getRandomMediaPost(randomSubreddit, env);
};
