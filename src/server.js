/**
 * The core server that runs on a Cloudflare worker.
 */

import ittyRouter from 'itty-router';
const { Router } = ittyRouter;
import {
  InteractionResponseType,
  InteractionType,
  verifyKey,
} from 'discord-interactions';
import { MEME_COMMAND } from './commands';
import { getMeme } from './reddit';
import JsonResponse from './jsonResponse';
import { REDIRECT_URL } from './constants';
import formatEmbed from './formatEmbed';

const router = Router();

/**
 * A simple :wave: hello page to verify the worker is working.
 */
// eslint-disable-next-line no-unused-vars
router.get('/', (request, env) => {
  return Response.redirect(REDIRECT_URL);
});

/**
 * Main route for all requests sent from Discord.  All incoming messages will
 * include a JSON payload described here:
 * https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object
 */
router.post('/', async (request, env) => {
  const message = await request.json();
  console.log(message);
  if (message.type === InteractionType.PING) {
    // The `PING` message is used during the initial webhook handshake, and is
    // required to configure the webhook in the developer portal.
    console.log('Handling Ping request');
    return new JsonResponse({
      type: InteractionResponseType.PONG,
    });
  }

  if (message.type === InteractionType.APPLICATION_COMMAND) {
    // Most user commands will come as `APPLICATION_COMMAND`.
    switch (message.data.name.toLowerCase()) {
      case MEME_COMMAND.name.toLowerCase(): {
        console.log('handling meme request');
        const data = await getMeme();
        console.log('got meme');
        const embed = formatEmbed(data);
        console.log('formatted embed');
        const respData = {
          type: 4,
          data: {
            embeds: [embed],
          },
        };
        console.log('constructed object');
        return new JsonResponse(respData);
      }
      default:
        console.error('Unknown Command');
        return new JsonResponse({ error: 'Unknown Type' }, { status: 400 });
    }
  }

  console.error('Unknown Type');
  return new JsonResponse({ error: 'Unknown Type' }, { status: 400 });
});
router.all('*', () => new Response('Not Found.', { status: 404 }));

export default {
  /**
   * Every request to a worker will start in the `fetch` method.
   * Verify the signature with the request, and dispatch to the router.
   * @param {*} request A Fetch Request object
   * @param {*} env A map of key/value pairs with env vars and secrets from the cloudflare env.
   * @returns
   */
  async fetch(request, env) {
    if (request.method === 'POST') {
      // Using the incoming headers, verify this request actually came from discord.
      const signature = request.headers.get('x-signature-ed25519');
      const timestamp = request.headers.get('x-signature-timestamp');
      console.log(signature, timestamp, env.DISCORD_PUBLIC_KEY);
      const body = await request.clone().arrayBuffer();
      const isValidRequest = verifyKey(
        body,
        signature,
        timestamp,
        env.DISCORD_PUBLIC_KEY
      );
      if (!isValidRequest) {
        console.error('Invalid Request');
        return new Response('Bad request signature.', { status: 401 });
      }
    }

    // Dispatch the request to the appropriate route
    return router.handle(request, env);
  },
};
