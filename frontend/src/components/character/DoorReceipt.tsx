import type { DoorReceipt as Receipt } from '../../api/types';

// The door's state, distance and reach all come from the server.
export function DoorReceipt({ receipt }: { readonly receipt: Receipt }) {
  return <section aria-label="Saved door receipt" className="space-y-2 border border-jade/20 p-3">
    <p>{receipt.door}: {receipt.before} → {receipt.after}.</p>
    <p>{receipt.distance_feet} ft away, within {receipt.reach_feet} ft reach. No action and no campaign time.</p>
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect door sources</summary>{receipt.source_ids.map((id, index) => <p key={index} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
