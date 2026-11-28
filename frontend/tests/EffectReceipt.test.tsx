import { render, screen } from '@testing-library/react';
import { expect, it } from 'vitest';
import type { EffectReceipt as Receipt } from '../src/api/types';
import { EffectReceipt } from '../src/components/character/EffectReceipt';

const faith = { id: 'faith-id', spell: 'shield-of-faith', started_at: 28800, expires_at: 29400 };
const detect = { id: 'detect-id', spell: 'detect-magic', started_at: 28860, expires_at: 29460 };
const base: Receipt = { rules_version: 'srd-5.2.1', kind: 'wait-minute', reason: 'One minute of quiet downtime elapsed.', time_before: 28800, time_after: 28860, slots_before: 1, slots_after: 1, hp_before: 10, hp_after: 10, hit_dice_before: 1, hit_dice_after: 1, source_ids: [] };

it('names the previous and current concentration spells', () => {
  render(<EffectReceipt receipt={{ ...base, before: { ...faith, spell_label: 'Shield of Faith' }, after: { ...detect, spell_label: 'Detect Magic' } }} />);
  const saved = screen.getByRole('region', { name: 'Saved caster effect receipt' });
  expect(saved).toHaveTextContent('Previous effect: Shield of Faith (expires at Day 1, 08:10:00).');
  expect(saved).toHaveTextContent('Concentration: Detect Magic until Day 1, 08:11:00.');
});

// An older API sends no spell_label: the lines show the ID in title case.
it('shows a title-cased spell ID when the server sends no name', () => {
  render(<EffectReceipt receipt={{ ...base, before: faith, after: detect }} />);
  const saved = screen.getByRole('region', { name: 'Saved caster effect receipt' });
  expect(saved).toHaveTextContent('Previous effect: Shield of Faith (expires at');
  expect(saved).toHaveTextContent('Concentration: Detect Magic until');
});
