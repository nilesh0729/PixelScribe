import type { SqlNullInt64, SqlNullInt32, SqlNullFloat64, SqlNullTime } from './common';

export interface PerformanceSummary {
    id: number;
    user_id: SqlNullInt64;
    dictation_id: SqlNullInt64;
    total_attempts: SqlNullInt32;
    best_accuracy: SqlNullFloat64;
    average_accuracy: SqlNullFloat64;
    average_time: SqlNullFloat64;
    last_attempt_at: SqlNullTime;
}
