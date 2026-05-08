// This is a mock implementation file for the Wails auto-generated module.
// It allows TypeScript to compile by providing concrete implementations.

// Helper to create a mock function that accepts any arguments but does nothing.
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const notImplemented = (..._args: any[]): Promise<any> => Promise.reject("Wails function not implemented in dev mode without running 'wails dev'");

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const Translate = (word: string): Promise<any> => notImplemented(word);
export const SaveWord = (word: object): Promise<void> => notImplemented(word);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const GetWordBook = (): Promise<any[]> => Promise.resolve([]);
export const UpdateNote = (word: string, note: string): Promise<void> => notImplemented(word, note);
export const MarkReviewed = (word: string): Promise<void> => notImplemented(word);
export const RemoveWord = (word: string): Promise<void> => notImplemented(word);
export const SaveArticle = (title: string, content: string): Promise<void> => notImplemented(title, content);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const GetArticles = (): Promise<any[]> => Promise.resolve([]);
export const DeleteArticle = (id: number): Promise<void> => notImplemented(id);
export const GetSetting = (key: string): Promise<string> => notImplemented(key);
export const SetSetting = (key: string, value: string): Promise<void> => notImplemented(key, value);

// Replace Wails Speak with Web Speech API
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
