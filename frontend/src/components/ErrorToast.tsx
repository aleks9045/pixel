import { useEffect } from 'react';

type Props = {
    message: string;
    onClose: () => void;
};

export function ErrorToast({ message, onClose }: Props) {
    useEffect(() => {
        const timeout = setTimeout(onClose, 3000);

        return () => clearTimeout(timeout);
    }, [message, onClose]);

    return (
        <div className="error-toast">
            <span>{message}</span>
            <button onClick={onClose}>×</button>
        </div>
    );
}