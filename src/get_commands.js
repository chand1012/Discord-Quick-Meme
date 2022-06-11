import fetch from 'node-fetch';
import dotenv from 'dotenv';

dotenv.config();

/**
 * This file is meant to be run from the command line, and is not used by the
 * application server.  It's allowed to use node.js primitives, and only needs
 * to be run once.
 */

/* eslint-disable no-undef */

const token = process.env.DISCORD_TOKEN;
const applicationId = process.env.DISCORD_APPLICATION_ID;
const testGuildId = process.env.DISCORD_TEST_GUILD_ID;

if (!token) {
  throw new Error('The DISCORD_TOKEN environment variable is required.');
}
if (!applicationId) {
  throw new Error(
    'The DISCORD_APPLICATION_ID environment variable is required.'
  );
}

async function getGlobalCommands() {
  const url = `https://discord.com/api/v10/applications/${applicationId}/commands`;
  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${token}`,
    },
    method: 'GET',
  });
  const json = await res.json();
  return json;
}

async function getGuildCommands() {
  const url = `https://discord.com/api/v10/applications/${applicationId}/guilds/${testGuildId}/commands`;
  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${token}`,
    },
    method: 'GET',
  });
  const json = await res.json();
  return json;
}

console.log({ global: await getGlobalCommands() });
console.log({ guild: await getGuildCommands() });
