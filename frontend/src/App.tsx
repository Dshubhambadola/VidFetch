import { useState, useEffect } from 'react'
import './App.css'
import { DownloadVideo, GetQueue, GetHistory } from "../wailsjs/wailsjs/go/main/App"
import { downloader } from "../wailsjs/wailsjs/go/models"
import { DownloadQueue } from "./components/DownloadQueue"
import { DownloadHistory } from "./components/DownloadHistory"

function App() {
    const [url, setUrl] = useState('')
    const [status, setStatus] = useState('')
    const [queue, setQueue] = useState<downloader.Download[]>([])
    const [history, setHistory] = useState<downloader.Download[]>([])

    const refreshData = async () => {
        try {
            const q = await GetQueue()
            setQueue(q || [])

            const h = await GetHistory()
            setHistory(h || [])
        } catch (e) {
            console.error("Failed to fetch data:", e)
        }
    }

    useEffect(() => {
        // Initial fetch
        refreshData()

        // Poll every second
        const interval = setInterval(refreshData, 1000)
        return () => clearInterval(interval)
    }, [])

    const handleDownload = async () => {
        if (!url) return
        setStatus('Requesting download...')
        try {
            const result = await DownloadVideo(url)
            setStatus(result)
            setUrl('') // Clear input
            refreshData()
        } catch (e) {
            setStatus('Error: ' + e)
        }
    }

    return (
        <div className="min-h-screen bg-slate-950 text-white p-8 font-sans">
            <div className="max-w-4xl mx-auto space-y-8">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-cyan-300 bg-clip-text text-transparent">VidFetch</h1>
                    <p className="text-slate-400">Local Privacy-First Video Downloader</p>
                </div>

                {/* Input Section */}
                <div className="bg-slate-900 p-6 rounded-xl border border-slate-800 shadow-xl">
                    <div className="flex gap-3">
                        <input
                            type="text"
                            value={url}
                            onChange={(e) => setUrl(e.target.value)}
                            placeholder="Paste video URL here (YouTube, Vimeo, etc.)"
                            className="flex-1 p-3 rounded-lg bg-slate-950 border border-slate-700 text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                        />
                        <button
                            onClick={handleDownload}
                            disabled={!url}
                            className="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed px-6 py-3 rounded-lg font-semibold transition-colors shadow-lg shadow-blue-900/20"
                        >
                            Download
                        </button>
                    </div>
                    {status && <div className="mt-3 text-sm text-blue-300/80 font-mono">{status}</div>}
                </div>

                {/* Content Area */}
                <div className="grid md:grid-cols-2 gap-8">
                    {/* Queue Column */}
                    <div className="bg-slate-900/50 p-6 rounded-xl border border-slate-800/50">
                        <DownloadQueue downloads={queue} />
                    </div>

                    {/* History Column */}
                    <div className="bg-slate-900/50 p-6 rounded-xl border border-slate-800/50">
                        <DownloadHistory history={history} />
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
