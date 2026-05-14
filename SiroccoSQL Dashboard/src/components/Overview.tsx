import { Activity, Database, Zap, Server, AlertTriangle, Clock } from "lucide-react";

export default function Overview() {
    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">SiroccoSQL Control Room</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Real-time proxy observability & cluster intelligence
                </p>
            </div>

            {/* STATS GRID */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">

                <Card
                    title="Queries/sec"
                    value="12,483"
                    icon={<Zap />}
                    trend="+8.2%"
                />

                <Card
                    title="Avg Latency"
                    value="4.8 ms"
                    icon={<Clock />}
                    trend="-2.1%"
                />

                <Card
                    title="Active Nodes"
                    value="18 / 20"
                    icon={<Server />}
                    trend="Stable"
                />

                <Card
                    title="Shard Load"
                    value="Balanced"
                    icon={<Database />}
                    trend="Optimal"
                />
            </div>

            {/* MAIN GRID */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

                {/* LIVE ACTIVITY */}
                <div className="lg:col-span-2 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="font-semibold flex items-center gap-2">
                            <Activity size={18} /> Live Query Stream
                        </h2>
                        <span className="text-xs text-green-500">● live</span>
                    </div>

                    <div className="space-y-3 text-sm">
                        <Log msg="SELECT * FROM users WHERE id = 42" status="cached" />
                        <Log msg="INSERT INTO orders VALUES (...)" status="replicated" />
                        <Log msg="UPDATE sessions SET active=1" status="routed shard-3" />
                        <Log msg="DELETE FROM logs WHERE expired=1" status="optimized" />
                    </div>
                </div>

                {/* HEALTH PANEL */}
                <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
                    <h2 className="font-semibold mb-4 flex items-center gap-2">
                        <AlertTriangle size={18} /> Cluster Health
                    </h2>

                    <Health label="Primary DB" status="healthy" />
                    <Health label="Shard Nodes" status="healthy" />
                    <Health label="Replication" status="degraded" />
                    <Health label="Cache Layer" status="healthy" />
                </div>
            </div>
        </div>
    );
}

/* -------------------- COMPONENTS -------------------- */

function Card({
    title,
    value,
    icon,
    trend,
}: {
    title: string;
    value: string;
    icon: React.ReactNode;
    trend: string;
}) {
    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
            <div className="flex items-center justify-between">
                <span className="text-gray-500 dark:text-gray-400 text-sm">{title}</span>
                <div className="text-blue-500">{icon}</div>
            </div>
            <div className="text-2xl font-bold mt-2">{value}</div>
            <div className="text-xs text-gray-500 mt-1">{trend}</div>
        </div>
    );
}

function Log({ msg, status }: { msg: string; status: string }) {
    return (
        <div className="flex items-start justify-between bg-gray-50 dark:bg-gray-800 p-2 rounded-lg">
            <span className="text-gray-700 dark:text-gray-200">{msg}</span>
            <span className="text-xs text-blue-500">{status}</span>
        </div>
    );
}

function Health({
    label,
    status,
}: {
    label: string;
    status: "healthy" | "degraded" | "down";
}) {
    const color =
        status === "healthy"
            ? "text-green-500"
            : status === "degraded"
                ? "text-yellow-500"
                : "text-red-500";

    return (
        <div className="flex items-center justify-between py-2 border-b border-gray-100 dark:border-gray-800">
            <span>{label}</span>
            <span className={`text-sm font-medium ${color}`}>{status}</span>
        </div>
    );
}