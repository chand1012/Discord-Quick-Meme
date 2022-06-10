export const constructProviderEmbed = (name = undefined, url = undefined) => {
  return { name, url };
};

export const constructAuthorEmbed = (
  name,
  url = undefined,
  iconUrl = undefined,
  proxyIconUrl = undefined
) => {
  if (name.length > 256) {
    throw new Error('Name must be less than 256 characters');
  }
  return {
    name,
    url,
    icon_url: iconUrl,
    proxy_icon_url: proxyIconUrl,
  };
};

export const constructImageEmbed = (
  url,
  proxyUrl = undefined,
  height = undefined,
  width = undefined
) => {
  return {
    url,
    proxy_url: proxyUrl,
    height,
    width,
  };
};

export const constructVideoEmbed = constructImageEmbed;
export const constructThumbnailEmbed = constructImageEmbed;

export const constructField = (name, value, inline = false) => {
  if (name.length > 256) {
    throw new Error('Name must be less than 256 characters');
  }
  if (value.length > 1024) {
    throw new Error('Value must be less than 1024 characters');
  }
  return {
    name,
    value,
    inline,
  };
};

export const constructFooter = (
  text,
  iconUrl = undefined,
  proxyIconUrl = undefined
) => {
  if (text.length > 2048) {
    throw new Error('Text must be less than 2048 characters');
  }
  return {
    text,
    icon_url: iconUrl,
    proxy_icon_url: proxyIconUrl,
  };
};

const constructEmbed = (
  title,
  type,
  description,
  url,
  objects = {
    image: {},
    video: {},
    footer: {},
    thumbnail: {},
    author: {},
    provider: {},
    fields: [],
  }
) => {
  // YOU CAN ONLY SPECIFY ONE OF THE OPTIONAL OBJECTS
  if (title.length > 256) {
    throw new Error('Title must be less than 256 characters');
  }
  if (description.length > 4096) {
    throw new Error('Description must be less than 4096 characters');
  }
  if (objects?.fields.length > 25) {
    throw new Error('25 field maximum');
  }
  return {
    title,
    type,
    description,
    url,
    ...objects,
  };
};

export default constructEmbed;
