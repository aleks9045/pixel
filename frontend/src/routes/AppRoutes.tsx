import { useAuth } from '../context/AuthContext';
import { LoginPage } from '../components/LoginPage';
import { PixelPage } from '../pages/PixelPage';
import { usePathname } from '../hooks/usePathname';
import { AdminPage } from '../pages/AdminPage';

export function AppRoutes() {
    const { user, loading } = useAuth();
    const { pathname, navigate } = usePathname();

    if (loading) {
        return <div className="loading-page">Загрузка...</div>;
    }

    if (!user) {
        return <LoginPage />;
    }

    if (pathname === '/admin') {
        if (!user.isAdmin) {
            navigate('/');
            return <PixelPage />;
        }

        return <AdminPage />;
    }

    return <PixelPage />;
}