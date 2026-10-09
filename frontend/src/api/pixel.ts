import { PixelField } from '../proto/pixel';
import config from '../config';
import { User } from '../proto/user';

export async function getPixelField() {

    const response = await fetch(`${config.host}/pixel/field`, {
        credentials: 'include',
    });

    if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
    }

    const buffer = await response.arrayBuffer();
    const message = PixelField.decode(new Uint8Array(buffer));

    return message;
}


export async function getAuthor(
    x: number,
    y: number
) {
    const response = await fetch(
        `${config.host}/pixel/author?x=${x}&y=${y}`,
        {
            method: 'GET',
            credentials: 'include',

        },
    );
    if (!response.ok) {
        return null;
    }

    const buffer = await response.arrayBuffer();
    const message = User.decode(new Uint8Array(buffer));

    return message;

}