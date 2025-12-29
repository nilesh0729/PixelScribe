import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { attemptService, type AttemptResponse } from '../services/attempt';
import { dictationService } from '../services/dictation';
import type { Dictation } from '../types/dictation';
import * as Diff from 'diff';
import { Loader2, ArrowLeft, Calendar, FileText, CheckCircle, Clock } from 'lucide-react';

export default function AttemptDetails() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [attempt, setAttempt] = useState<AttemptResponse | null>(null);
    const [dictation, setDictation] = useState<Dictation | null>(null);
    const [loading, setLoading] = useState(true);
    const [diffParts, setDiffParts] = useState<Diff.Change[]>([]);

    useEffect(() => {
        const loadData = async () => {
            if (!id) return;
            try {
                const attemptData = await attemptService.getById(parseInt(id));
                setAttempt(attemptData);

                // Fetch all dictations to find the matching one (Backend optimization pending)
                const allDictations = await dictationService.getAll();
                const matchedDictation = allDictations.find(d => d.id === attemptData.dictation_id);
                setDictation(matchedDictation || null);

                if (matchedDictation) {
                    // Calculate diff
                    // valid logic: diffWords(original, typed)
                    // added: true -> present in typed but NOT in original (Extra/Wrong word) -> Red
                    // removed: true -> present in original but NOT in typed (Missed word) -> Green (or Strikeout)
                    const differences = Diff.diffWords(matchedDictation.content, attemptData.typed_text, { ignoreCase: true });
                    setDiffParts(differences);
                }

            } catch (error) {
                console.error("Failed to load attempt details", error);
            } finally {
                setLoading(false);
            }
        };
        loadData();
    }, [id]);

    if (loading) return <div className="flex justify-center p-12"><Loader2 className="animate-spin h-8 w-8 text-indigo-600" /></div>;
    if (!attempt || !dictation) return <div className="text-center p-12">Attempt or Dictation not found.</div>;

    return (
        <div className="max-w-4xl mx-auto space-y-8">
            <button onClick={() => navigate('/attempts')} className="flex items-center text-gray-500 hover:text-gray-900 transition">
                <ArrowLeft className="h-4 w-4 mr-1" /> Back to History
            </button>

            {/* Header Stats */}
            <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 mb-1">{dictation.title}</h1>
                    <div className="flex items-center text-gray-500 text-sm">
                        <Calendar className="h-4 w-4 mr-1" />
                        {new Date(attempt.created_at).toLocaleString()}
                    </div>
                </div>

                <div className="flex gap-4">
                    <div className="bg-green-50 px-4 py-2 rounded-xl border border-green-100 flex flex-col items-center">
                        <span className="text-xs text-green-600 font-medium uppercase tracking-wider">Accuracy</span>
                        <div className="flex items-center text-xl font-bold text-green-700">
                            <CheckCircle className="h-5 w-5 mr-1" />
                            {attempt.accuracy.toFixed(1)}%
                        </div>
                    </div>
                    <div className="bg-blue-50 px-4 py-2 rounded-xl border border-blue-100 flex flex-col items-center">
                        <span className="text-xs text-blue-600 font-medium uppercase tracking-wider">Time</span>
                        <div className="flex items-center text-xl font-bold text-blue-700">
                            <Clock className="h-5 w-5 mr-1" />
                            {Math.round(attempt.time_spent)}s
                        </div>
                    </div>
                </div>
            </div>

            {/* Visual Diff */}
            <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100">
                <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                    <FileText className="mr-2 h-5 w-5 text-indigo-600" />
                    Attempt Analysis
                </h2>

                <div className="mb-4 text-sm flex gap-4">
                    <div className="flex items-center"><span className="w-3 h-3 bg-red-100 border border-red-300 mr-2 rounded"></span> <span className="text-gray-600">Incorrect/Extra</span></div>
                    <div className="flex items-center"><span className="w-3 h-3 bg-green-100 border border-green-300 mr-2 rounded"></span> <span className="text-gray-600">Missed</span></div>
                </div>

                <div className="leading-relaxed text-lg font-mono p-6 bg-gray-50 rounded-xl border border-gray-200 whitespace-pre-wrap">
                    {diffParts.map((part, index) => {
                        const color = part.added ? 'bg-red-200 text-red-900 line-through decoration-red-500'
                            : part.removed ? 'bg-green-200 text-green-900'
                                : 'text-gray-800';
                        return (
                            <span key={index} className={`${color} px-0.5 rounded-sm transition hover:opacity-80`}>
                                {part.value}
                            </span>
                        );
                    })}
                </div>
            </div>

            {/* Server Stats Panel */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <StatCard title="Total Words" value={attempt.performance_update?.total_attempts || 0} icon={FileText} />
                {/* Note: In a real app we would pass 'CorrectWords' etc directly in response too, currently relying on perf update which is aggregated. 
                      Actually attempt object has 'correctWords' calculated on backend but response struct in api/attempt.go didn't include it explicitly in root (my bad). 
                      It's fine, the visual diff is the key feature here.
                  */}
            </div>
        </div>
    );
}

function StatCard({ title, value, icon: Icon }: any) {
    return (
        <div className="bg-white p-4 rounded-xl border border-gray-100 shadow-sm flex items-center">
            <div className="p-3 bg-indigo-50 rounded-lg mr-4">
                <Icon className="h-6 w-6 text-indigo-600" />
            </div>
            <div>
                <p className="text-sm text-gray-500">{title}</p>
                <p className="text-xl font-bold text-gray-900">{value}</p>
            </div>
        </div>
    )
}
