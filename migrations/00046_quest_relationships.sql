-- +goose Up
-- link_quest_entity and branch_quest store quest-sourced and quest-to-quest
-- relationships, so both ends of a relationship may be a quest.
ALTER TABLE entity_relationships
  DROP CONSTRAINT entity_relationships_source_entity_type_check,
  ADD CONSTRAINT entity_relationships_source_entity_type_check
    CHECK (source_entity_type IN ('npc', 'location', 'faction', 'player_character', 'item', 'quest')),
  DROP CONSTRAINT entity_relationships_target_entity_type_check,
  ADD CONSTRAINT entity_relationships_target_entity_type_check
    CHECK (target_entity_type IN ('npc', 'location', 'faction', 'player_character', 'item', 'quest'));

-- +goose Down
DELETE FROM entity_relationships WHERE source_entity_type = 'quest' OR target_entity_type = 'quest';
ALTER TABLE entity_relationships
  DROP CONSTRAINT entity_relationships_source_entity_type_check,
  ADD CONSTRAINT entity_relationships_source_entity_type_check
    CHECK (source_entity_type IN ('npc', 'location', 'faction', 'player_character', 'item')),
  DROP CONSTRAINT entity_relationships_target_entity_type_check,
  ADD CONSTRAINT entity_relationships_target_entity_type_check
    CHECK (target_entity_type IN ('npc', 'location', 'faction', 'player_character', 'item'));
