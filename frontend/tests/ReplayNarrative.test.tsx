import { act, render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { ReplayNarrative } from '../src/components/replay/ReplayNarrative';
import type { SessionLogEntry } from '../src/api/types';

const response = `${'The road bends north. '.repeat(10)}The gate opens.`;
const entry = { turn_number: 1, player_input: '', llm_response: response } as SessionLogEntry;

it('finishes typing the latest turn after the replay page rerenders', () => {
  vi.useFakeTimers();
  // The replay engine passes a new array on every render.
  const view = render(<ReplayNarrative visibleEntries={[entry]} currentTurnIndex={0} isPlaying />);
  act(() => { vi.advanceTimersByTime(66); });
  expect(screen.queryByText(response)).not.toBeInTheDocument();
  view.rerender(<ReplayNarrative visibleEntries={[entry]} currentTurnIndex={0} isPlaying={false} />);
  act(() => { vi.advanceTimersByTime(1000); });
  expect(screen.getByText(response)).toBeInTheDocument();
});
