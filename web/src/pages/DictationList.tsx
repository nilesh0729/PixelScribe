import { useEffect, useState } from 'react';
import { dictationService } from '../services/dictation';
import type { Dictation } from '../types/dictation';
import { unwrap } from '../types/common';
import { Plus, Trash2, FileText, Play } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function DictationList() {
    const [dictations, setDictations] = useState<Dictation[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadDictations();
    }, []);

    const loadDictations = async () => {
        try {
            const data = await dictationService.getAll();
            setDictations(data);
        } catch (error) {
            console.error('Failed to load dictations', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm('Are you sure you want to delete this dictation?')) return;

        try {
            await dictationService.delete(id);
            setDictations(dictations.filter(d => d.id !== id));
        } catch (error) {
            console.error('Failed to delete dictation', error);
            alert('Failed to delete dictation');
        }
    };

    if (loading) {
        return <div className="p-8 text-center text-gray-500">Loading dictations...</div>;
    }

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Dictation Library</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        Choose a text to practice or create a new one.
                    </p>
                </div>
                <Link
                    to="/dictations/new"
                    className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                    <Plus className="h-4 w-4 mr-2" />
                    New Dictation
                </Link>
            </div>

            {dictations.length === 0 ? (
                <div className="text-center py-12 bg-white rounded-lg shadow">
                    <FileText className="mx-auto h-12 w-12 text-gray-400" />
                    <h3 className="mt-2 text-sm font-medium text-gray-900">No dictations</h3>
                    <p className="mt-1 text-sm text-gray-500">Get started by creating a new dictation.</p>
                    <div className="mt-6">
                        <Link
                            to="/dictations/new"
                            className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
                        >
                            <Plus className="h-4 w-4 mr-2" />
                            New Dictation
                        </Link>
                    </div>
                </div>
            ) : (
                <div className="bg-white shadow overflow-hidden sm:rounded-md">
                    <ul className="divide-y divide-gray-200">
                        {dictations.map((dictation) => (
                            <li key={dictation.id}>
                                <div className="px-4 py-4 flex items-center sm:px-6">
                                    <div className="min-w-0 flex-1 sm:flex sm:items-center sm:justify-between">
                                        <div className="truncate">
                                            <div className="flex text-sm">
                                                <p className="font-medium text-indigo-600 truncate">{unwrap(dictation.title, 'String', 'Untitled')}</p>
                                                <p className="ml-1 flex-shrink-0 font-normal text-gray-500">
                                                    in {unwrap(dictation.type, 'String', 'General')}
                                                </p>
                                            </div>
                                            <div className="mt-2 flex">
                                                <div className="flex items-center text-sm text-gray-500">
                                                    <p className="truncate">
                                                        {unwrap(dictation.content, 'String', '').substring(0, 100)}...
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="ml-5 flex-shrink-0 flex items-center space-x-2">
                                        <Link
                                            to={`/play/${dictation.id}`}
                                            className="inline-flex items-center p-2 border border-transparent rounded-full shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                                            title="Start Practice"
                                        >
                                            <Play className="h-4 w-4" />
                                        </Link>
                                        <button
                                            onClick={() => handleDelete(dictation.id)}
                                            className="inline-flex items-center p-2 border border-transparent rounded-full text-gray-400 hover:text-red-500 focus:outline-none"
                                            title="Delete"
                                        >
                                            <Trash2 className="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}
