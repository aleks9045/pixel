import { PixelCanvas } from '../components/pixel/PixelCanvas';
import { UserMenu } from '../components/user/UserMenu';

export function PixelPage() {
    return <>
        <UserMenu />
        <PixelCanvas />
    </>;
}
