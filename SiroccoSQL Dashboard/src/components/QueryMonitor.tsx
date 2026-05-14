import { useState } from "react";
import { Play, Filter, Search, Clock, Database } from "lucide-react";

type QueryStatus = "cached" | "routed" | "slow" | "error";

type Query = {
    id: string;
    sql: string;
    status: QueryStatus;
    latency: number;
    shard: string;
    timestamp: string;
};

const initialQueries: Query[] = [
    {
        id: "q1",
        sql: "SELECT * FROM users WHERE id = 42",
        status: "cached",
        latency: 1.2,
        shard: "shard-1",
        timestamp: "12:01:22",
    },
    {
        id: "q2",
        sql: "UPDATE orders SET status='paid'",
        status: "routed",
        latency: 6.8,
        shard: "shard-3",
        timestamp: "12:01:25",
    },
    {
        id: "q3",
        sql: "DELETE FROM logs WHERE expired=1",
        status: "slow",
        latency: 120.4,
        shard: "shard-2",
        timestamp: "12:01:30",
    },
];

export default function QueryMonitor() {
    const [queries] = useState(initialQueries);
    const [search, setSearch] = useState("");

    const filtered = queries.filter((q) =>
        q.sql.toLowerCase().includes(search.toLowerCase())
    );

    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Query Monitor</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Live proxy routing, latency tracking & execution trace
                </p>
            </div>

            {/* TOOLBAR */}
            <div className="flex flex-col md:flex-row gap-3 md:items-center justify-between">

                <div className="flex items-center gap-2 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 px-3 py-2 rounded-lg w-full md:w-1/2">
                    <Search size={16} className="text-gray-400" />
                    <input
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Search SQL queries..."
                        className="bg-transparent outline-none text-sm w-full"
                    />
                </div>

                <div className="flex gap-2">
                    <button className="flex items-center gap-2 px-3 py-2 rounded-lg bg-blue-500 text-white text-sm">
                        <Play size={16} /> Live Mode
                    </button>

                    <button className="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 text-sm">
                        <Filter size={16} /> Filter
                    </button>
                </div>
            </div>

            {/* QUERY LIST */}
            <div className="space-y-3">

                {filtered.map((q) => (
                    <QueryRow key={q.id} query={q} />
                ))}

            </div>
        </div>
    );
}

/* ---------------- QUERY ROW ---------------- */

function QueryRow({ query }: { query: Query }) {
    const color =
        query.status === "cached"
            ? "text-green-500"
            : query.status === "routed"
                ? "text-blue-500"
                : query.status === "slow"
                    ? "text-yellow-500"
                    : "text-red-500";

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4 hover:shadow transition">

            {/* TOP */}
            <div className="flex items-start justify-between gap-4">

                <div className="flex items-start gap-3">
                    <Database className="text-blue-500 mt-1" size={18} />

                    <div>
                        <p className="font-mono text-sm">{query.sql}</p>

                        <div className="flex items-center gap-3 mt-2 text-xs text-gray-500 dark:text-gray-400">
                            <span className={color}>{query.status}</span>
                            <span>Shard: {query.shard}</span>
                            <span className="flex items-center gap-1">
                                <Clock size={12} /> {query.latency}ms
                            </span>
                            <span>{query.timestamp}</span>
                        </div>
                    </div>
                </div>

                {/* ACTION */}
                <button className="text-xs px-3 py-1 rounded bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    Explain
                </button>
            </div>
        </div>
    );
}