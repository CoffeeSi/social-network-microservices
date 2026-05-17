'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import { clsx } from '@/shared/lib/clsx';

export function NavBar() {
  const pathname = usePathname();
  const { isAuthenticated, signOut, userId } = useAuth();

  function linkClass(href) {
    if (href === routes.users) {
      return clsx('nav__link', pathname === routes.users && 'nav__link--active');
    }
    if (href === routes.profile) {
      const active =
        pathname === routes.profile || Boolean(userId && pathname === routes.userProfile(userId));
      return clsx('nav__link', active && 'nav__link--active');
    }
    const active = pathname === href || (href !== routes.home && pathname.startsWith(href));
    return clsx('nav__link', active && 'nav__link--active');
  }

  return (
    <header className="nav">
      <Link href={routes.home} className="nav__brand">
        Social
      </Link>
      <nav className="nav__links">
        {isAuthenticated ? (
          <>
            <Link href={routes.feed} className={linkClass(routes.feed)}>
              Feed
            </Link>
            <Link href={routes.users} className={linkClass(routes.users)}>
              People
            </Link>
            <Link href={routes.profile} className={linkClass(routes.profile)}>
              Profile
            </Link>
            <Link href={routes.chat} className={linkClass(routes.chat)}>
              Chat
            </Link>
            <button type="button" className="nav__link nav__link--btn" onClick={signOut}>
              Sign out
            </button>
          </>
        ) : (
          <>
            <Link href={routes.login} className={linkClass(routes.login)}>
              Sign in
            </Link>
            <Link href={routes.register} className={linkClass(routes.register)}>
              Register
            </Link>
          </>
        )}
      </nav>
    </header>
  );
}
