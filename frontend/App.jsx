'use client';

import { NavBar } from '@/components/NavBar';

export function App({ children }) {
  return (
    <div className="shell">
      <NavBar />
      <main className="shell__main">{children}</main>
    </div>
  );
}
