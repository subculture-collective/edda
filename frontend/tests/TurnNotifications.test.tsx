import { act, render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { TurnNotifications } from '../src/components/layout/TurnNotifications';
import type { StateChange } from '../src/api/types';

const met: StateChange[] = [{ entity_type: 'npc', change_type: 'created', entity_id: 'n1', details: { name: 'Mara' } }];

it('still dismisses a notification when the next action clears the turn result first', () => {
  vi.useFakeTimers();
  const view = render(<TurnNotifications stateChanges={met} />);
  expect(screen.getByText('New NPC met: Mara')).toBeInTheDocument();
  act(() => { vi.advanceTimersByTime(1000); });
  view.rerender(<TurnNotifications stateChanges={[]} />);
  expect(screen.getByText('New NPC met: Mara')).toBeInTheDocument();
  act(() => { vi.advanceTimersByTime(3000); });
  expect(screen.queryByText('New NPC met: Mara')).not.toBeInTheDocument();
});
