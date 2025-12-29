import { createContext, useContext, useState, type ReactNode } from 'react';
import api from '../lib/axios';
import type { User, AuthResponse, LoginRequest, RegisterRequest } from '../types/auth';

interface AuthContextType {
    user: User | null;
    token: string | null;
    login: (data: LoginRequest) => Promise<void>;
    register: (data: RegisterRequest) => Promise<void>;
    logout: () => void;
    isLoading: boolean;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(() => {
        try {
            const storedUser = localStorage.getItem('user');
            return storedUser ? JSON.parse(storedUser) : null;
        } catch {
            return null;
        }
    });

    const [token, setToken] = useState<string | null>(() => {
        return localStorage.getItem('token');
    });

    const [isLoading, setIsLoading] = useState(false);

    const login = async (data: LoginRequest) => {
        setIsLoading(true);
        try {
            const response = await api.post<AuthResponse>('/users/login', data);
            const { access_token, user: userData } = response.data;

            setToken(access_token);
            setUser(userData);

            localStorage.setItem('token', access_token);
            localStorage.setItem('user', JSON.stringify(userData));
        } finally {
            setIsLoading(false);
        }
    };

    const register = async (data: RegisterRequest) => {
        setIsLoading(true);
        try {
            // First create user
            await api.post<User>('/users', data);
            // Then auto-login
            await login({ username: data.username, password: data.password });
        } finally {
            setIsLoading(false);
        }
    };

    const logout = () => {
        setToken(null);
        setUser(null);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                token,
                login,
                register,
                logout,
                isLoading,
                isAuthenticated: !!token,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
