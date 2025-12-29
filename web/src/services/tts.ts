import api from '../lib/axios';

export const ttsService = {
    generateAudio: async (text: string): Promise<Blob> => {
        const response = await api.post('/tts/generate', { text }, {
            responseType: 'blob'
        });
        return response.data;
    }
};
