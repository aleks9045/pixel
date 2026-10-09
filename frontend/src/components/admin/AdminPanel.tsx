import { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { usePathname } from '../../hooks/usePathname';
import type { User } from '../../proto/user';
import { getAllUsers, updateUser } from '../../api/user';

const PAGE_SIZE = 10;

export function AdminPanel() {
    const { user } = useAuth();
    const { navigate } = usePathname();

    const [users, setUsers] = useState<User[]>([]);
    const [page, setPage] = useState(0);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [name, setName] = useState('');
    const [studentGroup, setStudentGroup] = useState('');
    const [timeLimit, setTimeLimit] = useState(0);
    const [saving, setSaving] = useState(false);
    const [success, setSuccess] = useState('');

    useEffect(() => {
        if (!user?.isAdmin) {
            setLoading(false);
            return;
        }

        let cancelled = false;

        setLoading(true);
        setError('');
        setSelectedUserId(null);
        setIsEditing(false);

        getAllUsers(page, PAGE_SIZE)
            .then(data => {
                if (!cancelled) {
                    setUsers(data);
                }
            })
            .catch(error => {
                if (!cancelled) {
                    setError(error instanceof Error ? error.message : 'Ошибка загрузки');
                }
            })
            .finally(() => {
                if (!cancelled) {
                    setLoading(false);
                }
            });

        return () => {
            cancelled = true;
        };
    }, [page, user?.isAdmin]);

    const selectedUser = users.find(item => item.userId === selectedUserId) ?? null;

    const selectUser = (item: User) => {
        setSelectedUserId(item.userId);
        setName(item.name);
        setStudentGroup(item.studentGroup);
        setIsEditing(false);
        setTimeLimit(item.timeLimit);
        setError('');
        setSuccess('');
    };

    const cancelEditing = () => {
        if (!selectedUser) return;

        setName(selectedUser.name);
        setStudentGroup(selectedUser.studentGroup);
        setTimeLimit(selectedUser.timeLimit);
        setIsEditing(false);
    };

    const saveUser = async () => {
        if (!selectedUser) return;

        setSaving(true);
        setError('');
        setSuccess('');

        const updatedUser: User = {
            ...selectedUser,
            name,
            studentGroup,
            timeLimit,
        };

        try {
            const saved = await updateUser(updatedUser);

            if (!saved) {
                throw new Error('Не удалось сохранить пользователя');
            }

            setUsers(current =>
                current.map(item =>
                    item.userId === updatedUser.userId ? updatedUser : item,
                ),
            );

            setIsEditing(false);
            setSuccess('Изменения сохранены');
        } catch (error) {
            setError(error instanceof Error ? error.message : 'Ошибка сохранения');
        } finally {
            setSaving(false);
        }
    };

    if (!user?.isAdmin) {
        return (
            <div className="admin-page">
                <div className="admin-empty">
                    <h2>Доступ запрещён</h2>
                    <button onClick={() => navigate('/')}>На главную</button>
                </div>
            </div>
        );
    }

    return (
        <main className="admin-page">
            <header className="admin-topbar">
                <button
                    className="admin-back"
                    onClick={() => navigate('/')}
                >
                    <span>←</span>
                    На главную
                </button>

                <div className="admin-heading">
                    <h1>Пользователи</h1>
                    <p>Управление профилями пользователей</p>
                </div>

                <div className="admin-page-number">
                    Страница {page + 1}
                </div>
            </header>

            {error && <div className="admin-error">{error}</div>}

            <div className="admin-layout">
                <section className="admin-users">
                    <div className="admin-section-header">
                        <div>
                            <h2>Все пользователи</h2>
                            <span>{users.length} на этой странице</span>
                        </div>

                        <span className="admin-page-badge">
                            {String(page + 1).padStart(2, '0')}
                        </span>
                    </div>

                    {loading ? (
                        <div className="admin-loading">Загрузка пользователей...</div>
                    ) : users.length === 0 ? (
                        <div className="admin-empty">Пользователей нет</div>
                    ) : (
                        <div className="admin-user-list">
                            {users.map(item => (
                                <button
                                    key={item.userId}
                                    className={`admin-user-card ${selectedUserId === item.userId ? 'selected' : ''
                                        }`}
                                    onClick={() => selectUser(item)}
                                >
                                    <div className="admin-avatar">
                                        {item.avatarUrl ? (
                                            <img src={item.avatarUrl} alt="" />
                                        ) : (
                                            <span>👤</span>
                                        )}
                                    </div>

                                    <div className="admin-user-summary">
                                        <span className="admin-user-name">
                                            {item.name || item.username}
                                        </span>
                                        <span className="admin-username">
                                            @{item.username}
                                        </span>
                                        <span className="admin-user-group">
                                            {item.studentGroup || 'Без группы'}
                                        </span>
                                    </div>

                                    {item.isAdmin && (
                                        <span className="admin-role">ADMIN</span>
                                    )}

                                    <span className="admin-card-arrow">›</span>
                                </button>
                            ))}
                        </div>
                    )}

                    <div className="admin-pagination">
                        <button
                            disabled={page === 0 || loading}
                            onClick={() => setPage(value => value - 1)}
                        >
                            ← Назад
                        </button>

                        <span>{page + 1}</span>

                        <button
                            disabled={loading || users.length < PAGE_SIZE}
                            onClick={() => setPage(value => value + 1)}
                        >
                            Далее →
                        </button>
                    </div>
                </section>

                <aside className="admin-details">
                    {!selectedUser ? (
                        <div className="admin-details-empty">
                            <span className="admin-details-icon">◎</span>
                            <h2>Профиль пользователя</h2>
                            <p>Выбери пользователя из списка, чтобы посмотреть информацию и изменить профиль.</p>
                        </div>
                    ) : (
                        <>
                            <div className="admin-profile-header">
                                <div className="admin-profile-avatar">
                                    {selectedUser.avatarUrl ? (
                                        <img src={selectedUser.avatarUrl} alt="" />
                                    ) : (
                                        <span>👤</span>
                                    )}
                                </div>

                                <div className="admin-profile-title">
                                    <h2>{selectedUser.name || selectedUser.username}</h2>
                                    <span>@{selectedUser.username}</span>
                                    <span className={`admin-status ${selectedUser.isAdmin ? 'admin-status-admin' : ''}`}>
                                        {selectedUser.isAdmin ? 'Администратор' : 'Пользователь'}
                                    </span>
                                </div>
                            </div>

                            <div className="admin-detail-fields">
                                <label>
                                    <span>Имя</span>
                                    {isEditing ? (
                                        <input
                                            value={name}
                                            onChange={event => setName(event.target.value)}
                                            maxLength={100}
                                        />
                                    ) : (
                                        <div className="admin-field-value">{selectedUser.name || '—'}</div>
                                    )}
                                </label>

                                <label>
                                    <span>Username</span>
                                    <div className="admin-field-value">@{selectedUser.username}</div>
                                </label>

                                <label>
                                    <span>Группа</span>
                                    {isEditing ? (
                                        <input
                                            value={studentGroup}
                                            onChange={event => setStudentGroup(event.target.value)}
                                            maxLength={100}
                                        />
                                    ) : (
                                        <div className="admin-field-value">{selectedUser.studentGroup || '—'}</div>
                                    )}
                                </label>

                                <label>
                                    <span>Лимит размещения пикселя (мс)</span>
                                    {isEditing ? (
                                        <input
                                            type="string"
                                            min={0}
                                            step={1}
                                            value={timeLimit}
                                            onChange={event =>
                                                setTimeLimit(
                                                    Math.max(0, Math.floor(Number(event.target.value) || 0))
                                                )
                                            }
                                        />
                                    ) : (
                                        <div className="admin-field-value">
                                            {selectedUser.timeLimit} ms
                                        </div>
                                    )}
                                </label>

                                <label>
                                    <span>ID пользователя</span>
                                    <div className="admin-field-value admin-id">
                                        {selectedUser.userId}
                                    </div>
                                </label>
                            </div>

                            {success && <div className="admin-success">{success}</div>}

                            <div className="admin-detail-actions">
                                {isEditing ? (
                                    <>
                                        <button
                                            className="admin-save"
                                            disabled={saving}
                                            onClick={saveUser}
                                        >
                                            {saving ? 'Сохранение...' : 'Сохранить'}
                                        </button>

                                        <button
                                            className="admin-cancel"
                                            disabled={saving}
                                            onClick={cancelEditing}
                                        >
                                            Отмена
                                        </button>
                                    </>
                                ) : (
                                    <button
                                        className="admin-save"
                                        onClick={() => {
                                            setName(selectedUser.name);
                                            setStudentGroup(selectedUser.studentGroup);
                                            setIsEditing(true);
                                            setTimeLimit(selectedUser.timeLimit);
                                            setSuccess('');
                                            setError('');
                                        }}
                                    >
                                        Редактировать профиль
                                    </button>
                                )}
                            </div>
                        </>
                    )}
                </aside>
            </div>
        </main>
    );
}