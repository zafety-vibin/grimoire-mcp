---
name: grimoire-campaign-context
description: >
  Drive the Grimoire MCP server correctly for tabletop RPG campaign work. Use this
  skill whenever a Grimoire connector is available and the user asks anything about
  their campaign, world, NPCs, factions, locations, quests, sessions, plot threads,
  or wiki. Covers the session-start centering sequence (current_campaign,
  get_constitution, get_campaign_context, get_narrative_state) that must run before
  answering campaign questions, which read tool to reach for, how to write typed
  entities without dumping data into the wrong fields, how to search, read, and write
  the wiki without breaking its links or tables, how the visibility model constrains
  what you can see, and recipes for parsing prep notes, filing session recaps,
  relationship spelunking, wiki lookups, and pre-session prep pulls.
license: MIT
compatibility: Any MCP client that supports remote streamable HTTP connectors, including Claude, ChatGPT, Cursor, and Claude Code.
metadata:
  author: zafety-vibin
  version: "1.1.0"
---

# Working with Grimoire

Grimoire is a campaign manager for tabletop RPG gamemasters, and its MCP server
gives you live read and write access to one campaign's database. That database is
not a pile of notes. It is typed: 14 entity categories, typed relationships between
them, three knowledge graph projections, open narrative threads, and a block-based
wiki whose pages link to the records. The endpoint is `https://api.ttrpg.bot/mcp`,
transport is streamable HTTP, auth is OAuth 2.1, and access is scoped to the single
campaign the user approved.

Your job is to answer from that campaign's canon, not from generic tabletop RPG
knowledge. The homebrew is the authority. If the campaign says elves are extinct
and magic costs blood, that is true, and your training data is wrong.

Grimoire exposes 49 tools in 7 groups. The tools carry their own layer numbering in
their descriptions. Follow it. Every tool description is current; when this skill
and a tool description disagree, the tool description wins.

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

The wiki sits beside the layers rather than inside them. Reach for it when the
question is about prose the GM wrote (a handout, a house rule, a lore page, session
notes) rather than about a record's fields. See "The wiki" below.

If you skip step 2 and step 3, you will produce answers that sound like a tabletop
RPG and read like someone else's campaign. That is the failure mode this system
exists to prevent.

## If you are connected as a player

`tools/list` is the same 49 tools for everyone. Role is enforced when you **call**
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
`search_wiki`, `get_wiki_tree`, `get_wiki_page`. All of them return only what the
player's visibility admits, and a page the player cannot reach answers "page not
found" rather than "forbidden", so a not-found is not evidence that something is
being hidden.

So a player-role centering sequence is: `current_campaign`, then
`get_campaign_context` (if the portal allows it), then `get_wiki_tree`,
`search_wiki`, and `search_campaign` for everything else. You will not get the
constitution. Say so plainly if the user asks for world truths you cannot reach, and
point them at their GM. Every write tool is GM-only as well.

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
| Find a record by name or keyword, category unknown | `search_campaign` |
| Find prose by the words in it (a handout, a rule, a note) | `search_wiki`, then `get_wiki_page` |
| Everything in one category, paginated | `list_entities` |
| How things connect | `get_knowledge_graph` (political, timeline, geography) |
| One entity's edges | `get_relationships` |
| Full detail on a known entity | `get_entity` |
| Valid field values and tags before writing | `get_entity_schema`, `get_field_options`, `get_tag_options` |
| Wiki structure, then a page | `get_wiki_tree`, then `get_wiki_page` |

Etiquette that saves tokens and mistakes:

- `search_campaign` is full-text over entity name, description, and custom fields.
  It never touches wiki rows. It returns truncated descriptions and a relevance
  score. Use it to find ids, then deep-dive with `get_entity`. Filter with
  `categories` when you know the type.
- `search_wiki` is full-text over wiki page titles and block prose, and never
  touches entities. Results carry `pageId`, `pageTitle`, an optional `blockId` when
  the hit is inside a block, a snippet, and a rank. Follow up with `get_wiki_page`.
  When you do not know whether the thing is a record or a page, run both searches.
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
  level of the call. `name` is required on every create; a nameless entity is
  refused.
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
  environmental, technological. Tags you pass on create and update are kept.
- **Bulk.** `batch_create_entities` takes up to 10 entities of one category per
  call, as an array. It is **not** transactional: entities insert sequentially, so a
  failure at #3 leaves #1 and #2 committed. Always read the `failed` array. To undo
  a partial batch, `delete_entity` the ids from `created`.
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
- **Threads.** After a session, log progressions rather than rewriting thread text.
  `add_thread_progression` takes `thread_id`, `session_id`, `progression_type`
  (start, update, complication, resolution), and a `title`. Link it to a session key
  event with `key_event_id` when you have one; `key_event_index` is legacy and
  breaks if key events are reordered. GM-role connections only.
- **Foundations.** The six `*_foundation_*` tools and `update_world_foundations`
  change the frame every future answer is generated inside. Confirm with the GM
  before writing there.

## The wiki

The wiki is where the GM writes in sentences: lore pages, handouts, house rules,
session notes. Over MCP you reach the **campaign wiki space only**. The party wiki
and each player's private notebook are not listed, not searchable, and not writable
from here; their page ids answer "page not found".

**Read.** `get_wiki_tree` for structure, `search_wiki` for words, `get_wiki_page`
for content. A page comes back with its blocks, breadcrumb, children, and a token
estimate. Each block has a typed `content` object and a flat `text` rendering.

**Links are tokens, and they are the point.** In block text, an entity mention is
`@[Label](entity://<category>/<uuid>)` and a page link is
`@[Title](page://<uuid>)`. Reads emit them; writes resolve valid tokens into real
links, which is what makes a page show up under an entity's "Mentioned in" list and
a page's "Linked from" rail in the app. Two rules follow:

- When you rewrite prose that contains tokens, edit around them and write the
  tokens back verbatim. Retyping a name as plain text silently breaks the link.
- When you write new prose about a record whose id you already hold, write the
  token, not the bare name. `@[Abbess Mirel](entity://npcs/<uuid>)` links the page
  to her record; "Abbess Mirel" does not. Category is the plural table name from
  `get_entity_schema` or the catalog: `npcs`, `factions`, `locations`, and so on.
  Page ids come from `get_wiki_tree`. A token with an id that does not exist in this
  campaign is stored as plain text, so use ids you have actually read.

**Blocks.** Types are text, heading (level 1 to 6), bullet, numbered, quote, callout
(info, warning, success, error), divider, page, image, and table. Text is wrapped as
`{"content": [{"text": "..."}]}`. `update_wiki_block` takes a content object of the
same shape, never the flat `text` string. Block type cannot be changed; delete and
recreate instead, with the cascade caveat below.

**Ordering.** `batch_create_wiki_blocks` (up to 50) chains blocks in array order,
so stored order matches your array with no verification pass. Sequential single
creates are ordered too; only concurrent create calls can interleave. Use
`batch_reorder_wiki_blocks` when you actually want a different order, passing every
block id on the page.

**Tables.** A table block is written and read as a plain markdown grid:
`{"markdown": "| Name | Role |\n| --- | --- |\n| Varka | Chief |"}`. Cells are
paragraph-only (no spans, no nested blocks), tokens work inside cells, and reads
render the grid back in `text`. An update replaces the whole table. Tables are not
allowed on the Campaign Bible and cannot be moved onto it. A person with the table
open in a browser tab re-saves their copy about two seconds after they click into
it, so an edit you make while it is on someone's screen can be overwritten; prefer
table edits between sessions.

**Pages.** `create_wiki_page` with `parent_id` nests; root pages cannot be
`inherit`. After creating a page, make it reachable: either a `page` block on the
parent (or the home page) or an inline `@[Title](page://<uuid>)` token in a
sentence. Pages you create are in the campaign space.

**Deleting.** `delete_wiki_block` is a hard delete with no undo; only pages go to
the wiki trash. Deleting a **page block** (type `page`) also moves its linked page
and every sub-page to the trash, exactly like deleting the page card in the app.
Never delete a page block as a way to "edit" it; page blocks accept only a
visibility change through `update_wiki_block`. The Campaign Bible's title block
cannot be deleted.

**The Campaign Bible** is a wiki page with its own tools. `get_campaign_bible`
returns blocks in whichever shape they were written in; when you rewrite, pass the
simple `{"content": [{"text": "..."}]}` shape to `update_campaign_bible` and let the
server normalize. Do not hand-build TipTap documents. The update replaces every
block except the title, allows only text, heading, bullet, numbered, quote,
callout, divider, and image, and is create-then-delete, so a failed write leaves the
old content in place.

## Visibility is not yours to negotiate

Every entity, wiki page, wiki block, and graph node carries a visibility value:

- `common-knowledge`: everyone at the table knows it.
- `player-knowledge`: on an **entity**, every player in the campaign can see it in
  the portal (it is the "the party has learned this" tier). On a **wiki page or
  block** it is, today, a GM tracking marker only: no player connection can read it,
  because the admission list it expects has no control that sets it. If the party
  should read a wiki page, mark it `common-knowledge` or `system`; do not rely on
  `player-knowledge` to share wiki content.
- `dm-secret`: GM only. This is the "GM secret" tier.
- `system` (wiki only): rules and meta notes; players see it like common knowledge.
- `inherit` (wiki only): a block follows its page; a child page follows the page
  block that links it, then its parent chain. The sane default for body content.
  Not allowed on a root page.

This is enforced in the backend query layer, not in the response text. A player-role
connection never receives GM-secret rows in the first place, and a page it cannot
reach is absent from the tree, search, mentions, and backlinks. There is no flag,
phrasing, or tool call that gets around it, and trying is a bug in your reasoning,
not a puzzle to solve.

What this means for you:

- If a player-role user asks about something and you find nothing, say you cannot
  see it in this campaign's data from this connection. Do not fill the gap with
  invention and do not speculate about what is being hidden.
- When you write, set visibility deliberately. Anything the GM told you in
  confidence, anything that would spoil a reveal, anything a villain knows and the
  party does not: `dm-secret`. Public reputation, geography, common rumor:
  `common-knowledge`. When unsure on a secret, choose `dm-secret`; over-hiding is
  recoverable, leaking is not.
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
6. If the notes are also worth keeping as prose (a handout, a read-aloud box, a
   timeline of the evening), file them as a wiki page whose sentences carry entity
   tokens for the records you just created, so the page and the records point at
   each other.
7. Report back: what you created, what you updated, what you skipped as a
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
6. If the table keeps a written recap in the wiki, add it as a page under the
   sessions page with tokens for the NPCs and places it mentions.

### Find it in the wiki

1. `search_wiki` with the words you remember. If a hit has a `blockId`, the match is
   inside that block; if not, the page title matched.
2. `get_wiki_page` on the best hit. Read the `text` of the relevant blocks.
3. Follow the tokens: `get_entity` on any `entity://` id and `get_wiki_page` on any
   `page://` id the prose points at, as far as the question needs.
4. If nothing matches and the user is a player, say it is not in the wiki you can
   see, without speculating.

### Write a wiki page for an entity

1. `get_entity` for the record and `get_wiki_tree` for where the page belongs.
2. `create_wiki_page` with a `parent_id` under the right section, `inherit`
   visibility unless the GM says otherwise.
3. `batch_create_wiki_blocks` with a heading, prose that carries the entity's own
   token and tokens for every related record you hold ids for, and a `dm-secret`
   callout for anything the party must not read.
4. Make it reachable: a `page` block on the parent page, or a page token in the
   parent's prose.
5. Tell the GM which record the page is linked to, so they can confirm it shows
   under "Mentioned in".

### Who knows what about X

1. `search_campaign` for X, get the id and category.
2. `get_relationships` on it for direct edges.
3. `get_knowledge_graph("political")` to see it inside the wider web of factions and
   NPCs rather than as an isolated node.
4. `search_wiki` for X to find the prose that mentions it, then `get_wiki_page`.
5. `get_entity` on the neighbours that matter.
6. Answer with the visibility tier attached. "The city guard treats this as common
   knowledge; the Duke's involvement is a GM secret" is a more useful answer than a
   flat list.

### Prep pull before a session

1. `get_narrative_state` for open threads and recent arcs.
2. `get_knowledge_graph("timeline", attention=true)` for what the GM has already
   flagged as live.
3. `get_knowledge_graph("geography")` if the party is travelling, to get the
   location hierarchy right.
4. `get_entity` on the handful of NPCs and locations that will actually appear, and
   `search_wiki` for any handout or rule the scene will need.
5. Offer the GM a short prep brief, then use `toggle_graph_attention` to flag the
   entities you both agreed matter, so the next session's Layer 3 view is already
   correct. If the brief is worth keeping, file it as a `dm-secret` wiki page or as
   a table block on the session prep page.

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
- Searching for a handout with `search_campaign`, or for an NPC with `search_wiki`,
  then concluding it does not exist.
- Retyping an entity or page name as plain text when the block already carried a
  token, or when you hold the id. That is how links quietly disappear.
- Deleting a page block to "edit" it. That trashes the page and its sub-pages.
- Editing a table while it is on someone's screen.
- Sharing wiki content with players by marking it `player-knowledge`.
- Treating a visibility filter as an obstacle, or a "page not found" as a clue.
- Rewriting a thread's text instead of logging a progression, which destroys the
  session-over-session history the timeline graph is built on.

## Reference

- `references/tool-groups.md`: all 49 tools by group, with what each is for.
- `references/context-architecture.md`: why the layers exist and how to pick a
  traversal.
- Setup and client instructions: https://www.ttrpg.bot/docs/mcp/
- The wiki, as the GM sees it: https://www.ttrpg.bot/dnd-campaign-wiki/
- Architecture write-up: https://www.ttrpg.bot/blog/grimoire-mcp-secret-sauce/
