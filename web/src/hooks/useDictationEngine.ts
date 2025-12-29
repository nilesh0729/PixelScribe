import { useState, useEffect, useRef, useCallback } from 'react';

export interface GameState {
    status: 'idle' | 'playing' | 'paused' | 'completed';
    currentText: string;
    targetText: string;
    wpm: number;
    accuracy: number;
    startTime: number | null;
    endTime: number | null;
    errors: number;
}

export function useDictationEngine(targetContent: string, language: string = 'en-US') {
    const [status, setStatus] = useState<'idle' | 'playing' | 'paused' | 'completed'>('idle');
    const [currentText, setCurrentText] = useState('');
    const [startTime, setStartTime] = useState<number | null>(null);
    const [endTime, setEndTime] = useState<number | null>(null);

    // TTS State
    const synth = window.speechSynthesis;
    const utterance = useRef<SpeechSynthesisUtterance | null>(null);

    useEffect(() => {
        utterance.current = new SpeechSynthesisUtterance(targetContent);
        utterance.current.lang = language;
        utterance.current.rate = 0.9; // Slightly slower for dictation

        // Cleanup
        return () => {
            if (synth.speaking) {
                synth.cancel();
            }
        };
    }, [targetContent, language]);

    const start = useCallback(() => {
        if (status === 'completed') {
            setCurrentText('');
            setStartTime(Date.now());
            setStatus('playing');
            synth.cancel();
            synth.speak(utterance.current!);
            return;
        }

        if (status === 'idle') {
            setStartTime(Date.now());
        }

        setStatus('playing');
        if (utterance.current) {
            // If paused, resume logic implies either resume() or re-speak from index.
            // Web Speech API pause/resume is flaky. Better to just speak.
            if (synth.paused) {
                synth.resume();
            } else if (!synth.speaking) {
                synth.speak(utterance.current);
            }
        }
    }, [status, targetContent]);

    const pause = useCallback(() => {
        setStatus('paused');
        synth.pause();
    }, []);

    const stop = useCallback(() => {
        setStatus('idle');
        synth.cancel();
        setCurrentText('');
        setStartTime(null);
    }, []);

    const handleInput = useCallback((text: string) => {
        if (status !== 'playing') return;

        setCurrentText(text);

        // Auto-complete check
        if (text.length >= targetContent.length) {
            // Simple finish condition
            setStatus('completed');
            setEndTime(Date.now());
            synth.cancel();
        }
    }, [status, targetContent]);

    // Calculations
    const durationInMinutes = ((endTime || Date.now()) - (startTime || Date.now())) / 60000;
    // Standard WPM = (characters / 5) / minutes
    const wpm = durationInMinutes > 0 ? (currentText.length / 5) / durationInMinutes : 0;

    // Simple Levenshtein or character matching could go here, 
    // currently just simple length/match ratio for "accuracy" placeholder
    const calculateAccuracy = () => {
        if (currentText.length === 0) return 100;
        let correctChars = 0;
        const len = Math.min(currentText.length, targetContent.length);
        for (let i = 0; i < len; i++) {
            if (currentText[i] === targetContent[i]) correctChars++;
        }
        return (correctChars / currentText.length) * 100;
    };

    return {
        status,
        currentText,
        start,
        pause,
        stop,
        handleInput,
        stats: {
            wpm: Math.round(wpm),
            accuracy: Math.round(calculateAccuracy()),
            progress: (currentText.length / targetContent.length) * 100
        }
    };
}
