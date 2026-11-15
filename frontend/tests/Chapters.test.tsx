import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, expect, it, vi } from 'vitest';
import { ChapterSection } from '../src/components/journal/ChapterSection';
import { generateChapter, listChapters } from '../src/api/journal';
import type { ChapterResponse, ChapterRevisionResponse } from '../src/api/types';

vi.mock('../src/api/journal', () => ({ generateChapter: vi.fn(), listChapters: vi.fn() }));
const revision: ChapterRevisionResponse = { id: 'revision', chapter_id: 'chapter', request_id: 'key', revision: 1, status: 'completed', content: 'The gate stayed shut.', error: '', prompt_version: 'faithful-chapter-v1', created_at: '2026-01-01T00:00:00Z' };
const chapter: ChapterResponse = { id: 'chapter', campaign_id: 'alpha', from_turn: 1, to_turn: 2, source_hash: 'hash', source_turns: [], knowledge_boundary: 'player-visible-transcript-through-to-turn-v1', rules_version: 'srd-5.2.1', rules_mode: 'strict', created_at: revision.created_at, revisions: [revision] };
beforeEach(() => { vi.resetAllMocks(); vi.mocked(listChapters).mockResolvedValue([]); });
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><ChapterSection campaignId="alpha" /></QueryClientProvider>);
}
it('retains the next request key after a lost response and shows the persisted chapter', async () => {
  vi.mocked(generateChapter).mockRejectedValueOnce(new Error('connection lost')).mockImplementationOnce(async () => {
    vi.mocked(listChapters).mockResolvedValue([chapter]); return revision;
  });
  mount();
  await userEvent.click(screen.getByRole('button', { name: 'Generate next chapter' }));
  expect(await screen.findByText(/Could not finish the request/)).toBeVisible();
  const first = vi.mocked(generateChapter).mock.calls[0][1];
  await userEvent.click(screen.getByRole('button', { name: 'Generate next chapter' }));
  expect(await screen.findByText('The gate stayed shut.')).toBeVisible();
  expect(vi.mocked(generateChapter).mock.calls[1][1]).toEqual(first);
  expect(screen.getByText('Turns 1 – 2')).toBeVisible();
  await userEvent.click(screen.getByRole('button', { name: 'Regenerate chapter' }));
  const regen = vi.mocked(generateChapter).mock.calls[2][1];
  expect(regen.chapter_id).toBe('chapter'); expect(regen.request_id).not.toBe(first.request_id);
});
it('shows durable failure and retries that revision instead of creating another', async () => {
  vi.mocked(listChapters).mockResolvedValue([{ ...chapter, revisions: [{ ...revision, status: 'failed', content: '', error: 'Generation failed. Play remains available.' }] }]);
  vi.mocked(generateChapter).mockResolvedValue(revision);
  mount();
  expect(await screen.findByText('Generation failed. Play remains available.')).toBeVisible();
  await userEvent.click(screen.getByRole('button', { name: 'Retry revision 1' }));
  expect(generateChapter).toHaveBeenCalledWith('alpha', { chapter_id: 'chapter', request_id: 'key' });
});
it('announces an empty range and disables duplicate clicks during generation', async () => {
  let resolve!: (value: { message: string }) => void;
  vi.mocked(generateChapter).mockImplementation(() => new Promise(done => { resolve = done; }));
  mount();
  await userEvent.click(screen.getByRole('button', { name: 'Generate next chapter' }));
  expect(screen.getByRole('button', { name: 'Generating chapter...' })).toBeDisabled();
  resolve({ message: 'no completed turns available' });
  expect(await screen.findByText('No completed turns available for a new chapter.')).toBeVisible();
  await waitFor(() => expect(generateChapter).toHaveBeenCalledTimes(1));
});
