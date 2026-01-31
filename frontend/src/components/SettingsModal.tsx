import { downloader } from "../../wailsjs/wailsjs/go/models"
import { CheckForUpdates, GetYtdlpVersion } from "../../wailsjs/wailsjs/go/main/App"
import { X, Settings as SettingsIcon, Shield, Globe, Clock, Monitor, RefreshCw } from 'lucide-react'
import { useState, useEffect } from "react"
import toast from 'react-hot-toast'

interface SettingsModalProps {
    isOpen: boolean
    onClose: () => void
    options: downloader.DownloadOptions
    onChange: (opts: downloader.DownloadOptions) => void
}

export function SettingsModal({ isOpen, onClose, options, onChange }: SettingsModalProps) {
    if (!isOpen) return null

    const [version, setVersion] = useState<string>('Checking...')
    const [updating, setUpdating] = useState(false)

    useEffect(() => {
        if (isOpen) {
            GetYtdlpVersion().then(v => setVersion(v)).catch(() => setVersion("Unknown"))
        }
    }, [isOpen])

    const handleUpdate = async (mode: string) => {
        setUpdating(true)
        const toastId = toast.loading(`Updating to ${mode}...`)
        try {
            const res = await CheckForUpdates(mode)
            toast.success(res, { id: toastId })
            const v = await GetYtdlpVersion()
            setVersion(v)
        } catch (e: any) {
            toast.error("Update failed: " + e, { id: toastId })
        } finally {
            setUpdating(false)
        }
    }

    // Local state to avoid frequent parent updates, applied on change immediately though
    const update = (field: keyof downloader.DownloadOptions, value: any) => {
        onChange({ ...options, [field]: value } as downloader.DownloadOptions)
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
            <div className="bg-white dark:bg-slate-900 rounded-xl shadow-2xl w-full max-w-lg overflow-hidden border border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200">
                {/* Header */}
                <div className="flex justify-between items-center p-6 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg text-blue-600 dark:text-blue-400">
                            <SettingsIcon size={20} />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-slate-900 dark:text-white">Settings</h2>
                            <p className="text-sm text-slate-500 dark:text-slate-400">Advanced download configuration</p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                        <X size={20} />
                    </button>
                </div>

                {/* Body */}
                <div className="p-6 space-y-6 max-h-[70vh] overflow-y-auto custom-scrollbar">

                    {/* Anti-Blocking Section */}
                    <div className="space-y-4">
                        <h3 className="flex items-center gap-2 font-semibold text-slate-900 dark:text-white border-b pb-2 border-slate-100 dark:border-slate-800">
                            <Shield size={18} className="text-emerald-500" />
                            Anti-Blocking Strategy
                        </h3>

                        {/* Cookies */}
                        <div className="space-y-3 p-4 bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-700/50">
                            <div className="flex items-start justify-between">
                                <div>
                                    <label className="font-medium text-slate-900 dark:text-slate-200">Use Browser Cookies</label>
                                    <p className="text-xs text-slate-500 mt-1">Extract cookies from your browser to verify identity (useful for age-gated or premium content).</p>
                                </div>
                                <input
                                    type="checkbox"
                                    checked={options.use_cookies}
                                    onChange={(e) => update('use_cookies', e.target.checked)}
                                    className="w-5 h-5 rounded border-slate-300 text-blue-600 focus:ring-blue-500 mt-1 cursor-pointer"
                                />
                            </div>

                            {options.use_cookies && (
                                <div className="animate-in slide-in-from-top-2 duration-200">
                                    <label className="text-sm font-medium text-slate-700 dark:text-slate-300 block mb-1.5">Browser Source</label>
                                    <select
                                        value={options.browser_name || 'chrome'}
                                        onChange={(e) => update('browser_name', e.target.value)}
                                        className="w-full p-2.5 rounded-lg bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-700 text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                                    >
                                        <option value="chrome">Google Chrome</option>
                                        <option value="firefox">Firefox</option>
                                        <option value="safari">Safari</option>
                                        <option value="edge">Microsoft Edge</option>
                                        <option value="opera">Opera</option>
                                        <option value="brave">Brave</option>
                                        <option value="vivaldi">Vivaldi</option>
                                    </select>
                                </div>
                            )}
                        </div>

                        {/* User Agent */}
                        <div>
                            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5 flex items-center gap-2">
                                <Monitor size={16} />
                                Custom User Agent
                            </label>
                            <input
                                type="text"
                                value={options.user_agent || ''}
                                onChange={(e) => update('user_agent', e.target.value)}
                                placeholder="Mozilla/5.0..."
                                className="w-full p-2.5 rounded-lg bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-700 text-sm focus:ring-2 focus:ring-blue-500 outline-none placeholder:text-slate-400"
                            />
                            <p className="text-xs text-slate-500 mt-1">Leave empty to use the default optimized User-Agent.</p>
                        </div>
                    </div>

                    {/* Network Section */}
                    <div className="space-y-4">
                        <h3 className="flex items-center gap-2 font-semibold text-slate-900 dark:text-white border-b pb-2 border-slate-100 dark:border-slate-800">
                            <Globe size={18} className="text-blue-500" />
                            Network
                        </h3>

                        {/* Proxy */}
                        <div>
                            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">Proxy URL</label>
                            <input
                                type="text"
                                value={options.proxy_url || ''}
                                onChange={(e) => update('proxy_url', e.target.value)}
                                placeholder="http://user:pass@host:port"
                                className="w-full p-2.5 rounded-lg bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-700 text-sm font-mono focus:ring-2 focus:ring-blue-500 outline-none placeholder:text-slate-400"
                            />
                        </div>

                        {/* Rate Limit */}
                        <div>
                            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5 flex items-center gap-2">
                                <Clock size={16} />
                                Rate Limit
                            </label>
                            <input
                                type="text"
                                value={options.rate_limit || ''}
                                onChange={(e) => update('rate_limit', e.target.value)}
                                placeholder="e.g. 5M (5 MB/s)"
                                className="w-full p-2.5 rounded-lg bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-700 text-sm focus:ring-2 focus:ring-blue-500 outline-none placeholder:text-slate-400"
                            />
                        </div>
                    </div>

                    {/* Impersonate */}
                    <div className="pt-2 border-t border-slate-100 dark:border-slate-800">
                        <div className="flex items-start justify-between">
                            <div>
                                <label className="font-medium text-slate-900 dark:text-slate-200">TLS Impersonation (Experimental)</label>
                                <p className="text-xs text-slate-500 mt-1">Use if getting "Connection Reset" errors. Tries to mimic a real browser's TLS fingerprint.</p>
                            </div>
                            <select
                                value={options.impersonate || ''}
                                onChange={(e) => update('impersonate', e.target.value)}
                                className="p-2 rounded-lg bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                            >
                                <option value="">Disabled</option>
                                <option value="chrome">Chrome</option>
                                <option value="firefox">Firefox</option>
                                <option value="safari">Safari</option>
                            </select>
                        </div>
                    </div>
                </div>

                {/* Updates Section */}
                <div className="space-y-4 pt-2 border-t border-slate-100 dark:border-slate-800">
                    <h3 className="flex items-center gap-2 font-semibold text-slate-900 dark:text-white">
                        <RefreshCw size={18} className="text-blue-500" />
                        Engine Updates
                    </h3>

                    <div className="p-4 bg-slate-50 dark:bg-slate-800/50 rounded-lg border border-slate-100 dark:border-slate-700/50 space-y-4">
                        <div className="flex justify-between items-center">
                            <div>
                                <p className="text-sm font-medium text-slate-700 dark:text-slate-300">yt-dlp Version</p>
                                <p className="text-xs text-slate-500 font-mono mt-1">{version}</p>
                            </div>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => handleUpdate("stable")}
                                    disabled={updating}
                                    className="px-3 py-1.5 text-xs bg-slate-200 dark:bg-slate-700 hover:bg-slate-300 dark:hover:bg-slate-600 rounded-md transition-colors"
                                >
                                    Update Stable
                                </button>
                                <button
                                    onClick={() => handleUpdate("nightly")}
                                    disabled={updating}
                                    className="px-3 py-1.5 text-xs bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300 hover:bg-indigo-200 dark:hover:bg-indigo-900/60 rounded-md transition-colors font-medium border border-indigo-200 dark:border-indigo-800"
                                >
                                    Update Nightly
                                </button>
                            </div>
                        </div>

                        <div className="flex items-center justify-between pt-2 border-t border-slate-200 dark:border-slate-700/50">
                            <div>
                                <label className="text-sm font-medium text-slate-900 dark:text-slate-200">Auto-Update (Daily)</label>
                            </div>
                            <input
                                type="checkbox"
                                checked={options.auto_update_ytdlp}
                                onChange={(e) => update('auto_update_ytdlp', e.target.checked)}
                                className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500 cursor-pointer"
                            />
                        </div>
                    </div>
                </div>
            </div>

            {/* Footer */}
            <div className="p-6 border-t border-slate-100 dark:border-slate-800 flex justify-end">
                <button
                    onClick={onClose}
                    className="px-6 py-2.5 bg-slate-900 dark:bg-white text-white dark:text-slate-900 font-medium rounded-lg hover:opacity-90 transition-opacity"
                >
                    Done
                </button>
            </div>
        </div>
    )
}
