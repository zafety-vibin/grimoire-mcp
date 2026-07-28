# Grimoire MCP tool reference

48 tools across 7 groups. 20 read, 28 write. Categories are: npcs, locations,
factions, quests, items, lore_entries, session_recaps, creatures,
player_characters, world_rules, planar_forces, custom_mechanics, session_preps,
vehicles.

## Entities (10)

| Tool | Use it for |
|---|---|
| `search_campaign` | Full-text across name, description, and custom fields. Filter with `categories`. Returns id, category, name, truncated description, status, tags, relevance. Your default "find the thing" call. |
| `get_entity` | Full detail on one entity by id. Layer 5. |
| `list_entities` | Enumerate one category with offset and limit. Status filter: active (default), draft, hidden, archived. |
| `get_field_options` | Existing values already in use for a select field in this campaign. |
| `get_tag_options` | Tags already in use. Call before tagging. |
| `get_entity_schema` | The writable shape of a category: global fields, default fields, select option hints, campaign custom fields. Call before every create or update. |
| `create_entity` | One entity. Default fields go inside `custom_fields` and land in real columns. |
| `update_entity` | Partial update. Send only what changed. |
| `delete_entity` | Transactional, with graph cleanup. |
| `batch_create_entities` | Up to 10, one category, sequential and not transactional. Always read the `failed` array. |

## Relationships (3)

| Tool | Use it for |
|---|---|
| `add_relationship` | Typed edge between two entities. Routes to an FK field when one exists, otherwise a junction row labelled with `relationship_type`. Auto-visualizes in the knowledge graphs. |
| `get_relationships` | One entity's edges. |
| `delete_relationship` | Remove an edge. |

FK fields worth knowing: `npcs.faction_id`, `npcs.superior_npc_id`,
`npcs.location_ids[]`, `locations.parent_location_id`,
`locations.connected_location_ids[]`, `factions.leader_id`,
`factions.allied_faction_ids[]`, `factions.rival_faction_ids[]`,
`quests.quest_giver_id`, `items.owner_npc_id`, `items.location_id`.

To draw faction membership, use
`add_relationship(npcs -> factions, "primary_faction")`.

## Knowledge graphs (8)

| Tool | Use it for |
|---|---|
| `get_constitution` | Layer 1. Campaign summary, Campaign Bible summary, foundation nodes. Call first. |
| `get_entity_catalog` | Layer 2. Every active entity as id plus name, by type, plus compact graph node and edge lists and per-type counts. Cheap. |
| `get_knowledge_graph` | Layer 3 with `attention=true`, Layer 4 without. Types: political, timeline, geography. Nodes carry entity data, edges include database relationships automatically. |
| `list_entity_graphs` | Custom graphs beyond the three built-ins. |
| `get_entity_graph` | One custom graph in full. |
| `add_to_entity_graph` | Put an entity on a custom graph. Only for graphs, not database links. |
| `create_entity_graph_edge` | A custom edge with no database equivalent. Never for something `add_relationship` already covers. |
| `toggle_graph_attention` | Promote or demote an entity between Layer 4 and the Layer 3 session view. |

The three built-in projections:

- **political**: factions, NPCs, player characters, alliances, rivalries,
  memberships. The live relationship status of the campaign.
- **timeline**: sessions and events in sequence, with entities connected to the
  sessions they appeared in. This is what establishes "recent" versus "old."
- **geography**: locations in a hierarchy (plane, world, continent, downward) with
  spatial relationships between siblings, plus where NPCs and items sit inside it.

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
| `get_open_threads` | Unresolved narrative obligations, grouped major and minor. Types: consequence, promise, mystery, foreshadowing, callback_opportunity. |
| `create_open_thread` | New loose end. |
| `resolve_open_thread` / `unresolve_open_thread` | Close or reopen. |
| `update_open_thread` | Edit thread metadata. |
| `get_thread_progressions` | How a thread has moved session over session. |
| `add_thread_progression` | Log a movement: start, update, complication, resolution. Link with `key_event_id` when possible. GM role only. |

## Wiki (9)

Block-based collaborative pages. `get_wiki_tree`, `get_wiki_page`,
`create_wiki_page`, `create_wiki_block`, `update_wiki_block`, `move_wiki_block`,
`delete_wiki_block`, `batch_create_wiki_blocks`, `batch_reorder_wiki_blocks`.

Block types: text, heading (level 1 to 6), bullet, numbered, quote, callout (info,
warning, success, error), divider, page, image. Text is wrapped as
`{"content": [{"text": "..."}]}`.

Multi-block flow: `batch_create_wiki_blocks`, then `get_wiki_page` to verify order,
then `batch_reorder_wiki_blocks`. Single-block creates have a positioning race.

Block visibility: `inherit` (default), `common-knowledge`, `dm-secret`,
`player-knowledge`.

## Campaign meta (5)

| Tool | Use it for |
|---|---|
| `current_campaign` | Campaign, user, and role bound to this session. First call, always. |
| `get_campaign_context` | Genre, ruleset, setting technology and magic levels, which categories exist and what this campaign calls them. |
| `get_narrative_state` | Layer 1.5. Recent sessions with summaries and key events, open threads by weight, canonical facts, active arcs, recent observations. `recent_session_count` defaults to 3, max 10. |
| `get_campaign_bible` | Full Campaign Bible blocks. `get_constitution` returns a summary; this returns the blocks and their ids. |
| `update_campaign_bible` | Edit Bible blocks. Pass the simple content shape and let the server normalize. |

## Resources

The remote server exposes one MCP resource, `campaign://wiki` ("Wiki Pages"),
which returns the campaign's wiki page list and hierarchy as JSON. GM-secret pages
are filtered out for player-role connections. Use `get_wiki_tree` when you want the
same information through the tool interface.

## Role gating at a glance

`tools/list` returns the same 48 tools regardless of role. Gating happens at call
time.

| Availability | Tools |
|---|---|
| Player: permanently denied | `get_constitution`, `get_campaign_bible`, `get_relationships`, `get_open_threads`, `get_thread_progressions`, `list_entity_graphs`, `get_entity_graph` |
| Player: allowed only when the portal policy is `revealed` or `open` | `get_campaign_context`, `get_narrative_state`, `get_entity_catalog`, `get_knowledge_graph` |
| Player: always allowed | `current_campaign`, `search_campaign`, `get_wiki_tree`, `get_wiki_page` |
| GM only | every write tool, and everything in the first two rows |

Player-visible results are additionally visibility-filtered in the query layer, so
a permitted tool still returns less data for a player than for the GM.
