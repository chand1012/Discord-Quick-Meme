import constructEmbed, {
  constructImageEmbed,
  constructVideoEmbed,
} from './embed';
import { IMAGE_FILE_ENDINGS } from './constants';

const formatEmbed = (post) => {
  // check if the post url is an image
  if (post?.media_url && !post?.content) {
    if (IMAGE_FILE_ENDINGS.some((ending) => post.media_url.endsWith(ending))) {
      const image = constructImageEmbed(post.media_url, post.media_url);
      return constructEmbed(
        post.title,
        'image',
        `From r/${post.sub}`,
        post.permalink,
        { image }
      );
    } else if (post?.is_video) {
      const video = constructVideoEmbed(post.media_url, post.media_url);
      return constructEmbed(
        post.title,
        'video',
        `From r/${post.sub}`,
        post.permalink,
        { video }
      );
    }
  }
};

export default formatEmbed;
