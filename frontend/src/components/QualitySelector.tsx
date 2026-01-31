import { downloader } from "../../wailsjs/wailsjs/go/models";

interface QualitySelectorProps {
    options: downloader.DownloadOptions;
    onChange: (opts: downloader.DownloadOptions) => void;
}

export function QualitySelector({ options, onChange }: QualitySelectorProps) {
    const handleChange = (key: keyof downloader.DownloadOptions, value: any) => {
        const newOptions = new downloader.DownloadOptions(options);
        (newOptions as any)[key] = value;

        // Logic to update format string based on selection
        if (key === 'format' || key === 'audio_only') {
            if (newOptions.audio_only) {
                newOptions.format = "bestaudio/best";
            } else {
                // Simplified mapping for now
                newOptions.format = "bestvideo+bestaudio/best";
            }
        }

        onChange(newOptions);
    };

    return (
        <div className="flex flex-wrap gap-4 text-sm text-slate-300">
            <div className="flex items-center gap-2">
                <label>Quality:</label>
                <select
                    className="bg-slate-800 border border-slate-700 rounded px-2 py-1 focus:outline-none focus:border-blue-500"
                    onChange={(e) => {
                        const val = e.target.value;
                        if (val === 'audio') {
                            handleChange('audio_only', true);
                        } else {
                            handleChange('audio_only', false);
                            // We could set specific format resolution strings here later using video_format options
                            // For v1 we stick to "best"
                        }
                    }}
                >
                    <option value="best">Best Available</option>
                    <option value="audio">Audio Only</option>
                </select>
            </div>

            <div className="flex items-center gap-2">
                <label className="flex items-center gap-2 cursor-pointer select-none">
                    <input
                        type="checkbox"
                        checked={options.download_subs}
                        onChange={(e) => {
                            const val = e.target.checked;
                            const newOptions = new downloader.DownloadOptions(options);
                            newOptions.download_subs = val;
                            newOptions.embed_subtitles = val; // Auto embed if downloading
                            onChange(newOptions);
                        }}
                        className="rounded bg-slate-800 border-slate-700 text-blue-600 focus:ring-offset-slate-900"
                    />
                    Download Subtitles
                </label>
            </div>
        </div>
    );
}
