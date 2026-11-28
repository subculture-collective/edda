import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { expect, it } from 'vitest';
import type { ZeroHPReceipt } from '../src/api/types';
import { NarrativeEntry } from '../src/components/narrative/NarrativeEntry';

// The server sends an encounter turn that starts at 0 HP as one srd_death_save
// event per receipt: the repair, if due, then the Death Saving Throw.
const repair: ZeroHPReceipt = { rules_version: 'srd-5.2.1', event: 'unconscious-repair', reason: 'At 0 HP without the Unconscious condition: now Unconscious and Prone.', hp_before: 0, hp_after: 0, max_hp: 12, before: {}, after: {}, source_ids: ['falling-unconscious'] };
const save: ZeroHPReceipt = { rules_version: 'srd-5.2.1', event: 'death-save', reason: 'Death Saving Throw: d20 12, success (1 successes, 0 failures).', die: 12, outcome: 'success', hp_before: 0, hp_after: 0, max_hp: 12, before: {}, after: { successes: 1 }, source_ids: ['death-saving-throws'] };

it('shows the Unconscious repair receipt before the Death Saving Throw receipt', async () => {
  render(<NarrativeEntry entry={{ id: 'turn', kind: 'gm', text: 'Player turn ended.', timestamp: '2026-11-28T12:00:00Z', resolutionEvents: [repair, save].map(receipt => ({ type: 'srd_death_save', label: '0 HP', outcome: receipt.reason, details: { receipt } })) }} />);
  await userEvent.click(screen.getByText('Turn details'));
  const [first, second, ...rest] = screen.getAllByRole('region', { name: 'Saved 0 HP receipt' });
  expect(rest).toHaveLength(0);
  expect(first).toHaveTextContent(repair.reason);
  expect(first).toHaveTextContent('Death Saving Throws: 0 successes, 0 failures.');
  expect(first).toHaveTextContent('falling-unconscious');
  expect(second).toHaveTextContent(save.reason);
  expect(second).toHaveTextContent('Death Saving Throws: 1 successes, 0 failures.');
});
