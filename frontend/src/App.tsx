// Erasmo Cardoso - Dev
import { useState } from 'react';
import './App.css';
import { CheckSites } from "../wailsjs/go/main/App";

interface SiteStatus {
    url: string;
    status: string;
    message: string;
    isOnline: boolean;
}

function App() {
    const [sites, setSites] = useState<SiteStatus[]>([]);
    const [loading, setLoading] = useState(false);

    async function checkSites() {
        if (loading) return;

        setLoading(true);
        setSites([]);

        try {
            const results = await CheckSites();
            setSites(results);
        } catch (error) {
            console.error('Erro ao verificar sites:', error);
        } finally {
            setLoading(false);
        }
    }

    return (
        <div id="App">
            <div className="container">
                <h1 className="title">Monitor de Status</h1>
                <p className="subtitle">Verificação de disponibilidade de sites</p>

                <button
                    className="check-btn"
                    onClick={checkSites}
                    disabled={loading}
                >
                    {loading ? 'Verificando...' : 'Verificar Sites'}
                </button>

                {sites.length > 0 && (
                    <div className="sites-grid">
                        {sites.map((site, idx) => (
                            <div
                                key={idx}
                                className={`site-card ${site.isOnline ? 'online' : 'offline'}`}
                            >
                                <div className="site-header">
                                    <div className={`status-indicator ${site.isOnline ? 'online' : 'offline'}`}></div>
                                    <span className="status-text">{site.message}</span>
                                </div>
                                <div className="site-url">{site.url}</div>
                                <div className="site-status">{site.status}</div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default App;
