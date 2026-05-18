"use client";

import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { verifyUser } from '@/features/auth/api/authApi';

export default function ConfirmEmailClient() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [code, setCode] = useState('');
    const [status, setStatus] = useState('idle');
    const [message, setMessage] = useState('');
    const [pending, setPending] = useState(false);

    useEffect(() => {
        const emailParam = searchParams.get('email');

        if (emailParam) {
            setEmail(emailParam);
            setMessage(`Введите код подтверждения, отправленный на ${emailParam}`);
        } else {
            setMessage('Email не найден.');
            setStatus('error');
        }
    }, [searchParams]);

    async function verifyCode(emailVal, codeVal) {
        setPending(true);
        try {
            const res = await verifyUser({ email: emailVal, token: codeVal });
            if (res && res.ok === false) {
                setStatus('error');
                setMessage(res.message || 'Ошибка подтверждения.');
                return;
            }
            setStatus('success');
            setMessage('Email успешно подтверждён! Перенаправление на вход...');
            setTimeout(() => router.push('/login'), 2500);
        } catch (err) {
            setStatus('error');
            setMessage(err?.message || 'Серверная ошибка при подтверждении.');
        } finally {
            setPending(false);
        }
    }

    async function onSubmit(e) {
        e.preventDefault();
        if (!code.trim()) {
            setMessage('Введите код.');
            return;
        }
        await verifyCode(email, code);
    }

    return (
        <div className="page page--narrow">
            <div className="card">
                <h1 className="card__title">Подтверждение Email</h1>
                {message && <p className={`form__error`} style={{ marginBottom: '1rem' }}>{message}</p>}

                {status === 'success' && (
                    <p>Если вы не были перенаправлены, <a href="/login">войдите</a>.</p>
                )}

                {status !== 'success' && (
                    <form className="form" onSubmit={onSubmit}>
                        <label className="field">
                            <span className="field__label">Код подтверждения</span>
                            <input
                                type="text"
                                className="input"
                                value={code}
                                onChange={(e) => setCode(e.target.value)}
                                placeholder="Введите код из письма"
                                disabled={pending}
                                required
                            />
                        </label>
                        <button
                            type="submit"
                            className="btn"
                            disabled={pending || !code.trim()}
                        >
                            {pending ? 'Проверка...' : 'Подтвердить'}
                        </button>
                    </form>
                )}
            </div>
        </div>
    );
}
