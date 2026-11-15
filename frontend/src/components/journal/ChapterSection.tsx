import { useRef, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { generateChapter, listChapters, type GenerateChapterRequest } from '../../api/journal';

export function ChapterSection({ campaignId }: { readonly campaignId: string }) {
  const client = useQueryClient();
  // A transport failure must retain the original key, even if more turns commit.
  const pendingRequests = useRef(new Map<string, GenerateChapterRequest>());
  const [message, setMessage] = useState('');
  const queryKey = ['campaign', campaignId, 'journal', 'chapters'];
  const chapters = useQuery({
    queryKey,
    queryFn: () => listChapters(campaignId),
    refetchInterval: query => query.state.data?.some(c => c.revisions.some(r => r.status === 'generating')) ? 2000 : false,
  });
  const mutation = useMutation({
    mutationFn: (request: GenerateChapterRequest) => generateChapter(campaignId, request),
    onSuccess: (result, request) => {
      // A server response makes the request durable; failed/in-progress
      // revisions expose their key through the list, including after reload.
      for (const [key, pending] of pendingRequests.current) {
        if (pending.request_id === request.request_id) pendingRequests.current.delete(key);
      }
      setMessage('message' in result ? 'No completed turns available for a new chapter.' : result.status === 'generating' ? 'Chapter generation is in progress.' : result.status === 'failed' ? 'Generation failed. You can retry; play remains available.' : 'Chapter saved.');
      void client.invalidateQueries({ queryKey: ['campaign', campaignId, 'journal', 'chapters'] });
    },
    onError: () => {
      setMessage('Could not finish the request. Retry uses the same request key; play remains available.');
      void client.invalidateQueries({ queryKey: ['campaign', campaignId, 'journal', 'chapters'] });
    },
  });
  function generate(chapterId?: string, requestId?: string) {
    const key = `${campaignId}:${chapterId ?? 'next'}`;
    const request = requestId ? { request_id: requestId, chapter_id: chapterId } : pendingRequests.current.get(key) ?? { request_id: crypto.randomUUID(), chapter_id: chapterId };
    pendingRequests.current.set(key, request);
    setMessage('');
    mutation.mutate(request);
  }
  return (
    <section className="space-y-3" aria-label="Chapters">
      <p className="text-sm text-pewter">Create a chapter from the next completed turns. Chapters use only the recorded player-visible story.</p>
      <button type="button" className="hud-btn hud-text-button" disabled={mutation.isPending} onClick={() => generate()}>
        {mutation.isPending ? 'Generating chapter...' : 'Generate next chapter'}
      </button>
      {message && <p role="status" className="text-sm text-pewter">{message}</p>}
      {chapters.isPending && <p className="text-sm text-pewter">Loading chapters...</p>}
      {chapters.isError && <p role="alert">Could not load chapters. <button type="button" className="hud-btn hud-text-button" onClick={() => void chapters.refetch()}>Retry loading chapters</button></p>}
      {chapters.data?.length === 0 && <p className="text-sm text-pewter">No chapters yet.</p>}
      {chapters.data?.map(chapter => (
        <article key={chapter.id} className="space-y-3 border-2 border-midnight/40 bg-obsidian p-5">
          <h3 className="font-heading text-sm text-gold">Turns {chapter.from_turn} – {chapter.to_turn}</h3>
          <p className="text-xs text-pewter">{chapter.rules_version} · Player knowledge through turn {chapter.to_turn}</p>
          {chapter.revisions.map(revision => (
            <div key={revision.id} className="space-y-2">
              <p className="text-xs text-pewter">Revision {revision.revision} · {revision.status}</p>
              {revision.content && <p className="whitespace-pre-wrap text-sm leading-7 text-champagne/80">{revision.content}</p>}
              {revision.error && <p className="text-sm text-ruby">{revision.error}</p>}
              {revision.status !== 'completed' && <button type="button" className="hud-btn hud-text-button" disabled={mutation.isPending} onClick={() => generate(chapter.id, revision.request_id)}>Retry revision {revision.revision}</button>}
            </div>
          ))}
          <button type="button" className="hud-btn hud-text-button" disabled={mutation.isPending} onClick={() => generate(chapter.id)}>Regenerate chapter</button>
        </article>
      ))}
    </section>
  );
}
