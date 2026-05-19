# Grimoire MCP

> Your campaign, organized by magic.

Grimoire is a structured TTRPG campaign brain for game masters. This MCP gives any AI client you connect a live, structured view of your campaign data and current situation in-game. Upload and parse notes into structured database entities, query your world by relationship, drive prep and recaps from the conversation you're already having.

| | |
|---|---|
| **Endpoint** | `https://api.ttrpg.bot/mcp` |
| **Transport** | Streamable HTTP (MCP spec 2025-06-18) |
| **Auth** | OAuth 2.1, one-click consent flow |
| **Tools** | 41 across 6 groups |
| **Web app** | https://www.ttrpg.bot |

## Three clicks to connect

1. **Add the connector.** In Claude.ai (or any MCP client that supports custom connectors), open Settings → Connectors → Add custom connector and paste `https://api.ttrpg.bot/mcp`.
2. **Sign in and pick a campaign.** Your client opens the Grimoire consent screen. Sign in with the account that owns the campaign, pick which campaign the AI should see, click Approve. Access is scoped to that one campaign, enforced server-side.
3. **Ask it something about your world.** *"Who knows about the cult?" "What was the queen's reaction in session 4?" "Generate three rumors that fit my current city."* The client reads live and answers from your data.

Full setup walkthrough with client-specific instructions: https://www.ttrpg.bot/docs/mcp/

## What you can do with it

### The fast wins
Parse a paragraph of session notes into NPCs, factions, and threads. Dig up the location you mentioned three sessions ago. Update a faction's status straight from a recap.

### The deeper plays
Model an NPC reaction grounded in motivations they actually have and a political web they actually sit inside. Project the next move House Vale makes, if it moves at all. Generate three rumors that fit your city, not generic fantasy.

### The MC-Possibilities
Chain Grimoire MCP to other MCPs. Pipe your last session into a slideshow MCP and a voiceover MCP for an automated play-by-play. Hook a calendar MCP and the bot reminds you to prep.

The pitch is flexibility, not a baked-in workflow. MCP is a protocol; once Grimoire is connected to a frontier model, what you ask the model to do with it is the workflow.

## The toolbox

41 tools across six groups.

**Entities** (9): `search_campaign`, `get_entity`, `list_entities`, `get_field_options`, `get_tag_options`, `create_entity`, `update_entity`, `delete_entity`, `batch_create_entities`. Covers all 13 entity types: NPCs, Locations, Factions, Quests, Items, Player Characters, Creatures, Lore Entries, World Rules, Planar Forces, Session Recaps, Session Preps, Custom Mechanics.

**Relationships** (3): `add_relationship`, `get_relationships`, `delete_relationship`. Typed edges between any two entities.

**Knowledge graphs** (8): `get_constitution`, `get_entity_catalog`, `get_knowledge_graph`, `list_entity_graphs`, `get_entity_graph`, `add_to_entity_graph`, `create_entity_graph_edge`, `toggle_graph_attention`. Political, timeline, and geography projections with visibility filtering.

**Open threads** (7): `get_open_threads`, `create_open_thread`, `resolve_open_thread`, `unresolve_open_thread`, `update_open_thread`, `get_thread_progressions`, `add_thread_progression`. Loose ends and how they evolve session over session.

**Wiki** (9): `get_wiki_tree`, `get_wiki_page`, `create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`, `delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`. Block-based collaborative pages; the AI reads and edits them like you do.

**Campaign meta** (5): `current_campaign`, `get_campaign_context`, `get_narrative_state`, `get_campaign_bible`, `update_campaign_bible`. `get_narrative_state` aggregates recent sessions, open threads, canonical facts, and active arcs into one view.

## Visibility and safety

Everything respects Grimoire's visibility model: `common-knowledge`, `player-knowledge`, `dm-secret`. The AI sees what you've marked it can see and nothing more. Player accounts connecting through MCP get filtered views; GM accounts get the full picture.

OAuth scopes are campaign-scoped. Granting access to one campaign doesn't grant access to your other campaigns.

## Pricing

End-to-end free to try. Free tier on Grimoire pairs with Claude.ai's free tier with no card required. Pro and Pro+ unlock larger campaign sizes and additional features.

Pricing: https://www.ttrpg.bot/#pricing

## Links

- Product: https://www.ttrpg.bot
- MCP overview and connect flow: https://www.ttrpg.bot/mcp
- Full setup guide: https://www.ttrpg.bot/docs/mcp/
- Changelog: https://www.ttrpg.bot/changelog
- Support: zmanlevelup@gmail.com

## About this repository

This is the public documentation and registry submission artifact for the Grimoire MCP. The MCP server itself is operated as a hosted service; source for the server is not public.

`server.json` in this repo is what's published to `registry.modelcontextprotocol.io` as [`bot.ttrpg/grimoire`](https://registry.modelcontextprotocol.io/v0.1/servers?search=bot.ttrpg/grimoire).

## License

MIT (for this documentation repository).
