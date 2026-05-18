'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { routes } from '@/app/router/routes';
import { registerUser } from '@/features/auth/api/authApi';
import { Button } from '@/shared/ui/Button';
import { Input } from '@/shared/ui/Input';

export function RegisterForm() {
  const router = useRouter();
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [dob, setDob] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    setError('');
    setPending(true);
    try {
      await registerUser({
        first_name: firstName,
        last_name: lastName,
        dob: dob || '1970-01-01',
        email,
        password,
      });
      router.replace(`${routes.confirmEmail}?email=${encodeURIComponent(email)}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setPending(false);
    }
  }

  return (
    <form className="form" onSubmit={onSubmit}>
      {error ? <p className="form__error">{error}</p> : null}
      <div className="form__row">
        <Input
          label="First name"
          name="first_name"
          value={firstName}
          onChange={(ev) => setFirstName(ev.target.value)}
          required
        />
        <Input
          label="Last name"
          name="last_name"
          value={lastName}
          onChange={(ev) => setLastName(ev.target.value)}
          required
        />
      </div>
      <Input label="Date of birth" name="dob" type="date" value={dob} onChange={(ev) => setDob(ev.target.value)} />
      <Input
        label="Email"
        name="email"
        type="email"
        value={email}
        onChange={(ev) => setEmail(ev.target.value)}
        required
      />
      <Input
        label="Password"
        name="password"
        type="password"
        value={password}
        onChange={(ev) => setPassword(ev.target.value)}
        required
        minLength={8}
      />
      <Button type="submit" disabled={pending}>
        {pending ? 'Creating account…' : 'Create account'}
      </Button>
    </form>
  );
}
