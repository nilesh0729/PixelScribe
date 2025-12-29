export interface User {
    id: number;
    username: string;
    full_name: string;
    email: string;
    password_changed_at: string;
    created_at: string;
}

export interface AuthResponse {
    access_token: string;
    user: User;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface RegisterRequest {
    username: string;
    password: string;
    full_name: string;
    email: string;
}
