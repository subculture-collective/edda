import type { FallReceipt as Receipt } from '../../api/types';
import { CheckReceipt } from './CheckReceipt';
import { EffectReceipt } from './EffectReceipt';
import { ZeroHPReceipt } from './ZeroHPReceipt';

const proneText = (prone: boolean) => prone ? 'Prone' : 'not Prone';

// Every die, the damage and the condition come from the server.
export function FallReceipt({ receipt }: { readonly receipt: Receipt }) {
  return <section aria-label="Saved fall receipt" className="space-y-2 border border-jade/20 p-3">
    <p>Fell {receipt.height_feet} feet onto {receipt.landing === 'water' ? 'water' : 'solid ground'}.</p>
    {receipt.water_check && <CheckReceipt receipt={receipt.water_check} />}
    <p>{receipt.dice_count ? `${receipt.dice_count}d6: ${receipt.dice.join(', ')} = ${receipt.rolled}${receipt.halved ? `, halved to ${receipt.damage}` : ''}.` : 'Under 10 feet: no dice.'}</p>
    <p>{receipt.damage} {receipt.damage_type} damage. HP {receipt.hp_before} → {receipt.hp_after}. Condition: {proneText(receipt.prone_before)} → {proneText(receipt.prone_after)}.</p>
    {receipt.concentration && <EffectReceipt receipt={receipt.concentration} />}
    {receipt.zero_hp && <ZeroHPReceipt receipt={receipt.zero_hp} />}
    <details><summary className="min-h-11 cursor-pointer py-2">Inspect fall sources</summary>{receipt.source_ids.map((id, index) => <p key={index} className="break-all text-pewter">{id}</p>)}</details>
  </section>;
}
