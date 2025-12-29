import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { attemptService, type AttemptResponse } from '../services/attempt';
import { dictationService } from '../services/dictation';
import type { Dictation } from '../types/dictation';
import { Calendar, Clock, CheckCircle, ArrowRight, Activity } from 'lucide-react';

export default function AttemptHistory() {
    const navigate = useNavigate();
    const [attempts, setAttempts] = useState<AttemptResponse[]>([]);
    const [dictations, setDictations] = useState<Record<number, Dictation>>({});
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            try {
                const [attemptsData, dictationsData] = await Promise.all([
                    attemptService.getAll(),
                    dictationService.getAll()
                ]);

                // Create a map of dictations for easy lookup
                const dictMap: Record<number, Dictation> = {};
                dictationsData.forEach(d => dictMap[d.id] = d);
                setDictations(dictMap);

                setAttempts(attemptsData);
            } catch (error) {
                console.error("Failed to load history", error);
            } finally {
                setLoading(false);
            }
        };
        loadData();
    }, []);

    if (loading) return <div className="p-8 text-center text-gray-500">Loading history...</div>;

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold text-gray-900 flex items-center">
                <Activity className="mr-3 h-8 w-8 text-indigo-600" />
                Attempt History
            </h1>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Dictation</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Accuracy</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Time</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {attempts.map((attempt) => (
                            <tr key={attempt.id} className="hover:bg-gray-50 transition">
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    <div className="flex items-center">
                                        <Calendar className="h-4 w-4 mr-2 text-gray-400" />
                                        {new Date(attempt.created_at).toLocaleDateString()}
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                    {dictations[attempt.dictation_id]?.title || `Dictation #${attempt.dictation_id}`}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    <div className="flex items-center">
                                        <CheckCircle className={`h-4 w-4 mr-2 ${attempt.accuracy >= 90 ? 'text-green-500' : 'text-yellow-500'}`} />
                                        <span className={attempt.accuracy >= 90 ? 'text-green-700 font-semibold' : 'text-gray-700'}>
                                            {attempt.accuracy.toFixed(1)}%
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    <div className="flex items-center">
                                        <Clock className="h-4 w-4 mr-2 text-gray-400" />
                                        {Math.round(attempt.time_spent)}s
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                    <button
                                        onClick={() => navigate(`/attempts/${attempt.id}`)}
                                        className="text-indigo-600 hover:text-indigo-900 flex items-center justify-end"
                                    >
                                        View Details <ArrowRight className="h-4 w-4 ml-1" />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {attempts.length === 0 && (
                    <div className="p-8 text-center text-gray-500">
                        No attempts found. Start practicing!
                    </div>
                )}
            </div>
        </div>
    );
}
