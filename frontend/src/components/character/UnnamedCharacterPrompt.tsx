import { useId, useLayoutEffect, useRef, useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import { renameCampaignCharacter } from '../../api/characters';
import { APIError } from '../../api/client';
import type { CharacterResponse } from '../../api/types';

// The description the broken TUI guided form saved while it ignored the player's choices.
const guidedFormDefaults = 'human barbarian. Acolyte';

interface UnnamedCharacterPromptProps {
  readonly campaignId: string;
  readonly character: CharacterResponse | undefined;
}

// Asks the player to name a saved character that has none. It never blocks play:
// the player can skip it, and it asks again the next time the campaign loads.
// The web client deploys ahead of the API, so it asks only when the server sets
// needs_name, which a server without the rename route never does.
export function UnnamedCharacterPrompt({ campaignId, character }: UnnamedCharacterPromptProps) {
  const [outcome, setOutcome] = useState<'asking' | 'closed' | 'unsupported'>('asking');
  const queryClient = useQueryClient();
  const rename = useMutation({
    mutationFn: ({ characterId, name }: { characterId: string; name: string }) =>
      renameCampaignCharacter(campaignId, { character_id: characterId, name }),
    onSuccess: () => {
      setOutcome('closed');
      void queryClient.invalidateQueries({ queryKey: ['campaign', campaignId, 'character'] });
    },
    onError: (error) => {
      // The API may have been rolled back since it flagged the character.
      if (error instanceof APIError && (error.status === 404 || error.status === 405)) setOutcome('unsupported');
    },
  });

  if (outcome === 'unsupported') {
    return (
      <div className="flex shrink-0 flex-wrap items-center justify-between gap-3 border-2 border-gold/30 bg-charcoal px-4 py-2">
        <p role="status" className="text-sm text-champagne">The name was not saved. This server cannot rename the character right now. You can keep playing.</p>
        <button type="button" className="hud-btn hud-btn-secondary min-h-11" onClick={() => setOutcome('closed')}>Dismiss notice</button>
      </div>
    );
  }

  if (outcome === 'closed' || !character?.needs_name) return null;

  return (
    <NameCharacterDialog
      description={character.description}
      saving={rename.isPending}
      saveError={rename.error?.message ?? null}
      onSave={(name) => rename.mutate({ characterId: character.id, name })}
      onSkip={() => setOutcome('closed')}
    />
  );
}

function NameCharacterDialog({ description, saving, saveError, onSave, onSkip }: {
  readonly description: string;
  readonly saving: boolean;
  readonly saveError: string | null;
  readonly onSave: (name: string) => void;
  readonly onSkip: () => void;
}) {
  const id = useId();
  const dialogRef = useRef<HTMLDialogElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [name, setName] = useState('');
  const [blank, setBlank] = useState(false);
  const error = blank ? 'character name is required' : saveError;

  useLayoutEffect(() => {
    const dialog = dialogRef.current;
    dialog?.showModal();
    inputRef.current?.focus();
    // Closing before removal returns focus to where it was.
    return () => dialog?.close();
  }, []);

  return (
    <dialog
      ref={dialogRef}
      aria-labelledby={`${id}-title`}
      aria-describedby={`${id}-about`}
      onCancel={(event) => { event.preventDefault(); onSkip(); }}
      className="m-auto max-h-[90dvh] w-[calc(100vw-2rem)] max-w-md overflow-y-auto border-2 border-gold/30 bg-charcoal p-0 text-champagne backdrop:bg-obsidian/80"
    >
      <form
        noValidate
        onSubmit={(event) => {
          event.preventDefault();
          if (saving) return;
          const trimmed = name.trim();
          setBlank(!trimmed);
          if (trimmed) onSave(trimmed);
        }}
      >
        <div className="space-y-4 p-6">
          <h2 id={`${id}-title`} className="font-heading text-lg font-semibold uppercase tracking-wide">Name your character</h2>
          <div id={`${id}-about`} className="space-y-2 text-sm leading-6 text-pewter">
            <p>This saved character has no name.</p>
            {description === guidedFormDefaults ? (
              <p>The guided form that saved it also replaced your race, class, background and alignment with Human, Barbarian, Acolyte and Lawful Good. Naming does not restore your choices.</p>
            ) : null}
            {description.trim() ? <p className="line-clamp-3 break-words">Saved as: {description}</p> : null}
          </div>
          <label className="block space-y-2 text-sm font-semibold uppercase tracking-wide">
            <span>Character name</span>
            <input
              ref={inputRef}
              type="text"
              value={name}
              onChange={(event) => { setName(event.target.value); setBlank(false); }}
              aria-invalid={blank}
              aria-describedby={error ? `${id}-error` : undefined}
              autoComplete="off"
              className="min-h-11 w-full border-2 border-gold/20 bg-obsidian px-4 py-2 text-base normal-case tracking-normal text-champagne focus:border-gold focus:outline-none focus:ring-2 focus:ring-gold/40"
            />
          </label>
          {error ? <p id={`${id}-error`} role="alert" className="text-sm text-ruby">{error}</p> : null}
        </div>
        <div className="flex justify-end gap-3 border-t border-gold/15 px-6 py-4">
          <button type="button" className="hud-btn hud-btn-secondary min-h-11" onClick={onSkip}>Not now</button>
          <button type="submit" className="hud-btn hud-btn-primary min-h-11" disabled={saving}>{saving ? 'Saving…' : 'Save name'}</button>
        </div>
      </form>
    </dialog>
  );
}
