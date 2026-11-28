import type { DeathSaves } from '../../api/types';

// deathSaveText states the saved 0-HP counts, Stable or dead.
export const deathSaveText = (saves: DeathSaves) => saves.dead ? 'Dead.' : saves.stable ? 'Stable at 0 HP; still Unconscious.' : `Death Saving Throws: ${saves.successes ?? 0} successes, ${saves.failures ?? 0} failures.`;
