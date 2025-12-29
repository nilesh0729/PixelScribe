import { useState, useEffect, useRef, useCallback } from 'react';
import { ttsService } from '../services/tts';

export type DictationPhase = 'loading' | 'listening' | 'typing' | 'completed';

export function useDictationEngine(targetContent: string) {
    const [phase, setPhase] = useState<DictationPhase>('loading');
    const [currentText, setCurrentText] = useState('');
    const [startTime, setStartTime] = useState<number | null>(null);
    const [endTime, setEndTime] = useState<number | null>(null);

    // Audio State
    const audioRef = useRef<HTMLAudioElement | null>(null);
    const [audioUrl, setAudioUrl] = useState<string | null>(null);
    const [isPlaying, setIsPlaying] = useState(false);

    // Load Audio
    useEffect(() => {
        let active = true;
        const loadAudio = async () => {
            if (!targetContent) return;
            setPhase('loading');
            try {
                const blob = await ttsService.generateAudio(targetContent);
                if (active) {
                    const url = URL.createObjectURL(blob);
                    setAudioUrl(url);
                    audioRef.current = new Audio(url);

                    audioRef.current.onended = () => setIsPlaying(false);
                    audioRef.current.onplay = () => setIsPlaying(true);
                    audioRef.current.onpause = () => setIsPlaying(false);

                    setPhase('listening');
                }
            } catch (error) {
                console.error("Failed to load TTS", error);
                // Even if audio fails, allow typing? Or show error?
                setPhase('listening');
            }
        };
        loadAudio();

        return () => {
            active = false;
            if (audioUrl) URL.revokeObjectURL(audioUrl);
            if (audioRef.current) {
                audioRef.current.pause();
                audioRef.current = null;
            }
        };
    }, [targetContent]);

    // Audio Controls
    const playAudio = useCallback(() => {
        audioRef.current?.play();
    }, []);

    const pauseAudio = useCallback(() => {
        audioRef.current?.pause();
    }, []);

    const replayAudio = useCallback(() => {
        if (audioRef.current) {
            audioRef.current.currentTime = 0;
            audioRef.current.play();
        }
    }, []);

    // Workflow Controls
    const startTyping = useCallback(() => {
        if (audioRef.current) {
            audioRef.current.pause();
            audioRef.current.currentTime = 0;
        }
        setPhase('typing');
        setStartTime(Date.now());
        setCurrentText('');
    }, []);

    const handleInput = useCallback((text: string) => {
        setCurrentText(text);
    }, []);

    const complete = useCallback(() => {
        setEndTime(Date.now());
        setPhase('completed');
    }, []);

    // Simple time tracking (seconds)
    const timeSpentSeconds = startTime ? ((endTime || Date.now()) - startTime) / 1000 : 0;

    return {
        phase,
        currentText,
        audioState: {
            isPlaying,
            play: playAudio,
            pause: pauseAudio,
            replay: replayAudio
        },
        controls: {
            startTyping,
            handleInput,
            complete
        },
        stats: {
            timeSpent: timeSpentSeconds
        }
    };
}
