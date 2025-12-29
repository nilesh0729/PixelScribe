

export interface PerformanceSummary {
    id: number;
    user_id: number;
    dictation_id: number;
    total_attempts: number;
    best_accuracy: number;
    average_accuracy: number;
    average_time: number;
    last_attempt_at: string;
}
