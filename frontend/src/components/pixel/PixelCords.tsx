import { useEffect } from 'react';

type Props = {
    x: number;
    y: number;
};

export function PixelCords({
    x,
    y,
}: Props) {

    useEffect(() => {
    }, [x, y]);


    return (
        <div className="cords">
            x:{x} y:{y}
        </div>
    );
}