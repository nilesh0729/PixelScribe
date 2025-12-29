import api from '../lib/axios';


export interface AttemptRequest {
    dictation_id: number;
    typed_text: string;
    time_spent: number; // in seconds
}

export interface AttemptResponse {
    id: number;
    user_id: number;
    dictation_id: number;
    typed_text: string;
    attempt_no: number;
    accuracy: number;
    time_spent: number;
    created_at: string;
    performance_update?: {
        total_attempts: number;
        best_accuracy: number;
        average_accuracy: number;
        average_time: number;
    };
}

export const attemptService = {
    submit: async (data: AttemptRequest): Promise<AttemptResponse> => {
        const response = await api.post<AttemptResponse>('/attempts', data);
        return response.data;
    },

    getAll: async () => {
        const response = await api.get<AttemptResponse[]>('/attempts');
        return response.data;
    },

    getById: async (id: number) => {
        const response = await api.get<AttemptResponse>(`/attempts/${id}`);
        return response.data;
    }
};
