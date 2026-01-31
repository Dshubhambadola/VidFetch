import { downloader } from "../../wailsjs/wailsjs/go/models";

interface DownloadHistoryProps {
    history: downloader.Download[];
}

export function DownloadHistory({ history }: DownloadHistoryProps) {
    if (!history || history.length === 0) {
        return <div className="text-slate-500 text-sm mt-8">No download history</div>;
    }

    return (
        <div className="mt-8 space-y-4">
            <h2 className="text-xl font-semibold text-slate-200">History</h2>
            <div className="space-y-2">
                {history.map((dl) => (
                    <div key={dl.id} className="bg-slate-800/50 p-3 rounded-lg border border-slate-700/50 flex justify-between items-center hover:bg-slate-800 transition-colors">
                        <div className="overflow-hidden">
                            <div className="font-medium text-slate-200 truncate">{dl.title || dl.url}</div>
                            <div className="text-xs text-slate-500 flex gap-2">
                                <span>{dl.file_size > 0 ? (dl.file_size / 1024 / 1024).toFixed(1) + ' MB' : ''}</span>
                                <span>•</span>
                                <span className={dl.status === 'completed' ? 'text-green-500' : 'text-red-500'}>
                                    {dl.status}
                                </span>
                            </div>
                        </div>
                        <button className="text-slate-400 hover:text-white p-2">
                            📂
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
}
