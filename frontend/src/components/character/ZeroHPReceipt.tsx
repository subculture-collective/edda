import type { DeathSaves, ZeroHPReceipt as Receipt } from '../../api/types';

// deathSaveText states the saved 0-HP counts, Stable or dead.
export const deathSaveText = (saves: DeathSaves) => saves.dead ? 'Dead.' : saves.stable ? 'Stable at 0 HP; still Unconscious.' : `Death Saving Throws: ${saves.successes ?? 0} successes, ${saves.failures ?? 0} failures.`;

// The die, the counts and the HP come from the server.
export function ZeroHPReceipt({ receipt }: { readonly receipt: Receipt }) {
  return <section aria-label="Saved 0 HP receipt" className="space-y-2 border border-ruby/30 p-3">
    <p>{receipt.reason}</p>
    {receipt.hp_before !== receipt.hp_after && <p>HP {receipt.hp_before} → {receipt.hp_after}.</p>}
    {!receipt.after.dead && (receipt.hp_after === 0) && <p>{deathSaveText(receipt.after)}</p>}
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect 0 HP sources</summary>{receipt.source_ids.map((id, index) => <p key={index} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
