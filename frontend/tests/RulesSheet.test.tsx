import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { expect, it } from 'vitest';
import { RulesSheet } from '../src/components/character/RulesSheet';
import type { CharacterSheetProjection } from '../src/api/types';
it('displays server-derived values and inspectable terms without a client XP formula', async () => {
 const value = { value: 5, derivation: [{ label: 'Strength', value: 3, source_id: 'srd-5.2.1:content/character/ability-modifiers' }, { label: 'Proficiency', value: 2, source_id: 'srd-5.2.1:content/character/proficiency' }] };
 const sheet: CharacterSheetProjection = { rules_version: 'srd-5.2.1', scope: 'fixture', abilities: {}, proficiency: value, saves: {}, skills: { athletics: value }, hp: { ...value, value: 9 }, max_hp: { ...value, value: 12 }, defenses: {}, resources: [{ id: 'second-wind', label: 'Second Wind', current: 1, maximum: 2, source_id: 'srd-5.2.1:content/fighter/second-wind' }], capabilities: [{ id: 'attack', label: 'Attack', playable: false, source_id: 'srd-5.2.1:content/actions/attack' }], progression: { current_xp: 125, level_xp: 0, next_level_xp: 300, source_id: 'srd-5.2.1:content/character/level-advancement' } };
 render(<RulesSheet sheet={sheet} />);
 expect(screen.getByText('HP 9 / 12 · XP 125 / 300')).toBeVisible();
 expect(screen.getByText('Second Wind: 1 / 2')).toBeVisible();
 expect(screen.getByText('Attack (read only)')).toBeVisible();
 await userEvent.click(screen.getAllByText('+5')[0]);
 expect(screen.getAllByText('Strength: 3')[0]).toBeVisible();
 expect(screen.getAllByText('srd-5.2.1:content/character/proficiency')[0]).toBeVisible();
});
