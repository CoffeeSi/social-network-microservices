import { clsx } from '../lib/clsx';

export function Card({ title, children, className = '' }) {
  return (
    <section className={clsx('card', className)}>
      {title ? <h2 className="card__title">{title}</h2> : null}
      {children}
    </section>
  );
}
