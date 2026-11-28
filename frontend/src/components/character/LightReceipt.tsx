import type { LightReceipt as Receipt } from '../../api/types';
import { formatCampaignSeconds } from '../../lib/campaignTime';

export function LightReceipt({ receipt }: { readonly receipt: Receipt }) {
  return <section aria-label="Saved Light receipt" className="space-y-2 border border-jade/20 p-3">
    <p>{receipt.reason}</p>
    <p>{receipt.after ? `Light (${receipt.after.color}) lasts until ${formatCampaignSeconds(receipt.after.expires_at)}${receipt.after.covered ? '; covered' : ''}.` : 'No Light remains.'}</p>
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect Light sources</summary>{receipt.before && <p className="break-all">Previous effect ID: {receipt.before.id}</p>}{receipt.after && <p className="break-all">Current effect ID: {receipt.after.id}</p>}{receipt.source_ids.map(id => <p key={id} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
