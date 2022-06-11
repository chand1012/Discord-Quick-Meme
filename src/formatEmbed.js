import constructEmbed, {
  constructImageEmbed,
  constructVideoEmbed,
} from './embed';

const formatEmbed = (post) => {
  let objects = {
    image: {},
    video: {},
    footer: {},
    thumbnail: {},
    author: {},
    provider: {},
    fields: [],
  };

  if (post.hint === 'image') {
    objects['image'] = constructImageEmbed(post.media_url, post.media_url);
  } else if (post.hint.includes('video')) {
    objects['video'] = constructVideoEmbed(post.media_url, post.media_url);
  }

  return constructEmbed(
    post.title,
    'rich',
    `From ${post.sub}`,
    post.permalink,
    objects
  );
};

export default formatEmbed;
