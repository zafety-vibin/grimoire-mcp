# Grimoire MCP

> Your campaign, organized by magic.

| | |
|---|---|
| **Endpoint** | `https://api.ttrpg.bot/mcp` |
| **Transport** | Streamable HTTP (MCP spec 2025-06-18) |
| **Auth** | OAuth 2.1, one-click consent flow |
| **Tools** | 48 across 7 groups (20 read, 28 write) |
| **Web app** | https://www.ttrpg.bot |

## What is Grimoire?

Grimoire is a campaign manager for D&D and tabletop RPG gamemasters with a hosted MCP server built in. A campaign in Grimoire is structured data, not pages: typed NPCs, factions, locations, quests, items, sessions, and lore, connected by relationships and knowledge graphs. This MCP gives any AI client you connect a live, structured view of that database and your current situation in-game, so the AI answers from your actual canon instead of guessing. Upload and parse notes into structured database entities, query your world by relationship, drive prep and recaps from the conversation you're already having.

## How to use Grimoire?

Three clicks, nothing to install:

1. **Add the connector.** In Claude.ai or ChatGPT (or any MCP client that supports custom connectors), open Settings → Connectors → Add custom connector and paste `https://api.ttrpg.bot/mcp`.
2. **Sign in and pick a campaign.** Your client opens the Grimoire consent screen. Sign in with the account that owns the campaign, pick which campaign the AI should see, click Approve. Access is scoped to that one campaign, enforced server-side.
3. **Ask it something about your world.** *"Who knows about the cult?" "What was the queen's reaction in session 4?" "Generate three rumors that fit my current city."* The client reads live and answers from your data.

Full setup walkthrough with client-specific instructions: https://www.ttrpg.bot/docs/mcp/

For clients configured by JSON file instead of a connectors UI:

```json
{
  "mcpServers": {
    "grimoire": {
      "type": "streamable-http",
      "url": "https://api.ttrpg.bot/mcp"
    }
  }
}
```

Auth is OAuth 2.1 with PKCE; the client opens a browser sign-in on first connect. A free Grimoire account at https://www.ttrpg.bot provides the campaign the server reads.

## Key features of Grimoire

- 48 tools (20 read, 28 write) over a typed campaign database with 14 entity schemas
- Knowledge graph queries: political, timeline, and geography projections with visibility filtering
- Narrative state in one call: recent sessions, open plot threads, active arcs, and canon facts
- Full write path: create and update entities, log thread progressions, edit wiki pages and blocks
- GM secrets enforced server-side with three visibility tiers (common knowledge, player knowledge, GM secrets)
- Remote streamable HTTP with OAuth 2.1; nothing to install or self-host
- Free tier includes full MCP access, and Grimoire never charges for AI tokens

## Use cases of Grimoire

### The fast wins
Parse a paragraph of session notes into NPCs, factions, and threads. Dig up the location you mentioned three sessions ago. Update a faction's status straight from a recap.

### The deeper plays
Model an NPC reaction grounded in motivations they actually have and a political web they actually sit inside. Project the next move House Vale makes, if it moves at all. Generate three rumors that fit your city, not generic fantasy.

### The MC-Possibilities
Chain Grimoire MCP to other MCPs. Pipe your last session into a slideshow MCP and a voiceover MCP for an automated play-by-play. Hook a calendar MCP and the bot reminds you to prep.

The pitch is flexibility, not a baked-in workflow. MCP is a protocol; once Grimoire is connected to a frontier model, what you ask the model to do with it is the workflow.

## The toolbox

48 tools across seven groups.

**Entities** (10): `search_campaign`, `get_entity`, `list_entities`, `get_field_options`, `get_tag_options`, `get_entity_schema`, `create_entity`, `update_entity`, `delete_entity`, `batch_create_entities`. Covers all 14 entity types: NPCs, Locations, Factions, Quests, Items, Player Characters, Creatures, Vehicles, Lore Entries, World Rules, Planar Forces, Session Recaps, Session Preps, Custom Mechanics. `get_entity_schema` reports the exact fields a category accepts — native columns, select options, and the campaign's custom fields — so writes land in the right place instead of a catch-all bag.

**Relationships** (3): `add_relationship`, `get_relationships`, `delete_relationship`. Typed edges between any two entities.

**Knowledge graphs** (8): `get_constitution`, `get_entity_catalog`, `get_knowledge_graph`, `list_entity_graphs`, `get_entity_graph`, `add_to_entity_graph`, `create_entity_graph_edge`, `toggle_graph_attention`. Political, timeline, and geography projections with visibility filtering.

**World foundations** (6): `create_foundation_node`, `update_foundation_node`, `delete_foundation_node`, `create_foundation_edge`, `delete_foundation_edge`, `update_world_foundations`. Write tools for the campaign's constitution — the always-loaded foundation layer. Foundation nodes are world-defining concepts (a cataclysm, a hidden power, a cosmological law) connected by labeled edges, and `update_world_foundations` sets the structured World Foundations document (magic system, divine hierarchy, fundamental laws). These are the mechanical and cosmological truths of the world, distinct from database entities like NPCs and locations.

**Open threads** (7): `get_open_threads`, `create_open_thread`, `resolve_open_thread`, `unresolve_open_thread`, `update_open_thread`, `get_thread_progressions`, `add_thread_progression`. Loose ends and how they evolve session over session.

**Wiki** (9): `get_wiki_tree`, `get_wiki_page`, `create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`, `delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`. Block-based collaborative pages; the AI reads and edits them like you do.

**Campaign meta** (5): `current_campaign`, `get_campaign_context`, `get_narrative_state`, `get_campaign_bible`, `update_campaign_bible`. `get_narrative_state` aggregates recent sessions, open threads, canonical facts, and active arcs into one view.

## Visibility and safety

Everything respects Grimoire's visibility model: `common-knowledge`, `player-knowledge`, `dm-secret`. The AI sees what you've marked it can see and nothing more. Player accounts connecting through MCP get filtered views; GM accounts get the full picture.

OAuth scopes are campaign-scoped. Granting access to one campaign doesn't grant access to your other campaigns.

## Pricing

End-to-end free to try. Grimoire's free tier pairs with the free tiers on Claude.ai and ChatGPT, both of which support MCP connectors with no card required. Pro unlocks unlimited campaigns, custom fields on every category, and more storage.

Pricing: https://www.ttrpg.bot/#pricing

## FAQ from Grimoire

### Where is the source code for Grimoire?

The MCP server is operated as a hosted service and its source is not public. This repository is the public documentation and registry artifact: `server.json` here is what's published to `registry.modelcontextprotocol.io` as [`bot.ttrpg/grimoire`](https://registry.modelcontextprotocol.io/?q=grimoire).

### Does Grimoire include a standard MCP config?

Yes. Use the streamable HTTP block above, or skip config entirely: clients with custom connector support only need the URL `https://api.ttrpg.bot/mcp`.

### Is it free?

Yes, end to end. Grimoire's free tier covers full MCP access, and the free tiers of Claude.ai and ChatGPT each support a custom connector with no card required.

### Which clients work?

Claude (web and desktop), ChatGPT (plans with connector support), Cursor, and any MCP-compatible client that can talk to remote streamable HTTP servers.

### Does connecting share my campaign with an AI vendor?

You bring your own client. Grimoire exposes the one campaign you approve to that client and nothing else, the scopes are shown on the consent screen, and access is revocable at any time.

## Links

- Product: https://www.ttrpg.bot
- MCP overview and connect flow: https://www.ttrpg.bot/mcp
- Full setup guide: https://www.ttrpg.bot/docs/mcp/
- Changelog: https://www.ttrpg.bot/changelog
- Support: zmanlevelup@gmail.com

## License

MIT (for this documentation repository).
