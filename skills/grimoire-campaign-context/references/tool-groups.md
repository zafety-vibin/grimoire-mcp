# Grimoire MCP tool reference

49 tools across 7 groups. 21 read, 28 write. Categories are: npcs, locations,
factions, quests, items, lore_entries, session_recaps, creatures,
player_characters, world_rules, planar_forces, custom_mechanics, session_preps,
vehicles.

## Entities (10)

| Tool | Use it for |
|---|---|
| `search_campaign` | Full-text over every free-text field of an entity; GM tokens also match `dm_*` fields and custom field values, player tokens never do. Filter with `categories`. Returns id, category, name, truncated description, status, tags, relevance. A hit on a GM-only field may not show in the truncated description, so follow up with `get_entity`. Your default "find the thing" call. |
| `get_entity` | Full detail on one entity by id. Layer 5. |
| `list_entities` | Enumerate one category with offset and limit. Status filter: active (default), draft, hidden, archived. |
| `get_field_options` | Existing values already in use for a select field in this campaign. |
| `get_tag_options` | Tags already in use. Call before tagging. |
| `get_entity_schema` | The writable shape of a category: global fields, default fields, select option hints, campaign custom fields. json-type default fields carry a `shape` string describing the element shape (for example `session_recaps.key_events`). Call before every create or update. |
| `create_entity` | One entity. Default fields go inside `custom_fields` and land in real columns. |
| `update_entity` | Partial update. Send only what changed. |
| `delete_entity` | Transactional, with graph cleanup. |
| `batch_create_entities` | Up to 10, one category, sequential and not transactional. Always read the `failed` array. |

## Relationships (3)

| Tool | Use it for |
|---|---|
| `add_relationship` | Typed edge between two entities. Routes to an FK field when one exists, otherwise a junction row labelled with `relationship_type`. Appears as an edge in `get_knowledge_graph`; political relationships also put both endpoints on the app's Political Web, and so do the three political FK writes `npcs.faction_id` (member_of), `npcs.superior_npc_id` (reports_to) and `factions.leader_id` (led_by). Writing a session recap does the same for what it references: the NPCs and player characters it names join the Political Web, the locations join the geography graph. |
| `get_relationships` | One entity's edges, junction rows and FK rows, both real. On an incoming junction row `targetCategory` is the queried entity's category and `sourceCategory` is the other side. A table reachable from both sides (`faction_members`) is listed once per call. An FK plus a roster row for the same faction is one membership, not two. |
| `delete_relationship` | Remove an edge. |

FK fields worth knowing: `npcs.faction_id`, `npcs.superior_npc_id`,
`locations.parent_location_id`, `factions.leader_id`, `quests.quest_giver_id`,
`quests.started_session_id`, `quests.completed_session_id`, `items.owner_npc_id`,
`items.owner_pc_id`, `items.location_id`, `items.quest_id`,
`planar_forces.high_priest_id`, `custom_mechanics.world_rule_id`.

To draw faction membership, use
`add_relationship(npcs -> factions, "primary_faction")`. `add_relationship` with any
other `relationship_type` (member, agent, captain) writes a `faction_members` roster
row instead; it never sets both.

## Knowledge graphs (8)

| Tool | Use it for |
|---|---|
| `get_constitution` | Layer 1. Campaign summary, Campaign Bible summary, World Foundations nodes, the campaignContext block, plus attention-only compact node and edge maps (`graphStructures`) for the political, geography, and timeline projections; political edges use `get_knowledge_graph`'s vocabulary. Call first. |
| `get_entity_catalog` | Layer 2. Every active entity as id plus name, by type, plus compact graph node and edge lists and per-type counts. Compact edges use the same types as `get_knowledge_graph` and carry source, target and type only. Cheap. |
| `get_knowledge_graph` | Layer 3 with `attention=true`, Layer 4 without. Types: political, timeline, geography, foundations (foundations is portal-capped like the other three). Nodes carry entity data; edges are `member_of`, `ally`, `rival` and the NPC-to-NPC relationship types for political, `located_in` (child to parent) for geography, and `followed_by` plus the session-to-entity types for timeline. |
| `list_entity_graphs` | Custom graphs beyond the three built-ins. |
| `get_entity_graph` | One custom graph in full. |
| `add_to_entity_graph` | Put an entity on a custom graph. Only for graphs, not database links. |
| `create_entity_graph_edge` | A custom edge with no database equivalent. Never for something `add_relationship` already covers. |
| `toggle_graph_attention` | Promote or demote an entity between Layer 4 and the Layer 3 session view. |

The three built-in projections:

- **political**: factions, NPCs, alliances, rivalries, memberships, NPC-to-NPC
  relationships. The live relationship status of the campaign. Player characters are
  not on it today; read their ties with `get_relationships`.
- **timeline**: sessions and events in sequence, with entities connected to the
  sessions they appeared in. This is what establishes "recent" versus "old."
  Session-to-entity edges are typed `npc_encountered`, `location_visited`,
  `quest_progressed`, `loot_acquired`, `creature_encountered`, `pc_present`,
  `planar_force_mentioned` and `lore_learned`; entity endpoints are stub nodes
  carrying the `get_entity` category, so follow with `get_entity`.
- **geography**: locations in a hierarchy (plane, world, continent, downward), drawn
  as `located_in` edges from child to parent. Sibling links, and where NPCs and items
  sit inside it, are not drawn today; read those with `get_relationships`.

## World foundations (6)

The constitution's write path. Foundation nodes are world-defining concepts, a
cataclysm, a hidden power, a cosmological law, connected by labelled edges. These
are mechanical and cosmological truths, distinct from database entities.

`create_foundation_node`, `update_foundation_node`, `delete_foundation_node`,
`create_foundation_edge`, `delete_foundation_edge`, `update_world_foundations`
(sets the structured World Foundations document: magic system, divine hierarchy,
fundamental laws).

Treat these as high-stakes. Changing a foundation changes the frame every future
answer is generated inside. Confirm with the GM before writing here.

## Open threads (7)

| Tool | Use it for |
|---|---|
| `get_open_threads` | Unresolved narrative obligations, grouped major and minor. Types: consequence, promise, mystery, foreshadowing, callback_opportunity. Every thread carries a `status` of open or resolved; `status=all` interleaves both kinds inside the groups, so read the field rather than testing for `resolvedSessionId`. |
| `create_open_thread` | New loose end. With `created_session_id` a start progression is logged for you. |
| `resolve_open_thread` / `unresolve_open_thread` | Close (logs a resolution progression at that session unless one exists) or reopen (progression history is kept). |
| `update_open_thread` | Edit thread metadata. |
| `get_thread_progressions` | How a thread has moved session over session. |
| `add_thread_progression` | Log a movement: start, update, complication, resolution. Link with `key_event_id` when possible. GM role only. |

## Wiki (10)

Block-based collaborative pages. `get_wiki_tree`, `get_wiki_page`, `search_wiki`,
`create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`,
`delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`.

`search_wiki` is full-text over page titles and block prose; `search_campaign`
never touches wiki rows. Results carry `pageId`, `pageTitle`, an optional
`blockId` when the hit is inside a block, a snippet, and a rank. Follow up with
`get_wiki_page` for the full page.

Block types: text, heading (level 1 to 6), bullet, numbered, quote, callout (info,
warning, success, error), divider, page, image, table. Text is wrapped as
`{"content": [{"text": "..."}]}`. A table is written and read as a plain markdown
grid: `{"markdown": "| Name | Role |\n| --- | --- |\n| Varka | Chief |"}`.

Links inside any block, including table cells, travel as tokens the server resolves
to real records: `@[Name](entity://category/uuid)` for entities and
`@[Title](page://uuid)` for pages. Read a block's `text`, edit around the tokens,
and write it back verbatim to keep the links.

Multi-block flow: `batch_create_wiki_blocks` (up to 50) chains blocks in array order,
so stored order matches your array; no verification or reorder pass is needed.
Sequential single creates are ordered too; only concurrent create calls can
interleave. `batch_reorder_wiki_blocks` is for when you actually want a different
order, passing every block id on the page.

Block visibility: `inherit` (default; not allowed on a root page),
`common-knowledge`, `dm-secret`, `player-knowledge` (on wiki content a GM tracking
marker today, see the Visibility section of SKILL.md), `system` (rules and meta
notes; players read it like common-knowledge). All wiki tools work in the campaign
wiki space only; a page the connected role cannot reach answers "page not found"
everywhere.

## Campaign meta (5)

| Tool | Use it for |
|---|---|
| `current_campaign` | Campaign, user, and role bound to this session. First call, always. |
| `get_campaign_context` | Genre, ruleset, setting technology and magic levels, which categories exist and what this campaign calls them. Also embedded verbatim as `campaignContext` in `get_constitution`, `get_narrative_state` and `get_entity_catalog`; fetch it standalone only when none of those are loaded. |
| `get_narrative_state` | Layer 1.5. Recent sessions with summaries and key events, open threads by weight, canonical facts, active arcs, recent observations. GM connections also get `dmConsequences` and `dmBehindScenes` per session. `recent_session_count` defaults to 3, max 10. |
| `get_campaign_bible` | Full Campaign Bible blocks. `get_constitution` returns a summary; this returns the blocks and their ids. |
| `update_campaign_bible` | Edit Bible blocks. Pass the simple content shape and let the server normalize. |

## Resources

The remote server exposes one MCP resource, `campaign://wiki` ("Wiki Pages"),
which returns the campaign's wiki page list and hierarchy as JSON. GM-secret pages
are filtered out for player-role connections. Use `get_wiki_tree` when you want the
same information through the tool interface.

## Role gating at a glance

`tools/list` returns the same 49 tools regardless of role. Gating happens at call
time.

| Availability | Tools |
|---|---|
| Player: permanently denied | `get_constitution`, `get_campaign_bible`, `get_relationships`, `get_open_threads`, `get_thread_progressions`, `list_entity_graphs`, `get_entity_graph` |
| Player: allowed only when the portal policy is `revealed` or `open` | `get_campaign_context`, `get_narrative_state`, `get_entity_catalog`, `get_knowledge_graph` |
| Player: always allowed | `current_campaign`, `search_campaign`, `search_wiki`, `get_wiki_tree`, `get_wiki_page`, `get_entity`, `list_entities`, `get_entity_schema`, `get_field_options`, `get_tag_options` (`get_entity` and `list_entities` apply the dm-secret row filter; the other three carry no player-role gate today) |
| GM only | every write tool, and everything in the first two rows |

Player-visible results are additionally visibility-filtered in the query layer, so
a permitted tool still returns less data for a player than for the GM.
