import type { ZeroHPReceipt as Receipt } from '../../api/types';
import { deathSaveText } from './deathSaveText';

// The die, the counts and the HP come from the server. After an ally's aid,
// the Death Saving Throw that followed holds the final state.
export function ZeroHPReceipt({ receipt }: { readonly receipt: Receipt }) {
  const final = receipt.death_save ?? receipt;
  const sources = [...receipt.source_ids, ...(receipt.death_save?.source_ids ?? [])];
  return <section aria-label="Saved 0 HP receipt" className="space-y-2 border border-ruby/30 p-3">
    <p>{receipt.reason}</p>
    {receipt.hp_before !== final.hp_after && <p>HP {receipt.hp_before} → {final.hp_after}.</p>}
    {!final.after.dead && (final.hp_after === 0) && <p>{deathSaveText(final.after)}</p>}
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect 0 HP sources</summary>{sources.map((id, index) => <p key={index} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
