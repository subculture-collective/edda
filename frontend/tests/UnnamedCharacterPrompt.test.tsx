import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeAll, beforeEach, expect, it, vi } from 'vitest';
import { UnnamedCharacterPrompt } from '../src/components/character/UnnamedCharacterPrompt';
import type { CharacterResponse } from '../src/api/types';
import { character } from './fixtures/data';

// jsdom does not implement modal dialogs.
beforeAll(() => {
  HTMLDialogElement.prototype.showModal ??= function (this: HTMLDialogElement) { this.setAttribute('open', ''); };
  HTMLDialogElement.prototype.close ??= function (this: HTMLDialogElement) { this.removeAttribute('open'); };
});

const request = vi.fn<typeof fetch>();
beforeEach(() => {
  request.mockReset();
  vi.stubGlobal('fetch', request);
});

const unnamed: CharacterResponse = { ...character(), name: '', description: 'A scout', needs_name: true };

// null stands for a character that has not loaded.
function mount(value: CharacterResponse | null = unnamed) {
  const client = new QueryClient();
  const invalidate = vi.spyOn(client, 'invalidateQueries');
  const view = render(<QueryClientProvider client={client}><UnnamedCharacterPrompt campaignId="alpha" character={value ?? undefined} /></QueryClientProvider>);
  return { ...view, invalidate };
}

const dialog = () => screen.queryByRole('dialog', { name: 'Name your character' });
const save = async (name: string) => {
  const input = screen.getByRole('textbox', { name: 'Character name' });
  await userEvent.clear(input);
  if (name) await userEvent.type(input, name);
  await userEvent.click(screen.getByRole('button', { name: 'Save name' }));
};

it('asks only when the server flags the character, so an older API never prompts', () => {
  for (const value of [null, character(), { ...character(), name: '' }, { ...character(), name: '  ' }]) {
    const view = mount(value);
    expect(dialog()).toBeNull();
    expect(view.container).toBeEmptyDOMElement();
    view.unmount();
  }
});

it('refuses a blank name without a request, then saves the trimmed name and refreshes the sheet', async () => {
  request.mockResolvedValue(new Response(null, { status: 204 }));
  const { invalidate } = mount();
  expect(dialog()).toBeVisible();
  expect(screen.getByRole('textbox', { name: 'Character name' })).toHaveFocus();
  for (const blank of ['', '   ']) {
    await save(blank);
    expect(screen.getByRole('alert')).toHaveTextContent('character name is required');
  }
  expect(request).not.toHaveBeenCalled();

  await save('  Iona  ');

  await waitFor(() => expect(dialog()).toBeNull());
  expect(request).toHaveBeenCalledOnce();
  const [url, init] = request.mock.calls[0];
  expect(url).toBe('/api/v1/campaigns/alpha/character/name');
  expect(init?.method).toBe('PUT');
  expect(JSON.parse(String(init?.body))).toEqual({ character_id: 'hero', name: 'Iona' });
  expect(invalidate).toHaveBeenCalledWith({ queryKey: ['campaign', 'alpha', 'character'] });
});

it.each([404, 405])('closes with a dismissible notice when the server answers %i', async status => {
  request.mockResolvedValue(new Response('404 page not found', { status }));
  mount();
  await save('Iona');
  expect(await screen.findByRole('status')).toHaveTextContent('The name was not saved. This server cannot rename the character right now. You can keep playing.');
  expect(dialog()).toBeNull();
  await userEvent.click(screen.getByRole('button', { name: 'Dismiss notice' }));
  expect(screen.queryByRole('status')).toBeNull();
  expect(dialog()).toBeNull();
});

it('keeps the typed name after another failure and lets the player skip', async () => {
  request.mockResolvedValue(Response.json({ error: 'failed to rename character' }, { status: 500 }));
  mount();
  await save('Iona');
  expect(await screen.findByRole('alert')).toHaveTextContent('failed to rename character');
  expect(screen.getByRole('textbox', { name: 'Character name' })).toHaveValue('Iona');
  expect(screen.getByRole('button', { name: 'Save name' })).toBeEnabled();
  await userEvent.click(screen.getByRole('button', { name: 'Not now' }));
  expect(dialog()).toBeNull();
});

it('says which choices the broken guided form replaced', () => {
  const note = /also replaced your race, class, background and alignment/;
  const other = mount();
  expect(screen.queryByText(note)).toBeNull();
  expect(screen.getByText('Saved as: A scout')).toBeVisible();
  other.unmount();
  mount({ ...unnamed, description: 'human barbarian. Acolyte' });
  expect(screen.getByText(note)).toBeVisible();
});
