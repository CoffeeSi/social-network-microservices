import { clsx } from '../lib/clsx';

export function Button({ children, variant = 'primary', className = '', ...props }) {
  return (
    <button
      type="button"
      className={clsx('btn', variant === 'secondary' && 'btn--secondary', className)}
      {...props}
    >
      {children}
    </button>
  );
}
