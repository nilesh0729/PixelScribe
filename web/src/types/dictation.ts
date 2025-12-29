

export interface Dictation {
    id: number;
    user_id: number;
    title: string;
    type: string;
    content: string;
    audio_url: string;
    language: string;
    created_at: string;
    updated_at: string;
}

export interface CreateDictationRequest {
    title: string;
    type?: string;
    content: string;
    audio_url?: string;
    language?: string;
}
