import api from '../lib/axios';


export interface AttemptRequest {
    dictation_id: number;
    typed_text: string;
    total_words: number;
    correct_words: number;
    // We can add errors breakdown later if we implement detailed diffing
    accuracy: number;
    time_spent: number; // in seconds
}

export const attemptService = {
    submit: async (data: AttemptRequest) => {
        const response = await api.post('/attempts', data);
        return response.data;
    }
};
