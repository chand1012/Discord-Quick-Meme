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

// eslint-disable-next-line no-unused-vars
async function removeGlobalCommands() {
  const url = `https://discord.com/api/v10/applications/${applicationId}/commands`;
  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${token}`,
    },
    method: 'GET',
  });
  const json = await res.json();
  json.forEach(async (cmd) => {
    const response = await fetch(
      `https://discord.com/api/v10/applications/${applicationId}/commands/${cmd.id}`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bot ${token}`,
        },
        method: 'DELETE',
      }
    );
    if (!response.ok) {
      console.error(`Problem removing command ${cmd.id}`);
    }
  });
}

// eslint-disable-next-line no-unused-vars
async function removeGuildCommands() {
  const url = `https://discord.com/api/v10/applications/${applicationId}/guilds/${testGuildId}/commands`;
  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bot ${token}`,
    },
    method: 'GET',
  });
  const json = await res.json();
  json.forEach(async (cmd) => {
    const response = await fetch(
      `https://discord.com/api/v10/applications/${applicationId}/guilds/${testGuildId}/commands/${cmd.id}`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bot ${token}`,
        },
        method: 'DELETE',
      }
    );
    if (!response.ok) {
      console.error(`Problem removing command ${cmd.id}`);
    }
  });
}

// await removeGlobalCommands();
await removeGuildCommands();
