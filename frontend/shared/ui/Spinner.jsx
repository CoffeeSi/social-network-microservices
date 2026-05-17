export function Spinner({ className = '' }) {
  return <div className={`spinner ${className}`.trim()} role="status" aria-label="Loading" />;
}
