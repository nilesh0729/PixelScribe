import { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { dictationService } from '../services/dictation';
import { attemptService } from '../services/attempt';
import { useDictationEngine } from '../hooks/useDictationEngine';
import type { Dictation } from '../types/dictation';
import { unwrap } from '../types/common';
import { Play, Pause, RotateCcw, Save, ArrowLeft } from 'lucide-react';

export default function PlayDictation() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [dictation, setDictation] = useState<Dictation | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

    // Load Dictation Data
    useEffect(() => {
        if (!id) return;
        const fetchDictation = async () => {
            try {
                // We need a getById method in dictationService which we missed. 
                // We can either add it or filter from getAll (less efficient but works for now).
                // Let's implement getById in service properly next step. 
                // For now, I'll filter from getAll just to unblock the UI build.
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

    const targetContent = unwrap(dictation?.content, 'String', '');
    const language = unwrap(dictation?.language, 'String', 'en-US');

    const { status, currentText, start, pause, stop, handleInput, stats } = useDictationEngine(
        targetContent,
        language
    );

    const inputRef = useRef<HTMLTextAreaElement>(null);

    // Focus input on start
    useEffect(() => {
        if (status === 'playing' && inputRef.current) {
            inputRef.current.focus();
        }
    }, [status]);

    const handleSave = async () => {
        if (!dictation) return;
        setSubmitting(true);
        try {
            await attemptService.submit({
                dictation_id: dictation.id,
                typed_text: currentText,
                total_words: targetContent.split(' ').length,
                correct_words: Math.floor((targetContent.split(' ').length * stats.accuracy) / 100), // Approximate
                accuracy: stats.accuracy,
                time_spent: 0, // Need to implement time tracking properly in hook
            });
            navigate('/');
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
            <div className="flex items-center justify-between">
                <button onClick={() => navigate('/dictations')} className="flex items-center text-gray-500 hover:text-gray-700">
                    <ArrowLeft className="h-4 w-4 mr-1" /> Back
                </button>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${status === 'playing' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                    {status.toUpperCase()}
                </span>
            </div>

            {/* Header Stats */}
            <div className="grid grid-cols-3 gap-4 bg-white p-4 rounded-lg shadow">
                <div className="text-center">
                    <p className="text-xs text-gray-500 uppercase">WPM</p>
                    <p className="text-2xl font-bold text-indigo-600">{stats.wpm}</p>
                </div>
                <div className="text-center">
                    <p className="text-xs text-gray-500 uppercase">Accuracy</p>
                    <p className="text-2xl font-bold text-green-600">{stats.accuracy}%</p>
                </div>
                <div className="text-center">
                    <p className="text-xs text-gray-500 uppercase">Progress</p>
                    <p className="text-2xl font-bold text-blue-600">{Math.round(stats.progress)}%</p>
                </div>
            </div>

            {/* Game Area */}
            <div className="space-y-4">
                {/* Target Text (Blurred or Visible based on design choices - usually blurred or read only) 
            For Dictation, we might want to hide it or show it. 
            PixelScribe implies "Scribing" what you hear. Let's hide it by default or show it faded?
            Let's show it for now as reference, maybe allow toggling later.
        */}
                <div className="bg-gray-50 p-4 rounded-md border text-gray-400 select-none">
                    {targetContent}
                </div>

                <textarea
                    ref={inputRef}
                    value={currentText}
                    onChange={(e) => handleInput(e.target.value)}
                    disabled={status === 'completed'}
                    className="w-full h-64 p-4 text-lg border-2 border-indigo-200 rounded-lg focus:border-indigo-500 focus:ring-0 resize-none font-mono leading-relaxed"
                    placeholder="Press play and type what you hear..."
                />
            </div>

            {/* Controls */}
            <div className="flex justify-center space-x-4">
                {status === 'playing' ? (
                    <button onClick={pause} className="flex items-center px-6 py-3 bg-yellow-500 text-white rounded-full shadow-lg hover:bg-yellow-600 transition">
                        <Pause className="h-6 w-6 mr-2" /> Pause
                    </button>
                ) : (
                    <button onClick={start} className="flex items-center px-6 py-3 bg-green-600 text-white rounded-full shadow-lg hover:bg-green-700 transition">
                        <Play className="h-6 w-6 mr-2" /> {status === 'idle' ? 'Start' : 'Resume'}
                    </button>
                )}

                <button onClick={stop} className="flex items-center px-6 py-3 bg-gray-600 text-white rounded-full shadow-lg hover:bg-gray-700 transition">
                    <RotateCcw className="h-6 w-6 mr-2" /> Reset
                </button>
            </div>

            {/* Result Modal / Overlay */}
            {status === 'completed' && (
                <div className="fixed inset-0 bg-gray-900/50 flex items-center justify-center z-50">
                    <div className="bg-white p-8 rounded-lg shadow-2xl max-w-md w-full text-center">
                        <h2 className="text-3xl font-bold text-gray-900 mb-4">Good Job! 🎉</h2>
                        <div className="grid grid-cols-2 gap-6 mb-8">
                            <div>
                                <p className="text-gray-500">Speed</p>
                                <p className="text-3xl font-bold text-indigo-600">{stats.wpm} WPM</p>
                            </div>
                            <div>
                                <p className="text-gray-500">Accuracy</p>
                                <p className="text-3xl font-bold text-green-600">{stats.accuracy}%</p>
                            </div>
                        </div>
                        <div className="flex flex-col space-y-3">
                            <button
                                onClick={handleSave}
                                disabled={submitting}
                                className="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700"
                            >
                                <Save className="h-5 w-5 mr-2" /> {submitting ? 'Saving...' : 'Save Result'}
                            </button>
                            <button
                                onClick={() => stop()}
                                className="w-full px-4 py-3 text-gray-700 hover:bg-gray-100 rounded-md"
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
