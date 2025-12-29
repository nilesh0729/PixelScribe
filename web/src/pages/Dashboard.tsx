import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { performanceService } from '../services/performance';
import type { PerformanceSummary } from '../types/performance';
import { unwrap } from '../types/common';
import { Activity, Trophy, Clock, Target } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function Dashboard() {
    const { user } = useAuth();
    const [summaries, setSummaries] = useState<PerformanceSummary[]>([]);
    const [recent, setRecent] = useState<PerformanceSummary[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (user?.id) {
            const fetchData = async () => {
                try {
                    const [perfData, recentData] = await Promise.all([
                        performanceService.getPerformance(user.id),
                        performanceService.getRecentActivity(user.id, 5)
                    ]);
                    setSummaries(perfData);
                    setRecent(recentData);
                } catch (error) {
                    console.error("Failed to fetch dashboard data", error);
                } finally {
                    setLoading(false);
                }
            };
            fetchData();
        }
    }, [user]);

    // Calculate generic stats
    const totalAttempts = summaries.reduce((acc, curr) => acc + unwrap(curr.total_attempts, 'Int32', 0), 0);

    const totalAccuracy = summaries.reduce((acc, curr) => acc + unwrap(curr.average_accuracy, 'Float64', 0), 0);
    const avgAccuracy = summaries.length > 0 ? totalAccuracy / summaries.length : 0;

    const bestAccuracyEver = summaries.reduce((max, curr) => Math.max(max, unwrap(curr.best_accuracy, 'Float64', 0)), 0);

    return (
        <div className="space-y-6">
            {/* Header */}
            <div>
                <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Overview of your dictation performance.
                </p>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
                <div className="bg-white overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <Activity className="h-6 w-6 text-indigo-600" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Total Attempts</dt>
                                    <dd className="text-lg font-medium text-gray-900">{totalAttempts}</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="bg-white overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <Target className="h-6 w-6 text-green-600" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Avg. Accuracy</dt>
                                    <dd className="text-lg font-medium text-gray-900">{avgAccuracy.toFixed(1)}%</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="bg-white overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <Trophy className="h-6 w-6 text-yellow-500" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Best Accuracy</dt>
                                    <dd className="text-lg font-medium text-gray-900">{bestAccuracyEver.toFixed(1)}%</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Recent Activity */}
            <div className="bg-white shadow rounded-lg">
                <div className="px-4 py-5 sm:px-6 flex justify-between items-center">
                    <h3 className="text-lg leading-6 font-medium text-gray-900">Recent Activity</h3>
                    <Link to="/dictations" className="text-sm text-indigo-600 hover:text-indigo-500">Start new dictation &rarr;</Link>
                </div>
                <div className="border-t border-gray-200">
                    {loading ? (
                        <div className="p-4 text-center text-gray-500">Loading activity...</div>
                    ) : recent.length === 0 ? (
                        <div className="p-4 text-center text-gray-500">
                            No recent activity. <Link to="/dictations" className="text-indigo-600">Start your first dictation!</Link>
                        </div>
                    ) : (
                        <ul className="divide-y divide-gray-200">
                            {recent.map((item) => (
                                <li key={item.id} className="px-4 py-4 sm:px-6 hover:bg-gray-50">
                                    <div className="flex items-center justify-between">
                                        <div className="flex flex-col">
                                            <p className="text-sm font-medium text-indigo-600 truncate">
                                                Dictation #{unwrap(item.dictation_id, 'Int64', 0)}
                                            </p>
                                            <p className="flex items-center text-sm text-gray-500 mt-1">
                                                <Clock className="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" />
                                                {unwrap(item.last_attempt_at, 'Time', '') ? new Date(unwrap(item.last_attempt_at, 'Time', '')).toLocaleDateString() : 'Unknown date'}
                                            </p>
                                        </div>
                                        <div className="flex items-center">
                                            <div className="flex flex-col items-end">
                                                <p className="text-sm font-medium text-gray-900">
                                                    {unwrap(item.average_accuracy, 'Float64', 0).toFixed(1)}% Acc
                                                </p>
                                                <p className="text-xs text-gray-500 mt-1">
                                                    {unwrap(item.total_attempts, 'Int32', 0)} attempts
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
            </div>
        </div>
    );
}
