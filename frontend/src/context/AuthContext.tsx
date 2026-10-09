import {
    createContext,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from 'react';

import { getCurrentUser } from '../api/auth';
import type { User } from '../proto/user';

type AuthState = {
    user: User | null;
    loading: boolean;
    updateUser: (user: User | null) => void;
};

const AuthContext = createContext<AuthState>({
    user: null,
    loading: true,
    updateUser: () => { },
});

type AuthProviderProps = {
    children: ReactNode;
};

export function AuthProvider({
    children,
}: AuthProviderProps) {
    const [user, setUser] =
        useState<User | null>(null);

    const [loading, setLoading] =
        useState(true);

    useEffect(() => {
        getCurrentUser()
            .then(setUser)
            .finally(() => setLoading(false));
    }, []);

    const updateUser = (user: User | null) => {
        setUser(user);
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                loading,
                updateUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}