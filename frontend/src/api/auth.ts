import config from '../config';
import { User } from '../proto/user';

export async function getCurrentUser() {
    const response = await fetch(
        `${config.host}/users/me`, {
        credentials: 'include',
    }
    );

    if (!response.ok) {
        return null;
    }

    const buffer = await response.arrayBuffer();
    const message = User.decode(new Uint8Array(buffer));

    return message;

}

export function loginWithGithub() {
    window.location.href =
        `${config.host}/oauth/github/login`;
}

export async function logout(): Promise<boolean> {
    const response = await fetch(
        `${config.host}/auth/logout`,
        {
            method: 'POST',
            credentials: 'include',
        }
    );
    return response.ok;
}