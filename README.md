# Grimoire MCP: D&D and TTRPG Campaign Manager for Claude and ChatGPT

> The official hosted MCP server for Grimoire for Gamemasters. Your D&D campaign as a live, typed database your AI reads and writes, not markdown files you maintain.

| | |
|---|---|
| **Endpoint** | `https://api.ttrpg.bot/mcp` |
| **Transport** | Streamable HTTP (MCP spec 2025-06-18) |
| **Auth** | OAuth 2.1 with PKCE, one-click consent flow |
| **Tools** | 52 across 8 groups (21 read, 31 write) |
| **Registry** | `bot.ttrpg/grimoire` in the official MCP Registry |
| **Web app** | https://www.ttrpg.bot |

## What is Grimoire?

Grimoire for Gamemasters is a campaign manager for Dungeons and Dragons and other tabletop RPGs, built for the Dungeon Master or GM who runs the world. This repository documents its official MCP server, which lets Claude, ChatGPT, Cursor, or any MCP client read and write your campaign directly. A campaign in Grimoire is a typed database and knowledge graph, not a folder of notes: NPCs, factions, locations, quests, session recaps, and lore, connected by relationships, with a collaborative wiki whose pages link to those records. The AI answers from your current canon and stays consistent as the campaign changes between sessions.

The difference from a folder of campaign files is what the AI gets back. A file search returns prose that may be stale. Grimoire returns typed state: a faction's current leader, the open plot threads and how many sessions they have been hanging, the wiki pages that mention an NPC, the facts the table has marked canon. Writes go back into the same structure, so nothing you file from a chat drifts out of date.

It works for D&D 5e, Pathfinder 2e, Call of Cthulhu, Shadowdark, homebrew settings, and any other system: the records are the world, not a rules engine.

## How to use Grimoire?

Three clicks, nothing to install:

1. **Add the connector.** In Claude.ai (the free plan includes one custom connector), ChatGPT (Plus or Pro with Developer Mode), or any MCP client that supports custom connectors, open Settings, then Connectors, then Add custom connector, and paste `https://api.ttrpg.bot/mcp`.
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

- **Live typed state, not flat files.** 52 tools (21 read, 31 write) over a typed campaign database with 14 entity schemas. One narrative-state call returns recent sessions, open plot threads, active arcs, and canon facts.
- **Read and write, so state stays current.** Create and update NPCs, factions, locations, and quests, log thread progressions, and write wiki pages, blocks, and tables from the conversation, so canon updates as you prep instead of going stale in a file.
- **The wiki is part of the world.** Full-text wiki search, pages with breadcrumbs and children, tables that round-trip as markdown grids, and entity and page links that resolve to real records instead of string matches.
- **Knowledge-graph queries** across political, geographic, and timeline projections, with GM visibility filtering applied server-side.
- **GM secrets enforced server-side** through three visibility tiers (common knowledge, player knowledge, GM secret), so what the AI can surface respects what players are allowed to see.
- **Player connections too.** When the GM turns on player AI for a campaign, a player can connect their own AI and read exactly what the player portal shows them, nothing the GM has kept secret.
- **Hosted, one OAuth click, no token markup.** Remote streamable HTTP with OAuth 2.1; nothing to install or self-host. The free tier includes MCP access, and Grimoire never charges for AI tokens.

## Use cases of Grimoire

- Ask what changed since session 12 before you write session 13, answered from tracked state rather than your memory of the notes.
- Mid-session lookups that resolve from campaign canon instead of the model's best guess.
- Parse tonight's prep notes into typed NPCs, factions, and plot threads, instead of pasting the same lore into every new chat.
- Query the world by relationship: which factions have a stake in this city, and which NPCs would notice the party.
- Search the wiki for the handout that mentioned the drowned abbey, then read the page and the entities it links to.
- Draft a session recap and file it as a structured entity your players can read in the portal, scoped to their visibility tier.
- Chain it with other MCP servers: pipe the last session into a slideshow or voiceover server, or hook a calendar server so the bot reminds you to prep.

## The toolbox

52 tools across eight groups.

**Entities** (10): `search_campaign`, `get_entity`, `list_entities`, `get_field_options`, `get_tag_options`, `get_entity_schema`, `create_entity`, `update_entity`, `delete_entity`, `batch_create_entities`. Covers all 14 entity types: NPCs, Locations, Factions, Quests, Items, Player Characters, Creatures, Vehicles, Lore Entries, World Rules, Planar Forces, Session Recaps, Session Preps, Custom Mechanics. `get_entity_schema` reports the exact fields a category accepts, native columns, select options, and the campaign's custom fields, so writes land in the right place instead of a catch-all bag.

**Relationships** (3): `add_relationship`, `get_relationships`, `delete_relationship`. Typed edges between any two entities.

**Knowledge graphs** (8): `get_constitution`, `get_entity_catalog`, `get_knowledge_graph`, `list_entity_graphs`, `get_entity_graph`, `add_to_entity_graph`, `create_entity_graph_edge`, `toggle_graph_attention`. Political, timeline, geography, and foundations projections with visibility filtering. `get_constitution` is the Layer 1 read for both roles: a GM connection gets the full constitutional layer, a player connection gets the parts the GM chose to share, each flagged in the response.

**World foundations** (6): `create_foundation_node`, `update_foundation_node`, `delete_foundation_node`, `create_foundation_edge`, `delete_foundation_edge`, `update_world_foundations`. Write tools for the campaign's foundation layer, which is always loaded. Foundation nodes are world-defining concepts (a cataclysm, a hidden power, a cosmological law) connected by labeled edges, and `update_world_foundations` sets the structured World Foundations document (magic system, divine hierarchy, fundamental laws). These are the mechanical and cosmological truths of the world, distinct from database entities like NPCs and locations.

**Open threads** (7): `get_open_threads`, `create_open_thread`, `resolve_open_thread`, `unresolve_open_thread`, `update_open_thread`, `get_thread_progressions`, `add_thread_progression`. Loose ends and how they evolve session over session.

**Wiki** (10): `get_wiki_tree`, `get_wiki_page`, `search_wiki`, `create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`, `delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`. Block-based collaborative pages the AI reads and edits like you do. `search_wiki` is full-text over page titles and block prose (`search_campaign` covers entities only). Tables read and write as plain markdown grids, and entity links (`@[Name](entity://category/uuid)`) and page links (`@[Title](page://uuid)`) round-trip inside any block, including table cells, so a page the AI writes is linked to your records the same way one you typed would be.

**Campaign meta** (5): `current_campaign`, `get_campaign_context`, `get_narrative_state`, `get_campaign_bible`, `update_campaign_bible`. `get_narrative_state` aggregates recent sessions, open threads, canonical facts, and active arcs into one view. `get_campaign_bible` returns every block to a GM connection, and only the Common Knowledge blocks to a player connection, and then only when the GM has turned Bible sharing on.

**Session observations** (3): `create_session_observation`, `update_session_observation`, `delete_session_observation`. GM-only notes that colour a session: a typed, rated moment worth coming back to. Reading them needs no tool of its own, because `get_narrative_state` returns the observations on the sessions it summarizes and `get_knowledge_graph` returns the ones on the timeline it draws.

The server also exposes one MCP resource, `campaign://wiki`, which lists the campaign's wiki pages and their hierarchy.

## A skill for driving it well

`skills/grimoire-campaign-context/` is a Claude Skill (also usable as plain instructions in any client) that teaches an assistant the session-start centering sequence, which read tool to reach for, how to write typed entities without dumping data into the wrong fields, and recipes for parsing prep notes, filing recaps, and pre-session pulls. Point your client at it, or paste it into a project's instructions.

## Visibility and safety

Everything respects Grimoire's visibility model: `common-knowledge`, `player-knowledge`, `dm-secret`, and (wiki only) `system`, with `inherit` resolving through the block that links a wiki page and then its parent chain. The AI sees what the connected role is allowed to see and nothing more. Player connections get filtered views, and a page a player cannot reach does not exist for that connection: search, the tree, and page reads all answer not found. GM connections get the full picture.

The Campaign Bible is GM Secret by default. Mark a block `common-knowledge` when players may know it (world mechanics, everyday life); that records it as player-safe. A GM can then turn Bible sharing on, and a player connection reads those Common Knowledge blocks through `get_campaign_bible` and `get_constitution`, and nothing else. With sharing off, `get_campaign_bible` answers the normal result `{"shared": false, "message": "..."}` rather than an error, and `get_constitution` reports `shared.bible` as false with the Bible section absent, so an assistant tells the player plainly instead of retrying. The Bible cannot be renamed.

World Foundations share the same way and behind their own flag: `get_knowledge_graph(graph_type='foundations')` and the `foundationNodes` of `get_constitution` reach a player connection only when the GM turned Foundations sharing on, and then only as the player-safe projection, with GM Secret nodes and every edge touching them dropped. Both flags sit behind the campaign's portal policy: a wiki-only portal refuses these reads outright, whatever the flags say.

OAuth scopes are campaign-scoped. Granting access to one campaign does not grant access to your other campaigns, and access is revocable at any time.

## Pricing

Free to try end to end. Grimoire's free tier includes MCP access for one campaign with up to 100 entities, and it pairs with Claude's free plan, which supports one custom connector with no card required. ChatGPT custom connectors need Plus or Pro with Developer Mode. Pro unlocks unlimited campaigns and entities, custom fields on every category, and more storage. Grimoire never charges for AI tokens: you bring your own client.

Pricing: https://www.ttrpg.bot/#pricing

## FAQ from Grimoire

### Is there an MCP server that lets Claude or ChatGPT read and edit my D&D campaign?

Yes. Grimoire's hosted MCP server at `https://api.ttrpg.bot/mcp` gives Claude, ChatGPT, Cursor, and other MCP clients read and write access to a Grimoire campaign: NPCs, factions, locations, quests, session recaps, lore, open plot threads, and the campaign wiki. You connect with one OAuth click and choose which campaign the AI can see.

### Where is the source code for Grimoire?

The MCP server is operated as a hosted service and its source is not public. This repository is the public documentation and registry artifact: `server.json` here is what's published to `registry.modelcontextprotocol.io` as [`bot.ttrpg/grimoire`](https://registry.modelcontextprotocol.io/?q=grimoire), `lhm.plugin.json` is the LobeHub Marketplace manifest, and `skills/` holds the Claude Skill.

### Does Grimoire include a standard MCP config?

Yes. Use the streamable HTTP block above, or skip config entirely: clients with custom connector support only need the URL `https://api.ttrpg.bot/mcp`.

### Is it free?

Yes, to try. Grimoire's free tier includes MCP access for one campaign with up to 100 entities, and Claude's free plan supports one custom connector with no card required. ChatGPT needs Plus or Pro with Developer Mode for custom connectors.

### Which clients work?

Claude (web and desktop), ChatGPT (Plus or Pro with Developer Mode), Cursor, and any MCP-compatible client that can talk to remote streamable HTTP servers.

### Does it work for Pathfinder, Call of Cthulhu, or my homebrew system?

Yes. Grimoire is system-agnostic. Campaigns can be set up for D&D 5e, Pathfinder 2e, Call of Cthulhu, Shadowdark, and many other systems or a homebrew setting, and the genre changes the labels and defaults, not what the AI can do. The tools read and write the world (people, places, factions, history, threads), not character-sheet rules.

### Can it read and write my wiki?

Yes. `search_wiki` finds pages and blocks by text, `get_wiki_page` returns a page with its blocks, breadcrumb, and children, and the write tools create pages, blocks, and tables. Tables travel as markdown grids, and entity and page links are preserved as tokens the server resolves to real records.

### How is this different from keeping campaign notes in markdown files?

Files hold prose; the AI has to guess which file is current and re-read it every session. Grimoire holds typed state: records with fields and relationships, threads with ages, facts marked canon, pages linked to records. The AI queries what is true now and writes changes back into the same structure, so nothing drifts between sessions.

### How does it compare to World Anvil, Obsidian, or Notion for worldbuilding?

Those tools store pages, and an AI connected to them reads pages. Grimoire stores typed records and the relationships between them, so an AI can traverse the campaign: from a city to the factions with a stake in it, to the NPCs who lead them, to the open threads they touch, with GM secrets filtered on the server. If you already have notes elsewhere, the AI can read them in a chat and file them into Grimoire as typed entities.

### Does connecting share my campaign with an AI vendor?

You bring your own client. Grimoire exposes the one campaign you approve to that client and nothing else, the scopes are shown on the consent screen, and access is revocable at any time.

## Links

- Product: https://www.ttrpg.bot
- MCP overview and connect flow: https://www.ttrpg.bot/mcp
- Full setup guide: https://www.ttrpg.bot/docs/mcp/
- Campaign knowledge graph: https://www.ttrpg.bot/campaign-knowledge-graph/
- Campaign wiki: https://www.ttrpg.bot/dnd-campaign-wiki/
- Changelog: https://www.ttrpg.bot/changelog
- Support: zmanlevelup@gmail.com

## License

MIT (for this documentation repository).
