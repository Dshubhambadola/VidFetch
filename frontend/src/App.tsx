import { useState } from 'react'
import './App.css'
import { DownloadVideo } from "../wailsjs/wailsjs/go/main/App"

function App() {
    const [url, setUrl] = useState('')
    const [status, setStatus] = useState('')

    const handleDownload = async () => {
        if (!url) return
        setStatus('Requesting download...')
        try {
            const result = await DownloadVideo(url)
            setStatus(result)
        } catch (e) {
            setStatus('Error: ' + e)
        }
    }

    return (
        <div className="min-h-screen bg-slate-900 text-white p-8">
            <h1 className="text-3xl font-bold mb-6">VidFetch</h1>
            <div className="flex gap-2">
                <input
                    type="text"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    placeholder="Enter video URL"
                    className="flex-1 p-2 rounded bg-slate-800 border border-slate-700 text-white"
                />
                <button
                    onClick={handleDownload}
                    className="bg-blue-600 hover:bg-blue-500 px-4 py-2 rounded font-medium"
                >
                    Download
                </button>
            </div>
            {status && <div className="mt-4 text-slate-300">{status}</div>}
        </div>
    )
}

export default App
