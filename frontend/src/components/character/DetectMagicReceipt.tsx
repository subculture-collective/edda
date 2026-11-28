import type { DetectMagicReceipt as Receipt } from '../../api/types';
import { formatCampaignSeconds } from '../../lib/campaignTime';

// The reason is the whole disclosure: a cast says only whether magic is sensed,
// and a study words each visible aura.
export function DetectMagicReceipt({ receipt }: { readonly receipt: Receipt }) {
  return <section aria-label="Saved Detect Magic receipt" className="space-y-2 border border-jade/20 p-3">
    <p>{receipt.reason}</p>
    {receipt.kind === 'spell-cast' && <p>Magic Initiate uses: {receipt.magic_initiate_uses_before} → {receipt.magic_initiate_uses_after}. Spell slots: {receipt.slots_before} → {receipt.slots_after}.</p>}
    {receipt.after && <p>Detect Magic concentration lasts until {formatCampaignSeconds(receipt.after.expires_at)}.</p>}
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect Detect Magic sources</summary>{receipt.before && receipt.before.id !== receipt.after?.id && <p className="break-all">Previous effect ID: {receipt.before.id}</p>}{receipt.after && <p className="break-all">Current effect ID: {receipt.after.id}</p>}{receipt.source_ids.map(id => <p key={id} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
