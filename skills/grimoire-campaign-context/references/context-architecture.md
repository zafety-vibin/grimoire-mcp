# Why Grimoire layers context

Dumping an entire campaign into a prompt does not work. A model given every session
note, recap, and prep document at once cannot tell what is recent from what is old,
or what is central from what is incidental, or what has been retconned. It produces
answers that sound plausible and contradict the GM's canon.

Grimoire's server splits the campaign into layers, each with a purpose and a
criterion for when it gets pulled in. Your traversal through those layers is what
constrains the answer to the world the GM actually built.

## The layers

**Layer 1, Constitution.** Always loaded. `get_constitution` returns the world
foundations graph and the Campaign Bible summary. The foundations graph holds the
biggest concepts of the worldbuilding: what is uniquely different about this world,
which single historical events define everything else, anything non-negotiable, and
any nuance that breaks from default genre assumptions. The Campaign Bible holds what
does not fit a graph, including meta rules such as content boundaries. This layer is
the enforcement mechanism on your behavior. Load it before anything, including
something as simple as parsing notes.

Pair it with `get_campaign_context`, which tells you the genre, the ruleset, the
technology and magic levels, and this campaign's own vocabulary for its categories.
Together they answer "what kind of world is this and how do I talk about it."

**Layer 1.5, Narrative state.** `get_narrative_state`. Where the story is right now:
last three sessions with summaries and key events, open threads split into
story-critical and callback-worthy, canonical facts, and active arcs. Entities the
GM has marked with attention get pulled in here too, which is how a distant villain
stays in scope without being on the party's current path.

**Layer 2, The Map.** `get_entity_catalog`. Every entity and how it connects, ids
and names only, no detail. You can see that a character is from a place with "home"
as the relationship, and nothing else. This is how you survey a large campaign
without burning the context window. See the shape, decide where to look, then pull
detail from there. Much better than one over-broad search that misses something and
answers anyway.

**Layer 3, Session context.** `get_knowledge_graph(type, attention=true)`. The
GM-curated slice: only entities flagged as relevant right now. On a campaign of any
size this should be your first graph call.

**Layer 4, A full graph.** `get_knowledge_graph(type)`. One projection, complete.
Every user gets political, timeline, and geography; higher tiers can define custom
entity graphs.

**Layer 5, Entity deep dive.** `get_entity`. One record, complete. Call it as many
times as the question needs.

**Beside the layers, the wiki.** The wiki is the campaign's prose: handouts, lore
pages, house rules, session notes, in blocks whose sentences carry tokens that
resolve to real records and pages. It is not numbered because it is reached by
subject rather than by breadth: `search_wiki` for the words, `get_wiki_page` for the
page, then the tokens lead you back into the layers (`get_entity` on an
`entity://` id, another page on a `page://` id). When the question is "what did the
GM write about X", start here; when it is "what is true about X", start in the
layers and use the wiki to read the prose around the facts.

**The result.** The context you end up holding is assembled from the layers you
chose, in the order you chose them. It is unique per question, and that trajectory
is what keeps the answer inside the world's nuance instead of drifting toward genre
defaults.

## Picking a traversal

| Question shape | Traversal |
|---|---|
| "Who is X" | constitution, context, search, entity |
| "What happened last session" | constitution, context, narrative state |
| "How would X react to Y" | constitution, context, narrative state, political graph, entity on X and the relevant neighbours |
| "Where is X relative to Y" | constitution, context, geography graph |
| "What should I prep" | constitution, context, narrative state, timeline graph with attention, entity on the shortlist |
| "Turn these notes into entities" | constitution, context, catalog, schema, batch create, relationships, then a wiki page with tokens if the prose is worth keeping |
| "Who knows about Z" | constitution, context, search, relationships, political graph, wiki search for the prose |
| "What did the GM write about the abbey" | constitution, context, wiki search, wiki page, then entity on the tokens the page carries |
| "Write up X for the players" | constitution, context, entity on X, wiki tree, create page under the right section, blocks with tokens, a page block on the parent |

## The point

The AI is meant to be a steward of the GM's world, not a generator that produces
things that sound right but do not fit. Loading the constitution first, then only
what the question needs, is what makes the difference.
