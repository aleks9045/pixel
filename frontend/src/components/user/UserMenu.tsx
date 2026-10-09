import {
    useState,
} from 'react';
import { useAuth } from '../../context/AuthContext';
import type { User } from '../../proto/user';
import { updateUser } from '../../api/user';
import { logout } from '../../api/auth';

export function UserMenu() {
    const {
        user,
        loading,
        updateUser: setUser,
    } = useAuth();

    const [
        isOpen,
        setIsOpen,
    ] = useState(false);

    const [
        isEditing,
        setIsEditing,
    ] = useState(false);

    const [
        name,
        setName,
    ] = useState('');

    const [
        studentGroup,
        setStudentGroup,
    ] = useState('');

    if (loading || !user) {
        return null;
    }

    const startEditing = () => {
        setName(user.name);
        setStudentGroup(user.studentGroup);
        setIsEditing(true);
    };

    const cancelEditing = () => {
        setName(user.name);
        setStudentGroup(user.studentGroup);
        setIsEditing(false);
    };

    const saveUser = async () => {
        const updatedUser: User = {
            ...user,
            name,
            studentGroup,
        };

        const success =
            await updateUser(updatedUser);

        if (!success) {
            return;
        }

        setUser(updatedUser);
        setIsEditing(false);
    };

    const handleLogout = async () => {
        const success = await logout();

        if (!success) {
            return;
        }

        setIsOpen(false);
        setUser(null);
    };

    return (
        <div className="user-menu">
            <button
                className="user-menu-button"
                onClick={() =>
                    setIsOpen(!isOpen)
                }
            >
                {user.avatarUrl ? (
                    <img
                        src={user.avatarUrl}
                        alt={user.username}
                    />
                ) : (
                    <span>👤</span>
                )}
            </button>

            {isOpen && (
                <div className="user-menu-dropdown">
                    <div className="user-menu-header">
                        {user.avatarUrl ? (
                            <img
                                src={user.avatarUrl}
                                alt={user.username}
                            />
                        ) : (
                            <span>👤</span>
                        )}

                        <div>
                            {isEditing ? (
                                <input
                                    className="user-menu-input"
                                    value={name}
                                    onChange={e =>
                                        setName(
                                            e.target.value
                                        )
                                    }
                                />
                            ) : (
                                <div>
                                    {user.name}
                                </div>
                            )}
                        </div>
                    </div>

                    <div className="user-menu-info">
                        <div>
                            <span>Username</span>
                            <span>
                                {user.username}
                            </span>
                        </div>

                        <div>
                            <span>Group</span>

                            {isEditing ? (
                                <input
                                    className="user-menu-input"
                                    value={studentGroup}
                                    onChange={e =>
                                        setStudentGroup(
                                            e.target.value
                                        )
                                    }
                                />
                            ) : (
                                <span>
                                    {user.studentGroup || '—'}
                                </span>
                            )}
                        </div>

                        <div>
                            <span>Time limit</span>
                            <span>
                                {user.timeLimit} ms
                            </span>
                        </div>
                    </div>

                    <div className="user-menu-actions">
                        {isEditing ? (
                            <>
                                <button
                                    className="user-menu-save"
                                    onClick={saveUser}
                                >
                                    Save
                                </button>

                                <button
                                    className="user-menu-cancel"
                                    onClick={cancelEditing}
                                >
                                    Cancel
                                </button>

                                <button
                                    className="user-menu-logout"
                                    onClick={handleLogout}
                                >
                                    Logout
                                </button>
                            </>
                        ) : (
                            <button
                                className="user-menu-edit"
                                onClick={startEditing}
                            >
                                Edit
                            </button>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}