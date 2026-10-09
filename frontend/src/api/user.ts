import config from '../config';
import {
    User,
    Users
} from '../proto/user';

export async function getAllUsers(
    page: number,
    size: number,
): Promise<User[]> {
    const response = await fetch(
        `${config.host}/users/all?page=${page}&size=${size}`,
        {
            credentials: 'include',
        },
    );

    if (!response.ok) {
        throw new Error('Не удалось загрузить пользователей');
    }

    const data = new Uint8Array(await response.arrayBuffer());

    return Users.decode(data).user;
}

export async function updateUser(
    user: User,
): Promise<boolean> {
    const response = await fetch(
        `${config.host}/users/${user.userId}`,
        {
            method: 'PATCH',
            credentials: 'include',
            headers: {
                'Content-Type':
                    'application/octet-stream',
            },
            body: User.encode(user).finish(),
        },
    );

    return response.ok;
}
