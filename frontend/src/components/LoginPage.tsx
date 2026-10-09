
import config from '../config';

export function LoginPage() {
    const handleLogin = () => {
        window.location.href =
            `${config.host}/oauth/github/login`;
    };

    return (
        <div className="login-page">
            <div className="login-card">
                <h1>Pixels</h1>

                <p>
                    Войдите через GitHub, чтобы продолжить
                </p>

                <button
                    className="github-login-button"
                    onClick={handleLogin}
                >
                    Войти через GitHub
                </button>
            </div>
        </div>
    );
}
