import type { User } from "../../proto/user";

type UserCardProps = {
    user: User;
};

export function UserCard({
    user,
}: UserCardProps) {
    return (
        <div className="user-card">
            <div className="user-card-header">
                {user.avatarUrl ? (
                    <img
                        src={user.avatarUrl}
                        alt={user.username}
                    />
                ) : (
                    <span className="user-card-avatar">
                        👤
                    </span>
                )}

                <div className="user-card-name">
                    {user.name}
                </div>
            </div>

            <div className="user-card-info">
                <div className="user-card-row">
                    <span>Username</span>
                    <span>{user.username}</span>
                </div>

                <div className="user-card-row">
                    <span>Group</span>
                    <span>
                        {user.studentGroup || '—'}
                    </span>
                </div>

                <div className="user-card-row">
                    <span>Time limit</span>
                    <span>
                        {user.timeLimit} ms
                    </span>
                </div>
            </div>
        </div>
    );
}