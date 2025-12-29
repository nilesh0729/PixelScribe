import { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { dictationService } from '../services/dictation';
import { attemptService } from '../services/attempt';
import type { AttemptResponse } from '../services/attempt';
import { useDictationEngine } from '../hooks/useDictationEngine';
import type { Dictation } from '../types/dictation';
import { Play, Pause, RotateCcw, Check, ArrowLeft, Ear, Keyboard, Loader2 } from 'lucide-react';

export default function PlayDictation() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [dictation, setDictation] = useState<Dictation | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [result, setResult] = useState<AttemptResponse | null>(null);

    // Load Dictation Data
    useEffect(() => {
        if (!id) return;
        const fetchDictation = async () => {
            try {
                // Keep the filter logic for now as getById is not yet implemented in service
                const all = await dictationService.getAll();
                const found = all.find(d => d.id === Number(id));
                if (found) setDictation(found);
                else alert('Dictation not found');
            } catch (error) {
                console.error("Failed to load dictation", error);
            } finally {
                setLoading(false);
            }
        };
        fetchDictation();
    }, [id]);

    const targetContent = dictation?.content || '';

    const {
        phase,
        currentText,
        audioState,
        controls,
        stats
    } = useDictationEngine(targetContent);

    const inputRef = useRef<HTMLTextAreaElement>(null);

    // Focus input when separate typing phase starts
    useEffect(() => {
        if (phase === 'typing' && inputRef.current) {
            setTimeout(() => inputRef.current?.focus(), 100);
        }
    }, [phase]);

    const handleSave = async () => {
        if (!dictation) return;
        setSubmitting(true);
        try {
            const response = await attemptService.submit({
                dictation_id: dictation.id,
                typed_text: currentText,
                time_spent: stats.timeSpent,
            });
            setResult(response);
            controls.complete(); // Mark engine as completed
        } catch (error) {
            console.error("Failed to submit attempt", error);
            alert("Failed to save result.");
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading dictation...</div>;
    if (!dictation) return <div className="p-8 text-center text-gray-500">Dictation not found.</div>;

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <button onClick={() => navigate('/dictations')} className="flex items-center text-gray-500 hover:text-gray-700">
                    <ArrowLeft className="h-4 w-4 mr-1" /> Back to Library
                </button>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${phase === 'listening' ? 'bg-blue-100 text-blue-800' :
                    phase === 'typing' ? 'bg-yellow-100 text-yellow-800' :
                        phase === 'completed' ? 'bg-green-100 text-green-800' :
                            'bg-gray-100 text-gray-800'
                    }`}>
                    {phase.toUpperCase()}
                </span>
            </div>

            <div className="bg-white p-8 rounded-lg shadow-lg min-h-[400px] flex flex-col items-center justify-center text-center space-y-8">

                {/* LISTENING PHASE */}
                {phase === 'listening' || phase === 'loading' ? (
                    <div className="space-y-6 w-full max-w-md">
                        <div className="mx-auto bg-blue-50 p-6 rounded-full w-24 h-24 flex items-center justify-center">
                            {phase === 'loading' ? <Loader2 className="h-10 w-10 text-blue-500 animate-spin" /> :
                                <Ear className="h-10 w-10 text-blue-600" />}
                        </div>

                        <div>
                            <h2 className="text-2xl font-bold text-gray-900">Listen Carefully</h2>
                            <p className="text-gray-500 mt-2">Listen to the audio segment. When you are ready, click 'Start Typing'.</p>
                        </div>

                        <div className="flex justify-center space-x-4">
                            {audioState.isPlaying ? (
                                <button onClick={audioState.pause} className="flex items-center px-6 py-3 bg-yellow-500 text-white rounded-full shadow-lg hover:bg-yellow-600 transition transform hover:scale-105">
                                    <Pause className="h-6 w-6 mr-2" /> Pause
                                </button>
                            ) : (
                                <button onClick={audioState.play} disabled={phase === 'loading'} className="flex items-center px-6 py-3 bg-indigo-600 text-white rounded-full shadow-lg hover:bg-indigo-700 transition transform hover:scale-105 disabled:opacity-50">
                                    <Play className="h-6 w-6 mr-2" /> Play Audio
                                </button>
                            )}

                            <button onClick={audioState.replay} disabled={phase === 'loading'} className="flex items-center px-4 py-3 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 transition">
                                <RotateCcw className="h-5 w-5" />
                            </button>
                        </div>

                        <div className="pt-8 border-t border-gray-100">
                            <button onClick={controls.startTyping} disabled={phase === 'loading'} className="w-full flex items-center justify-center px-6 py-4 border-2 border-indigo-600 text-indigo-700 rounded-lg text-lg font-bold hover:bg-indigo-50 transition">
                                <Keyboard className="h-6 w-6 mr-2" /> I'm Ready to Type
                            </button>
                        </div>
                    </div>
                ) : null}

                {/* TYPING PHASE */}
                {phase === 'typing' && (
                    <div className="w-full space-y-6">
                        <div className="flex items-center justify-between text-sm text-gray-400">
                            <span><Keyboard className="inline h-4 w-4 mr-1" /> Typing Mode</span>
                            <span>Time: {Math.round(stats.timeSpent)}s</span>
                        </div>

                        <textarea
                            ref={inputRef}
                            value={currentText}
                            onChange={(e) => controls.handleInput(e.target.value)}
                            className="w-full h-80 p-6 text-lg text-gray-900 border-2 border-indigo-100 rounded-xl focus:border-indigo-500 focus:ring-0 resize-none font-mono leading-relaxed shadow-inner bg-gray-50 placeholder:text-gray-400"
                            placeholder="Type what you heard here..."
                        />

                        <div className="space-y-4">
                            <button
                                onClick={handleSave}
                                disabled={submitting || currentText.trim().length === 0}
                                className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {submitting ? <Loader2 className="h-5 w-5 mr-2 animate-spin" /> : <Check className="h-5 w-5 mr-2" />}
                                Save Attempt
                            </button>
                            {result && (
                                <button
                                    onClick={() => navigate(`/attempts/${result.id}`)}
                                    className="w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                                >
                                    View Analysis
                                </button>
                            )}
                        </div>
                    </div>
                )}
            </div>

            {/* Result Modal */}
            {result && (
                <div className="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50">
                    <div className="bg-white p-8 rounded-2xl shadow-2xl max-w-md w-full text-center animate-in fade-in zoom-in duration-200">
                        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100 mb-6">
                            <Check className="h-8 w-8 text-green-600" />
                        </div>

                        <h2 className="text-3xl font-bold text-gray-900 mb-2">Attempt Complete!</h2>
                        <p className="text-gray-500 mb-8">Here is how you performed.</p>

                        <div className="grid grid-cols-2 gap-4 mb-8">
                            <div className="bg-gray-50 p-4 rounded-xl">
                                <p className="text-sm text-gray-500 uppercase font-semibold">Accuracy</p>
                                <p className="text-3xl font-bold text-indigo-600">{result.accuracy.toFixed(1)}%</p>
                            </div>
                            <div className="bg-gray-50 p-4 rounded-xl">
                                <p className="text-sm text-gray-500 uppercase font-semibold">Time</p>
                                <p className="text-3xl font-bold text-gray-600">{Math.round(result.time_spent)}s</p>
                            </div>
                        </div>

                        <div className="flex flex-col space-y-3">
                            <button
                                onClick={() => navigate(`/attempts/${result.id}`)}
                                className="w-full px-4 py-3 bg-indigo-600 text-white rounded-lg font-bold hover:bg-indigo-700 transition"
                            >
                                View Analysis
                            </button>
                            <button
                                onClick={() => navigate('/dictations')}
                                className="w-full px-4 py-3 bg-white border border-gray-300 text-gray-700 rounded-lg font-bold hover:bg-gray-50 transition"
                            >
                                Back to Library
                            </button>
                            {/* Retry logic would need page reload or state reset */}
                            <button
                                onClick={() => window.location.reload()}
                                className="w-full px-4 py-3 text-indigo-600 hover:bg-indigo-50 rounded-lg font-medium transition"
                            >
                                Try Again
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
