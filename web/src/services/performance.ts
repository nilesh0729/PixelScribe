import api from '../lib/axios';
import type { PerformanceSummary } from '../types/performance';

export const performanceService = {
    // Get all performance summaries for the current user
    getPerformance: async (userId: number | string) => {
        const response = await api.get<PerformanceSummary[]>('/performance', {
            params: { user_id: userId }
        });
        return response.data;
    },

    // Get recent attempts (feed)
    // Note: Backend endpoint for recent is /performance/recent
    getRecentActivity: async (userId: number | string, limit = 5) => {
        const response = await api.get<PerformanceSummary[]>('/performance/recent', {
            params: { user_id: userId, limit }
        });
        return response.data;
    }
};
