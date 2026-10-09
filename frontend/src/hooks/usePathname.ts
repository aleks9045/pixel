import { useCallback, useEffect, useState } from 'react';

export function usePathname() {
    const [pathname, setPathname] = useState(window.location.pathname);

    useEffect(() => {
        const handleChange = () => {
            setPathname(window.location.pathname);
        };

        window.addEventListener('popstate', handleChange);

        return () => {
            window.removeEventListener('popstate', handleChange);
        };
    }, []);

    const navigate = useCallback((path: string) => {
        window.history.pushState(null, '', path);
        window.dispatchEvent(new PopStateEvent('popstate'));
    }, []);

    return { pathname, navigate };
}