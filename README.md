# Grimoire MCP

> Your campaign as a live database your AI reads and writes, not markdown files you maintain.

| | |
|---|---|
| **Endpoint** | `https://api.ttrpg.bot/mcp` |
| **Transport** | Streamable HTTP (MCP spec 2025-06-18) |
| **Auth** | OAuth 2.1 with PKCE, one-click consent flow |
| **Tools** | 49 across 7 groups (21 read, 28 write) |
| **Registry** | `bot.ttrpg/grimoire` in the official MCP Registry |
| **Web app** | https://www.ttrpg.bot |

## What is Grimoire?

Grimoire is a campaign manager for D&D and tabletop RPG gamemasters that gives your AI live campaign state instead of flat markdown files to keep in sync. A campaign in Grimoire is a typed database, not a folder of notes: NPCs, factions, locations, quests, sessions, and lore, connected by relationships and knowledge graphs, with a collaborative wiki whose pages link to those records. The hosted MCP server lets any connected client read and write that database directly, so the AI answers from your current canon and stays consistent as the campaign changes between sessions.

The difference from a folder of campaign files is what the AI gets back. A file search returns prose that may be stale. Grimoire returns typed state: a faction's current leader, the open plot threads and how many sessions they have been hanging, the wiki pages that mention an NPC, the facts the table has marked canon. Writes go back into the same structure, so nothing you file from a chat drifts out of date.

## How to use Grimoire?

Three clicks, nothing to install:

1. **Add the connector.** In Claude.ai or ChatGPT (or any MCP client that supports custom connectors), open Settings, then Connectors, then Add custom connector, and paste `https://api.ttrpg.bot/mcp`.
2. **Sign in and pick a campaign.** Your client opens the Grimoire consent screen. Sign in with the account that owns the campaign, pick which campaign the AI should see, click Approve. Access is scoped to that one campaign, enforced server-side.
3. **Ask it something about your world.** *"Who knows about the cult?" "What was the queen's reaction in session 4?" "What changed since session 12?"* The client reads live and answers from your data.

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

- **Live typed state, not flat files.** 49 tools (21 read, 28 write) over a typed campaign database with 14 entity schemas. One narrative-state call returns recent sessions, open plot threads, active arcs, and canon facts.
- **Read and write, so state stays current.** Create and update entities, log thread progressions, and write wiki pages, blocks, and tables from the conversation, so canon updates as you prep instead of going stale in a file.
- **The wiki is part of the world.** Full-text wiki search, pages with breadcrumbs and children, tables that round-trip as markdown grids, and entity and page links that resolve to real records instead of string matches.
- **Knowledge-graph queries** across political, geographic, and timeline projections, with GM visibility filtering applied server-side.
- **GM secrets enforced server-side** through three visibility tiers (common knowledge, player knowledge, GM secret), so what the AI can surface respects what players are allowed to see.
- **Hosted, one OAuth click, no token markup.** Remote streamable HTTP with OAuth 2.1; nothing to install or self-host. The free tier includes full MCP access, and Grimoire never charges for AI tokens.

## Use cases of Grimoire

- Ask what changed since session 12 before you write session 13, answered from tracked state rather than your memory of the notes.
- Mid-session lookups that resolve from campaign canon instead of the model's best guess.
- Parse tonight's prep notes into typed NPCs, factions, and plot threads, instead of pasting the same lore into every new chat.
- Query the world by relationship: which factions have a stake in this city, and which NPCs would notice the party.
- Search the wiki for the handout that mentioned the drowned abbey, then read the page and the entities it links to.
- Draft a session recap and file it as a structured entity your players can read in the portal, scoped to their visibility tier.
- Chain it with other MCP servers: pipe the last session into a slideshow or voiceover server, or hook a calendar server so the bot reminds you to prep.

## The toolbox

49 tools across seven groups.

**Entities** (10): `search_campaign`, `get_entity`, `list_entities`, `get_field_options`, `get_tag_options`, `get_entity_schema`, `create_entity`, `update_entity`, `delete_entity`, `batch_create_entities`. Covers all 14 entity types: NPCs, Locations, Factions, Quests, Items, Player Characters, Creatures, Vehicles, Lore Entries, World Rules, Planar Forces, Session Recaps, Session Preps, Custom Mechanics. `get_entity_schema` reports the exact fields a category accepts, native columns, select options, and the campaign's custom fields, so writes land in the right place instead of a catch-all bag.

**Relationships** (3): `add_relationship`, `get_relationships`, `delete_relationship`. Typed edges between any two entities.

**Knowledge graphs** (8): `get_constitution`, `get_entity_catalog`, `get_knowledge_graph`, `list_entity_graphs`, `get_entity_graph`, `add_to_entity_graph`, `create_entity_graph_edge`, `toggle_graph_attention`. Political, timeline, and geography projections with visibility filtering.

**World foundations** (6): `create_foundation_node`, `update_foundation_node`, `delete_foundation_node`, `create_foundation_edge`, `delete_foundation_edge`, `update_world_foundations`. Write tools for the campaign's foundation layer, which is always loaded. Foundation nodes are world-defining concepts (a cataclysm, a hidden power, a cosmological law) connected by labeled edges, and `update_world_foundations` sets the structured World Foundations document (magic system, divine hierarchy, fundamental laws). These are the mechanical and cosmological truths of the world, distinct from database entities like NPCs and locations.

**Open threads** (7): `get_open_threads`, `create_open_thread`, `resolve_open_thread`, `unresolve_open_thread`, `update_open_thread`, `get_thread_progressions`, `add_thread_progression`. Loose ends and how they evolve session over session.

**Wiki** (10): `get_wiki_tree`, `get_wiki_page`, `search_wiki`, `create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`, `delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`. Block-based collaborative pages the AI reads and edits like you do. `search_wiki` is full-text over page titles and block prose (`search_campaign` covers entities only). Tables read and write as plain markdown grids, and entity links (`@[Name](entity://category/uuid)`) and page links (`@[Title](page://uuid)`) round-trip inside any block, including table cells, so a page the AI writes is linked to your records the same way one you typed would be.

**Campaign meta** (5): `current_campaign`, `get_campaign_context`, `get_narrative_state`, `get_campaign_bible`, `update_campaign_bible`. `get_narrative_state` aggregates recent sessions, open threads, canonical facts, and active arcs into one view.

The server also exposes one MCP resource, `campaign://wiki`, which lists the campaign's wiki pages and their hierarchy.

## A skill for driving it well

`skills/grimoire-campaign-context/` is a Claude Skill (also usable as plain instructions in any client) that teaches an assistant the session-start centering sequence, which read tool to reach for, how to write typed entities without dumping data into the wrong fields, and recipes for parsing prep notes, filing recaps, and pre-session pulls. Point your client at it, or paste it into a project's instructions.

## Visibility and safety

Everything respects Grimoire's visibility model: `common-knowledge`, `player-knowledge`, and `dm-secret`, with `inherit` resolving through the block that links a wiki page and then its parent chain. The AI sees what the connected role is allowed to see and nothing more. Player connections get filtered views, and a page a player cannot reach does not exist for that connection: search, the tree, and page reads all answer not found. GM connections get the full picture.

OAuth scopes are campaign-scoped. Granting access to one campaign does not grant access to your other campaigns, and access is revocable at any time.

## Pricing

End-to-end free to try. Grimoire's free tier includes full MCP access and pairs with the free tiers on Claude.ai and ChatGPT, both of which support MCP connectors with no card required. Pro unlocks unlimited campaigns and entities, custom fields on every category, and more storage.

Pricing: https://www.ttrpg.bot/#pricing

## FAQ from Grimoire

### Where is the source code for Grimoire?

The MCP server is operated as a hosted service and its source is not public. This repository is the public documentation and registry artifact: `server.json` here is what's published to `registry.modelcontextprotocol.io` as [`bot.ttrpg/grimoire`](https://registry.modelcontextprotocol.io/?q=grimoire), `lhm.plugin.json` is the LobeHub Marketplace manifest, and `skills/` holds the Claude Skill.

### Does Grimoire include a standard MCP config?

Yes. Use the streamable HTTP block above, or skip config entirely: clients with custom connector support only need the URL `https://api.ttrpg.bot/mcp`.

### Is it free?

Yes, end to end. Grimoire's free tier covers full MCP access, and the free tiers of Claude.ai and ChatGPT each support a custom connector with no card required.

### Which clients work?

Claude (web and desktop), ChatGPT (plans with connector support), Cursor, and any MCP-compatible client that can talk to remote streamable HTTP servers.

### Can it read and write my wiki?

Yes. `search_wiki` finds pages and blocks by text, `get_wiki_page` returns a page with its blocks, breadcrumb, and children, and the write tools create pages, blocks, and tables. Tables travel as markdown grids, and entity and page links are preserved as tokens the server resolves to real records.

### How is this different from keeping campaign notes in markdown files?

Files hold prose; the AI has to guess which file is current and re-read it every session. Grimoire holds typed state: records with fields and relationships, threads with ages, facts marked canon, pages linked to records. The AI queries what is true now and writes changes back into the same structure, so nothing drifts between sessions.

### Does connecting share my campaign with an AI vendor?

You bring your own client. Grimoire exposes the one campaign you approve to that client and nothing else, the scopes are shown on the consent screen, and access is revocable at any time.

## Links

- Product: https://www.ttrpg.bot
- MCP overview and connect flow: https://www.ttrpg.bot/mcp
- Full setup guide: https://www.ttrpg.bot/docs/mcp/
- Campaign wiki: https://www.ttrpg.bot/dnd-campaign-wiki/
- Changelog: https://www.ttrpg.bot/changelog
- Support: zmanlevelup@gmail.com

## License

MIT (for this documentation repository).
