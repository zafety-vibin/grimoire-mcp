---
name: grimoire-campaign-context
description: >
  Drive the Grimoire MCP server correctly for tabletop RPG campaign work. Use this
  skill whenever a Grimoire connector is available and the user asks anything about
  their campaign, world, NPCs, factions, locations, quests, sessions, plot threads,
  or wiki. Covers the session-start centering sequence (current_campaign,
  get_constitution, get_campaign_context, get_narrative_state) that must run before
  answering campaign questions, which read tool to reach for, how to write typed
  entities without dumping data into the wrong fields, how the three-tier visibility
  model constrains what you can see, and recipes for parsing prep notes, filing
  session recaps, relationship spelunking, and pre-session prep pulls.
license: MIT
compatibility: Any MCP client that supports remote streamable HTTP connectors, including Claude, ChatGPT, Cursor, and Claude Code.
metadata:
  author: zafety-vibin
  version: "1.0.0"
---

# Working with Grimoire

Grimoire is a campaign manager for tabletop RPG gamemasters, and its MCP server
gives you live read and write access to one campaign's database. That database is
not a pile of notes. It is typed: 14 entity categories, typed relationships between
them, three knowledge graph projections, open narrative threads, and a block-based
wiki. The endpoint is `https://api.ttrpg.bot/mcp`, transport is streamable HTTP,
auth is OAuth 2.1, and access is scoped to the single campaign the user approved.

Your job is to answer from that campaign's canon, not from generic tabletop RPG
knowledge. The homebrew is the authority. If the campaign says elves are extinct
and magic costs blood, that is true, and your training data is wrong.

Grimoire exposes 48 tools in 7 groups. The tools carry their own layer numbering in
their descriptions. Follow it.

## Center yourself first

Do not answer a campaign question, and do not write anything, until you have
grounded yourself in this campaign. The whole point of the layered design is that
context gets loaded deliberately instead of dumped all at once.

**Every session, before the first substantive answer:**

1. `current_campaign` (free, instant). Returns the campaign, user, and **role** bound
   to this MCP session. Call it first, every time. Everything after this depends on
   it, and the role decides which of the next two steps you are even allowed to
   make. See "If you are connected as a player" below.
2. `get_constitution` (Layer 1). GM role only. The world's non-negotiable truths:
   campaign summary, Campaign Bible summary, and the world foundations graph. This
   is where the setting's defining events, cosmology, and the GM's content
   boundaries live. It colors every downstream answer.
3. `get_campaign_context`. Genre, ruleset, setting technology and magic levels,
   which entity categories this campaign uses, and what this campaign calls them.
   Call it early. It is how you stop saying "NPC" when this table says "contact,"
   and how you stop applying D&D 5e assumptions to a Blades-adjacent horror game.

**Then, if the question touches the present tense of the story:**

4. `get_narrative_state` (Layer 1.5). Recent sessions with summaries and key events,
   open threads split into major and minor, canonical facts, active arcs, and recent
   observations. Call it for anything about "now": what the party is doing, what
   just happened, how an NPC would react today, prep, recaps, consequences. Skip it
   for a pure static lookup like "what is the population of Vethmoor."

**Then, only if you still need more:**

5. `get_entity_catalog` (Layer 2). Every active entity as id plus name, grouped by
   type, plus compact node and edge lists for the political, geography, timeline,
   and foundations graphs. Cheap. This is the map. Call it when you need to know
   what exists, and always call it before creating anything so you do not make a
   duplicate.
6. `get_knowledge_graph(graph_type, attention=true)` (Layer 3). The GM-curated
   subset: only entities the GM has flagged as relevant right now. This is the right
   first graph call on a large campaign.
7. `get_knowledge_graph(graph_type)` (Layer 4). The full projection. Use when the
   attention-filtered view is genuinely not enough.
8. `get_entity(...)` (Layer 5). Full detail on one entity, by id. Call it as many
   times as the question warrants, on ids you already hold.

If you skip step 2 and step 3, you will produce answers that sound like a tabletop
RPG and read like someone else's campaign. That is the failure mode this system
exists to prevent.

## If you are connected as a player

`tools/list` is the same 48 tools for everyone. Role is enforced when you **call**
a tool, not when you list them. So do not assume a tool works just because you can
see it. Check `current_campaign` first and branch.

**Permanently unavailable to a player-role connection** (returns "not available to
player-role MCP connections"): `get_constitution`, `get_campaign_bible`,
`get_relationships`, `get_open_threads`, `get_thread_progressions`,
`list_entity_graphs`, `get_entity_graph`. Do not retry these, and do not treat the
error as a transient failure.

**Gated by the campaign's player portal policy**: `get_campaign_context`,
`get_narrative_state`, `get_entity_catalog`, `get_knowledge_graph`. Under a
`wiki-only` portal they refuse outright. Under `revealed` or `open` they work and
return visibility-filtered results. The refusal message names the policy, so read
it rather than guessing.

**Always available to a player**: `current_campaign`, `search_campaign`,
`get_wiki_tree`, `get_wiki_page`.

So a player-role centering sequence is: `current_campaign`, then
`get_campaign_context` (if the portal allows it), then `get_wiki_tree` and
`search_campaign` for everything else. You will not get the constitution. Say so
plainly if the user asks for world truths you cannot reach, and point them at their
GM. Every write tool is GM-only as well.

## Choosing a read tool

| You want | Call |
|---|---|
| Which campaign and role am I in | `current_campaign` |
| World truths, tone, content boundaries | `get_constitution`, then `get_campaign_bible` for full blocks |
| Genre, ruleset, category vocabulary | `get_campaign_context` |
| Where the story is right now | `get_narrative_state` |
| Just the loose ends | `get_open_threads` |
| How a thread has moved over time | `get_thread_progressions` |
| A cheap map of everything | `get_entity_catalog` |
| Find something by name or keyword, category unknown | `search_campaign` |
| Everything in one category, paginated | `list_entities` |
| How things connect | `get_knowledge_graph` (political, timeline, geography) |
| One entity's edges | `get_relationships` |
| Full detail on a known entity | `get_entity` |
| Valid field values and tags before writing | `get_entity_schema`, `get_field_options`, `get_tag_options` |
| Wiki structure, then a page | `get_wiki_tree`, then `get_wiki_page` |

Etiquette that saves tokens and mistakes:

- `search_campaign` is full-text over name, description, and custom fields. It
  returns truncated descriptions and a relevance score. Use it to find ids, then
  deep-dive with `get_entity`. Filter with `categories` when you know the type.
- `list_entities` is for enumeration, not lookup. Do not page through 200 NPCs
  hunting for a name; search instead.
- Reach for `get_knowledge_graph` when the question is about relationships rather
  than attributes. "Who does House Vale hate" is a graph question. "How old is Lady
  Vale" is an entity question.
- `get_entity_catalog` before any create. Duplicates are the most common damage an
  AI does to a campaign database.

## Writing without making a mess

Grimoire's categories have real columns. Writing blind means dumping structured
data into a generic blob where nothing can query it. Do not do that.

**The rule: `get_entity_schema(category)` before `create_entity` or
`update_entity`.** It returns:

- `globalFields`: name, description, status, player_knowledge, tags, set at the top
  level of the call.
- `defaultFields`: the category's native fields. These go **inside** `custom_fields`
  and are written to real columns. Examples: `npcs.faction_id`,
  `locations.location_type`, `factions.allied_faction_ids`, `quests.quest_giver_id`.
- `selectOptions`: suggested values for select-type fields. These are hints, not a
  closed enum. A select-type default field also accepts a brand-new free-text value,
  and the new value becomes a campaign option. Do not avoid the right word just
  because it is not listed.
- `customFields`: the GM's own defined fields for this campaign.

Any key in `custom_fields` that is neither a default field nor a defined custom
field is rejected outright. Relation fields take arrays of entity UUIDs and sync
junction tables, so you need real ids, which means catalog or search first.

Other write rules:

- **Tags.** Call `get_tag_options` first. The standard set is political, personal,
  economic, religious, magical, historical, secret, prophetic, combat, social,
  environmental, technological.
- **Bulk.** `batch_create_entities` takes up to 10 entities of one category per
  call. It is **not** transactional: entities insert sequentially, so a failure at
  #3 leaves #1 and #2 committed. Always read the `failed` array. To undo a partial
  batch, `delete_entity` the ids from `created`.
- **Relationships.** `add_relationship` routes automatically: it writes an FK field
  when one exists (`faction_id`, `parent_location_id`, `owner_npc_id`) and a
  junction row otherwise, using `relationship_type` as the edge label. Database
  relationships auto-visualize in the knowledge graphs. Do **not** also call
  `add_to_entity_graph` or `create_entity_graph_edge` for the same link; those two
  exist only for custom edges with no database equivalent, and using them here
  creates duplicate edges.
- **Attention.** `toggle_graph_attention` promotes an entity from Layer 4 into the
  Layer 3 session-context view. Use it when prep identifies what matters next
  session. Do not mass-flag; the value of Layer 3 is that it is small.
- **Wiki.** `create_wiki_block` has a known positioning race. For more than one
  block use `batch_create_wiki_blocks`, then `get_wiki_page` to verify order, then
  `batch_reorder_wiki_blocks` to fix it. Block content is a typed shape per block
  type (text, heading, bullet, numbered, quote, callout, divider, page, image); read
  the tool description before composing. When you read a Bible block and want to
  rewrite it, pass the simple `{"content":[{"text":"..."}]}` shape back and let the
  server normalize. Do not hand-build TipTap documents.
- **Threads.** After a session, log progressions rather than rewriting thread text.
  `add_thread_progression` takes `thread_id`, `session_id`, `progression_type`
  (start, update, complication, resolution), and a `title`. Link it to a session key
  event with `key_event_id` when you have one; `key_event_index` is legacy and
  breaks if key events are reordered. GM-role connections only.

## Visibility is not yours to negotiate

Every entity, wiki block, and graph node carries one of three visibility values:

- `common-knowledge`: everyone at the table knows it.
- `player-knowledge`: specific players know it.
- `dm-secret`: GM only. This is the "GM secret" tier.

This is enforced in the backend query layer, not in the response text. A player-role
connection never receives GM-secret rows in the first place. There is no flag,
phrasing, or tool call that gets around it, and trying is a bug in your reasoning,
not a puzzle to solve.

What this means for you:

- If a player-role user asks about something and you find nothing, say you cannot
  see it in this campaign's data from this connection. Do not fill the gap with
  invention and do not speculate about what is being hidden.
- When you write, set `player_knowledge` deliberately. Anything the GM told you in
  confidence, anything that would spoil a reveal, anything a villain knows and the
  party does not: `dm-secret`. Public reputation, geography, common rumor:
  `common-knowledge`. When unsure on a secret, choose `dm-secret`; over-hiding is
  recoverable, leaking is not.
- Wiki blocks additionally support `inherit`, which takes the page's visibility.
  That is the sane default for ordinary body content.
- Category-specific GM-only fields exist and are named for it: `npcs.dm_secrets`,
  `factions.dm_true_agenda`, `quests.dm_true_objective`,
  `items.dm_secret_properties`, `session_recaps.dm_consequences`. Put the twist
  there, not in the public description.

## Recipes

### Parse prep notes into entities

1. `get_constitution`, `get_campaign_context` if not already loaded this session.
2. `get_entity_catalog`. Match every proper noun in the notes against it.
3. For each name that already exists, `get_entity`, then `update_entity` with only
   the new information. Do not recreate.
4. For each genuinely new thing, `get_entity_schema(category)`, then
   `batch_create_entities` in groups of 10 per category.
5. `add_relationship` for every link the notes imply. Do not add graph edges by
   hand afterwards.
6. Report back: what you created, what you updated, what you skipped as a
   duplicate, and any contradiction you hit between the notes and existing canon.
   Contradictions get surfaced, not silently resolved.

### File a session recap

1. `get_narrative_state` for what was already open going in.
2. `get_entity_schema("session_recaps")`, then `create_entity` with
   `session_number`, `summary`, `key_events`, `pcs_present_ids`, and
   `dm_consequences` for anything the party has not realized yet.
3. `get_open_threads`. For each thread the session touched,
   `add_thread_progression` against the new session id with the right
   `progression_type`, linked to a key event where one applies.
4. `resolve_open_thread` for anything actually closed. `create_open_thread` for new
   loose ends, with the type that fits: consequence, promise, mystery,
   foreshadowing, callback_opportunity.
5. `update_entity` on NPCs, factions, and locations the session changed.

### Who knows what about X

1. `search_campaign` for X, get the id and category.
2. `get_relationships` on it for direct edges.
3. `get_knowledge_graph("political")` to see it inside the wider web of factions and
   NPCs rather than as an isolated node.
4. `get_entity` on the neighbours that matter.
5. Answer with the visibility tier attached. "The city guard treats this as common
   knowledge; the Duke's involvement is a GM secret" is a more useful answer than a
   flat list.

### Prep pull before a session

1. `get_narrative_state` for open threads and recent arcs.
2. `get_knowledge_graph("timeline", attention=true)` for what the GM has already
   flagged as live.
3. `get_knowledge_graph("geography")` if the party is travelling, to get the
   location hierarchy right.
4. `get_entity` on the handful of NPCs and locations that will actually appear.
5. Offer the GM a short prep brief, then use `toggle_graph_attention` to flag the
   entities you both agreed matter, so the next session's Layer 3 view is already
   correct.

## Anti-patterns

- Answering a campaign question before `get_constitution` and
  `get_campaign_context`. This is the single biggest failure mode.
- Filling a gap with generic fantasy or sci-fi content instead of saying the data
  does not cover it, or asking.
- Creating an entity without checking `get_entity_catalog` for a near-duplicate.
- Writing structured data into `custom_fields` keys that are not real fields,
  instead of calling `get_entity_schema` first.
- Calling `get_knowledge_graph` with `attention=false` as the opening move on a
  large campaign, then drowning in it.
- Adding a graph edge with `create_entity_graph_edge` for a relationship that
  `add_relationship` already visualized.
- Treating a visibility filter as an obstacle.
- Rewriting a thread's text instead of logging a progression, which destroys the
  session-over-session history the timeline graph is built on.

## Reference

- `references/tool-groups.md`: all 48 tools by group, with what each is for.
- `references/context-architecture.md`: why the layers exist and how to pick a
  traversal.
- Setup and client instructions: https://www.ttrpg.bot/docs/mcp/
- Architecture write-up: https://www.ttrpg.bot/blog/grimoire-mcp-secret-sauce/
