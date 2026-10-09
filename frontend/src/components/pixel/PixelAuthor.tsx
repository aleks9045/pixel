import {
    useEffect,
    useState,
} from 'react';

import type { User } from '../../proto/user';
import { getAuthor } from '../../api/pixel';
import { UserCard } from '../user/UserCard';

type PixelAuthorProps = {
    x: number;
    y: number;
    left: number;
    top: number;
};

export function PixelAuthor({
    x,
    y,
    left,
    top,
}: PixelAuthorProps) {
    const [user, setUser] =
        useState<User | null>(null);

    useEffect(() => {
        getAuthor(x, y)
            .then(setUser);
    }, [x, y]);

    if (!user) {
        return null;
    }

    return (
        <div
            className="pixel-author"
            style={{
                left,
                top,
            }}
        >
            <UserCard user={user} />
        </div>
    );
}