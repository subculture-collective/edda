import type { TravelReceipt as Receipt } from '../../api/types';
import { formatCampaignSeconds } from '../../lib/campaignTime';

const duration = (seconds: number) => `${Math.floor(seconds / 60)} min ${String(seconds % 60).padStart(2, '0')} s`;

// The route facts, travel time and clock all come from the server.
export function TravelReceipt({ receipt }: { readonly receipt: Receipt }) {
  const { plan } = receipt;
  return <section aria-label="Saved travel receipt" className="space-y-2 border border-jade/20 p-3 wrap-break-word">
    <p>Travelled from {receipt.from} to {receipt.to} at a {plan.pace} pace.</p>
    <p>{(plan.distance_feet / 5280).toFixed(2)} miles ({plan.distance_feet} feet), {plan.terrain} terrain, {plan.good_road ? 'good road' : 'no road'}: {duration(plan.seconds)}.</p>
    <p>Campaign time {formatCampaignSeconds(receipt.time_before)} → {formatCampaignSeconds(receipt.time_after)}. Travel today {duration(receipt.day_travel_before)} → {duration(receipt.day_travel_after)}.</p>
    <p>Pace effect while travelling: {plan.effects.join('; ')}. No checks were made during the segment.</p>
    {receipt.expired && <p>{receipt.expired.spell === 'detect-magic' ? 'Detect Magic' : 'Shield of Faith'} expired on the way.</p>}
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect travel sources</summary>{receipt.source_ids.map((id, index) => <p key={index} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
