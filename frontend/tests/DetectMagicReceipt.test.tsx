import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { expect, it } from 'vitest';
import type { DetectMagicReceipt } from '../src/api/types';
import { NarrativeEntry } from '../src/components/narrative/NarrativeEntry';

// The server's public receipt (srd521.DetectMagicReceipt). A cast says only
// whether magic is sensed; a study names the visible auras in its reason.
const concentration = { id: 'cast-id', spell: 'detect-magic', started_at: 28860, expires_at: 29460 };
const base = { rules_version: 'srd-5.2.1', spell: 'detect-magic', slots_before: 1, slots_after: 1, time_before: 28860, time_after: 28860, presence: true, source_ids: ['spells/detect-magic'] };
const cast: DetectMagicReceipt = { ...base, kind: 'spell-cast', use: 'magic-initiate-free', magic_initiate_uses_before: 1, magic_initiate_uses_after: 0, before: { ...concentration, id: 'faith-id', spell: 'shield-of-faith' }, after: concentration, auras: [], reason: 'Detect Magic cast with the Magic Initiate (Cleric) free use. Adjudicated magic is sensed within 30 feet; take the Magic action to study visible auras.' };
const study: DetectMagicReceipt = { ...base, kind: 'detect-magic-study', use: '', magic_initiate_uses_before: 0, magic_initiate_uses_after: 0, before: concentration, after: concentration, auras: [{ subject_id: 'ward', subject: 'Warded door', distance_feet: 20, school: 'abjuration', origin: 'spell' }], reason: 'Detect Magic study; Magic action; no slot or feat use. Faint aura around Warded door (20 feet): abjuration spell.' };

async function mountHistory(receipt: DetectMagicReceipt) {
  render(<NarrativeEntry entry={{ id: 'turn', kind: 'gm', text: 'Saved.', timestamp: '2026-11-28T12:00:00Z', resolutionEvents: [{ type: 'srd_detect_magic', label: 'Detect Magic', outcome: receipt.reason, details: { receipt } }] }} />);
  await userEvent.click(screen.getByText('Turn details'));
  return screen.getByRole('region', { name: 'Saved Detect Magic receipt' });
}

it('shows a saved Detect Magic cast in the turn history', async () => {
  const saved = await mountHistory(cast);
  expect(saved).toHaveTextContent(cast.reason);
  expect(saved).toHaveTextContent('Magic Initiate uses: 1 → 0. Spell slots: 1 → 1.');
  expect(saved).toHaveTextContent('Detect Magic concentration lasts until Day 1, 08:11:00.');
  expect(saved).toHaveTextContent('Previous effect ID: faith-id');
  expect(saved).toHaveTextContent('Current effect ID: cast-id');
  expect(saved).toHaveTextContent('spells/detect-magic');
});

it('shows a saved Detect Magic study in the turn history', async () => {
  const saved = await mountHistory(study);
  // The aura appears once, as the receipt words it.
  expect(saved.textContent?.split('Faint aura around Warded door (20 feet): abjuration spell.')).toHaveLength(2);
  expect(saved).toHaveTextContent('Detect Magic concentration lasts until Day 1, 08:11:00.');
  // A study spends no feat use and replaces no effect.
  expect(saved).not.toHaveTextContent('Magic Initiate uses');
  expect(saved).not.toHaveTextContent('Previous effect ID');
});
