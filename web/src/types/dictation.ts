import type { SqlNullInt64, SqlNullString } from './common';

export interface Dictation {
    id: number;
    user_id: SqlNullInt64;
    title: SqlNullString;
    type: SqlNullString;
    content: SqlNullString;
    audio_url: SqlNullString;
    language: SqlNullString;
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
