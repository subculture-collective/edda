import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import type { CharacterSheetProjection, Concentration, EquipmentReceipt as Equipment, TravelReceipt as Travel } from '../src/api/types';
import { EquipmentReceipt } from '../src/components/character/EquipmentReceipt';
import { RulesSheet } from '../src/components/character/RulesSheet';
import { TravelReceipt } from '../src/components/character/TravelReceipt';
import fixture from './fixtures/srd-check-sheet.json';

const effect = (spell: string, spell_label?: string): Concentration => ({ id: 'cast', spell, spell_label, started_at: 28800, expires_at: 29400 });
// A spell the client does not know, once with the server's name and once from an older API.
const cases: [string, Concentration, string][] = [
  ["the server's name", effect('hunters-mark', 'Hunter’s Mark'), 'Hunter’s Mark'],
  ['no name, as from an older API', effect('protection-from-evil-and-good'), 'Protection from Evil and Good'],
  ['no name for Detect Magic', effect('detect-magic'), 'Detect Magic'],
  ['no name for Shield of Faith', effect('shield-of-faith'), 'Shield of Faith'],
];

describe.each(cases)('with %s', (_, concentration, name) => {
  const other = name !== 'Shield of Faith';

  it('names the spell that expired during an equipment change', () => {
    const receipt: Equipment = { rules_version: 'srd-5.2.1', kind: 'equipment-equip', item_id: 'item', equipment: 'chain-mail', equipped_before: false, equipped_after: true, ac_before: 11, ac_after: 16, time_before: 28800, time_after: 29400, expired: concentration, source_ids: [] };
    render(<EquipmentReceipt receipt={receipt} />);
    const saved = screen.getByRole('region', { name: 'Saved equipment receipt' });
    expect(saved).toHaveTextContent(`${name} expired.`);
    if (other) expect(saved).not.toHaveTextContent('Shield of Faith');
  });

  it('names the spell that expired during travel', () => {
    const receipt: Travel = { rules_version: 'srd-5.2.1', kind: 'travel-hour', reason: '', route_id: 'r', from_id: 'a', from: 'Gate', to_id: 'b', to: 'Mill', plan: { pace: 'normal', maximum_pace: 'normal', distance_feet: 5280, terrain: 'forest', good_road: false, miles_per_hour: 3, seconds: 1200, effects: [] }, time_before: 28800, time_after: 30000, day_travel_before: 0, day_travel_after: 1200, expired: concentration, source_ids: [] };
    render(<TravelReceipt receipt={receipt} />);
    const saved = screen.getByRole('region', { name: 'Saved travel receipt' });
    expect(saved).toHaveTextContent(`${name} expired on the way.`);
    if (other) expect(saved).not.toHaveTextContent('Shield of Faith');
  });

  it('names the active concentration on the sheet', () => {
    render(<RulesSheet sheet={{ ...(fixture.sheet as CharacterSheetProjection), spellcasting: { ability: 'wisdom', modifier: { value: 3, derivation: [] }, spells: [] }, concentration }} />);
    const saved = screen.getByRole('region', { name: 'Concentration and recovery' });
    // Only the two supported spells have a described effect.
    const effects: Record<string, string> = { 'detect-magic': ': sensing magic within 30 feet', 'shield-of-faith': ': +2 AC' };
    expect(saved).toHaveTextContent(`${name}${effects[concentration.spell] ?? ''}. Started Day 1, 08:00:00; expires Day 1, 08:10:00.`);
    if (other) expect(saved).not.toHaveTextContent('+2 AC');
  });
});
