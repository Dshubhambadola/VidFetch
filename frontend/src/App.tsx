import { useState, useEffect } from 'react'
import './App.css'
import { DownloadVideoWithOptions, GetQueue, GetHistory } from "../wailsjs/wailsjs/go/main/App"
import { EventsOn } from "../wailsjs/wailsjs/runtime"
import { downloader } from "../wailsjs/wailsjs/go/models"
import { DownloadQueue } from "./components/DownloadQueue"
import { DownloadHistory } from "./components/DownloadHistory"
import { QualitySelector } from "./components/QualitySelector"
import { BatchDownload } from "./components/BatchDownload"
import { ThemeToggle } from "./components/ThemeToggle"
import { SettingsModal } from "./components/SettingsModal"
import { Settings } from "lucide-react"
import toast, { Toaster } from 'react-hot-toast';

function App() {
    const [url, setUrl] = useState('')
    const [status, setStatus] = useState('')
    const [queue, setQueue] = useState<downloader.Download[]>([])
    const [history, setHistory] = useState<downloader.Download[]>([])
    const [activeTab, setActiveTab] = useState<'single' | 'batch'>('single')
    const [isSettingsOpen, setIsSettingsOpen] = useState(false)

    // Download Options State
    const [options, setOptions] = useState<downloader.DownloadOptions>(new downloader.DownloadOptions({
        format: "bestvideo+bestaudio/best",
        download_subs: true,
        embed_subtitles: true,
        subtitle_langs: ["all"],
        output_dir: "./downloads",
        output_template: "%(title)s.%(ext)s"
    }))

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

        // Listen for completion events
        const unsub = EventsOn("download-complete", (dl: downloader.Download) => {
            toast.success(`Download complete: ${dl.title || 'Video'}`, {
                duration: 4000,
                position: 'bottom-right',
                style: {
                    background: '#1e293b',
                    color: '#fff',
                    border: '1px solid #334155',
                },
            });
            refreshData();
        });

        return () => {
            clearInterval(interval);
            unsub();
        }
    }, [])

    const handleDownload = async () => {
        if (!url) return
        setStatus('Requesting download...')
        try {
            const result = await DownloadVideoWithOptions(url, options)
            setStatus(result)
            setUrl('') // Clear input
            refreshData()
        } catch (e) {
            setStatus('Error: ' + e)
        }
    }

    const handleBatchDownload = async (urls: string[]) => {
        setStatus(`Queueing ${urls.length} downloads...`)
        for (const u of urls) {
            try {
                await DownloadVideoWithOptions(u, options)
            } catch (e) {
                console.error(`Failed to queue ${u}`, e)
            }
        }
        setStatus(`Queued ${urls.length} downloads.`)
        refreshData()
    }

    return (
        <div className="min-h-screen bg-slate-950 text-white p-8 font-sans transition-colors duration-300 dark:bg-slate-950 dark:text-white bg-gray-100 text-slate-900">
            <Toaster />
            <div className="max-w-4xl mx-auto space-y-8">
                {/* Header */}
                <div className="flex justify-between items-start">
                    <div>
                        <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-cyan-500 dark:from-blue-400 dark:to-cyan-300 bg-clip-text text-transparent">VidFetch</h1>
                        <p className="text-slate-500 dark:text-slate-400">Local Privacy-First Video Downloader</p>
                    </div>
                    <div className="flex gap-3">
                        <button
                            onClick={() => setIsSettingsOpen(true)}
                            className="p-2.5 rounded-lg bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-500 hover:text-blue-600 dark:text-slate-400 dark:hover:text-blue-400 transition-colors shadow-sm"
                            title="Settings"
                        >
                            <Settings size={20} />
                        </button>
                        <ThemeToggle />
                    </div>
                </div>

                <SettingsModal
                    isOpen={isSettingsOpen}
                    onClose={() => setIsSettingsOpen(false)}
                    options={options}
                    onChange={setOptions}
                />

                {/* Tab Switcher */}
                <div className="flex gap-4 border-b border-slate-300 dark:border-slate-800 pb-2">
                    <button
                        onClick={() => setActiveTab('single')}
                        className={`pb-2 px-2 transition-colors ${activeTab === 'single' ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400 font-medium' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}`}
                    >
                        Single URL
                    </button>
                    <button
                        onClick={() => setActiveTab('batch')}
                        className={`pb-2 px-2 transition-colors ${activeTab === 'batch' ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400 font-medium' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}`}
                    >
                        Batch Download
                    </button>
                </div>

                {/* Input Section */}
                {activeTab === 'single' ? (
                    <div className="bg-white dark:bg-slate-900 p-6 rounded-xl border border-slate-200 dark:border-slate-800 shadow-xl space-y-4">
                        <div className="flex gap-3">
                            <input
                                type="text"
                                value={url}
                                onChange={(e) => setUrl(e.target.value)}
                                placeholder="Paste video URL here (YouTube, Vimeo, etc.)"
                                className="flex-1 p-3 rounded-lg bg-gray-50 dark:bg-slate-950 border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-white placeholder-slate-500 dark:placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                            />
                            <button
                                onClick={handleDownload}
                                disabled={!url}
                                className="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed px-6 py-3 rounded-lg font-semibold text-white transition-colors shadow-lg shadow-blue-500/20"
                            >
                                Download
                            </button>
                        </div>

                        {/* Options */}
                        <QualitySelector options={options} onChange={setOptions} />

                        {status && <div className="mt-3 text-sm text-blue-600 dark:text-blue-300/80 font-mono">{status}</div>}
                    </div>
                ) : (
                    <BatchDownload onDownload={handleBatchDownload} />
                )}


                {/* Content Area */}
                <div className="grid md:grid-cols-2 gap-8">
                    {/* Queue Column */}
                    <div className="bg-white/50 dark:bg-slate-900/50 p-6 rounded-xl border border-slate-200 dark:border-slate-800/50">
                        <DownloadQueue downloads={queue} />
                    </div>

                    {/* History Column */}
                    <div className="bg-white/50 dark:bg-slate-900/50 p-6 rounded-xl border border-slate-200 dark:border-slate-800/50">
                        <DownloadHistory history={history} />
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
