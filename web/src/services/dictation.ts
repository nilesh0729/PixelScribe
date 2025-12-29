import api from '../lib/axios';
import type { Dictation, CreateDictationRequest } from '../types/dictation';

export const dictationService = {
    // Get all dictations (with optional pagination/filtering later)
    getAll: async () => {
        const response = await api.get<Dictation[]>('/dictations');
        return response.data;
    },

    // Create a new dictation
    create: async (data: CreateDictationRequest) => {
        const response = await api.post<Dictation>('/dictations', data);
        return response.data;
    },

    // Delete a dictation
    delete: async (id: number) => {
        await api.delete(`/dictations/${id}`);
    }
};
