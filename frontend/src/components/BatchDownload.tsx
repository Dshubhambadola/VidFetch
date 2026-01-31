import { useState } from "react";

interface BatchDownloadProps {
    onDownload: (urls: string[]) => void;
}

export function BatchDownload({ onDownload }: BatchDownloadProps) {
    const [text, setText] = useState("");

    const handleDownload = () => {
        const urls = text.split(/\r?\n/).map(u => u.trim()).filter(u => u.length > 0);
        if (urls.length > 0) {
            onDownload(urls);
            setText("");
        }
    };

    return (
        <div className="bg-slate-900 p-6 rounded-xl border border-slate-800 shadow-xl space-y-4">
            <h2 className="text-xl font-semibold text-slate-200">Batch Download</h2>
            <textarea
                className="w-full h-32 bg-slate-950 border border-slate-700 rounded-lg p-3 text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Paste multiple URLs here (one per line)..."
                value={text}
                onChange={(e) => setText(e.target.value)}
            />
            <div className="flex justify-end">
                <button
                    onClick={handleDownload}
                    disabled={!text.trim()}
                    className="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed px-6 py-2 rounded-lg font-semibold transition-colors"
                >
                    Download All
                </button>
            </div>
        </div>
    );
}
