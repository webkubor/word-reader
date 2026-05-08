export const Speak = (word: string): Promise<void> => {
    return new Promise((resolve) => {
        if ('speechSynthesis' in window) {
            const msg = new SpeechSynthesisUtterance(word);
            msg.lang = 'en-US';
            msg.onend = () => resolve();
            msg.onerror = () => resolve();
            window.speechSynthesis.speak(msg);
        } else {
            console.warn("TTS is not supported in this environment");
            resolve();
        }
    });
};
