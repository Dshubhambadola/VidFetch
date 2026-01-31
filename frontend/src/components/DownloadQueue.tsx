import { downloader } from "../../wailsjs/wailsjs/go/models";

interface DownloadQueueProps {
    downloads: downloader.Download[];
}

export function DownloadQueue({ downloads }: DownloadQueueProps) {
    if (downloads.length === 0) {
        return <div className="text-slate-500 text-sm">No active downloads</div>;
    }

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold text-slate-200">Active Downloads</h2>
            <div className="grid gap-3">
                {downloads.map((dl) => (
                    <div key={dl.id} className="bg-slate-800 p-4 rounded-lg border border-slate-700">
                        <div className="flex justify-between items-start mb-2">
                            <div>
                                <h3 className="font-medium text-white truncate max-w-md">{dl.title || dl.url}</h3>
                                <div className="text-xs text-slate-400 mt-1">
                                    {dl.status} • {dl.speed} • ETA: {dl.eta}
                                </div>
                            </div>
                            <div className="text-xs bg-blue-900 text-blue-200 px-2 py-1 rounded">
                                {dl.quality}
                            </div>
                        </div>

                        {/* Progress Bar */}
                        <div className="w-full bg-slate-700 rounded-full h-2.5">
                            <div
                                className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
                                style={{ width: `${(dl.progress || 0) * 100}%` }}
                            ></div>
                        </div>

                        <div className="flex justify-between mt-2 text-xs text-slate-400">
                            <span>{(dl.progress * 100).toFixed(1)}%</span>
                            <span>{dl.subtitle_langs?.length || 0} subs</span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
