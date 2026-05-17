import { clsx } from '../lib/clsx';

export function Input({ label, id, className = '', ...props }) {
  const inputId = id ?? props.name;
  return (
    <label className="field">
      {label ? <span className="field__label">{label}</span> : null}
      <input id={inputId} className={clsx('input', className)} {...props} />
    </label>
  );
}
