import { useEffect, useState } from 'react';

type Props = {
    time: number;
    trigger: number;
    onChange: (canPlace: boolean) => void;
};

export function PixelTimeout({ time, trigger, onChange }: Props) {
    const [remaining, setRemaining] = useState(0);

    useEffect(() => {
        if (trigger === 0) {
            setRemaining(0);
            return;
        }

        if (time <= 0) {
            setRemaining(0);
            onChange(true);
            return;
        }

        const endTime = Date.now() + time;
        setRemaining(time);
        onChange(false);

        const updateRemaining = () => {
            const nextRemaining = Math.max(0, endTime - Date.now());
            setRemaining(nextRemaining);

            if (nextRemaining === 0) {
                clearInterval(interval);
                onChange(true);
            }
        };

        const interval = setInterval(updateRemaining, 100);
        updateRemaining();

        return () => clearInterval(interval);
    }, [trigger, time, onChange]);

    if (trigger === 0 || remaining <= 0) {
        return null;
    }

    return (
        <div className="timeout">
            {(remaining / 1000).toFixed(1)} s
        </div>
    );
}