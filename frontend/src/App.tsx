import { useState, useEffect, useCallback, useRef } from 'react';
import styles from './App.module.css';
import {
    Translate,
    SaveWord,
    GetWordBook,
    RemoveWord,
    MarkReviewed,
    UpdateNote,
    GetSetting,
    SetSetting,
    GetArticles,
    SaveArticle,
    DeleteArticle
} from '../../wailsjs/go/main/App';
import { Speak } from './utils/tts';
import type { TranslateResult, WordBookItem } from './types.js';
import logo from './assets/images/logo-universal.png';

type Page = 'reading' | 'wordbook' | 'articles' | 'settings';

function App() {
    // Global state
    const [currentPage, setCurrentPage] = useState<Page>('reading');
    const [wordBook, setWordBook] = useState<WordBookItem[]>([]);
    const hoverTimeoutRef = useRef<number | null>(null);

    // Reading Page state
    const [text, setText] = useState('');
    const [articleTitle, setArticleTitle] = useState('');
    const [readingContent, setReadingContent] = useState<string[]>([]);
    
    // Translation Popup state
    const [popup, setPopup] = useState<{ data: TranslateResult | null; position: { x: number; y: number }, loading: boolean } | null>(null);

    // Settings state
    const [apiKey, setApiKey] = useState('');
    const [saveStatus, setSaveStatus] = useState('');

    // Articles state
    const [articles, setArticles] = useState<any[]>([]);

    // --- Debounce Hook ---
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    function useDebounce<T extends (...args: any[]) => void>(callback: T, delay: number) {
        const timeoutRef = useRef<number>();
        useEffect(() => {
            return () => clearTimeout(timeoutRef.current);
        }, []);
        return (...args: Parameters<T>) => {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = window.setTimeout(() => {
                callback(...args);
            }, delay);
        };
    }

    const fetchWordBook = useCallback(async () => {
        try {
            const words = await GetWordBook() || [];
            setWordBook(words);
        } catch (error) {
            console.error("Error fetching word book:", error);
        }
    }, []);

    const fetchSettings = useCallback(async () => {
        try {
            const key = await GetSetting("deepl_api_key");
            setApiKey(key || '');
        } catch(e) {
            console.error("Failed to load settings", e);
        }
    }, []);

    const fetchArticles = useCallback(async () => {
        try {
            const list = await GetArticles() || [];
            setArticles(list);
        } catch(e) {
            console.error("Failed to fetch articles", e);
        }
    }, []);

    useEffect(() => {
        fetchWordBook();
        fetchSettings();
    }, [fetchWordBook, fetchSettings]);

    useEffect(() => {
        if (currentPage === 'wordbook') fetchWordBook();
        if (currentPage === 'settings') fetchSettings();
        if (currentPage === 'articles') fetchArticles();
    }, [currentPage, fetchWordBook, fetchSettings, fetchArticles]);

    const handleUpdateNote = useDebounce(async (word: string, note: string) => {
        try {
            await UpdateNote(word, note);
        } catch(e) {
            console.error("Failed to update note:", e);
        }
    }, 500);

    const handleMarkReviewed = async (word: string) => {
        try {
            await MarkReviewed(word);
            setWordBook(prev => prev.map(item => item.word === word ? { ...item, reviewed: item.reviewed + 1 } : item));
        } catch (e) {
            console.error("Failed to mark as reviewed:", e);
        }
    };
    
    const handleRemoveWord = async (word: string) => {
        if (!confirm(`Are you sure you want to remove "${word}"?`)) return;
        try {
            await RemoveWord(word);
            setWordBook(prev => prev.filter(item => item.word !== word));
        } catch (e) {
            console.error("Failed to remove word:", e);
        }
    };

    const handleStartReading = () => {
        const segments = text.split(/([a-zA-Z0-9'-]+)/g);
        setReadingContent(segments);
    };

    const handleSaveArticle = async () => {
        if (!text || !articleTitle) {
            alert("Title and content cannot be empty.");
            return;
        }
        try {
            await SaveArticle(articleTitle, text);
            alert("Article saved!");
            setArticleTitle('');
        } catch(e) {
            alert("Failed to save article: " + e);
        }
    };

    const handleDeleteArticle = async (id: number) => {
        if (!confirm("Delete this article?")) return;
        try {
            await DeleteArticle(id);
            fetchArticles();
        } catch(e) {
            alert("Failed to delete: " + e);
        }
    };

    const handleLoadArticle = (content: string) => {
        setText(content);
        setCurrentPage('reading');
        setTimeout(() => {
            const segments = content.split(/([a-zA-Z0-9'-]+)/g);
            setReadingContent(segments);
        }, 100);
    };

    const handleWordHover = (e: React.MouseEvent<HTMLSpanElement>, word: string) => {
        const cleanWord = word.replace(/[^a-zA-Z'-]/g, '').toLowerCase();
        if (!cleanWord || cleanWord.length < 2) return;

        if (hoverTimeoutRef.current) clearTimeout(hoverTimeoutRef.current);

        const target = e.currentTarget;
        const rect = target.getBoundingClientRect();

        hoverTimeoutRef.current = window.setTimeout(async () => {
            setPopup({ data: null, position: { x: rect.left, y: rect.bottom + 5 }, loading: true });
            try {
                const result = await Translate(cleanWord);
                if (result) {
                    Speak(cleanWord).catch(err => console.error("Speak error:", err));
                    setPopup({
                        data: result,
                        position: { x: rect.left, y: rect.bottom + 5 },
                        loading: false
                    });
                } else {
                    setPopup(null);
                }
            } catch (error: any) {
                console.error(`Error processing word ${cleanWord}:`, error);
                setPopup(null);
            }
        }, 300);
    };

    const handleMouseLeave = () => {
        if (hoverTimeoutRef.current) clearTimeout(hoverTimeoutRef.current);
        setPopup(null);
    };

    const handleSaveWord = async (wordData: TranslateResult) => {
        if (!wordData?.word) return;
        try {
            await SaveWord(wordData.word, wordData.translation, wordData.phonetic || '');
            fetchWordBook();
            setPopup(null);
        } catch (error) {
            console.error(`Error saving word ${wordData.word}:`, error);
        }
    };
    
    const isWordSaved = (word: string) => wordBook.some(item => item.word.toLowerCase() === word.toLowerCase());

    const saveSettings = async () => {
        try {
            await SetSetting("deepl_api_key", apiKey);
            setSaveStatus("Settings saved successfully!");
            setTimeout(() => setSaveStatus(''), 2000);
        } catch(e) {
            setSaveStatus("Failed to save settings.");
        }
    };

    const renderReadingPage = () => (
        <div className={styles.readingPage}>
            <div className={styles.readingForm}>
                <input 
                    className={styles.titleInput} 
                    value={articleTitle} 
                    onChange={e => setArticleTitle(e.target.value)} 
                    placeholder="Article Title (Optional)" 
                />
                <textarea
                    className={styles.textarea}
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    placeholder="Paste English text here..."
                />
                <div style={{display: 'flex', gap: '10px'}}>
                    <button className={styles.button} onClick={handleStartReading}>Start Reading</button>
                    <button className={`${styles.button} ${styles.secondaryButton}`} style={{background: 'var(--bg-3)', color: 'var(--text)'}} onClick={handleSaveArticle}>Save Article</button>
                </div>
            </div>
            {readingContent.length > 0 && (
                <div className={styles.readingArea} onMouseLeave={handleMouseLeave}>
                    {readingContent.map((segment, index) => {
                        const isWord = /[a-zA-Z'-]+/.test(segment);
                        if (isWord) {
                            return (
                                <span
                                    key={index}
                                    className={styles.wordSpan}
                                    onMouseEnter={(e) => handleWordHover(e, segment)}
                                >
                                    {segment}
                                </span>
                            );
                        }
                        return <span key={index}>{segment}</span>;
                    })}
                </div>
            )}
        </div>
    );

    const renderWordBookPage = () => (
        <div className={styles.wordBookPage}>
            <div className={styles.wordBookHeader}>
                <h2 style={{color: 'var(--text)'}}>Word Book ({wordBook.length})</h2>
                <button className={styles.actionButton} onClick={fetchWordBook}>Refresh</button>
            </div>
            <ul className={styles.wordList}>
                {wordBook.map((item) => (
                    <li key={item.word} className={styles.wordItem}>
                        <div className={styles.wordInfo} onClick={() => Speak(item.word)}>
                            <div className={styles.word}>{item.word}</div>
                            <div className={styles.translationSmall}>{item.translation.split(',')[0]}</div>
                        </div>
                        <input
                            type="text"
                            defaultValue={item.note}
                            onChange={(e) => handleUpdateNote(item.word, e.target.value)}
                            className={styles.noteInput}
                            placeholder="Add a note..."
                        />
                         <div className={styles.wordActions}>
                            <button className={styles.actionButton} onClick={() => handleMarkReviewed(item.word)}>
                                Review ({item.reviewed})
                            </button>
                            <button className={`${styles.actionButton} ${styles.removeButton}`} onClick={() => handleRemoveWord(item.word)}>
                                Remove
                            </button>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );

    const renderArticlesPage = () => (
        <div className={styles.wordBookPage}>
            <div className={styles.wordBookHeader}>
                <h2 style={{color: 'var(--text)'}}>Saved Articles ({articles.length})</h2>
                <button className={styles.actionButton} onClick={fetchArticles}>Refresh</button>
            </div>
            <ul className={styles.wordList}>
                {articles.map((item) => (
                    <li key={item.id} className={styles.wordItem} style={{flexDirection: 'column', alignItems: 'flex-start'}}>
                        <div style={{display: 'flex', justifyContent: 'space-between', width: '100%'}}>
                            <div style={{fontWeight: 'bold', cursor: 'pointer', color: 'var(--accent)'}} onClick={() => handleLoadArticle(item.content)}>
                                {item.title || 'Untitled'}
                            </div>
                            <button className={`${styles.actionButton} ${styles.removeButton}`} onClick={() => handleDeleteArticle(item.id)}>
                                Delete
                            </button>
                        </div>
                        <div style={{fontSize: '0.8em', color: 'var(--text-3)', marginTop: '5px'}}>
                            Saved at: {new Date(item.created_at).toLocaleString()}
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );

    const renderSettingsPage = () => (
        <div className={styles.settingsPage}>
            <h2 className={styles.settingsTitle}>Settings</h2>
            <div className={styles.settingGroup}>
                <label className={styles.settingLabel}>DeepL API Key (Free API):</label>
                <input 
                    type="password" 
                    className={styles.settingInput}
                    value={apiKey} 
                    onChange={e => setApiKey(e.target.value)} 
                    placeholder="Enter your DeepL API key"
                />
            </div>
            <button className={styles.button} onClick={saveSettings}>Save Settings</button>
            {saveStatus && <p style={{color: 'var(--green)', marginTop: '15px', fontSize: '14px'}}>{saveStatus}</p>}
        </div>
    );

    return (
        <div className={styles.app}>
            <nav className={styles.nav}>
                <div className={styles.logoSection}>
                    <img src={logo} alt="Logo" className={styles.logo} />
                    <span className={styles.brandName}>Word Reader</span>
                </div>
                <div className={styles.segmentedControl}>
                    <button
                        onClick={() => setCurrentPage('reading')}
                        className={`${styles.navButton} ${currentPage === 'reading' ? styles.active : ''}`}
                    >
                        Reading
                    </button>
                    <button
                        onClick={() => setCurrentPage('articles')}
                        className={`${styles.navButton} ${currentPage === 'articles' ? styles.active : ''}`}
                    >
                        Articles
                    </button>
                    <button
                        onClick={() => setCurrentPage('wordbook')}
                        className={`${styles.navButton} ${currentPage === 'wordbook' ? styles.active : ''}`}
                    >
                        Word Book
                    </button>
                    <button
                        onClick={() => setCurrentPage('settings')}
                        className={`${styles.navButton} ${currentPage === 'settings' ? styles.active : ''}`}
                    >
                        Settings
                    </button>
                </div>
            </nav>

            <main className={styles.content}>
                {currentPage === 'reading' && renderReadingPage()}
                {currentPage === 'articles' && renderArticlesPage()}
                {currentPage === 'wordbook' && renderWordBookPage()}
                {currentPage === 'settings' && renderSettingsPage()}
            </main>

            {popup && (
                <div className={styles.popup} style={{ top: popup.position.y, left: popup.position.x }}>
                    {popup.loading ? (
                        <div style={{color: 'var(--text-3)', fontSize: '14px'}}>Translating...</div>
                    ) : popup.data ? (
                        <>
                            <div className={styles.popupHeader}>
                                <div className={styles.popupWord}>{popup.data.word}</div>
                                <span className={styles.phonetic}>{popup.data.phonetic}</span>
                                <button
                                    className={`${styles.saveButton} ${isWordSaved(popup.data.word) ? styles.saved : ''}`}
                                    onClick={() => handleSaveWord(popup.data!)}
                                    disabled={isWordSaved(popup.data.word)}
                                >
                                    {isWordSaved(popup.data.word) ? '✓' : '★'}
                                </button>
                            </div>
                            <p className={styles.translation}>{popup.data.translation}</p>
                        </>
                    ) : null}
                </div>
            )}
        </div>
    );
}

export default App;
