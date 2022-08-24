/* eslint-disable no-undef */

export const getChannel = async (channelID, env) => {
  const applicationId = env.DISCORD_APPLICATION_ID;
  const url = `https://discord.com/api/v10/applications/${applicationId}/channels/${channelID}`;
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${env.DISCORD_TOKEN}`,
    },
    method: 'GET',
  });

  return await response.json();
};

export const isChannelNSFW = async (channelID, env) => {
  const channel = await getChannel(channelID, env);
  return channel?.nsfw || false;
};
