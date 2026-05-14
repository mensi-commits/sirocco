import { useState } from "react";
import { Search, Terminal, AlertCircle, Info, ShieldAlert } from "lucide-react";

type LogLevel = "info" | "warning" | "error";

type Log = {
    id: string;
    message: string;
    level: LogLevel;
    source: string;
    time: string;
};

const initialLogs: Log[] = [
    {
        id: "l1",
        message: "Query routed to shard-2 successfully",
        level: "info",
        source: "router",
        time: "12:01:02",
    },
    {
        id: "l2",
        message: "High latency detected on node-03",
        level: "warning",
        source: "monitor",
        time: "12:01:10",
    },
    {
        id: "l3",
        message: "Failed authentication attempt from 10.0.0.4",
        level: "error",
        source: "security",
        time: "12:01:20",
    },
    {
        id: "l4",
        message: "Cache warmed for user session queries",
        level: "info",
        source: "cache",
        time: "12:01:25",
    },
];

export default function Logs() {
    const [logs] = useState(initialLogs);
    const [search, setSearch] = useState("");

    const filtered = logs.filter((l) =>
        l.message.toLowerCase().includes(search.toLowerCase())
    );

    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Logs Explorer</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Real-time system logs, routing traces & security events
                </p>
            </div>

            {/* SEARCH BAR */}
            <div className="flex items-center gap-2 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 px-3 py-2 rounded-lg w-full md:w-1/2">
                <Search size={16} className="text-gray-400" />
                <input
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Search logs..."
                    className="bg-transparent outline-none text-sm w-full"
                />
            </div>

            {/* LOG STREAM */}
            <div className="space-y-2">

                {filtered.map((log) => (
                    <LogRow key={log.id} log={log} />
                ))}

            </div>
        </div>
    );
}

/* ---------------- LOG ROW ---------------- */

function LogRow({ log }: { log: Log }) {
    const color =
        log.level === "info"
            ? "text-blue-500"
            : log.level === "warning"
                ? "text-yellow-500"
                : "text-red-500";

    const icon =
        log.level === "info" ? (
            <Info size={16} />
        ) : log.level === "warning" ? (
            <ShieldAlert size={16} />
        ) : (
            <AlertCircle size={16} />
        );

    return (
        <div className="flex items-start justify-between bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-lg p-3 hover:shadow transition">

            {/* LEFT */}
            <div className="flex gap-3">

                <div className={color}>{icon}</div>

                <div>
                    <p className="text-sm">{log.message}</p>

                    <div className="text-xs text-gray-500 mt-1 flex gap-3">
                        <span>{log.source}</span>
                        <span>{log.time}</span>
                        <span className={color}>{log.level}</span>
                    </div>
                </div>
            </div>

            {/* RIGHT */}
            <Terminal size={16} className="text-gray-400" />
        </div>
    );
}