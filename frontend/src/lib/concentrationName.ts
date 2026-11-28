import type { Concentration } from '../api/types';

// The same words srd521.ConcentrationLabel keeps lowercase.
const minorWords = ['a', 'an', 'and', 'for', 'from', 'into', 'of', 'on', 'or', 'the', 'to', 'via', 'with', 'without'];

// The server's spell name. An older API sends none, so the ID is shown in title case.
export function concentrationName({ spell, spell_label }: Concentration): string {
  return spell_label || spell.split('-').map((word, index) => index > 0 && minorWords.includes(word) ? word : word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
}
