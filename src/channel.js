/* eslint-disable no-undef */

export const getChannel = async (channelID) => {
  const applicationId = process.env.DISCORD_APPLICATION_ID;
  const url = `https://discord.com/api/v10/applications/${applicationId}/channels/${channelID}`;
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${process.env.DISCORD_TOKEN}`,
    },
    method: 'GET',
  });

  return await response.json();
};

export const isChannelNSFW = async (channelID) => {
  const channel = await getChannel(channelID);
  return channel?.nsfw || false;
};
